package config

import (
	"database/sql"

	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Lancement de la BDD déjà existante.
func RunDB(pathDB string) {
	// Ouverture de la BDD
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Fatal("Error opening database :", err)
	}
	// Vérification de la connexion, car sql.Open ne le fais pas.
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to database", err)
	}
}

// Création de la BDD si elle n'existe pas.
func InitDB(pathDB string) {
	// Ouverture de la BDD
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Fatal("Error opening database :", err)
	}
	// Vérification de la connexion, car sql.Open ne le fais pas.
	if err = db.Ping(); err != nil {
		log.Fatal("Error connecting to database", err)
	}
	// Création de la BDD à partir des différentes tables.
}
