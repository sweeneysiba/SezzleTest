package middleware

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bcrypt Password", string(hashedPassword))
	return string(hashedPassword)
}

func VerifyBcrypt(db_password, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(db_password), []byte(password))
	if err != nil {
		return false
	}
	return true
}
