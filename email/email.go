package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

type Email interface {
}

type defaultEmail struct {
	auth     smtp.Auth
	from     string
	password string
	port     string
	host     string
}

func New(smtpFrom, smtpPassword, smtpHost, smtpPort string) (Email, error) {
	return &defaultEmail{
		auth:     smtp.PlainAuth("", smtpFrom, smtpPassword, smtpHost),
		from:     smtpFrom,
		password: smtpPassword,
		host:     smtpHost,
		port:     smtpPort,
	}, nil
}

func (e *defaultEmail) SendEmail(to, subject, templatePath string, data any) error {
	t, _ := template.ParseFiles(templatePath)
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))
	if err := t.Execute(&body, data); err != nil {
		return err
	}
	if err := smtp.SendMail(e.host+":"+e.port, e.auth, e.from, []string{to}, body.Bytes()); err != nil {
		return err
	}
	return nil
}
