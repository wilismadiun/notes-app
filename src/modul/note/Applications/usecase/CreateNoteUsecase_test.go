package usecase

import (
	"errors"
	"notes-app/src/modul/note/Domains/entities"
	"notes-app/src/modul/note/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateNoteUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGenerator := mocks.NewMockIdGenerator(ctrl)
	mockRepository := mocks.NewMockNotesRepository(ctrl)

	cn := CreateNote{
		Generator: mockGenerator,
		Repo:      mockRepository,
	}

	t.Run("should return error when owner is empty", func(t *testing.T) {
		id := "id-123"
		mockGenerator.EXPECT().Generator().Return(id)

		note := entities.Note{
			ID:      id,
			Title:   "test use case",
			Content: "test content use case",
		}

		noteResponse, err := cn.Execute(note)

		assert.Error(t, err)
		assert.EqualError(t, err, "note owner is required")
		assert.Empty(t, noteResponse)
	})

	t.Run("repository failed to create note", func(t *testing.T) {
		mockGenerator.EXPECT().Generator().Return("id-123")

		note := entities.Note{
			Title:   "test use case",
			Content: "test content use case",
			Owner:   "user-123",
		}

		mockRepository.EXPECT().CreateNote(gomock.Any()).Return(entities.CreateNoteResponse{}, errors.New("Gagal menambahkan note"))

		noteResponse, err := cn.Execute(note)

		assert.Error(t, err)
		assert.EqualError(t, err, "Gagal menambahkan note")
		assert.Empty(t, noteResponse)
	})

	t.Run("success", func(t *testing.T) {
		mockGenerator.EXPECT().Generator().Return("id-123")

		note := entities.Note{
			Title:   "test use case",
			Content: "test content use case",
			Owner:   "user-123",
		}

		expectedResponse := entities.CreateNoteResponse{
			ID:    "id-123",
			Title: "test use case",
		}

		mockRepository.EXPECT().CreateNote(gomock.Any()).Return(expectedResponse, nil)

		noteResponse, err := cn.Execute(note)

		assert.NoError(t, err)
		assert.Equal(t, nil, err)
		assert.Equal(t, expectedResponse, noteResponse)
	})
}
