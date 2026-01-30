package auth

import (
	"database/sql"
	"encoding/json"
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
			// Création de la session
			token, err := CreateSession(db, userID)
			if err != nil {
				http.Error(w, "Impossible de créer a session", http.StatusInternalServerError)
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
