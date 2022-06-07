package config_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func prepare(action int, config *go_block.Config) {
	if config == nil {
		return
	}
	switch action {
	case actionCreate:
		config.CreatedAt = ts.Now()
		config.UpdatedAt = ts.Now()
	case actionUpdate:
		config.UpdatedAt = ts.Now()
	}
}
