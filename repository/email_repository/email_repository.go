package email_repository

import (
	"context"
	"github.com/nuntiodev/block-proto/go_block"
	"time"
)

type Email struct {
	Id             string    `bson:"_id" json:"id"`
	LogoUrl        string    `bson:"logo_url" json:"logo_url"`
	WelcomeMessage string    `bson:"welcome_message" json:"welcome_message"`
	BodyMessage    string    `bson:"body_message" json:"body_message"`
	FooterMessage  string    `bson:"footer_message" json:"footer_message"`
	CreatedAt      time.Time `bson:"updated_at" json:"updated_at"`
	EncryptedAt    time.Time `bson:"encrypted_at" json:"encrypted_at"`
}

type EmailRepository interface {
	Create(ctx context.Context, email *go_block.Email) (*Email, error)
	Get(ctx context.Context, email *go_block.Email) (*Email, error)
	Update(ctx context.Context, email *go_block.Email) (*Email, error)
	Delete(ctx context.Context, email *go_block.Email) error
}
