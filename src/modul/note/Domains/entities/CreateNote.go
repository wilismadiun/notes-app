package entities

import (
	"errors"
	"strings"
)

func VerifyNote(note Note) error {
	if strings.TrimSpace(note.ID) == "" {
		return errors.New("note id is required")
	}

	if strings.TrimSpace(note.Title) == "" {
		return errors.New("note title is required")
	}

	if len([]rune(note.Title)) > 100 {
		return errors.New("note title must not exceed 100 characters")
	}

	if strings.TrimSpace(note.Content) == "" {
		return errors.New("note content is required")
	}

	if strings.TrimSpace(note.Owner) == "" {
		return errors.New("note owner is required")
	}

	// CreatedAt tidak boleh zero value
	if note.CreateAt.IsZero() {
		return errors.New("note created at is invalid")
	}

	// UpdatedAt tidak boleh zero value
	if note.UpdateAt.IsZero() {
		return errors.New("note updated at is invalid")
	}

	// UpdatedAt tidak boleh lebih awal dari CreatedAt
	if note.UpdateAt.Before(note.CreateAt) {
		return errors.New("updated at cannot be before created at")
	}

	return nil

}
