package runner

import (
	"context"
	"github.com/joho/godotenv"
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

func initialize() error {
	initializeSecrets, _ = strconv.ParseBool(os.Getenv("INITIALIZE_SECRETS"))
	initializeEngine = os.Getenv("INITIALIZE_ENGINE")
	return nil
}

func Run(ctx context.Context, zapLog *zap.Logger) error {
	if err := godotenv.Load(".env"); err != nil {
		zapLog.Warn("could not get .env")
	}
	if err := initialize(); err != nil {
		return err
	}
	if initializeSecrets {
		myInitializer, err := initializer.New(zapLog, initializeEngine)
		if err != nil {
			return err
		}
		if err := myInitializer.CreateSecrets(ctx); err != nil {
			return err
		}
	}
	zapLog.Info("runner is initializing the application...")
	serve, err := server.New(ctx, zapLog)
	if err != nil {
		return err
	}
	return serve.GrpcServer.Run()
}
