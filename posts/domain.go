package posts

import "time"

type Post struct {
	ID          int
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int    `json:"authorid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryIDs []int `json:"category_ids"`
}
