package models

import (
	"github.com/nuntiodev/hera-sdks/go_hera"
	"github.com/nuntiodev/x/cryptox"
)

type Hash struct {
	Variant go_hera.HasingAlgorithm    `json:"Variant"`
	Body    string                     `json:"body"`
	Params  map[string]cryptox.Stringx `json:"Params"`
}

func HashToProto(hash *Hash) *go_hera.Hash {
	if hash == nil {
		return nil
	}
	params := map[string]string{}
	for key, val := range hash.Params {
		params[key] = val.Body
	}
	return &go_hera.Hash{
		Variant: hash.Variant,
		Body:    hash.Body,
		Params:  params,
	}
}

func ProtoToHash(hash *go_hera.Hash) *Hash {
	if hash == nil {
		return nil
	}
	params := map[string]cryptox.Stringx{}
	for key, val := range hash.Params {
		params[key] = cryptox.Stringx{Body: val}
	}
	return &Hash{
		Variant: hash.Variant,
		Body:    hash.Body,
		Params:  params,
	}
}
