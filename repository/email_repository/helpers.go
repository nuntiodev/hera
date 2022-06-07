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
	email.Id = strings.TrimSpace(email.Id)
}
