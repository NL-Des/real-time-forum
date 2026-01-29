package posts

import "time"

type Post struct {
	ID          int
	Title       string
	Content     string
	AuthorID    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryIDs []int
}
