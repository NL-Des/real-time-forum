package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"real-time-forum/auth"
	"real-time-forum/posts"
	"real-time-forum/users"
)

// Commande de lancement depuis la racine : go run cmd/main.go

func Server(port string, db *sql.DB) {
	mux := http.NewServeMux() // Création d'un serveur mux vide.

	mux.HandleFunc("/auth/login", auth.LoginHandler(db))
	mux.HandleFunc("/post", posts.NewPostHandler(db))

	// Quand l'utilisateur arrive, affiche mainPage.
	mux.HandleFunc("/", users.MainPage)
	// servir les fichiers static
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", fs))
	// Lancement serveur Go
	fmt.Printf("Serveur démarré sur le port %s \n", port)
	fmt.Printf("Page d'accès : http://localhost:8081/ \n")
	err := http.ListenAndServe(port, mux) // A laisser à la fin, élément bloquant la lecture des instructions suivantes.
	if err != nil {
		fmt.Printf("Error lauching servor : %v \n", err)
	}
}
