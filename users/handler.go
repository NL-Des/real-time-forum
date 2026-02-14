package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

// ✅ Structure pour l'API (version simplifiée de User)
type UserAPI struct {
	ID         int    `json:"id"`
	UserName   string `json:"nickname"`
	UserOnline int    `json:"online"`
}

// ✅ Handler pour récupérer tous les utilisateurs
func GetAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("📡 Requête /api/users reçue")

		// ✅ Vérifier la méthode HTTP
		if r.Method != http.MethodGet {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		// ✅ Requête SQL
		rows, err := db.Query(`
            SELECT id, UserName, userOnline 
            FROM users 
            ORDER BY userOnline DESC, UserName ASC
        `)
		if err != nil {
			log.Printf("❌ Erreur SQL: %v\n", err)
			http.Error(w, "Erreur serveur", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// ✅ Parcourir les résultats
		var users []UserAPI
		for rows.Next() {
			var user UserAPI
			err := rows.Scan(&user.ID, &user.UserName, &user.UserOnline)
			if err != nil {
				log.Printf("⚠️ Erreur scan: %v\n", err)
				continue
			}
			users = append(users, user)
		}

		log.Printf("✅ %d utilisateurs trouvés\n", len(users))

		// ✅ Renvoyer le JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("❌ Erreur encodage JSON: %v\n", err)
			http.Error(w, "Erreur encodage", http.StatusInternalServerError)
		}
	}
}
