package hash

import (
	"errors"
	"fmt"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type Hash interface {
	Generate(password string) (*go_hera.Hash, error)
	Compare(password string, hashModel *go_hera.Hash) error
}

type hasher struct {
	config *go_hera.Config
}

func New(config *go_hera.Config) Hash {
	return &hasher{config: config}
}

func (h *hasher) Generate(password string) (*go_hera.Hash, error) {
	if h.config == nil {
		return nil, errors.New("config is nil")
	}
	switch h.config.HasingAlgorithm {
	case go_hera.HasingAlgorithm_BCRYPT:
		bcryptCost := bcrypt.DefaultCost
		if h.config.Bcrypt != nil {
			bcryptCost = int(h.config.Bcrypt.Cost)
		}
		return (&bcryptHash{cost: bcryptCost}).Generate(password)
	case go_hera.HasingAlgorithm_SCRYPT:
		signerKey := []byte("")
		saltSeparator := []byte("Bw==")
		rounds := 8
		memCost := 14
		p := 1
		keyLen := 32
		if h.config.Scrypt != nil {
			saltSeparator = []byte(h.config.Scrypt.SaltSeparator)
			signerKey = []byte(h.config.Scrypt.SignerKey)
			rounds = int(h.config.Scrypt.Rounds)
			memCost = int(h.config.Scrypt.MemCost)
			p = int(h.config.Scrypt.P)
			keyLen = int(h.config.Scrypt.KeyLen)
		}
		return (&scryptHash{
			signerKey:     signerKey,
			saltSeparator: saltSeparator,
			rounds:        rounds,
			memCost:       memCost,
			p:             p,
			keyLen:        keyLen,
		}).generate(password)
	}
	return nil, fmt.Errorf("hash generate: invalid hash %s", h.config.GetName())
}

func (h *hasher) Compare(password string, hashModel *go_hera.Hash) error {
	switch hashModel.Variant {
	case go_hera.HasingAlgorithm_BCRYPT:
		return (&bcryptHash{}).Compare(password, hashModel.Body)
	case go_hera.HasingAlgorithm_SCRYPT:
		salt, ok := hashModel.Params["salt"]
		if !ok || salt == "" {
			return errors.New("missing required salt")
		}
		signerKey, ok := hashModel.Params["signer_key"]
		if !ok || len(signerKey) == 0 {
			return errors.New("missing required signer key")
		}
		saltSeparator, ok := hashModel.Params["salt_separator"]
		if !ok || len(saltSeparator) == 0 {
			return errors.New("missing required salt separator")
		}
		roundsString, ok := hashModel.Params["rounds"]
		if !ok {
			return errors.New("missing required rounds")
		}
		rounds, err := strconv.Atoi(roundsString)
		if err != nil {
			return err
		}
		memCostString, ok := hashModel.Params["mem_cost"]
		if !ok {
			return errors.New("missing required mem cost")
		}
		memCost, err := strconv.Atoi(memCostString)
		if err != nil {
			return err
		}
		pString, ok := hashModel.Params["p"]
		if !ok {
			return errors.New("missing required p")
		}
		p, err := strconv.Atoi(pString)
		if err != nil {
			return err
		}
		keyLenString, ok := hashModel.Params["key_len"]
		if !ok {
			return errors.New("missing required key_len")
		}
		keyLen, err := strconv.Atoi(keyLenString)
		if err != nil {
			return err
		}
		return (&scryptHash{
			signerKey:     []byte(signerKey),
			saltSeparator: []byte(saltSeparator),
			rounds:        rounds,
			memCost:       memCost,
			p:             p,
			keyLen:        keyLen,
		}).compare(password, salt, hashModel.Body)
	}
	return fmt.Errorf("hash compare: invalid hash %s", hashModel.Variant)
}
