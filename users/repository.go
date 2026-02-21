package users

import "database/sql"

// Exporté : Repository
type Repository struct {
	DB *sql.DB
}

// Exporté : méthode GetOnlineUsers
func (r *Repository) GetOnlineUsers() ([]string, error) {
	rows, err := r.DB.Query(`
        SELECT u.UserName
        FROM users u
        JOIN session s ON s.UserID = u.id
        WHERE s.ExpiresAt > CURRENT_TIMESTAMP
    `)
	if err != nil {
		return nil, err
	}

	var users []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		users = append(users, name)
	}

	return users, nil
}
