package posts

import (
	"real-time-forum/comments"
	"time"
)

type Post struct {
	ID          int
	Title       string `json:"title"`
	Content     string `json:"content"`
	AuthorID    int    `json:"authorid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryIDs []int `json:"category_ids"`
}

type PostWithCommentsResponse struct {
	Post     Post               `json:"post"`
	Comments []comments.Comment `json:"comments"`
}
