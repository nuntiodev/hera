package measurement_repository

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func getUserHash(userId string) string {
	userId = strings.TrimSpace(userId)
	hash := sha256.New()
	hash.Write([]byte(userId))
	userShaHash := hash.Sum(nil)
	return hex.EncodeToString(userShaHash)
}
