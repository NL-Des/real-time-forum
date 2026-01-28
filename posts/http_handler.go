package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var newPost Post
	referer := r.Referer()
	if referer == "" {
		referer = "/"
	}

	// -- récupération données --
	if r.Method == "POST" {
		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// -- gestion d'erreurs --
	// ajouter verif post existant

	if err := IsValid(newPost); err != nil {
		log.Println("<NewPostHandler> Error invalid post: ", err)
		RespondError(w, http.StatusBadRequest, err)
	}

	// -- sauvegarde db --

	if err := SavePost(db, &newPost); err != nil {
		log.Println("<NewPostHandler> Error cannot save post: ", err)
		RespondError(w, http.StatusInternalServerError, err)
	}

	//notif post ajouté
	log.Println("New post saved into db")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

//Prochaine étape : faire un post.js pour vérifier la bonne récéption des données ici

// Fonction renvoi message d'erreur dans le frontq
func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
