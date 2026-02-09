package posts

import (
	"fmt"
)

func IsValidFormat(post Post) error {
	if post.Title == "" {
		return fmt.Errorf("empty title")
	} else if post.Content == "" {
		return fmt.Errorf("no content provided")
	}
	return nil
}
