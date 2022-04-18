package initializer

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

const (
	BLOCK_USER_RSA_SECRET_NAME        = "NUNTIO_USER_BLOCK_RSA"
	BLOCK_USER_ENCRYPTION_SECRET_NAME = "NUNTIO_USER_BLOCK_ENCRYPTION"
	Kubernetes                        = "kubernetes"
	Memory                            = "memory"
)

type Initializer interface {
	CreateRSASecrets(ctx context.Context) error
	CreateEncryptionSecret(ctx context.Context) error
}

type defaultInitializer struct {
	namespace string
	zapLog    *zap.Logger
	k8s       *kubernetes.Clientset
	crypto    cryptox.Crypto
}

func New(zapLog *zap.Logger, engine string) (*defaultInitializer, error) {
	zapLog.Info("Initializing system with encryption secrets and public/private keys")
	if engine != Kubernetes && engine != Memory {
		return nil, fmt.Errorf("invalid engine %s", engine)
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	if _, err := os.Stat(tokenPath); err != nil {
		return nil, err
	}
	config.BearerTokenFile = tokenPath
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	bytesNamespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, err
	}
	crypto, err := cryptox.New()
	if err != nil {
		return nil, err
	}
	return &defaultInitializer{
		zapLog:    zapLog,
		k8s:       clientSet,
		namespace: string(bytesNamespace),
		crypto:    crypto,
	}, nil
}

func (i *defaultInitializer) CreateRsaSecrets(ctx context.Context) error {
	if os.Getenv("PUBLIC_KEY") == "" || os.Getenv("PRIVATE_KEY") == "" {
		i.zapLog.Info("RSA keys is already provided internally by the system (the PUBLIC_KEY or PRIVATE_KEY variable is set).")
		return nil
	}
	// check if secret already exists
	if cryptoSecret, err := i.k8s.CoreV1().Secrets(i.namespace).Get(ctx, BLOCK_USER_RSA_SECRET_NAME, metav1.GetOptions{}); err != nil {
		i.zapLog.Info("Block user RSA secret does not exist... creating....")
		// create public private keys
		rsaPrivateKey, rsaPublicKey, err := i.crypto.GenerateRsaKeyPair(2048)
		if err != nil {
			return err
		}
		userPrivateKey := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		}))
		userPublicKey := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(rsaPublicKey),
		}))
		// create secret in the Kubernetes api
		secretData := map[string]string{
			"PRIVATE_KEY": userPrivateKey,
			"PUBLIC_KEY":  userPublicKey,
		}
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: i.namespace,
				Name:      BLOCK_USER_RSA_SECRET_NAME,
			},
			StringData: secretData,
			Type:       v1.SecretTypeOpaque,
		}
		if _, err := i.k8s.CoreV1().Secrets(i.namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return err
		}
		// set in memory
		if err := os.Setenv("PUBLIC_KEY", userPublicKey); err != nil {
			return err
		}
		if err := os.Setenv("PRIVATE_KEY", userPrivateKey); err != nil {
			return err
		}
		i.zapLog.Info("Successfully created RSA secret")
	} else {
		if err := os.Setenv("PUBLIC_KEY", string(cryptoSecret.Data["PUBLIC_KEY"])); err != nil {
			return err
		}
		if err := os.Setenv("PRIVATE_KEY", string(cryptoSecret.Data["PRIVATE_KEY"])); err != nil {
			return err
		}
		i.zapLog.Info("RSA secret already exists")
	}
	return nil
}

func (i *defaultInitializer) CreateEncryptionSecret(ctx context.Context) error {
	if os.Getenv("ENCRYPTION_KEYS") != "" {
		i.zapLog.Info("Encryption keys is already provided internally by the system (the ENCRYPTION_KEYS variable is set).")
		return nil
	}
	// check if secret already exists
	if cryptoSecret, err := i.k8s.CoreV1().Secrets(i.namespace).Get(ctx, BLOCK_USER_ENCRYPTION_SECRET_NAME, metav1.GetOptions{}); err != nil {
		i.zapLog.Info("Block user encryption secret does not exist... creating....")
		// create encryption secret
		encryptionSecret, err := i.crypto.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return err
		}
		// create secret in the Kubernetes api
		secretData := map[string]string{
			"ENCRYPTION_KEYS": encryptionSecret,
		}
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: i.namespace,
				Name:      BLOCK_USER_ENCRYPTION_SECRET_NAME,
			},
			StringData: secretData,
			Type:       v1.SecretTypeOpaque,
		}
		if _, err := i.k8s.CoreV1().Secrets(i.namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return err
		}
		// set in memory
		if err := os.Setenv("ENCRYPTION_KEYS", encryptionSecret); err != nil {
			return err
		}
		i.zapLog.Info("Successfully created encryption secret")
	} else {
		if err := os.Setenv("ENCRYPTION_KEYS", string(cryptoSecret.Data["ENCRYPTION_KEYS"])); err != nil {
			return err
		}
		i.zapLog.Info("Encryption secret already exists")
	}
	return nil
}
