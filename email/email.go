package email

import (
	"bytes"
	"errors"
	"github.com/keighl/postmark"
	"html/template"
	"os"
	"strings"
)

var (
	EmailSender Sender
	templates   = map[string]bool{}
	emailFrom   = ""
)

type TemplateData struct {
	LogoUrl        string
	Title          string
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

type Sender interface {
	Send(to, from, subject, data string, html bool) error
}

type PostmarkSender struct {
	client *postmark.Client
}

func (s *PostmarkSender) Send(to, from, subject string, data string, html bool) error {
	email := postmark.Email{
		From:       from,
		To:         to,
		Subject:    subject,
		TrackOpens: true,
	}
	if html {
		email.HtmlBody = data
	} else {
		email.TextBody = data
	}
	if _, err := s.client.SendEmail(email); err != nil {
		return err
	}
	return nil
}

func NewPostmarkSender(serverToken, accountToken string) (*PostmarkSender, error) {
	client := postmark.NewClient(serverToken, accountToken)
	return &PostmarkSender{
		client: client,
	}, nil
}

type defaultEmail struct {
	sender Sender
}

func initialize() error {
	var ok bool
	emailFrom, ok = os.LookupEnv("EMAIL_FROM")
	if !ok || emailFrom == "" {
		return errors.New("missing required EMAIL_FROM")
	}
	emailTemplatePaths, _ := os.LookupEnv("EMAIL_TEMPLATE_PATHS")
	for _, val := range strings.Fields(emailTemplatePaths) {
		templates[val] = true
	}
	return nil
}

func New(sender Sender) (Email, error) {
	if sender == nil {
		return nil, errors.New("invalid")
	}
	if err := initialize(); err != nil {
		return nil, err
	}
	for template, _ := range templates {
		if _, err := os.Stat(template); err != nil && errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}
	return &defaultEmail{
		sender: sender,
	}, nil
}

func (e *defaultEmail) SendEmail(to, subject, templatePath string, data *TemplateData) error {
	return nil
}

func (e *defaultEmail) SendVerificationEmail(to, subject, templatePath string, data *VerificationData) error {
	mail := ""
	html := false
	if templatePath != "" {
		if _, ok := templates[templatePath]; !ok {
			return errors.New("invalid path")
		}
		t, _ := template.ParseFiles(templatePath)
		var body bytes.Buffer
		if err := t.Execute(&body, data); err != nil {
			return err
		}
		mail = body.String()
		html = true
	} else {
		plaintext := data.WelcomeMessage + " " + data.NameOfUser + ",\n\n"
		plaintext += data.BodyMessage + "\n\n"
		plaintext += data.Code + "\n\n"
		plaintext += data.FooterMessage
		mail = plaintext
	}
	if err := e.sender.Send(to, emailFrom, subject, mail, html); err != nil {
		return err
	}
	return nil
}
