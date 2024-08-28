package encrypt

import "golang.org/x/crypto/bcrypt"

func GenPasswordHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func ValidatePasswordHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}
	return true
}
