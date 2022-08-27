package utils

import "golang.org/x/crypto/bcrypt"

func CompareHashedKeys(hashedKey, plainKey string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedKey), []byte(plainKey)); err != nil {
		return false
	}
	return true
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
