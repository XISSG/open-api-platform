package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func AccessKeyGenerator() (string, error) {
	accessKey, err := generateRandomKey(16)
	if err != nil {
		return "", err
	}
	return accessKey, nil
}

func SecretKeyGenerator() (string, error) {
	secretKey, err := generateRandomKey(32)
	if err != nil {
		return "", err
	}
	return secretKey, nil
}

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
