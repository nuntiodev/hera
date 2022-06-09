package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nuntiodev/nuntio-user-block/models"
	"github.com/nuntiodev/nuntio-user-block/repository/config_repository"
	"github.com/nuntiodev/nuntio-user-block/repository/user_repository"
	"golang.org/x/sync/errgroup"
	"time"

	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

/*
	Login - this method is used to authenticate a user and returns an access and refresh token used to validate a user afterwards.
*/
func (h *defaultHandler) Login(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	var (
		configRepo    config_repository.ConfigRepository
		userRepo      user_repository.UserRepository
		config        *models.Config
		user          *models.User
		refreshToken  string
		refreshClaims *go_block.CustomClaims
		accessToken   string
		accessClaims  *go_block.CustomClaims
		errGroup      = &errgroup.Group{}
		err           error
	)
	// async action 1 - get namespace config.
	errGroup.Go(func() error {
		configRepo, err = h.repository.Config(ctx, req.Namespace, req.EncryptionKey)
		if err != nil {
			return err
		}
		config, err = configRepo.GetNamespaceConfig(ctx)
		return err
	})
	// async action 2 - validate a users credentials and fetch user info.
	errGroup.Go(func() error {
		userRepo, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).SetEncryptionKey(req.EncryptionKey).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepo.Get(ctx, req.User)
		if err != nil {
			return err
		}
		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.User.Password))
	})
	if err := errGroup.Wait(); err != nil {
		return &go_block.UserResponse{}, err
	}
	// validate if email is verified
	// if email validation is required and email is not verified; return error
	if config.RequireEmailVerification && user.EmailIsVerified == false {
		// check if we should send a new email
		if user.VerificationEmailExpiresAt.Sub(time.Now()).Seconds() <= 0 {
			// sent new email
			verificationEmail, err := h.SendVerificationEmail(ctx, req)
			if err != nil {
				return &go_block.UserResponse{}, fmt.Errorf("could not send email with err: %v", err)
			}
			user = models.ProtoUserToUser(verificationEmail.User)
		}
		return &go_block.UserResponse{
			LoginSession: &go_block.LoginSession{
				LoginStatus:    go_block.LoginStatus_EMAIL_IS_NOT_VERIFIED,
				EmailSentAt:    ts.New(user.VerificationEmailSentAt),
				EmailExpiresAt: ts.New(user.VerificationEmailExpiresAt),
			},
		}, nil
	}
	// step 3: generate and save refresh and access tokens
	tokens, err := h.repository.Tokens(ctx, req.Namespace, req.EncryptionKey)
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
	// async action 3 - generate refresh token.
	errGroup.Go(func() error {
		refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, refreshTokenId, user.Id, "", token.TokenTypeRefresh, refreshTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate refresh token with err: %v", err)
		}
		// create refresh token info in database
		_, err := tokens.Create(ctx, &go_block.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_block.TokenType_TOKEN_TYPE_REFRESH,
		})
		return err
	})
	// async action 4 - generate access token.
	errGroup.Go(func() error {
		accessToken, accessClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), user.Id, refreshTokenId, token.TokenTypeAccess, accessTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate access token with err: %v", err)
		}
		// create access token info in database
		_, err = tokens.Create(ctx, &go_block.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_block.TokenType_TOKEN_TYPE_ACCESS,
		})
		return err
	})
	err = errGroup.Wait()
	return &go_block.UserResponse{
		Token: &go_block.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: models.UserToProtoUser(user),
	}, err
}
