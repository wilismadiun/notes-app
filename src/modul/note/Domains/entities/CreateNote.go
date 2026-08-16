package entities

import (
	"errors"
	"strings"
	"time"
)

func VerifyNote(id, title, content, owner string, createAt, updateAt time.Time) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("note id is required")
	}

	if strings.TrimSpace(title) == "" {
		return errors.New("note title is required")
	}

	if len([]rune(title)) > 100 {
		return errors.New("note title must not exceed 100 characters")
	}

	if strings.TrimSpace(content) == "" {
		return errors.New("note content is required")
	}

	if strings.TrimSpace(owner) == "" {
		return errors.New("note owner is required")
	}

	// CreatedAt tidak boleh zero value
	if createAt.IsZero() {
		return errors.New("note created at is invalid")
	}

	// UpdatedAt tidak boleh zero value
	if updateAt.IsZero() {
		return errors.New("note updated at is invalid")
	}

	// UpdatedAt tidak boleh lebih awal dari CreatedAt
	if updateAt.Before(createAt) {
		return errors.New("updated at cannot be before created at")
	}

	return nil

}
