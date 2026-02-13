package posts

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/shared"
	"strconv"
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
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "post created",
			"post":    newPost,
		})
	}

}

func DisplayPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		referer := r.Header.Get("Referer")
		if referer == "" {
			referer = "/"
		}

		//reception id post
		idStr := r.URL.Query().Get("id")

		postId, err := strconv.Atoi(idStr)
		if err != nil || postId < 1 {
			log.Printf("Erreur conversion ou id invalide: %v, retour à la page précédente", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}

		// récuperer toutes les données du post

		currentPost, currentCmts, err := GetPostData(db, postId)
		if err != nil {
			log.Printf("Erreur récupération des données du post: %v", err)
			shared.RespondError(w, http.StatusInternalServerError, err)
			return
		}

		//renvoyer les données à afficher
		response := PostWithCommentsResponse{
			Post:     currentPost,
			Comments: currentCmts,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
