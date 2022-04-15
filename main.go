package main

import (
	"context"
	"github.com/nuntio-dev/nuntio-user-block/runner"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zapLog.Fatal(runner.Run(context.Background(), zapLog).Error())
}
