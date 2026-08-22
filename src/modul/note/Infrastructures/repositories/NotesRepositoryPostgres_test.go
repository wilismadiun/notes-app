package repositories

import (
	"fmt"
	"notes-app/src/commons/database"
	"notes-app/src/modul/note/Domains/entities"
	entitiesUser "notes-app/src/modul/user/Domains/entities"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var repo *NoteRepository
var userId = fmt.Sprintf("user-%s", uuid.New().String())
var noteId = fmt.Sprintf("note-%s", uuid.New().String())
var now = time.Now()

func TestMain(m *testing.M) {
	database.ConnectPostgresql(".test.env")

	repo = &NoteRepository{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func DatabaseHelper() (entitiesUser.User, entities.Note) {
	userId := "user-123"
	noteId := "note-123"
	now := time.Now().Truncate(time.Microsecond)

	user := entitiesUser.User{
		ID:       userId,
		Username: "Jaya123",
		Password: "12345678",
	}

	note := entities.Note{
		ID:       noteId,
		Title:    "title",
		Content:  "content",
		CreateAt: now,
		UpdateAt: now,
		Owner:    userId,
	}

	database.DB.Create(&user)
	database.DB.Create(&note)

	return user, note
}

func Test_CreateNote(t *testing.T) {
	user, _ := DatabaseHelper()

	noteId := "note ID"
	now := time.Now()

	note := entities.Note{
		ID:       noteId,
		Title:    "test database",
		Content:  "test content database",
		CreateAt: now,
		UpdateAt: now,
		Owner:    user.ID,
	}

	err := repo.CreateNote(note)
	assert.NoError(t, err)

	exisistNote := entities.Note{
		ID: noteId,
	}
	err = repo.DB.First(&exisistNote).Error

	assert.NoError(t, err)
	assert.Equal(t, note.ID, exisistNote.ID)
	assert.Equal(t, note.Title, exisistNote.Title)
	assert.Equal(t, note.Owner, exisistNote.Owner, gorm.ErrRecordNotFound)

	repo.DB.Exec("DELETE FROM notes")
	repo.DB.Exec("DELETE FROM users")
}

func Test_GetAllNote(t *testing.T) {
	user, _ := DatabaseHelper()

	now := time.Now()

	notes := []entities.Note{
		{
			ID:       "id-1",
			Title:    "title 1",
			Content:  "content 1",
			CreateAt: now,
			UpdateAt: now,
			Owner:    user.ID,
		},
		{
			ID:       "id-2",
			Title:    "title 2",
			Content:  "content 2",
			CreateAt: now,
			UpdateAt: now,
			Owner:    user.ID,
		},
	}

	database.DB.Create(&notes)

	t.Run("when owner not found", func(t *testing.T) {
		exisistNote := repo.GetAllNote("user-13")

		assert.Empty(t, exisistNote)
		assert.Equal(t, 0, len(exisistNote))
	})

	t.Run("success", func(t *testing.T) {
		exisistNote := repo.GetAllNote("user-123")

		assert.NotEmpty(t, exisistNote)
		assert.Equal(t, 3, len(exisistNote))
	})
	repo.DB.Exec("DELETE FROM notes")
	repo.DB.Exec("DELETE FROM users")
}

func Test_FindNoteById(t *testing.T) {
	user, note := DatabaseHelper()

	t.Run("should be error when id not found", func(t *testing.T) {
		note := entities.Note{
			ID:    "Note ID",
			Owner: "User ID",
		}

		err := repo.FindNoteById(&note)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("success", func(t *testing.T) {
		exisistNote := entities.Note{
			ID:    note.ID,
			Owner: user.ID,
		}

		err := repo.FindNoteById(&exisistNote)
		assert.NoError(t, err)
		assert.Equal(t, note.ID, exisistNote.ID)
		assert.Equal(t, note.Title, exisistNote.Title)
		assert.Equal(t, note.Content, exisistNote.Content)

	})
	repo.DB.Exec("DELETE FROM notes")
	repo.DB.Exec("DELETE FROM users")
}

func Test_EditNoteById(t *testing.T) {
	user, notePrepare := DatabaseHelper()

	now := time.Now().Truncate(time.Microsecond)

	note := entities.Note{
		ID:       "Note ID",
		Title:    "title",
		Content:  "content",
		CreateAt: now,
		UpdateAt: now,
		Owner:    user.ID,
	}

	editNote := map[string]any{
		"title":    "new title",
		"content":  "new content",
		"updateAt": time.Now().Truncate(time.Microsecond),
	}

	t.Run("should be error when id not found", func(t *testing.T) {
		err := repo.EditNoteById(note, editNote)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("edit note success", func(t *testing.T) {
		err := repo.EditNoteById(notePrepare, editNote)

		assert.NoError(t, err)

		exisistNote := entities.Note{
			ID:    notePrepare.ID,
			Owner: user.ID,
		}

		err = repo.FindNoteById(&exisistNote)
		assert.NoError(t, err)

		assert.Equal(t, notePrepare.ID, exisistNote.ID)
		assert.Equal(t, editNote["title"], exisistNote.Title)
		assert.Equal(t, editNote["content"], exisistNote.Content)
		assert.Equal(t, editNote["updateAt"], exisistNote.UpdateAt)

	})

	database.DB.Exec("DELETE FROM notes")
	database.DB.Exec("DELETE FROM users")
}

func Test_DeleteNote(t *testing.T) {
	user, note := DatabaseHelper()

	t.Run("should be error when id not found", func(t *testing.T) {
		exisistNote := entities.Note{
			ID:    "note",
			Owner: user.ID,
		}

		err := repo.DeleteNoteById(exisistNote)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("delete note success", func(t *testing.T) {
		exisistNote := entities.Note{
			ID:    note.ID,
			Owner: user.ID,
		}

		err := repo.DB.First(&exisistNote).Error
		assert.NoError(t, err)
		assert.Equal(t, note.Title, exisistNote.Title)
		assert.Equal(t, note.Content, exisistNote.Content)

		err = repo.DeleteNoteById(exisistNote)

		assert.NoError(t, err)

		newExisistNote := entities.Note{
			ID: noteId,
		}

		err = repo.DB.First(&newExisistNote).Error
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		repo.DB.Exec("DELETE FROM users")
	})
}
