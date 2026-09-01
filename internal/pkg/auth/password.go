package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash password: %w", err)
	}

	return string(hashed), nil
}

func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func IsHashedPassword(password string) bool {
	_, err := bcrypt.Cost([]byte(password))
	return err == nil
}
