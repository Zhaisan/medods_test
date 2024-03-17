package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRefreshToken(userID string) string {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	token := fmt.Sprintf("%s.%s", userID, base64.URLEncoding.EncodeToString(randomBytes))
	return token
}
