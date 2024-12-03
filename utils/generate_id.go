package utils

import (
	"crypto/rand"
	"encoding/base32"
	"log"
)

func GenerateID() (string, error) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("⚠️ Error generating random bytes: %v", err)
		return "", err
	}

	id := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	if len(id) < 26 {
		log.Printf("❌ Generated ID is too short: %s (length: %d)", id, len(id))
		return "", err
	}

	return id[:26], nil
}
