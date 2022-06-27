package runner

import (
	"context"
	"fmt"
	"github.com/nuntiodev/hera/initializer"
	"github.com/nuntiodev/hera/server"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var (
	initializeSecrets = false
	initializeEngine  = ""
)

const (
	version = 1
)

func initialize() error {
	initializeSecrets, _ = strconv.ParseBool(os.Getenv("INITIALIZE_SECRETS"))
	initializeEngine = os.Getenv("INITIALIZE_ENGINE")
	return nil
}

func Run(ctx context.Context, logger *zap.Logger) error {
	logger.Info(fmt.Sprintf("running Hera version: %d", version))
	if err := initialize(); err != nil {
		return err
	}
	if initializeSecrets {
		myInitializer, err := initializer.New(logger, initializeEngine)
		if err != nil {
			return err
		}
		if err := myInitializer.CreateSecrets(ctx); err != nil {
			return err
		}
	}
	logger.Info("runner is initializing the application...")
	serve, err := server.New(ctx, logger)
	if err != nil {
		return err
	}
	return serve.GrpcServer.Run()
}
