package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func NewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost Post

		// -- récupération données --
		if r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		log.Println(newPost.CategoryIDs)
		// Récupérer l'auteur grâce au cookie

		// -- gestion d'erreurs --
		// ajouter verif post existant

		if err := IsValid(newPost); err != nil {
			log.Println("<NewPostHandler> Error invalid post: ", err)
			RespondError(w, http.StatusBadRequest, err)
			return
		}

		// -- sauvegarde db --

		if err := SavePost(db, &newPost); err != nil {
			log.Println("<NewPostHandler> Error cannot save post: ", err)
			RespondError(w, http.StatusInternalServerError, err)
			return
		}

		if err := SavePostCategories(db, &newPost); err != nil {
			log.Println("<NewPostHandler> Error cannot save post categories: ", err)
			RespondError(w, http.StatusInternalServerError, err)
			return
		}

		// --- notification post ajouté ---
		log.Println("New post saved into db")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "post created",
		})
	}

}

// Fonction renvoi message d'erreur au front
func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
