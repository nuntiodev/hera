package email_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func prepare(action int, email *go_block.Email) {
	if email == nil {
		return
	}
	switch action {
	case actionCreate:
		email.CreatedAt = ts.Now()
		email.UpdatedAt = ts.Now()
	case actionUpdate:
		email.UpdatedAt = ts.Now()
	}
	email.Subject = strings.TrimSpace(email.Subject)
	email.Logo = strings.TrimSpace(email.Logo)
	email.WelcomeMessage = strings.TrimSpace(email.WelcomeMessage)
	email.BodyMessage = strings.TrimSpace(email.BodyMessage)
	email.FooterMessage = strings.TrimSpace(email.FooterMessage)
	email.Id = strings.TrimSpace(email.Id)
}
