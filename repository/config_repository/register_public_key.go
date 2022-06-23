package config_repository

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (c *defaultConfigRepository) RegisterPublicKey(ctx context.Context, rsaPublicKey string) error {
	pubBlock, _ := pem.Decode([]byte(rsaPublicKey))
	_, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return err
	}
	mongoUpdate := bson.M{
		"$set": bson.M{
			"public_key": rsaPublicKey,
			"updated_at": time.Now(),
		},
	}
	if _, err := c.collection.UpdateOne(ctx, bson.M{"_id": namespaceConfigName}, mongoUpdate); err != nil {
		return err
	}
	return nil
}
