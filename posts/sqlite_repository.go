package posts

import "database/sql"

func SavePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO posts (title, content, category_id) VALUES (?, ?, ?)`
	_, err := db.Exec(query, post.Title, post.Content, post.CategoryId)
	return err
}
