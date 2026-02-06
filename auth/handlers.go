package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Json reçu depuis le frontend
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// représente le JSON renvoyé au client
type LoginResponse struct {
	Success bool `json:"success"`
	User    struct {
		Nickname string `json:"nickname"`
	} `json:"user"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Vérification des identifiants avec authenticator.go
		ok, userID, err := CheckCredentials(db, req.Login, req.Password)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("DB error:", err)
			return
		}

		res := LoginResponse{}

		if ok {
			// Vérifier si une session active existe déjà
			existingToken := GetActiveSessionToken(db, userID)

			var token string

			if existingToken != "" {
				token = existingToken
				fmt.Println("Token déjà existant pour cet utilisateur")
			} else {
				// Création de la session
				token, err = CreateSession(db, userID)
				if err != nil {
					http.Error(w, "Impossible de créer a session", http.StatusInternalServerError)
					return
				}
			}

			_, err = db.Exec("UPDATE users SET userOnline = 1 WHERE id = ?", userID)
			if err != nil {
				http.Error(w, "Erreur mise à jour statut utilisateur", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    token,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})

			res.Success = true
			res.User.Nickname = req.Login

		} else {
			res.Success = false
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Lire le cookie session_token
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// Cookie absent = l'utilisateur est déjà déconnecté
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Déconnecté"))
			return
		}

		// Récupérer le UserID associé au token
		var userID int
		err = db.QueryRow("SELECT UserID FROM session WHERE Token = ?", cookie.Value).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Session inexistante = déjà déconnecté
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Déconnecté"))
				return
			}
			// Erreur BDD
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}

		// Mettre userOnline à 0
		_, err = db.Exec("UPDATE users SET userOnline = 0 WHERE id = ?", userID)
		if err != nil {
			http.Error(w, "Impossible de mettre à jour le statut", http.StatusInternalServerError)
			return
		}

		// Supprimer la session
		/* _, err = db.Exec("DELETE FROM session WHERE Token = ?", cookie.Value)
		if err != nil {
			http.Error(w, "Impossible de supprimer la session", http.StatusInternalServerError)
			return
		} */

		// Supprimer le cookie côté client
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // expiré
			HttpOnly: true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Déconnecté"))
	}
}
