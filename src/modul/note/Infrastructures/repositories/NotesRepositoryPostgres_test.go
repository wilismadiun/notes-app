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

func Test_CreateNote(t *testing.T) {
	user := entitiesUser.User{
		ID:       userId,
		Username: "Jaya 123",
		Password: "12345678",
	}

	err := repo.DB.Create(&user).Error
	assert.NoError(t, err)

	note := entities.Note{
		ID:       noteId,
		Title:    "test database",
		Content:  "test content database",
		CreateAt: now,
		UpdateAt: now,
		Owner:    userId,
	}

	err = repo.CreateNote(note)
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

func Test_DeleteNote(t *testing.T) {
	t.Run("should be error when id not found", func(t *testing.T) {
		err := repo.DeleteNoteById("note-123")

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("delete note success", func(t *testing.T) {
		user := entitiesUser.User{
			ID:       userId,
			Username: "Jaya 123",
			Password: "12345678",
		}

		err := repo.DB.Create(&user).Error
		assert.NoError(t, err)

		note := entities.Note{
			ID:       noteId,
			Title:    "test database",
			Content:  "test content database",
			CreateAt: now,
			UpdateAt: now,
			Owner:    userId,
		}

		err = repo.CreateNote(note)
		assert.NoError(t, err)

		exisistNote := entities.Note{
			ID: noteId,
		}

		err = repo.DB.First(&exisistNote).Error
		assert.NoError(t, err)
		assert.Equal(t, note.Title, exisistNote.Title)
		assert.Equal(t, note.Content, exisistNote.Content)

		err = repo.DeleteNoteById(noteId)

		assert.NoError(t, err)

		newExisistNote := entities.Note{
			ID: noteId,
		}

		err = repo.DB.First(&newExisistNote).Error
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

		repo.DB.Exec("DELETE FROM notes")
		repo.DB.Exec("DELETE FROM users")
	})
}
