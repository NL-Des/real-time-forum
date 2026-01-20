package main

import (
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
	}
	// Attribution du port du serveur.
	// port := os.Getenv("server_port")
	// Attribution du chemin de la database.
	// pathDB := os.Getenv("path_to_database")

	// Lancement de la BDD.
	// Vérification si la BDD existe déjà.
	pathDB := "internal/config-database/initdb.go"
	_, err := os.ReadFile(pathDB)
	if err != nil {
		config.InitDB(pathDB) // Création BDD
	} else {
		config.RunDB(pathDB) // Lancement BDD
	}

	// Lancement du serveur GO.
	port := ":8080"
	server.Server(port) // server = nom du package | Server() = nom de la fonction

}
