package hash

import (
	"github.com/nuntiodev/hera-sdks/go_hera"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type bcryptHash struct {
	cost int
}

func (bh *bcryptHash) Generate(password string) (*go_hera.Hash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bh.cost)
	if err != nil {
		return nil, err
	}
	return &go_hera.Hash{
		Variant: go_hera.HasingAlgorithm_BCRYPT,
		Body:    string(hash),
		Params: map[string]string{
			"cost": strconv.Itoa(bh.cost),
		},
	}, nil
}

func (bh *bcryptHash) Compare(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
