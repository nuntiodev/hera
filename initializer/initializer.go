package initializer

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/fatih/color"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strings"
)

const (
	blockUserSecretName = "user-block-secret"
	EngineKubernetes    = "kubernetes"
	EngineMemory        = "memory"
)

type Initializer interface {
	CreateSecrets(ctx context.Context) error
}

type defaultInitializer struct {
	namespace string
	zapLog    *zap.Logger
	redLog    *color.Color
	k8s       *kubernetes.Clientset
}

func New(zapLog *zap.Logger, engine string) (*defaultInitializer, error) {
	zapLog.Info("initializing system with encryption secrets and public/private keys")
	if engine != EngineKubernetes && engine != EngineMemory {
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
	return &defaultInitializer{
		zapLog:    zapLog,
		k8s:       clientSet,
		redLog:    color.New(color.FgRed),
		namespace: string(bytesNamespace),
	}, nil
}

func (i *defaultInitializer) CreateSecrets(ctx context.Context) error {
	if os.Getenv("PUBLIC_KEY") != "" && os.Getenv("PRIVATE_KEY") != "" && os.Getenv("ENCRYPTION_KEYS") != "" {
		i.zapLog.Info("secrets are already provided internally by the system (the PUBLIC_KEY, ENCRYPTION_KEYS and PRIVATE_KEY variable is set).")
		return nil
	}
	// check if secret already exists
	if cryptoSecret, err := i.k8s.CoreV1().Secrets(i.namespace).Get(ctx, blockUserSecretName, metav1.GetOptions{}); err != nil {
		i.zapLog.Info("block user secret does not exist... creating....")
		// create public private keys
		rsaPrivateKey, rsaPublicKey, err := cryptox.GenerateRsaKeyPair(4096)
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
		// create encryption secret
		encryptionSecret, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return err
		}
		// create secret in the EngineKubernetes api
		secretData := map[string]string{
			"PRIVATE_KEY":     userPrivateKey,
			"PUBLIC_KEY":      userPublicKey,
			"ENCRYPTION_KEYS": encryptionSecret,
		}
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: i.namespace,
				Name:      blockUserSecretName,
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
		if err := os.Setenv("ENCRYPTION_KEYS", encryptionSecret); err != nil {
			return err
		}
		i.zapLog.Info("successfully created secret")
	} else {
		i.zapLog.Info("block user secret already exists")
		newEncryptionKey := strings.TrimSpace(os.Getenv("NEW_ENCRYPTION_KEY"))
		newKeyAlreadyExists := false
		encryptionKeys := strings.Fields(string(cryptoSecret.Data["ENCRYPTION_KEYS"]))
		for i, key := range encryptionKeys {
			encryptionKeys[i] = strings.TrimSpace(key)
			if newEncryptionKey == encryptionKeys[i] && !newKeyAlreadyExists {
				newKeyAlreadyExists = true
			}
		}
		if newEncryptionKey != "" && !newKeyAlreadyExists {
			i.zapLog.Info("new key added... updating existing secret...")
			encryptionKeys = append(encryptionKeys, newEncryptionKey)
			secretData := map[string]string{
				"ENCRYPTION_KEYS": strings.Join(encryptionKeys, " "),
				"PRIVATE_KEY":     string(cryptoSecret.Data["PUBLIC_KEY"]),
				"PUBLIC_KEY":      string(cryptoSecret.Data["PUBLIC_KEY"]),
			}
			if _, err := i.k8s.CoreV1().Secrets(i.namespace).Update(ctx, &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: i.namespace,
					Name:      blockUserSecretName,
				},
				StringData: secretData,
				Type:       v1.SecretTypeOpaque,
			}, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}
		// set in memory
		if err := os.Setenv("PUBLIC_KEY", string(cryptoSecret.Data["PUBLIC_KEY"])); err != nil {
			return err
		}
		if err := os.Setenv("PRIVATE_KEY", string(cryptoSecret.Data["PRIVATE_KEY"])); err != nil {
			return err
		}
		if err := os.Setenv("ENCRYPTION_KEYS", strings.Join(encryptionKeys, " ")); err != nil {
			return err
		}
	}
	return nil
}
