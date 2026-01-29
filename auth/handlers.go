package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// Json reçu depuis le frontend
type LoginRequest struct {
	Username string `json:"username"`
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
		ok, err := CheckCredentials(db, req.Username, req.Password)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("DB error:", err)
			return
		}

		res := LoginResponse{}

		if ok {
			res.Success = true
			res.User.Nickname = req.Username
		} else {
			res.Success = false
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
