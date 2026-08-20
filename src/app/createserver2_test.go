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
	"notes-app/src/modul/user/Infrastructures/security"
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

func databaseHelper() (entities.User, string, entitiesNote.Note, string, string) {
	now := time.Now()
	userId := "idUser-123"
	username := "John123"
	password := "8764321"
	noteId := "idNote-123"

	hasher := security.HashPasswordBcrypt{}
	hashPassword, _ := hasher.HashingPassword(password)

	user := entities.User{
		ID:       userId,
		Username: username,
		Password: hashPassword,
	}
	database.DB.Create(&user)

	note := entitiesNote.Note{
		ID:       noteId,
		Title:    "title note preparatin",
		Content:  "content note preparatin",
		CreateAt: now,
		UpdateAt: now,
		Owner:    userId,
	}

	database.DB.Create(&note)

	authToken, _ := tokenService.GenerateToken(userId)

	path := fmt.Sprintf("/api/note/%s", noteId)

	return user, password, note, authToken, path
}

func Test_CreateUser(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	t.Run("should response 400 when username available", func(t *testing.T) {
		user := entities.User{
			ID:       uuid.New().String(),
			Username: "jaya",
			Password: "12345678",
		}

		database.DB.Table("users").Create(&user)

		body := []byte(`{
				"username": "jaya",
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
	})

	t.Run("Create user success", func(t *testing.T) {
		body := []byte(`{
			"username": "Jaya123",
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

	database.DB.Exec("DELETE FROM users")
}

func Test_UserLogin(t *testing.T) {
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

	user, password, _, _, _ := databaseHelper()

	t.Run("There should be an error when the password does not match the password hash.", func(t *testing.T) {
		requestBody := fmt.Sprintf(`{
				"username": "%s",
				"password": "123456789"
			}`, user.Username)

		body := []byte(requestBody)

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
		requestBody := fmt.Sprintf(`{
			"username": "%s",
			"password": "%s"
			}`, user.Username, password)

		body := []byte(requestBody)

		req := httptest.NewRequest(
			http.MethodPost,
			"/login",
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		type LoginResponse struct {
			Data struct {
				Token    string `json:"token"`
				Username string `json:"username"`
			} `json:"data"`
			Message string `json:"message"`
		}

		var dataResponse LoginResponse

		err := json.Unmarshal(w.Body.Bytes(), &dataResponse)

		assert.NoError(t, err)
		assert.Equal(t, "login berhasil", dataResponse.Message)
		assert.NotEmpty(t, dataResponse.Data)
	})

	database.DB.Exec("DELETE FROM users")
}

func Test_CreateNote(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	t.Run("response 401 when authorization is missing", func(t *testing.T) {
		body := []byte(`{
				"title": "Note1",
				"content": "content note1"
			}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/note",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{
				"message": "authorization header tidak ditemukan"
			}`, w.Body.String())
	})

	_, _, _, authToken, _ := databaseHelper()

	t.Run("Add note success", func(t *testing.T) {
		body := []byte(`{
				"title": "Note1",
				"content": "content note1"
			}`)

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/note",
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var note entitiesNote.Note

		err := database.DB.Where("title = ?", "Note1").First(&note).Error

		assert.NoError(t, err)
		expectedResponse := fmt.Sprintf(`{
				"message": "note berhasil ditambahkan",
				"data": {
					"ID": "%s",
					"Title": "%s"
				}
			}`, note.ID, note.Title)

		assert.JSONEq(t, expectedResponse, w.Body.String())

	})

	database.DB.Exec("DELETE FROM notes")
	database.DB.Exec("DELETE FROM users")
}

func Test_GetNoteById(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	_, _, _, authToken, path := databaseHelper()

	t.Run("should response 400 when authorization is missing", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodGet,
			path,
			nil,
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{
				"message": "authorization header tidak ditemukan"
			}`, w.Body.String())
	})

	t.Run("get note by id success", func(t *testing.T) {
		var exisistNote entitiesNote.Note
		err := database.DB.First(&exisistNote).Error
		assert.NoError(t, err)

		req := httptest.NewRequest(
			http.MethodGet,
			path,
			nil,
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Body.String())
	})

	database.DB.Exec("DELETE FROM notes")
	database.DB.Exec("DELETE FROM users")
}

func Test_EditNoteById(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	_, _, _, authToken, path := databaseHelper()

	t.Run("should response 400 when authorization is missing", func(t *testing.T) {
		body := []byte(`{
				"title": "new title",
				"content": "new content"
			}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			path,
			bytes.NewBuffer(body),
		)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{
				"message": "authorization header tidak ditemukan"
			}`, w.Body.String())
	})

	t.Run("should be error when note to edit not found", func(t *testing.T) {
		body := []byte(`{
				"title": null,
				"content": null
			}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			path,
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

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
		body := []byte(`{
				"title": "new title"
				}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			path,
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var note entitiesNote.Note
		database.DB.First(&note)

		responseExpected := fmt.Sprintf(`{
				"message": "Note berhasil diperbarui",
				"data": "%s"
			}`, note.ID)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())
		assert.Equal(t, "new title", note.Title)
	})

	t.Run("should success update note when updating title and content", func(t *testing.T) {
		body := []byte(`{
				"title": "update new title",
				"content": "new content"
			}`)

		req := httptest.NewRequest(
			http.MethodPatch,
			path,
			bytes.NewBuffer(body),
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		var note entitiesNote.Note

		database.DB.First(&note)

		responseExpected := fmt.Sprintf(`{
				"message": "Note berhasil diperbarui",
				"data": "%s"
			}`, note.ID)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, responseExpected, w.Body.String())
		assert.Equal(t, "update new title", note.Title)
		assert.Equal(t, "new content", note.Content)
	})

	database.DB.Exec("DELETE FROM notes")
	database.DB.Exec("DELETE FROM users")
}

func Test_DeleteNoteById(t *testing.T) {
	router := gin.New()
	Router(router, database.DB)

	_, _, note, authToken, path := databaseHelper()

	t.Run("should response 400 when id not found", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodDelete,
			"/api/note/:id",
			nil,
		)

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{
				"message": "gagal menghapus note",
				"error": "id tidak ditemukan"
			}`, w.Body.String())
	})

	t.Run("delete success", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodDelete,
			path,
			nil,
		)

		log.Println("=====================INI ADALAHH DELETE==================")
		log.Println(path)
		log.Println("=====================INI ADALAHH DELETE==================")

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+authToken)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		response := fmt.Sprintf(`{
				"message": "berhasil menghapus note",
				"data": "%s"
			}`, note.ID)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, response, w.Body.String())
	})

	database.DB.Exec("DELETE FROM notes")
	database.DB.Exec("DELETE FROM users")
}
