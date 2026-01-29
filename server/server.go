package server

import (
	"fmt"
	"net/http"
	"real-time-forum/categories"
	"real-time-forum/posts"

	_ "github.com/mattn/go-sqlite3"
)

// Commande de lancement depuis la racine : go run cmd/main.go

func Server(port string) {
	mux := http.NewServeMux() // Création d'un serveur mux vide.

	// Quand l'utilisateur arrive, affiche mainPage.
	mux.HandleFunc("/", categories.MainPage)
	// servir les fichiers static
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", fs))
	mux.HandleFunc("/newpost", posts.NewPostHandler)
	// Lancement serveur Go
	fmt.Printf("Serveur démarré sur le port %s \n", port)
	fmt.Printf("Page d'accès : http://localhost:8080/ \n")
	err := http.ListenAndServe(port, mux) // A laisser à la fin, élément bloquant la lecture des instructions suivantes.
	if err != nil {
		fmt.Printf("Error lauching servor : %v \n", err)
	}
}
