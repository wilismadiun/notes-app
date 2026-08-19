package usecase

import (
	"errors"
	"notes-app/src/modul/note/Domains/entities"
	"notes-app/src/modul/note/mocks"
	"testing"

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

	noteId := "note-123"
	userId := "user-123"

	note := entities.Note{
		ID:    noteId,
		Owner: userId,
	}

	t.Run("should be error when id not found", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(&note).Return(errors.New("id tidak ditemukan"))

		result, err := fn.Execute(noteId, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
		assert.Empty(t, result)
	})

	t.Run("get note success", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(&note).Return(nil)

		exisistNote, err := fn.Execute(noteId, userId)

		assert.NoError(t, err)
		assert.Equal(t, exisistNote.ID, exisistNote.ID)
		assert.Equal(t, exisistNote.Content, exisistNote.Content)
		assert.Equal(t, exisistNote.Title, exisistNote.Title)
		assert.Equal(t, exisistNote.Owner, exisistNote.Owner)
	})
}
