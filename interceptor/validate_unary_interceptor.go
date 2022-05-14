package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nuntiodev/block-proto/go_block"
	"google.golang.org/grpc"
)

const (
	ProjectName             = "/BlockUser.UserService/"
	Heartbeat               = "Heartbeat"
	Create                  = "Create"
	UpdatePassword          = "UpdatePassword"
	UpdateMetadata          = "UpdateMetadata"
	UpdateImage             = "UpdateImage"
	UpdateEmail             = "UpdateEmail"
	UpdatePhoneNumber       = "UpdatePhoneNumber"
	UpdateName              = "UpdateName"
	UpdateBirthdate         = "UpdateBirthdate"
	UpdateUsername          = "UpdateUsername"
	UpdatePreferredLanguage = "UpdatePreferredLanguage"
	UpdateSecurity          = "UpdateSecurity"
	Get                     = "Get"
	GetAll                  = "GetAll"
	ValidateCredentials     = "ValidateCredentials"
	Login                   = "Login"
	ValidateToken           = "ValidateToken"
	BlockToken              = "BlockToken"
	BlockTokenById          = "BlockTokenById"
	RefreshToken            = "RefreshToken"
	GetTokens               = "GetTokens"
	PublicKeys              = "PublicKeys"
	RecordActiveMeasurement = "RecordActiveMeasurement"
	UserActiveHistory       = "UserActiveHistory"
	NamespaceActiveHistory  = "NamespaceActiveHistory"
	SendVerificationEmail   = "SendVerificationEmail"
	VerifyEmail             = "VerifyEmail"
	SendResetPasswordEmail  = "SendResetPasswordEmail"
	ResetPassword           = "ResetPassword"
	Delete                  = "Delete"
	DeleteBatch             = "DeleteBatch"
	DeleteNamespace         = "DeleteNamespace"
	CreateNamespaceConfig   = "CreateNamespaceConfig"
	UpdateConfigSettings    = "UpdateConfigSettings"
	UpdateConfigDetails     = "UpdateConfigDetails"
	GetConfig               = "GetConfig"
	DeleteConfig            = "DeleteConfig"
	CreateText              = "CreateText"
	UpdateGeneralText       = "UpdateGeneralText"
	UpdateWelcomeText       = "UpdateWelcomeText"
	UpdateRegisterText      = "UpdateRegisterText"
	UpdateLoginText         = "UpdateLoginText"
	UpdateProfileText       = "UpdateProfileText"
	DeleteText              = "DeleteText"
	InitializeApplication   = "InitializeApplication"
)

var (
	TokenIsNil          = errors.New("token is nil")
	TokenPointerIsEmpty = errors.New("token pointer is nil")
	UpdateIsNil         = errors.New("update is nil")
	UserIsNil           = errors.New("user is nil")
	NamespaceIsEmpty    = errors.New("namespace is empty")
	UserBatchIsNil      = errors.New("user batch is nil")
	MeasurementIsNil    = errors.New("measurement is nil")
	ConfigIsNil         = errors.New("config is nil")
	AuthConfigIsNil     = errors.New("auth config is nil")
	TextIsNil           = errors.New("text is nil")
)

func (i *DefaultInterceptor) WithValidateUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info == nil {
		return nil, errors.New("invalid request")
	}
	translatedReq, ok := req.(*go_block.UserRequest)
	if !ok {
		translatedReq = &go_block.UserRequest{}
	}
	method := strings.Split(info.FullMethod, ProjectName)
	if len(method) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid method call: %s", info.FullMethod))
	}
	switch method[1] {
	case Heartbeat, GetAll, PublicKeys,
		NamespaceActiveHistory, GetConfig, DeleteConfig,
		InitializeApplication:
		break
	case BlockToken:
		if translatedReq.TokenPointer == "" {
			return nil, TokenPointerIsEmpty
		}
	case RecordActiveMeasurement, UserActiveHistory:
		if translatedReq.ActiveMeasurement == nil {
			return &go_block.UserResponse{}, MeasurementIsNil
		}
	case Get, Create, VerifyEmail,
		SendVerificationEmail, SendResetPasswordEmail, ResetPassword:
		if translatedReq.User == nil {
			return &go_block.UserResponse{}, UserIsNil
		}
	case UpdatePassword, UpdateMetadata,
		UpdateImage, UpdateEmail, UpdateUsername,
		UpdateName, UpdateBirthdate, UpdatePhoneNumber,
		UpdatePreferredLanguage:
		if translatedReq.Update == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		} else if translatedReq.User == nil {
			return &go_block.UserResponse{}, UpdateIsNil
		}
	case ValidateCredentials, Delete, UpdateSecurity, Login:
		if translatedReq.User == nil {
			return nil, UserIsNil
		}
	case RefreshToken, GetTokens, BlockTokenById:
		if translatedReq.Token == nil {
			return nil, TokenIsNil
		}
	case ValidateToken:
		if translatedReq.TokenPointer == "" {
			return nil, errors.New("token pointer is nil")
		}
	case DeleteNamespace:
		if translatedReq.Namespace == "" {
			return nil, NamespaceIsEmpty
		}
	case DeleteBatch:
		if translatedReq.UserBatch == nil {
			return nil, UserBatchIsNil
		}
	case CreateNamespaceConfig, UpdateConfigSettings,
		UpdateConfigDetails:
		if translatedReq.Config == nil {
			return nil, ConfigIsNil
		}
	case UpdateGeneralText, UpdateWelcomeText, UpdateRegisterText,
		UpdateLoginText, CreateText, UpdateProfileText,
		DeleteText:
		if translatedReq.Text == nil {
			return nil, TextIsNil
		}
	default:
		return &go_block.UserResponse{}, errors.New(fmt.Sprintf("invalid request: %s", info.FullMethod))
	}
	h, err := handler(ctx, req) // make actual request
	return h, err
}
