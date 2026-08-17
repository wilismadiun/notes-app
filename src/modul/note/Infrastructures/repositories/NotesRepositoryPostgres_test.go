package repositories

import (
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

func TestMain(m *testing.M) {
	database.ConnectPostgresql(".test.env")

	repo = &NoteRepository{
		DB: database.DB,
	}

	user := entitiesUser.User{
		ID:       "user-123",
		Username: "Jaya123",
		Password: "pass-123",
	}

	err := repo.DB.Create(&user).Error
	if err != nil {
		panic(err)
	}

	code := m.Run()

	repo.DB.Exec("DELETE FROM users")
	os.Exit(code)
}

func Test_CreateNote(t *testing.T) {
	now := time.Now()
	id := uuid.New().String()

	note := entities.Note{
		ID:       id,
		Title:    "test database",
		Content:  "test content database",
		CreateAt: now,
		UpdateAt: now,
		Owner:    "user-123",
	}

	err := repo.CreateNote(note)
	assert.NoError(t, err)

	expectedNote := entities.Note{ID: id}
	err = repo.DB.First(&expectedNote).Error

	assert.NoError(t, err)
	assert.Equal(t, note.ID, expectedNote.ID)
	assert.Equal(t, note.Title, expectedNote.Title)
	assert.Equal(t, note.Owner, expectedNote.Owner, gorm.ErrRecordNotFound)

	repo.DB.Exec("DELETE FROM notes")
}

func Test_DeleteNote(t *testing.T) {
	t.Run("should be error when id not found", func(t *testing.T) {
		err := repo.DeleteNoteById("note-123")

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("delete note success", func(t *testing.T) {
		now := time.Now()
		id := uuid.New().String()

		note := entities.Note{
			ID:       id,
			Title:    "test database",
			Content:  "test content database",
			CreateAt: now,
			UpdateAt: now,
			Owner:    "user-123",
		}

		err := repo.CreateNote(note)
		assert.NoError(t, err)

		exisistNote := entities.Note{ID: id}

		err = repo.DB.First(&exisistNote).Error
		assert.NoError(t, err)
		assert.Equal(t, note.Title, exisistNote.Title)
		assert.Equal(t, note.Content, exisistNote.Content)

		err = repo.DeleteNoteById(id)

		assert.NoError(t, err)

		newExisistNote := entities.Note{ID: id}

		err = repo.DB.First(&newExisistNote).Error
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
