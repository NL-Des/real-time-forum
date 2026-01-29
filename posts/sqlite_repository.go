package posts

import "database/sql"

func SavePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO post (title, content, authorid) VALUES (?, ?, ?)`
	_, err := db.Exec(query, post.Title, post.Content, post.AuthorID)
	return err
}
