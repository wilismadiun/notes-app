package entities

import (
	"strings"
	"testing"
	"time"
)

func TestVerifyNote(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		id        string
		title     string
		content   string
		owner     string
		createAt  time.Time
		updateAt  time.Time
		wantError string
	}{
		{
			name:      "id kosong",
			id:        "",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note id is required",
		},
		{
			name:      "id hanya spasi",
			id:        "   ",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note id is required",
		},
		{
			name:      "title kosong",
			id:        "note-123",
			title:     "",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note title is required",
		},
		{
			name:      "title hanya spasi",
			id:        "note-123",
			title:     "   ",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note title is required",
		},
		{
			name:      "title lebih dari 100 karakter",
			id:        "note-123",
			title:     strings.Repeat("a", 101),
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note title must not exceed 100 characters",
		},
		{
			name:      "content kosong",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note content is required",
		},
		{
			name:      "content hanya spasi",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "   ",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "note content is required",
		},
		{
			name:      "owner kosong",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "",
			createAt:  now,
			updateAt:  now,
			wantError: "note owner is required",
		},
		{
			name:      "owner hanya spasi",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "   ",
			createAt:  now,
			updateAt:  now,
			wantError: "note owner is required",
		},
		{
			name:      "created at invalid",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  time.Time{},
			updateAt:  now,
			wantError: "note created at is invalid",
		},
		{
			name:      "updated at invalid",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  time.Time{},
			wantError: "note updated at is invalid",
		},
		{
			name:      "updated at sebelum created at",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now.Add(-time.Hour),
			wantError: "updated at cannot be before created at",
		},
		{
			name:      "valid note",
			id:        "note-123",
			title:     "Belajar Go",
			content:   "Belajar Go untuk backend",
			owner:     "user-123",
			createAt:  now,
			updateAt:  now,
			wantError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyNote(
				tt.id,
				tt.title,
				tt.content,
				tt.owner,
				tt.createAt,
				tt.updateAt,
			)

			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantError)
			}

			if err.Error() != tt.wantError {
				t.Fatalf(
					"expected error %q, got %q",
					tt.wantError,
					err.Error(),
				)
			}
		})
	}
}
