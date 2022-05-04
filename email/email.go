package email

import (
	"errors"
	"github.com/nuntiodev/x/emailx"
	"os"
	"strings"
)

var (
	smtpFrom     = ""
	smtpPassword = ""
	smtpHost     = ""
	smtpPort     = ""
	templates    map[string]bool
)

type TemplateData struct {
	LogoUrl        string
	WelcomeMessage string
	NameOfUser     string
	BodyMessage    string
	FooterMessage  string
}

type VerificationData struct {
	Code string
	TemplateData
}

type Email interface {
	SendEmail(to, subject, templatePath string, data *TemplateData) error
	SendVerificationEmail(to, subject, templatePath string, data *VerificationData) error
}

type defaultEmail struct {
	emailx emailx.Email
}

func initialize() error {
	var ok bool
	smtpFrom, ok = os.LookupEnv("SMTP_FROM")
	if !ok || smtpFrom == "" {
		return errors.New("missing required SMTP_FROM")
	}
	smtpPassword, ok = os.LookupEnv("SMTP_PASSWORD")
	if !ok || smtpPassword == "" {
		return errors.New("missing required SMTP_PASSWORD")
	}
	smtpHost, ok = os.LookupEnv("SMTP_HOST")
	if !ok || smtpHost == "" {
		return errors.New("missing required SMTP_HOST")
	}
	smtpPort, ok = os.LookupEnv("SMTP_PORT")
	if !ok || smtpPort == "" {
		return errors.New("missing required SMTP_PORT")
	}
	emailTemplatePaths, _ := os.LookupEnv("EMAIL_TEMPLATE_PATHS")
	for _, val := range strings.Fields(emailTemplatePaths) {
		templates[val] = true
	}
	return nil
}

func New() (Email, error) {
	if err := initialize(); err != nil {
		return nil, err
	}
	myEmail, err := emailx.New(smtpFrom, smtpPassword, smtpHost, smtpPort)
	if err != nil {
		return nil, err
	}
	return &defaultEmail{
		emailx: myEmail,
	}, nil
}

func (e *defaultEmail) SendEmail(to, subject, templatePath string, data *TemplateData) error {
	if templatePath != "" {
		if _, ok := templates[templatePath]; !ok {
			return errors.New("invalid path")
		}
	}
	if err := e.emailx.SendEmail(to, subject, templatePath, data); err != nil {
		return err
	}
	return nil
}

func (e *defaultEmail) SendVerificationEmail(to, subject, templatePath string, data *VerificationData) error {
	if templatePath != "" {
		if _, ok := templates[templatePath]; !ok {
			return errors.New("invalid path")
		}
	}
	if err := e.emailx.SendEmail(to, subject, templatePath, data); err != nil {
		return err
	}
	return nil
}
