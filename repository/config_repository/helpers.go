package config_repository

import (
	"github.com/nuntiodev/hera-proto/go_hera"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

func prepare(action int, config *go_hera.Config) {
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
