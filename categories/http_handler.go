package categories

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("frontend/index.html"))

func MainPage(w http.ResponseWriter, r *http.Request) {
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
