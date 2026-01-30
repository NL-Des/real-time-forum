package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// MARK: DB opening
// Lancement de la BDD déjà existante.
func RunDB(pathDB string) (*sql.DB, error) {
	// Ouverture ou création de la BDD
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		log.Fatalf("Error opening database : %v", err)
	}
	// Vérification de la connexion, car sql.Open ne le fais pas.
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database : %v", err)
	}
	return db, nil
}

// MARK: DB Creation
// Création de la BDD si elle n'existe pas.
func InitDB(pathDB string, db *sql.DB) (*sql.DB, error) {
	// Lecture des fichiers contenant les tables.
	// Ouverture du fichier.
	testTables, err := os.ReadFile("./internal/config-database/001_create_tables.sql")
	if err != nil {
		log.Fatalf("Error with testTables in internal/config-database/001_create_tables.sql : %v", err)
	}
	// Ecriture des tables dans la BDD.
	_, err = db.Exec(string(testTables))
	if err != nil {
		log.Fatalf("Error writing testTables in DB : %v", err)
	}
	// Importation des informations de bases de la BDD (Catégories,...)
	seedTables, err := os.ReadFile("./internal/config-database/002_seed.sql")
	if err != nil {
		log.Fatalf("Error with seedTables in internal/config-database/002_seed.sql : %v", err)
	}
	_, err = db.Exec(string(seedTables))
	if err != nil {
		log.Fatalf("Error writing seedTables in DB : %v", err)
	}
	return db, nil
}
