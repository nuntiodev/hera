package crypto

import (
	"github.com/softcorp-io/block-proto/go_block"
)

type Crypto interface {
	Encrypt(stringToEncrypt string, keyString string) (string, error)
	Decrypt(encryptedString string, keyString string) (string, error)
	EncryptUser(key string, user *go_block.User) error
	DecryptUser(key string, user *go_block.User) error
}

type defaultCrypto struct{}

func New() (Crypto, error) {
	return &defaultCrypto{}, nil
}
