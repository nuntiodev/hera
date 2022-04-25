package config_repository

import (
	"github.com/nuntiodev/block-proto/go_block"
	ts "google.golang.org/protobuf/types/known/timestamppb"
	"strings"
)

func prepare(action int, config *go_block.Config) {
	if config == nil {
		return
	}
	switch action {
	case actionCreate:
		config.CreatedAt = ts.Now()
		config.UpdatedAt = ts.Now()
	}
	config.Name = strings.TrimSpace(config.Name)
	config.Website = strings.TrimSpace(config.Website)
	config.About = strings.TrimSpace(config.About)
	config.Email = strings.TrimSpace(config.Email)
	config.Logo = strings.TrimSpace(config.Logo)
	config.Terms = strings.TrimSpace(config.Terms)
}
