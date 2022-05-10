package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *defaultHandler) Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	resp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could not get user with err: %v", err)
	}
	// if email validation is required and email is not verified; return error
	if resp.User.RequireEmailVerification && resp.User.EmailIsVerified == false {
		// check if we should send a new email
		if resp.User.VerificationEmailExpiresAt.AsTime().Sub(time.Now()).Seconds() <= 0 {
			// sent new email
			verificationEmail, err := h.SendVerificationEmail(ctx, req)
			if err != nil {
				return &go_block.UserResponse{}, fmt.Errorf("could not send email with err: %v", err)
			}
			resp.User = verificationEmail.User
		}
		return &go_block.UserResponse{
			LoginSession: &go_block.LoginSession{
				LoginStatus:    go_block.LoginStatus_EMAIL_IS_NOT_VERIFIED,
				EmailSentAt:    resp.User.VerificationEmailSentAt,
				EmailExpiresAt: resp.User.VerificationEmailExpiresAt,
			},
		}, nil //status.Error(codes.Code(go_block.ErrorType_ERROR_EMAIL_IS_NOT_VERIFIED), "user has not verified his/her email")
	}
	if resp.User.Password == "" {
		return &go_block.UserResponse{}, errors.New("please update the user with a non-empty password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(resp.User.Password), []byte(req.User.Password)); err != nil {
		return &go_block.UserResponse{}, err
	}
	// issue access and refresh token pair
	refreshToken, refreshClaims, err := h.token.GenerateToken(privateKey, resp.User.Id, "", token.TokenTypeRefresh, refreshTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could generate refresh token with err: %v", err)
	}
	accessToken, accessClaims, err := h.token.GenerateToken(privateKey, resp.User.Id, refreshClaims.UserId, token.TokenTypeAccess, accessTokenExpiry)
	if err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could generate access token with err: %v", err)
	}
	// setup token database
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
	// create refresh token in database
	if _, err := tokens.Create(ctx, &go_block.Token{
		Id:           refreshClaims.Id,
		UserId:       refreshClaims.UserId,
		ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
		LoggedInFrom: loggedInFrom,
		DeviceInfo:   deviceInfo,
		Type:         go_block.TokenType_TOKEN_TYPE_REFRESH,
	}); err != nil {
		return &go_block.UserResponse{}, fmt.Errorf("could create refresh token with err: %v", err)
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
		return &go_block.UserResponse{}, fmt.Errorf("could create refresh token with err: %v", err)
	}
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: resp.User,
	}, nil
}
