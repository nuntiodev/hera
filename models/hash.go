package models

import (
	"github.com/nuntiodev/hera-sdks/go_hera"
)

type Hash struct {
	Variant go_hera.HasingAlgorithm `json:"Variant"`
	Body    string                  `json:"body"`
	Params  map[string]string       `json:"Params"`
}

func HashToProto(hash *Hash) *go_hera.Hash {
	if hash == nil {
		return nil
	}
	return &go_hera.Hash{
		Variant: hash.Variant,
		Body:    hash.Body,
		Params:  hash.Params,
	}
}

func ProtoToHash(hash *go_hera.Hash) *Hash {
	if hash == nil {
		return nil
	}
	return &Hash{
		Variant: hash.Variant,
		Body:    hash.Body,
		Params:  hash.Params,
	}
}
