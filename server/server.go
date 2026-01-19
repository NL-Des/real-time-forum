package server

import (
	"fmt"
	"net/http"
	"text/template"
)

func Server(port string) { // Première lettre en majuscule pour que la fonction soit publique.

	// Définition de la route de l'unique page.
	mainPage := func(w http.ResponseWriter, r *http.Request) {
		// Lecture du fichier HTML
		tmpl, err := template.ParseFiles("../frontend/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // Affiche sur le site.
			fmt.Printf("Error Parsing HTML Page: %v \n", err)          // Affiche sur la console.
			return
		}
		// Envoi le fichier HTML au w http.ResponseWriter.
		// Ce qui l'envoie et l'affiche sur le site.
		tmpl.Execute(w, nil)
	}

	// Si quelqu'un va sur mainPage, affiche mainPage
	http.HandleFunc("/mainPage", mainPage)

	// Lancement serveur Go
	fmt.Printf("Serveur démarré sur le port %s \n", port)
	err := http.ListenAndServe(port, nil) // A laisser à la fin, élément bloquant la lecture des instructions suivantes.
	if err != nil {
		fmt.Printf("Error lauching servor : %v \n", err)
	}
}
