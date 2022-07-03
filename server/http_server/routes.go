package http_server

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/authenticator"
	"github.com/nuntiodev/hera/interceptor"
)

type HeraHttpRoute struct {
	Name   string
	Handle func(ctx context.Context, request *go_hera.HeraRequest) (*go_hera.HeraResponse, error)
}

const (
	userBasePath               = "/hera/user"
	userCreate                 = "/create"
	userGet                    = "/user"
	usersGet                   = "/users"
	userList                   = "/list"
	userDelete                 = "/delete"
	userLogin                  = "/login"
	userResetPassword          = "/reset-password"
	userSearch                 = "/search"
	userSendResetPasswordEmail = "/send-reset-password-email"
	userSendResetPasswordText  = "/send-reset-password-text"
	userSendVerificationEmail  = "/send-verification-email"
	userSendVerificationText   = "/send-verification-text"
	userUpdateContact          = "/update-contact"
	userUpdateProfile          = "/update-profile"
	userValidateCredentials    = "/validate-credentials"
	userVerifyEmail            = "/verify-email"
	userVerifyPhone            = "/verify-phone"
)

const (
	configBasePath        = "/hera/config"
	configCreateNamespace = "/create-namespace"
	configDelete          = "/delete"
	configDeleteNamespace = "/delete-namespace"
	configGet             = "/config"
	configRemovePublicKey = "/remove-public-key"
	configRegisterRsaKey  = "/register-rsa-key"
	configUpdate          = "/update"
	configPublicKeys      = "/public-keys"
)
const (
	tokenBasePath = "/hera/token"
	tokenCreate   = "/create"
	tokenValidate = "/validate"
	tokenBlock    = "/block"
	tokenRefresh  = "/refresh"
	tokenTokens   = "/tokens"
)

func (s *Server) performAction(name string, x func(ctx context.Context, req *go_hera.HeraRequest) (*go_hera.HeraResponse, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("received new request: " + r.RequestURI)
		var heraRequest go_hera.HeraRequest
		if err := jsonpb.Unmarshal(r.Body, &heraRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Debug("request failed with err: " + err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := s.authenticator.AuthenticateRequest(ctx, &heraRequest, &authenticator.Info{
			IsHttp: true,
			Name:   name,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			s.logger.Debug("request failed with err: " + err.Error())
			return
		}
		if err := interceptor.ValidateRequest(name, &heraRequest); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			s.logger.Debug("request failed with err: " + err.Error())
			return
		}
		heraResponse, err := x(ctx, &heraRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Debug("request failed with err: " + err.Error())
			return
		}
		if err := (&jsonpb.Marshaler{}).Marshal(w, heraResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			s.logger.Debug("request failed with err: " + err.Error())
			return
		}
	}
}

func (s *Server) routes() *mux.Router {
	router := mux.NewRouter()
	// user router
	userRouter := router.PathPrefix(userBasePath).Subrouter()
	userRouter.HandleFunc(userCreate, s.performAction(interceptor.CreateUser, s.handler.CreateUser))
	userRouter.HandleFunc(userGet, s.performAction(interceptor.GetUser, s.handler.GetUser))
	userRouter.HandleFunc(usersGet, s.performAction(interceptor.GetUser, s.handler.GetUsers))
	userRouter.HandleFunc(userList, s.performAction(interceptor.ListUsers, s.handler.ListUsers))
	userRouter.HandleFunc(userDelete, s.performAction(interceptor.DeleteUser, s.handler.DeleteUser))
	userRouter.HandleFunc(userLogin, s.performAction(interceptor.Login, s.handler.Login))
	userRouter.HandleFunc(userResetPassword, s.performAction(interceptor.ResetPassword, s.handler.ResetPassword))
	userRouter.HandleFunc(userSearch, s.performAction(interceptor.SearchForUser, s.handler.SearchForUser))
	userRouter.HandleFunc(userSendResetPasswordEmail, s.performAction(interceptor.SendResetPasswordEmail, s.handler.SendResetPasswordEmail))
	userRouter.HandleFunc(userSendResetPasswordText, s.performAction(interceptor.SendResetPasswordText, s.handler.SendResetPasswordText))
	userRouter.HandleFunc(userSendVerificationEmail, s.performAction(interceptor.SendVerificationEmail, s.handler.SendVerificationEmail))
	userRouter.HandleFunc(userSendVerificationText, s.performAction(interceptor.SendVerificationText, s.handler.SendVerificationText))
	userRouter.HandleFunc(userUpdateContact, s.performAction(interceptor.UpdateUserContact, s.handler.UpdateUserContact))
	userRouter.HandleFunc(userUpdateProfile, s.performAction(interceptor.UpdateUserProfile, s.handler.UpdateUserProfile))
	userRouter.HandleFunc(userValidateCredentials, s.performAction(interceptor.ValidateCredentials, s.handler.ValidateCredentials))
	userRouter.HandleFunc(userVerifyEmail, s.performAction(interceptor.VerifyEmail, s.handler.VerifyEmail))
	userRouter.HandleFunc(userVerifyPhone, s.performAction(interceptor.VerifyPhone, s.handler.VerifyPhone))
	// config router
	configRouter := router.PathPrefix(configBasePath).Subrouter()
	configRouter.HandleFunc(configCreateNamespace, s.performAction(interceptor.CreateNamespace, s.handler.CreateNamespace))
	configRouter.HandleFunc(configDelete, s.performAction(interceptor.DeleteConfig, s.handler.DeleteConfig))
	configRouter.HandleFunc(configDeleteNamespace, s.performAction(interceptor.DeleteNamespace, s.handler.DeleteNamespace))
	configRouter.HandleFunc(configGet, s.performAction(interceptor.GetConfig, s.handler.GetConfig))
	configRouter.HandleFunc(configRemovePublicKey, s.performAction(interceptor.RemovePublicKey, s.handler.RemovePublicKey))
	configRouter.HandleFunc(configRegisterRsaKey, s.performAction(interceptor.RegisterRsaKey, s.handler.RegisterRsaKey))
	configRouter.HandleFunc(configUpdate, s.performAction(interceptor.UpdateConfig, s.handler.UpdateConfig))
	configRouter.HandleFunc(configPublicKeys, s.performAction(interceptor.PublicKeys, s.handler.PublicKeys))
	// token router
	tokenRouter := router.PathPrefix(tokenBasePath).Subrouter()
	tokenRouter.HandleFunc(tokenCreate, s.performAction(interceptor.CreateTokenPair, s.handler.CreateTokenPair))
	tokenRouter.HandleFunc(tokenValidate, s.performAction(interceptor.ValidateToken, s.handler.ValidateToken))
	tokenRouter.HandleFunc(tokenBlock, s.performAction(interceptor.BlockToken, s.handler.BlockToken))
	tokenRouter.HandleFunc(tokenRefresh, s.performAction(interceptor.RefreshToken, s.handler.RefreshToken))
	tokenRouter.HandleFunc(tokenTokens, s.performAction(interceptor.GetTokens, s.handler.GetTokens))
	return router
}
