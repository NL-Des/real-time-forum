package posts

import (
	"database/sql"
	"fmt"
	"real-time-forum/comments"
)

// --- Fonctions Sauvegarde de Data

func SavePost(db *sql.DB, post *Post) error {
	query := `INSERT INTO post (title, content, authorid) VALUES (?, ?, ?)`
	_, err := db.Exec(query, post.Title, post.Content, post.AuthorID)
	return err
}

func SavePostCategories(db *sql.DB, post *Post) error {
	postWithId, err := GetPostID(db, post)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(
		`INSERT INTO post_categories (postID, categoryID) VALUES(?, ?)`,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, catId := range post.CategoryIDs {
		_, err := stmt.Exec(postWithId.ID, catId)
		fmt.Println(postWithId.ID, catId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// --- Fonctions Récupération de Data ---

func GetPostID(db *sql.DB, post *Post) (*Post, error) {
	query := `SELECT id FROM post WHERE title = ? AND content = ? AND authorid = ?`
	row := db.QueryRow(query, post.Title, post.Content, post.AuthorID)
	err := row.Scan(&post.ID)
	if err != nil {
		return &Post{}, err
	}
	return post, nil
}

func GetPostData(db *sql.DB, id int) (Post, []comments.Comment, error) {
	var post Post
	queryPost := `SELECT title, content, authorid FROM post WHERE id = ?`
	row := db.QueryRow(queryPost, id)
	err := row.Scan(&post.Title, &post.Content, &post.AuthorID)
	if err != nil {
		return Post{}, []comments.Comment{}, err
	}

	var coms []comments.Comment
	var comment comments.Comment
	queryComment := `SELECT id, authorid, content FROM comments WHERE postid = ?`
	rows, err := db.Query(queryComment, id)

	for rows.Next() {
		if err := rows.Scan(&comment.ID, &comment.AuthorID, &comment.Content); err != nil {
			if err == sql.ErrNoRows {
				return post, []comments.Comment{}, nil
			}
			return Post{}, []comments.Comment{}, err
		}

		coms = append(coms, comment)
	}

	return post, coms, nil
}
