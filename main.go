package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/nuntiodev/hera/runner"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	logMode = ""
)

func main() {
	var zapLog *zap.Logger
	var err error
	if err := godotenv.Load(".env"); err != nil {
		log.Println("could not get .env")
	}
	logMode = os.Getenv("LOG_MODE")
	if logMode == "prod" {
		zapLog, err = zap.NewProduction()
		if err != nil {
			log.Fatal(err)
		}
	} else if logMode == "" || logMode == "dev" {
		zapLog, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatalf("invalid log mode %s, supported log-modes: dev & prod", logMode)
	}
	zapLog.Fatal(runner.Run(context.Background(), zapLog).Error())
}
