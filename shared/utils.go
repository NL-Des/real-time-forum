package shared

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Chiffrement des mots de passes.
func Encryption(password []byte) (hashPassword []byte) {
	// Prend "password" pour le chiffrer.
	hashPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Error during encryption : %w", err)
	}

	// Test à effacer par la suite.
	fmt.Println(string(hashPassword))
	return hashPassword
}
