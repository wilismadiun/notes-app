package repositories

import (
	"log"
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
var id = uuid.New().String()

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
	log.Println("=============================ini adalah errornya yang pertama========================")
	log.Println(err)
	log.Println("=============================ini adalah errornya========================")

	expectedNote := entities.Note{
		ID: id,
	}
	err = repo.DB.First(&expectedNote).Error

	
	log.Println("=============================ini adalah errornya create note========================")
	log.Println(err)
	log.Println("=============================ini adalah errornya========================")

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
			ID:       id,
			Title:    "test database",
			Content:  "test content database",
			CreateAt: now,
			UpdateAt: now,
			Owner:    "user-123",
		}
	
		err := repo.CreateNote(note)
	
		assert.NoError(t, err)
		log.Println("=============================ini adalah errornya yang pertama delete note========================")
		log.Println(err)
		log.Println("=============================ini adalah errornya========================")

		
		var existNote entities.Note
		err = repo.DB.First(&existNote, "id = ?", id).Error
		if err != nil {
			log.Println("=============================================")
			log.Println("Data note tidak ditemukan")
			log.Println("=============================================")
		} else {
			log.Println("=============================================")
			log.Println("Data note ditemukan")
			log.Println(existNote.ID)
			log.Println("=============================================")
		}

		err = repo.DeleteNoteById(existNote.ID)

		assert.NoError(t, err)

		var deletedNote entities.Note
		err = repo.DB.First(&deletedNote, "id = ?", id).Error

		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}
