package posts

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	// referer := r.Referer()
	// if referer == "" {
	// 	referer = "/"
	// }

	// -- récupération données --
	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// -- gestion d'erreurs --

	if err := IsValid(newPost); err != nil {
		log.Println("<NewPostHandler> Error invalid post: ", err)
		RespondError(w, http.StatusBadRequest, err)
	}

	//sauvegarde db

	//notif post ajouté
}

// Fonction renvoi message d'erreur dans le front
func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
