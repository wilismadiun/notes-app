package usecase

import (
	"errors"
	"notes-app/src/modul/note/Domains/entities"
	"notes-app/src/modul/note/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_FindNoteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotesRepository(ctrl)

	fn := GetNoteById{
		Repo: mockRepo,
	}

	id := "note-123"

	t.Run("should be error when id not found", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(id).Return(entities.Note{}, errors.New("id tidak ditemukan"))

		note, err := fn.Execute(id)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
		assert.Empty(t, note)
	})

	t.Run("get note success", func(t *testing.T) {
		note := entities.Note{
			ID:       id,
			Title:    "Title",
			Content:  "Content",
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Owner:    "user-123",
		}
		mockRepo.EXPECT().FindNoteById(id).Return(note, nil)

		exisistNote, err := fn.Execute(id)

		assert.NoError(t, err)
		assert.Equal(t, note.ID, exisistNote.ID)
		assert.Equal(t, note.Content, exisistNote.Content)
		assert.Equal(t, note.Title, exisistNote.Title)
		assert.Equal(t, note.Owner, exisistNote.Owner)
	})
}
