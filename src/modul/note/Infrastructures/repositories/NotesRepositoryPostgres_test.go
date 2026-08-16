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

	var expectedNote entities.Note
	err = repo.DB.Where("id = ?", id).First(&expectedNote).Error

	assert.NoError(t, err)
	assert.Equal(t, note.ID, expectedNote.ID)
	assert.Equal(t, note.Title, expectedNote.Title)
	assert.Equal(t, note.Owner, expectedNote.Owner)

	// repo.db.Exec("DELETE FROM notes")
}

func Test_DeleteNote(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		err := repo.DeleteNoteById("id-123")

		assert.EqualError(t, err, "id tidak ditemukan")
	})

	t.Run("delete success", func(t *testing.T) {
		now := time.Now()

		note := entities.Note{
			ID:       "id-123",
			Title:    "delete note",
			Content:  "delete content note",
			CreateAt: now,
			UpdateAt: now,
			Owner:    "user-123",
		}

		err := repo.DB.Create(&note).Error
		assert.NoError(t, err)

		err = repo.DeleteNoteById("id-123")

		assert.NoError(t, err)

		var deletedNote entities.Note
		err = repo.DB.
			Where("id = ?", "id-123").
			First(&deletedNote).Error

		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
