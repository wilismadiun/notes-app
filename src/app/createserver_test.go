package app

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"notes-app/src/commons/database"
	entitiesNote "notes-app/src/modul/note/Domains/entities"
	"notes-app/src/modul/user/Domains/entities"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	database.ConnectPostgresql(".test.env")

	gin.SetMode(gin.TestMode)

	code := m.Run()

	os.Exit(code)
}

func TestUserModule(t *testing.T) {

	router := gin.New()

	Router(router, database.DB)

	t.Run("should response 400 when password less than 8 character", func(t *testing.T) {
		body := []byte(`{
			"username":"jaya",
			"email":"jaya@gmail.com",
			"password":"123456"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "PASSWORD_TOO_SHORT"
		}`, w.Body.String())
	})

	t.Run("should response 400 when username avalilable", func(t *testing.T) {
		user := entities.User{
			ID:       uuid.New().String(),
			Username: "jaya",
			Password: "12345678",
		}

		database.DB.Table("users").Create(&user)

		body := []byte(`{
			"username": "jaya",
			"email": "jaya@gmail.com",
			"password": "12345678"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message":"username sudah digunakan"
		}`, w.Body.String())

		database.DB.Exec("DELETE FROM users")
	})

	t.Run("Create user success", func(t *testing.T) {
		body := []byte(`{
		"username": "jaya",
		"password": "12345678",
		"email": "jaya@gmail.com"
	}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/register",
			bytes.NewBuffer(body),
		)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Ambil data dari database
		var user entities.User
		err := database.DB.
			Where("username = ?", "jaya").
			First(&user).Error

		assert.NoError(t, err)

		// Assert HTTP
		assert.Equal(t, http.StatusCreated, w.Code)

		// Password di database harus sudah di-hash
		assert.NotEqual(t, "12345678", user.Password)

		// Response harus sama dengan data di database
		userJson := fmt.Sprintf(`{
			"message": "Berhasil menambahkan user",
			"data": {
				"ID": "%s",
				"Username": "%s"
			}
		}`, user.ID, user.Username)

		assert.JSONEq(t, userJson, w.Body.String())

		database.DB.Exec("DELETE FROM users")
	})
}

func TestNoteModul(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	user := entities.User{
		ID:       "user-123",
		Username: "Jaya123",
		Password: "pass123",
	}

	err := database.DB.Create(&user).Error
	assert.NoError(t, err)

	t.Run("should response 400 when owner empty", func(t *testing.T) {
		body := []byte(`{
			"title": "Note1",
			"content": "content note1"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/note",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "gagal menambahkan note",
			"error": "note owner is required"
		}`, w.Body.String())
	})

	t.Run("Add note success", func(t *testing.T) {
		body := []byte(`{
			"title": "Note1",
			"content": "content note1",
			"owner": "user-123"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/note",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var note entitiesNote.Note

		err := database.DB.
			Where("title = ? AND owner = ?", "Note1", "user-123").
			First(&note).
			Error

		assert.NoError(t, err)
		expectedResponse := fmt.Sprintf(`{
			"message": "note berhasil ditambahkan",
			"data": {
				"ID": "%s",
				"Title": "%s"
			}
		}`, note.ID, note.Title)

		assert.JSONEq(t, expectedResponse, w.Body.String())

		// Cleanup
		database.DB.Exec("DELETE FROM notes")
	})
	database.DB.Exec("DELETE FROM users")
}

func Test_DeleteNoteModul(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	user := entities.User{
		ID:       "user-123",
		Username: "Jaya123",
		Password: "pass123",
	}

	err := database.DB.Create(&user).Error
	assert.NoError(t, err)

	t.Run("should response 400 when id not found", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodDelete,
			"/note/:id",
			nil,
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "gagal menghapus note",
			"error": "id tidak ditemukan"
		}`, w.Body.String())
	})

	t.Run("delete success", func(t *testing.T) {
		now := time.Now()

		note := entitiesNote.Note{
			ID:       "id-123",
			Title:    "delete note",
			Content:  "delete content note",
			CreateAt: now,
			UpdateAt: now,
			Owner:    "user-123",
		}

		err := database.DB.Create(&note).Error
		assert.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodDelete,
			"/note/id-123",
			nil,
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{
			"message": "berhasil menghapus note",
			"data": "id-123"
		}`, w.Body.String())

	})
	database.DB.Exec("DELETE FROM users")
}
