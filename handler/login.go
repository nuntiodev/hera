package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nuntiodev/nuntio-user-block/repository/token_repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *defaultHandler) Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		err             error
		namespaceConfig *go_block.Config
		user            *go_block.User
		accessToken     string
		refreshToken    string
	)
	// setup async err group
	g, ctx := errgroup.WithContext(ctx)
	// step 1: get config and valid user credentials
	g.Go(func() error {
		namespaceConfig, err = h.getNamespaceConfig(ctx, req.Namespace)
		return err
	})
	g.Go(func() error {
		user, err = h.validateUserCredentials(ctx, req)
		return err
	})
	if err := g.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// step 2: validate if email is verified
	// if email validation is required and email is not verified; return error
	if namespaceConfig.RequireEmailVerification && user.EmailIsVerified == false {
		// check if we should send a new email
		if user.VerificationEmailExpiresAt.AsTime().Sub(time.Now()).Seconds() <= 0 {
			// sent new email
			verificationEmail, err := h.SendVerificationEmail(ctx, req)
			if err != nil {
				return &go_block.UserResponse{}, fmt.Errorf("could not send email with err: %v", err)
			}
			user = verificationEmail.User
		}
		return &go_block.UserResponse{
			LoginSession: &go_block.LoginSession{
				LoginStatus:    go_block.LoginStatus_EMAIL_IS_NOT_VERIFIED,
				EmailSentAt:    user.VerificationEmailSentAt,
				EmailExpiresAt: user.VerificationEmailExpiresAt,
			},
		}, nil
	}
	// step 3: generate and save refresh and access tokens
	tokens, err := h.repository.Tokens(ctx, req.Namespace)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could setup tokens with err: %v", err)
	}
	// build data for tokens
	refreshTokenId := uuid.NewString()
	loggedInFrom := &go_block.Location{}
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	g.Go(func() error {
		refreshToken, err = h.generateRefreshToken(ctx, refreshTokenId, user.Id, deviceInfo, loggedInFrom, tokens)
		return err
	})
	g.Go(func() error {
		accessToken, err = h.generateAccessToken(ctx, refreshTokenId, user.Id, deviceInfo, loggedInFrom, tokens)
		return err
	})
	if err := g.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// step 4: return tokens
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: user,
	}, nil
}

func (h *defaultHandler) generateRefreshToken(ctx context.Context, refreshId, userId, deviceInfo string, loggedInFrom *go_block.Location, tokens token_repository.TokenRepository) (string, error) {
	// issue access and refresh token pair
	refreshToken, refreshClaims, err := h.token.GenerateToken(privateKey, refreshId, userId, "", token.TokenTypeRefresh, refreshTokenExpiry)
	if err != nil {
		return "", fmt.Errorf("could generate refresh token with err: %v", err)
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
		return "", fmt.Errorf("could create refresh token with err: %v", err)
	}
	return refreshToken, nil
}

func (h *defaultHandler) generateAccessToken(ctx context.Context, refreshTokenId, userId, deviceInfo string, loggedInFrom *go_block.Location, tokens token_repository.TokenRepository) (string, error) {
	accessToken, accessClaims, err := h.token.GenerateToken(privateKey, uuid.NewString(), userId, refreshTokenId, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return "", fmt.Errorf("could generate access token with err: %v", err)
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
		return "", fmt.Errorf("could create refresh token with err: %v", err)
	}
	return accessToken, nil
}

func (h *defaultHandler) validateUserCredentials(ctx context.Context, req *go_block.UserRequest) (*go_block.User, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("could not get user with err: %v", err)
	}
	if resp.User.Password == "" {
		return nil, errors.New("please update the user with a non-empty password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return nil, errors.New("please update the user with a non-empty password")
	}
	return resp.User, nil
}

func (h *defaultHandler) getNamespaceConfig(ctx context.Context, namespace string) (*go_block.Config, error) {
	// get config
	config, err := h.repository.Config(ctx, namespace)
	if err != nil {
		return nil, err
	}
	namespaceConfig, err := config.GetNamespaceConfig(ctx)
	if err != nil {
		return nil, err
	}
	return namespaceConfig, nil
}
