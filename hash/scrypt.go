package hash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"golang.org/x/crypto/scrypt"
	"strconv"
)

type scryptHash struct {
	signerKey     []byte
	saltSeparator []byte
	rounds        int
	memCost       int
	p             int
	keyLen        int
}

func (sh *scryptHash) generate(password string) (*go_hera.Hash, error) {
	salt := make([]byte, 12)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	hash, err := key([]byte(password), salt, sh.signerKey, sh.saltSeparator, sh.rounds, sh.memCost, sh.p, sh.keyLen)
	if err != nil {
		return nil, err
	}
	return &go_hera.Hash{
		Variant: go_hera.HasingAlgorithm_BCRYPT,
		Body:    base64.StdEncoding.EncodeToString(hash),
		Params: map[string]string{
			"salt":           base64.StdEncoding.EncodeToString(salt),
			"signer_key":     string(sh.signerKey),
			"salt_separator": string(sh.saltSeparator),
			"rounds":         strconv.Itoa(sh.rounds),
			"mem_cost":       strconv.Itoa(sh.memCost),
			"p":              strconv.Itoa(sh.p),
			"key_len":        strconv.Itoa(sh.keyLen),
		},
	}, nil
}

func (sh *scryptHash) compare(password, salt, hash string) error {
	_salt, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return err
	}
	res, err := key([]byte(password), _salt, sh.signerKey, sh.saltSeparator, sh.rounds, sh.memCost, sh.p, sh.keyLen)
	if err != nil {
		return err
	}
	if base64.StdEncoding.EncodeToString(res) != hash {
		return errors.New("scrypt: invalid hash")
	}
	return nil
}

func Key(password, salt []byte, signerKey, saltSeparator string, rounds, memCost, p, keyLen int) ([]byte, error) {
	var (
		sk, ss []byte
		err    error
	)

	if sk, err = base64.StdEncoding.DecodeString(signerKey); err != nil {
		return nil, err
	}
	if ss, err = base64.StdEncoding.DecodeString(saltSeparator); err != nil {
		return nil, err
	}

	return key(password, salt, sk, ss, rounds, memCost, p, keyLen)
}

func key(password, salt, signerKey, saltSeparator []byte, rounds, memCost, p, keyLen int) ([]byte, error) {
	ck, err := scrypt.Key(password, append(salt, saltSeparator...), 1<<memCost, rounds, p, keyLen)
	if err != nil {
		return nil, err
	}

	var block cipher.Block
	if block, err = aes.NewCipher(ck); err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(signerKey))
	stream := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	stream.XORKeyStream(cipherText[aes.BlockSize:], signerKey)
	return cipherText[aes.BlockSize:], nil
}
