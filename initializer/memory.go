package initializer

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/fatih/color"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
	"os"
	"strings"
)

type memoryInitializer struct {
	namespace string
	zapLog    *zap.Logger
	redLog    *color.Color
	blueLog   *color.Color
}

func (i *memoryInitializer) CreateSecrets(ctx context.Context) error {
	i.redLog.Println("running in memory and creating new secrets. All data will be lost when shut down")
	if os.Getenv("PUBLIC_KEY") == "" || os.Getenv("PRIVATE_KEY") == "" {
		// create public private keys
		rsaPrivateKey, rsaPublicKey, err := cryptox.GenerateRsaKeyPair(4096)
		if err != nil {
			return err
		}
		heraPrivateKey := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		}))
		heraPublicKey := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(rsaPublicKey),
		}))
		// set in memory
		if err := os.Setenv("PUBLIC_KEY", heraPublicKey); err != nil {
			return err
		}
		if err := os.Setenv("PRIVATE_KEY", heraPrivateKey); err != nil {
			return err
		}
	}
	// create encryption secret
	if strings.TrimSpace(os.Getenv("NEW_ENCRYPTION_KEY")) == "" {
		encryptionSecret, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return err
		}
		if err := os.Setenv("ENCRYPTION_KEYS", encryptionSecret); err != nil {
			return err
		}
	}
	return nil
}
