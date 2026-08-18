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

func Test_EditNoteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotesRepository(ctrl)

	en := EditNoteById{
		Repo: mockRepo,
	}

	noteId := "id-123"
	now := time.Now().Truncate(time.Microsecond)
	title := "title note"
	content := "content note"

	editNotePayload := entities.EditNoteRequest{
		Title:   nil,
		Content: nil,
	}

	editNoteMap := map[string]any{
		"title":    title,
		"content":  content,
		"updateAt": now,
	}

	t.Run("should be error when id not found", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(noteId).Return(entities.Note{}, errors.New("id tidak ditemukan"))

		id, err := en.Execute(noteId, editNotePayload)

		assert.Error(t, err)
		assert.EqualError(t, err, "id tidak ditemukan")
		assert.Equal(t, "", id)
	})

	note := entities.Note{
		ID:       noteId,
		Title:    title,
		Content:  content,
		CreateAt: now,
		UpdateAt: now,
		Owner:    "user-123",
	}

	t.Run("should be error when note to edit not found", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(noteId).Return(note, nil)

		id, err := en.Execute(noteId, editNotePayload)

		assert.Error(t, err)
		assert.EqualError(t, err, "Tidak ada data yang dikirim untuk diubah")
		assert.Equal(t, "", id)
	})

	t.Run("edit note success", func(t *testing.T) {
		mockRepo.EXPECT().FindNoteById(noteId).Return(note, nil)

		editNotePayload.Title = &title
		editNotePayload.Content = &content

		mockRepo.EXPECT().EditNoteById(note, editNoteMap).Return(nil)

		id, err := en.Execute(noteId, editNotePayload)

		assert.NoError(t, err)
		assert.Equal(t, noteId, id)
	})
}
