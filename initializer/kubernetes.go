package initializer

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"github.com/fatih/color"
	"github.com/nuntiodev/x/cryptox"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
)

const (
	heraSecretName = "hera-secret"
)

type kubernetesInitializer struct {
	namespace string
	zapLog    *zap.Logger
	redLog    *color.Color
	blueLog   *color.Color
	k8s       *kubernetes.Clientset
}

func (i *kubernetesInitializer) CreateSecrets(ctx context.Context) error {
	i.blueLog.Println("Hera is running in Kubernetes-engine mode...")
	if os.Getenv("PUBLIC_KEY") != "" && os.Getenv("PRIVATE_KEY") != "" && os.Getenv("ENCRYPTION_KEYS") != "" {
		i.blueLog.Println("secrets are already provided internally by the system (the PUBLIC_KEY, ENCRYPTION_KEYS and PRIVATE_KEY variable is set).")
		return nil
	}
	// check if secret already exists
	if cryptoSecret, err := i.k8s.CoreV1().Secrets(i.namespace).Get(ctx, heraSecretName, metav1.GetOptions{}); err != nil {
		i.zapLog.Info("hera secret does not exist... creating....")
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
		// create encryption secret
		encryptionSecret, err := cryptox.GenerateSymmetricKey(32, cryptox.AlphaNum)
		if err != nil {
			return err
		}
		// create secret in the EngineKubernetes api
		secretData := map[string]string{
			"PRIVATE_KEY":     heraPrivateKey,
			"PUBLIC_KEY":      heraPublicKey,
			"ENCRYPTION_KEYS": encryptionSecret,
		}
		secret := &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: i.namespace,
				Name:      heraSecretName,
			},
			StringData: secretData,
			Type:       v1.SecretTypeOpaque,
		}
		if _, err := i.k8s.CoreV1().Secrets(i.namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return err
		}
		// set in memory
		if err := os.Setenv("PUBLIC_KEY", heraPublicKey); err != nil {
			return err
		}
		if err := os.Setenv("PRIVATE_KEY", heraPrivateKey); err != nil {
			return err
		}
		if err := os.Setenv("ENCRYPTION_KEYS", encryptionSecret); err != nil {
			return err
		}
		i.zapLog.Info("successfully created secret")
	} else {
		i.zapLog.Info("hera secret already exists")
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
					Name:      heraSecretName,
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
