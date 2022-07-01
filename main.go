package main

import (
	"context"
	"github.com/nuntiodev/hera/runner"
	"log"
)

func main() {
	log.Fatal(runner.Run(context.Background()).Error())
}
