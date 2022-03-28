package repository_mock

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/softcorp-io/block-user-service/crypto"
	"github.com/softcorp-io/block-user-service/repository"
	database "github.com/softcorp-io/softcorp_db_helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net"
	"os"
	"strconv"
	"time"
)

/*
	NewUserRepoMock spin up a user repository
	by first creating a MongoDB Docker instance
	and connecting the user repository  to that
	MongoDB instance.
*/
func NewRepositoryMock(ctx context.Context, zapLog *zap.Logger, containerName string) (repository.Repository, *dockertest.Pool, *dockertest.Resource, error) {
	// create the pool (docker instance).
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, nil, err
	}
	// remove old containers
	if err := pool.RemoveContainerByName(containerName); err != nil {
		return nil, nil, nil, err
	}
	// get random free port
	mongoPort := ""
	for {
		ln, err := net.Listen("tcp", ":"+"0")
		if err == nil {
			mongoPort = strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)
			ln.Close()
			break
		}
		ln.Close()
	}
	// start the container.
	container, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "mongo",
		Name:         containerName,
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: mongoPort},
			},
		},
	})
	if err != nil {
		return nil, nil, nil, err
	}
	// setup environment
	mongoUri := "mongodb://localhost:" + mongoPort
	os.Setenv("MONGO_URI", mongoUri)
	// check db connection and create mongo client
	var mongoClient *mongo.Client
	if err = pool.Retry(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		db, err := database.CreateDatabase(zapLog)
		client, err := db.CreateMongoClient(ctx)
		if err != nil {
			return err
		}
		mongoClient = client
		return nil
	}); err != nil {
		if err := pool.Purge(container); err != nil {
			zapLog.Fatal(fmt.Sprintf("failed to purge pool with err: %s", err))
		}
		if err := pool.RemoveContainerByName(containerName); err != nil {
			zapLog.Fatal(fmt.Sprintf("failed to remove Docker container with err: %s", err))
		}
		return nil, pool, container, err
	}
	// create the repository_mock
	// generate private and public keys
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, pool, container, err
	}
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}
	publickey := &privatekey.PublicKey
	marshalPublicKey, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		return nil, pool, container, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: marshalPublicKey,
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyBlock)
	publicKeyBytes := pem.EncodeToMemory(publicKeyBlock)
	myCrypto, err := crypto.New(privateKeyBytes, publicKeyBytes)
	os.Setenv("JWT_PRIVATE_KEY", string(privateKeyBytes))
	os.Setenv("JWT_PUBLIC_KEY", string(publicKeyBytes))
	repo, err := repository.New(mongoClient, myCrypto, zapLog)
	if err != nil {
		if err := pool.Purge(container); err != nil {
			zapLog.Fatal(fmt.Sprintf("failed to purge pool with err: %s", err))
		}
		if err := pool.RemoveContainerByName(containerName); err != nil {
			zapLog.Fatal(fmt.Sprintf("failed to remove Docker container with err: %s", err))
		}
		return nil, pool, container, err
	}
	return repo, pool, container, nil
}
