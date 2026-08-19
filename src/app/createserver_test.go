package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func Test_CreateUserModule(t *testing.T) {

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

	t.Run("should response 400 when username available", func(t *testing.T) {
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
		"username": "Jaya123",
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
			Where("username = ?", "Jaya123").
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
	})
}

func Test_LoginUserModule(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	t.Run("should be an error when the username doesn't exist", func(t *testing.T) {
		body := []byte(`{
			"username": "",
			"password": "123456789"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "login gagal",
			"error": "username atau Password tidak ada"
		}`, w.Body.String())
	})

	t.Run("There should be an error if the username isn't registered yet", func(t *testing.T) {
		body := []byte(`{
			"username": "Jaya12",
			"password": "12345678"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "login gagal",
			"error": "username atau password salah"
		}`, w.Body.String())
	})

	t.Run("There should be an error when the password does not match the password hash.", func(t *testing.T) {
		body := []byte(`{
			"username": "Jaya123",
			"password": "123456789"
		}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
			"message": "login gagal",
			"error": "username atau password salah"
		}`, w.Body.String())
	})

	t.Run("login success", func(t *testing.T) {
		body := []byte(`{
		"username": "Jaya123",
		"password": "12345678"
		}`)

		var user entities.User
		database.DB.First(&user)
		log.Println("================================================================")
		log.Println(user)
		log.Println("================================================================")

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Message string `json:"message"`
			Data    string `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "login berhasil", response.Message)
		assert.NotEmpty(t, response.Data)
	})

	database.DB.Exec("DELETE FROM users")
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

func Test_GetNoteById(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	user := entities.User{
		ID:       "user-1234",
		Username: "Jaya1234",
		Password: "pass123",
	}

	err := database.DB.Create(&user).Error
	assert.NoError(t, err)

	t.Run("should response 400 when id not found", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			"/note/id-123",
			nil,
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		response := fmt.Sprintln(`{
			"message": "gagal menampilkan note",
			"error": "id tidak ditemukan"
		}`)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, response, w.Body.String())
	})

	t.Run("get note by id success", func(t *testing.T) {
		now := time.Now().Truncate(time.Microsecond)

		note := entitiesNote.Note{
			ID:       "id-1234",
			Title:    "delete note",
			Content:  "delete content note",
			CreateAt: now,
			UpdateAt: now,
			Owner:    "user-1234",
		}

		err := database.DB.Create(&note).Error
		assert.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			"/note/id-1234",
			nil,
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		noteJSON, _ := json.Marshal(note)

		response := fmt.Sprintf(`{
			"message": "note berhasil ditampilkan",
			"data": %s
		}`, noteJSON)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, response, w.Body.String())

		database.DB.Exec("DELETE FROM notes")
	})
	database.DB.Exec("DELETE FROM users")
}

func Test_EditNoteModul(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	user := entities.User{
		ID:       "user-1234",
		Username: "Jaya1234",
		Password: "pass123",
	}

	err := database.DB.Create(&user).Error
	assert.NoError(t, err)

	t.Run("should be error when id not found", func(t *testing.T) {
		body := []byte(`{
			"title": "new title",
			"content": "new content"
		}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/note/note-123",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		responseExpected := fmt.Sprintln(`{
			"message": "gagal memperbarui note",
			"error": "id tidak ditemukan"
		}`)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())
	})

	now := time.Now().Truncate(time.Microsecond)

	note := entitiesNote.Note{
		ID:       "note-123",
		Title:    "title",
		Content:  "content",
		CreateAt: now,
		UpdateAt: now,
		Owner:    user.ID,
	}

	err = database.DB.Create(&note).Error
	assert.NoError(t, err)

	t.Run("should be error when note to edit not found", func(t *testing.T) {
		body := []byte(`{
			"title": null,
			"content": null
		}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/note/note-123",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		responseExpected := fmt.Sprintln(`{
			"message": "gagal memperbarui note",
			"error": "Tidak ada data yang dikirim untuk diubah"
		}`)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())
	})

	t.Run("should success update note when only updating title", func(t *testing.T) {
		log.Println("===================================only title sebelum edit======================================")
		note1 := entitiesNote.Note{
			ID: "note-123",
		}
		database.DB.First(&note1)
		log.Println(note1)
		log.Println("===================================only title sebelum edit======================================")
		body := []byte(`{
			"title": "new title"
			}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/note/note-123",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		log.Println("===================================only title setelah edit======================================")
		note2 := entitiesNote.Note{
			ID: "note-123",
		}
		database.DB.First(&note2)
		log.Println(note2)
		log.Println("===================================only title sebelum edit======================================")

		responseExpected := fmt.Sprintln(`{
			"message": "Note berhasil diperbarui",
			"data": "note-123"
		}`)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())

		exisistNote := entitiesNote.Note{
			ID: "note-123",
		}
		err := database.DB.First(&exisistNote).Error
		assert.NoError(t, err)

		assert.Equal(t, "new title", exisistNote.Title)
	})

	t.Run("should success update note when updating title and content", func(t *testing.T) {
		log.Println("===================================title & content sebelum edit======================================")
		note3 := entitiesNote.Note{
			ID: "note-123",
		}
		database.DB.First(&note3)
		log.Println(note3)
		log.Println("===================================title & content sebelum edit======================================")

		body := []byte(`{
			"title": "update new title",
			"content": "new content"
		}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			"/note/note-123",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		log.Println("===================================title & content sebelum edit======================================")
		note4 := entitiesNote.Note{
			ID: "note-123",
		}
		database.DB.First(&note4)
		log.Println(note4)
		log.Println("===================================title & content sebelum edit======================================")

		responseExpected := fmt.Sprintln(`{
			"message": "Note berhasil diperbarui",
			"data": "note-123"
		}`)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())

		exisistNote := entitiesNote.Note{
			ID: "note-123",
		}

		err := database.DB.First(&exisistNote).Error
		assert.NoError(t, err)

		assert.Equal(t, "update new title", exisistNote.Title)
		assert.Equal(t, "new content", exisistNote.Content)
	})

	database.DB.Exec("DELETE FROM notes")
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
