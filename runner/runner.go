package runner

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/nuntiodev/nuntio-user-block/initializer"
	"github.com/nuntiodev/nuntio-user-block/server"
	"go.uber.org/zap"
	"os"
	"strconv"
)

var (
	initializeSecrets = false
	initializerEngine = ""
)

func initialize() error {
	initializeSecrets, _ = strconv.ParseBool(os.Getenv("INITIALIZE_SECRETS"))
	initializerEngine = os.Getenv("INITIALIZER_ENGINE")
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
		myInitializer, err := initializer.New(zapLog, initializerEngine)
		if err != nil {
			return err
		}
		if err := myInitializer.CreateRsaSecrets(ctx); err != nil {
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
