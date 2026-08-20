package usecase

import (
	"notes-app/src/modul/note/Domains/entities"
	"notes-app/src/modul/note/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_GetAllNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNotesRepository(ctrl)

	gn := GetAllNotes{
		Repo: mockRepo,
	}

	t.Run("empty note", func(t *testing.T) {
		mockRepo.EXPECT().GetAllNote("user-123").Return([]entities.Note{})

		notes, err := gn.Execute("user-123")

		assert.Error(t, err)
		assert.EqualError(t, err, "Tidak ada data yg bisa ditampilkan")
		assert.Empty(t, notes)
	})

	t.Run("note with content", func(t *testing.T) {
		now := time.Now()

		notes := []entities.Note{
			{
				ID:       "id-1",
				Title:    "title 1",
				Content:  "content 1",
				CreateAt: now,
				UpdateAt: now,
				Owner:    "user-123",
			},
			{
				ID:       "id-2",
				Title:    "title 2",
				Content:  "content 2",
				CreateAt: now,
				UpdateAt: now,
				Owner:    "user-123",
			},
		}

		mockRepo.EXPECT().GetAllNote("user-123").Return(notes)

		exisistNotes, err := gn.Execute("user-123")

		assert.NoError(t, err)
		assert.Equal(t, notes, exisistNotes)
	})
}
