package handler

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
	"github.com/nuntiodev/nuntio-user-block/email"
	"github.com/nuntiodev/nuntio-user-block/repository/email_repository"
	"github.com/nuntiodev/x/cryptox"
	"strings"
)

func (h *defaultHandler) SendVerificationEmail(ctx context.Context, req *go_block.UserRequest) (*go_block.UserResponse, error) {
	if !h.emailEnabled {
		return nil, errors.New("email provider is not enabled")
	}
	// get requested user and check if the email is already verified
	userResp, err := h.Get(ctx, req)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	get := userResp.User
	if get.Email == "" {
		return &go_block.UserResponse{}, errors.New("user do not have an email - set the email for the user")
	}
	if get.EmailIsVerified {
		return &go_block.UserResponse{}, errors.New("email is already verified")
	}
	// generate verification code and send it to the user
	emails, err := h.repository.Emails(ctx, req.Namespace) // email config containing text to send
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	randomCode, err := h.crypto.GenerateSymmetricKey(6, cryptox.Numeric)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	verificationCode, err := hex.DecodeString(randomCode)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	emailConfig, err := emails.Get(ctx, &go_block.Email{
		Id: email_repository.VerificationEmail,
	})
	if err != nil {
		return nil, err
	}
	nameOfUser := get.Email
	if get.FirstName != "" {
		nameOfUser = strings.TrimSpace(get.FirstName + " " + get.LastName)
	}
	if err := h.email.SendVerificationEmail(get.Email, emailConfig.Subject, emailConfig.TemplatePath, &email.VerificationData{
		Code: string(verificationCode),
		TemplateData: email.TemplateData{
			LogoUrl:        emailConfig.Logo,
			WelcomeMessage: emailConfig.WelcomeMessage,
			NameOfUser:     nameOfUser,
			BodyMessage:    emailConfig.BodyMessage,
			FooterMessage:  emailConfig.FooterMessage,
		},
	}); err != nil {
		if err != nil {
			return &go_block.UserResponse{}, err
		}
	}
	// set verification code and timestamp in repository
	users, err := h.repository.Users().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return &go_block.UserResponse{}, err
	}
	get.VerificationCode = string(verificationCode)
	if _, err := users.UpdateVerificationEmailSent(ctx, &go_block.User{
		VerificationCode: string(verificationCode),
		Id:               get.Id,
	}); err != nil {
		return &go_block.UserResponse{}, err
	}
	return &go_block.UserResponse{}, nil
}
