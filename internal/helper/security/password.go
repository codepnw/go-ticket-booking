package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GenenrateHashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password length be at least 6 characters long")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password hash failed")
	}

	return string(hashed), nil
}

func VerifyPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password does not matchh")
	}

	return nil
}
