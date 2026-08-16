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
		note      Note
		wantError string
	}{
		{
			name: "id kosong",
			note: Note{
				ID:       "",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note id is required",
		},
		{
			name: "id hanya spasi",
			note: Note{
				ID:       "   ",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note id is required",
		},
		{
			name: "title kosong",
			note: Note{
				ID:       "note-123",
				Title:    "",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note title is required",
		},
		{
			name: "title hanya spasi",
			note: Note{
				ID:       "note-123",
				Title:    "   ",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note title is required",
		},
		{
			name: "title lebih dari 100 karakter",
			note: Note{
				ID:       "note-123",
				Title:    strings.Repeat("a", 101),
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note title must not exceed 100 characters",
		},
		{
			name: "content kosong",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note content is required",
		},
		{
			name: "content hanya spasi",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "   ",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note content is required",
		},
		{
			name: "owner kosong",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note owner is required",
		},
		{
			name: "owner hanya spasi",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "   ",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "note owner is required",
		},
		{
			name: "created at invalid",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: time.Time{},
				UpdateAt: now,
			},
			wantError: "note created at is invalid",
		},
		{
			name: "updated at invalid",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: time.Time{},
			},
			wantError: "note updated at is invalid",
		},
		{
			name: "updated at sebelum created at",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now.Add(-time.Hour),
			},
			wantError: "updated at cannot be before created at",
		},
		{
			name: "valid note",
			note: Note{
				ID:       "note-123",
				Title:    "Belajar Go",
				Content:  "Belajar Go untuk backend",
				Owner:    "user-123",
				CreateAt: now,
				UpdateAt: now,
			},
			wantError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyNote(tt.note)

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
