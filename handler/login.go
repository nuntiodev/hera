package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nuntiodev/hera/models"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/hera/repository/user_repository"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"k8s.io/utils/strings/slices"
	"net/http"
	"time"

	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/token"
	"golang.org/x/crypto/bcrypt"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

/*
	Login - this method is used to authenticate a user and returns an access and refresh token used to validate a user afterwards.
*/
func (h *defaultHandler) Login(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		userRepository   user_repository.UserRepository
		config           *models.Config
		user             *models.User
		refreshToken     string
		refreshClaims    *go_hera.CustomClaims
		accessToken      string
		accessClaims     *go_hera.CustomClaims
		errGroup         = &errgroup.Group{}
	)
	// async action 1 - get namespace config.
	errGroup.Go(func() (err error) {
		configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		config, err = configRepository.Get(ctx)
		return err
	})
	// async action 2 - validate a users credentials and fetch user info.
	errGroup.Go(func() (err error) {
		userRepository, err = h.repository.UserRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
		if err != nil {
			return err
		}
		user, err = userRepository.Get(ctx, req.User)
		if err != nil {
			return err
		}
		return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.User.Password))
	})
	if err = errGroup.Wait(); err != nil {
		return nil, err
	}
	// validate if email is verified
	// if email validation is required and email is not verified; return error
	if req.User.GetEmail() != "" && config.VerifyEmail && slices.Contains(user.VerifiedEmails, user.EmailHash) == false {
		// check if we should send a new email
		if user.VerificationEmailExpiresAt.Sub(time.Now()).Seconds() <= 0 {
			// sent new email
			verificationEmail, err := h.SendVerificationEmail(ctx, req) // do not call directly on interface, but make same requests...
			if err != nil {
				return nil, fmt.Errorf("could not send email with err: %v", err)
			}
			user = models.ProtoUserToUser(verificationEmail.User)
		}
		return &go_hera.HeraResponse{
			LoginSession: &go_hera.LoginSession{
				LoginStatus:    go_hera.LoginStatus_EMAIL_IS_NOT_VERIFIED,
				EmailSentAt:    ts.New(user.VerificationEmailSentAt),
				EmailExpiresAt: ts.New(user.VerificationEmailExpiresAt),
			},
		}, nil
	}
	if req.User.GetPhone() != "" && config.VerifyPhone && slices.Contains(user.VerifiedPhoneNumbers, user.PhoneHash) == false {
		// check if we should send a new email
		// todo: implement
	}
	//  generate and save refresh and access tokenRepository
	tokenRepository, err := h.repository.TokenRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("could setup tokenRepository with err: %v", err)
	}
	// build data for tokenRepository
	refreshTokenId := uuid.NewString()
	loggedInFrom := ""
	deviceInfo := ""
	if req.Token != nil {
		loggedInFrom = req.Token.LoggedInFrom
		deviceInfo = req.Token.DeviceInfo
	}
	// async action 3 - generate refresh token.
	errGroup.Go(func() (err error) {
		refreshToken, refreshClaims, err = h.token.GenerateToken(privateKey, refreshTokenId, user.Id, "", token.RefreshToken, refreshTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate refresh token with err: %v", err)
		}
		// create refresh token info in database
		if err = tokenRepository.Create(ctx, &go_hera.Token{
			Id:           refreshClaims.Id,
			UserId:       refreshClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(refreshClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_hera.TokenType_TOKEN_TYPE_REFRESH,
		}); err != nil {
			return err
		}
		return nil
	})
	// async action 4 - generate access token.
	errGroup.Go(func() (err error) {
		accessToken, accessClaims, err = h.token.GenerateToken(privateKey, uuid.NewString(), user.Id, refreshTokenId, token.AccessToken, accessTokenExpiry)
		if err != nil {
			return fmt.Errorf("could generate access token with err: %v", err)
		}
		// create access token info in database
		if err = tokenRepository.Create(ctx, &go_hera.Token{
			Id:           accessClaims.Id,
			UserId:       accessClaims.UserId,
			ExpiresAt:    ts.New(time.Unix(accessClaims.ExpiresAt, 0)),
			LoggedInFrom: loggedInFrom,
			DeviceInfo:   deviceInfo,
			Type:         go_hera.TokenType_TOKEN_TYPE_ACCESS,
		}); err != nil {
			return err
		}
		return nil
	})
	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}
	// set access and refresh cookies for the browser
	accessCookie := &http.Cookie{
		Name:  HeraAccessTokenId,
		Value: accessToken,
		//HttpOnly: true,
	}
	refreshCookie := &http.Cookie{
		Name:  HeraRefreshTokenId,
		Value: refreshToken,
		//HttpOnly: true,
	}
	fmt.Println(accessCookie.String())
	// todo: check if header is also set for http server
	if err := grpc.SetHeader(ctx, map[string][]string{
		"Set-Cookie": {accessCookie.String(), refreshCookie.String()},
	}); err != nil {
		return nil, err
	}
	// return access and refresh token to the client
	return &go_hera.HeraResponse{
		Token: &go_hera.Token{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: models.UserToProtoUser(user),
	}, nil
}
