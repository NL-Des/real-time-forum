package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func IsValid(post Post) error {
	if post.Title == "" {
		return fmt.Errorf("empty title")
	} else if post.Content == "" {
		return fmt.Errorf("no content provided")
	}

	return nil
}

func RespondError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}
