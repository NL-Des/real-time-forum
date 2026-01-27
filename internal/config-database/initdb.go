package config

import (
	"database/sql"
	"fmt"

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
func InitDB(pathDB string, schemaSQL string) (*sql.DB, error) {
	// Ouverture de la BDD
	// Exécuter le schéma SQL complet

	// Ouvrir/créer la base de données
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de la base de données:", err)
		return nil, err
	}

	// Vérifier la connexion
	if err = db.Ping(); err != nil {
		db.Close()
		fmt.Println("Erreur de connexion à la base de données:", err)
		return nil, err
	}
	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		db.Close()
		fmt.Println("Erreur lors de l'exécution du schéma SQL:", err)
		return nil, err
	}

	fmt.Println("Base de données créée avec succès!")
	return db, nil
}
