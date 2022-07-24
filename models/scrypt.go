package models

import "github.com/nuntiodev/x/cryptox"

type Scrypt struct {
	SignerKey     cryptox.Stringx `bson:"signer_key" json:"signer_key"`
	SaltSeparator cryptox.Stringx `bson:"salt_separator" json:"salt_separator"`
	Rounds        int             `bson:"rounds" json:"rounds"`
	MemCost       int             `bson:"mem_cost" json:"mem_cost"`
	P             int             `bson:"p" json:"p"`
	KeyLen        int             `bson:"key_len" json:"key_len"`
}

type ScryptHera struct {
	SignerKey     string `bson:"signer_key" json:"signer_key"`
	SaltSeparator string `bson:"salt_separator" json:"salt_separator"`
	Rounds        int    `bson:"rounds" json:"rounds"`
	MemCost       int    `bson:"mem_cost" json:"mem_cost"`
	P             int    `bson:"p" json:"p"`
	KeyLen        int    `bson:"key_len" json:"key_len"`
}
