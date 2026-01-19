package main

import (
	"log"
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

	// Lancement du serveur GO.
	server.Server(":8080") // server = nom du package | Server() = nom de la fonction

}
