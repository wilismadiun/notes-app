package usecase

import (
	"errors"
	"notes-app/src/modul/note/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_DeleteNoteUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotesRepository(ctrl)

	dn := DeleteNote{
		Repo: mockRepo,
	}

	userId := "user-123"
	noteId := "note-123"

	t.Run("should return error when id not found", func(t *testing.T) {
		mockRepo.EXPECT().DeleteNoteById(gomock.Any()).Return(errors.New("Id tidak ditemukan"))

		_, err := dn.Execute(noteId, userId)

		assert.Error(t, err)
		assert.EqualError(t, err, "Id tidak ditemukan")
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().DeleteNoteById(gomock.Any()).Return(nil)

		id, err := dn.Execute(noteId, userId)

		assert.NoError(t, err)
		assert.Equal(t, nil, err)
		assert.Equal(t, noteId, id)
	})
}
