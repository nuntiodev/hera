package runner

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/softcorp-io/block-user-service/server"
	"go.uber.org/zap"
)

func Run(ctx context.Context, zapLog *zap.Logger) error {
	if err := godotenv.Load(".env"); err != nil {
		zapLog.Warn("could not get .env")
	}
	zapLog.Info("runner is initializing the application...")
	serve, err := server.New(ctx, zapLog)
	if err != nil {
		return err
	}
	return serve.GrpcServer.Run()
}
