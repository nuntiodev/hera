package main

import (
	"context"
	"github.com/fatih/color"
	"github.com/nuntiodev/hera/runner"
	"go.uber.org/zap"
	"log"
)

func main() {
	color.New(color.FgHiBlue).Println("Setting up your system in Kubernetes")
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zapLog.Fatal(runner.Run(context.Background(), zapLog).Error())
}
