package models

type Bcrypt struct {
	Cost int `bson:"cost" json:"cost"`
}
