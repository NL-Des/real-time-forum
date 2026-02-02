package auth

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID int) (string, error) {
	token := uuid.NewString()
	now := time.Now()
	expireAt := now.Add(time.Hour * 24)

	_, err := db.Exec("INSERT INTO session (UserID, Token, CreatedAt, ExpiresAt) VALUES (?, ?, ?, ?);", userID, token, now, expireAt)
	if err != nil {
		return "", err
	}

	return token, nil
}
