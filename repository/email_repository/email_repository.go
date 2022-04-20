package email_repository

import "time"

type Email struct {
	Id          string    `bson:"_id" json:"id"`
	SendAt      time.Time `bson:"updated_at" json:"updated_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	EncryptedAt time.Time `bson:"encrypted_at" json:"encrypted_at"`
}

type EmailRepository interface {
}
