package main

import (
	"context"
	"github.com/softcorp-io/block-user-service/runner"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zapLog.Fatal(runner.Run(context.Background(), zapLog).Error())
}
