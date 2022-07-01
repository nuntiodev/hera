package runner

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nuntiodev/hera/initializer"
	"github.com/nuntiodev/hera/server"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"strconv"
)

var (
	initializeSecrets = false
	initializeEngine  = ""
	logMode           = ""
)

const (
	version = 1
)

func initialize() error {
	initializeSecrets, _ = strconv.ParseBool(os.Getenv("INITIALIZE_SECRETS"))
	initializeEngine = os.Getenv("INITIALIZE_ENGINE")
	return nil
}

func Run(ctx context.Context) error {
	var logger *zap.Logger
	var err error
	if err := godotenv.Load(".env"); err != nil {
		log.Println("could not get .env")
	}
	logMode = os.Getenv("LOG_MODE")
	if logMode == "prod" {
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
	} else if logMode == "" || logMode == "dev" {
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("invalid log mode %s, supported log-modes: dev & prod", logMode)
	}
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
	errGroup := errgroup.Group{}
	if serve.HttpServer != nil {
		if serve.GrpcServer == nil {
			return serve.HttpServer.Run()
		} else {
			errGroup.Go(func() error {
				return serve.HttpServer.Run()
			})
		}
	}
	if serve.GrpcServer != nil {
		return serve.GrpcServer.Run()
	}
	return errGroup.Wait()
}
