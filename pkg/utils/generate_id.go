package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"orion/pkg/logger"
)

func GenerateID() (string, error) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		logger.Debug("⚠️ Error generating random bytes: " + err.Error())
		return "", err
	}

	id := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	if len(id) < 26 {
		logger.Error(fmt.Sprintf("❌ Generated ID is too short: %s (length: %d)", id, len(id)), nil)

		return "", fmt.Errorf("generated ID is too short: %s (length: %d)", id, len(id))
	}

	return id[:26], nil
}
