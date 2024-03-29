package util

import (
	"crypto/rand"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GenerateRandomHex(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return BytesToHex(bytes), nil
}

func BytesToHex(bytes []byte) string {
	return fmt.Sprintf("%x", bytes)
}

func EncryptPassword(password string) (string, error) {
	return argon2id.CreateHash("pa$$word", argon2id.DefaultParams)
}

func VerifyPassword(password, encryptedPassword string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, encryptedPassword)
}
