package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/shared"
)

func NewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPost Post

		// -- récupération données --
		if r.Method == "POST" {
			if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				shared.RespondError(w, http.StatusBadRequest, err)
				return
			}
		}

		// Récupérer l'auteur grâce au cookie

		// -- gestion d'erreurs --
		// ajouter verif post existant

		if err := IsValidFormat(newPost); err != nil {
			log.Println("<NewPostHandler> Error invalid post: ", err)
			shared.RespondError(w, http.StatusBadRequest, err)
			return
		}

		// -- sauvegarde db --

		if err := SavePost(db, &newPost); err != nil {
			log.Println("<NewPostHandler> Error cannot save post: ", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		if err := SavePostCategories(db, &newPost); err != nil {
			log.Println("<NewPostHandler> Error cannot save post categories: ", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
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
