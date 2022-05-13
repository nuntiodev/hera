package handler

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *defaultHandler) Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	// init variables
	var (
		err             error
		namespaceConfig *go_block.Config
		get             *go_block.User
		accessToken     string
		refreshToken    string
		refreshClaims   *go_block.CustomClaims
		accessClaims    *go_block.CustomClaims
	)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// get config
		config, err := h.repository.Config(ctx, req.Namespace)
		if err != nil {
			return err
		}
		namespaceConfig, err = config.GetNamespaceConfig(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		resp, err := h.Get(ctx, req)
		if err != nil {
			return fmt.Errorf("could not get user with err: %v", err)
		}
		if resp.User.Password == "" {
			return errors.New("please update the user with a non-empty password")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
			return err
		}
		get = resp.User
		return nil
	})
	if err := g.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// if email validation is required and email is not verified; return error
	if namespaceConfig.RequireEmailVerification && get.EmailIsVerified == false {
		// check if we should send a new email
		if get.VerificationEmailExpiresAt.AsTime().Sub(time.Now()).Seconds() <= 0 {
			// sent new email
			verificationEmail, err := h.SendVerificationEmail(ctx, req)
			if err != nil {
				return &go_block.UserResponse{}, fmt.Errorf("could not send email with err: %v", err)
			}
			get = verificationEmail.User
		}
		return &go_block.UserResponse{
			LoginSession: &go_block.LoginSession{
				LoginStatus:    go_block.LoginStatus_EMAIL_IS_NOT_VERIFIED,
				EmailSentAt:    get.VerificationEmailSentAt,
				EmailExpiresAt: get.VerificationEmailExpiresAt,
			},
		}, nil //status.Error(codes.Code(go_block.ErrorType_ERROR_EMAIL_IS_NOT_VERIFIED), "user has not verified his/her email")
	}
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could setup tokens with err: %v", err)
	}
	// build data for token
	loggedInFrom := &go_block.Location{}
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	g.Go(func() error {
		// issue access and refresh token pair
		refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, get.Id, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate refresh token with err: %v", err)
		}
		// create refresh token in database
		if _, err := tokens.Create(ctx, &go_block.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_block.TokenType_TOKEN_TYPE_REFRESH,
		}); err != nil {
			return fmt.Errorf("could create refresh token with err: %v", err)
		}
		return nil
	})
	g.Go(func() error {
		accessToken, accessClaims, err = h.token.GenerateToken(privateKey, get.Id, refreshClaims.UserId, token.TokenTypeAccess, accessTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate access token with err: %v", err)
		}
		// create access token in database
		if _, err := tokens.Create(ctx, &go_block.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_block.TokenType_TOKEN_TYPE_ACCESS,
		}); err != nil {
			return fmt.Errorf("could create refresh token with err: %v", err)
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: get,
	}, nil
}
