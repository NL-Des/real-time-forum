package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"real-time-forum/internal/config-database"
	"real-time-forum/server" // Ceci permet d'appeler la fonction qui se trouve dans le fichier.

	"github.com/joho/godotenv"
)

func main() {
	// MARK: .ENV
	// Chargement des données du .env
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
		return
	}
	// Attribution du port du serveur.
	// port := os.Getenv("server_port")
	// Attribution du chemin de la database.
	// pathDB := os.Getenv("path_to_database")

	// Lancement de la BDD.
	// Vérification si la BDD existe déjà.

	// Créer le dossier vault s'il n'existe pas
	if err := os.MkdirAll("./vault", 0755); err != nil {
		log.Fatalf("Error creating vault directory: %v", err)
	}
	pathDB := "./vault/real_time_forum_database.db"
	var db *sql.DB
	var err error

	_, statErr := os.Stat(pathDB) // Stat repère si le fichier existe, sans le charger.
	if statErr != nil {
		fmt.Println("Initialazing Database...")
		db, err = config.InitDB(pathDB) // Création BDD
		fmt.Println("Connection to Database...")
	} else {
		fmt.Println("Connection to Database...")
		db, err = config.RunDB(pathDB) // Lancement BDD
	}
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	// Lancement du serveur GO.
	port := ":8080"
	server.Server(port, db)
}
