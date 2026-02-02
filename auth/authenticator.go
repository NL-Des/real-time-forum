package auth

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

// CheckCredentials vérifie si le username et le mot de passe sont corrects
func CheckCredentials(db *sql.DB, username string, password string) (bool, error) {
	var hashedPassword string

	// On récupère le mot de passe hashé depuis la base
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		return false, err // utilisateur non trouvé ou erreur DB
	}

	// On compare le mot de passe fourni avec le hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil // mot de passe incorrect
	}

	return true, nil // identifiants corrects
}
