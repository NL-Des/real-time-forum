package server

import (
	"fmt"
	"net/http"
	"text/template"
)

// Commande de lancement depuis la racine : go run cmd/main.go

var tmpl = template.Must(template.ParseFiles("frontend/index.html"))

func Server(port string) {
	mux := http.NewServeMux() // Création d'un serveur mux vide.

	mainPage := func(w http.ResponseWriter, r *http.Request) {
		// La déclaration à l'extérieur de "Server" évite de refaire à chaque fois ParseFiles.
		// Mais nous perdons la déclaration d'erreur personnalisée.
		/* 		tmpl, err := template.ParseFiles("../frontend/index.html")
		   		if err != nil {
		   			http.Error(w, err.Error(), http.StatusInternalServerError) // Affiche sur le site.
		   			fmt.Printf("Error Parsing HTML Page: %v \n", err)          // Affiche sur la console.
		   			return
		   		} */
		tmpl.Execute(w, nil)
	}
	// Quand l'utilisateur arrive, affiche mainPage.
	mux.HandleFunc("/", mainPage)

	// Lancement serveur Go
	fmt.Printf("Serveur démarré sur le port %s \n", port)
	fmt.Printf("Page d'accès : http://localhost:8080/ \n")
	err := http.ListenAndServe(port, mux) // A laisser à la fin, élément bloquant la lecture des instructions suivantes.
	if err != nil {
		fmt.Printf("Error lauching servor : %v \n", err)
	}
}
