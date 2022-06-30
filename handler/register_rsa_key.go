package handler

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/hera/repository/config_repository"
	"github.com/nuntiodev/x/cryptox"
)

func (h *defaultHandler) RegisterRsaKey(ctx context.Context, req *go_hera.HeraRequest) (resp *go_hera.HeraResponse, err error) {
	var (
		configRepository config_repository.ConfigRepository
		p                *rsa.PublicKey
		s                *rsa.PrivateKey
		publicKey        string
		privateKey       string
		config           *go_hera.Config
	)
	configRepository, err = h.repository.ConfigRepositoryBuilder().SetNamespace(req.Namespace).Build(ctx)
	if err != nil {
		return nil, err
	}
	config, err = configRepository.Get(ctx)
	if err != nil {
		return nil, err
	}
	if config.PublicKey != "" {
		return nil, errors.New("a public key is already present. Please call RemovePublicKey before adding a new one")
	}
	// generate rsa key pair
	s, p, err = cryptox.GenerateRsaKeyPair(2048)
	if err != nil {
		return nil, err
	}
	privateKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s),
	}))
	publicKey = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(p),
	}))
	if err := configRepository.RegisterPublicKey(ctx, publicKey); err != nil {
		return nil, err
	}
	return &go_hera.HeraResponse{
		PrivateKey: privateKey,
	}, nil
}
