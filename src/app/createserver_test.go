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
	"reflect"
	"testing"

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

func Test_CreateServer(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

	t.Run("Create User", func(t *testing.T) {
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

		var hashPassword string
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
			hashPassword = user.Password
			log.Println("===============================ini adalah hash password=============================")
			log.Println(hashPassword)
			log.Println("===============================ini adalah hash password=============================")

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
	})

	var authToken string
	t.Run("Login User", func(t *testing.T) {
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

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBuffer(body),
			)

			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			log.Println("=========================")
			log.Println(reflect.TypeOf(w.Body.String()))
			log.Println("=========================")
			log.Println(reflect.TypeOf(w.Body.Bytes()))

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

			// Mendapatkan token authentication
			authToken = dataResponse.Data.Token
		})
	})

	var path string
	t.Run("Create Note", func(t *testing.T) {
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

			fmt.Println("==========================auth token di add note==============================")
			fmt.Println(authToken)
			fmt.Println("==========================auth token di add note==============================")

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

			path = fmt.Sprintf("/api/note/%s", note.ID)
		})
	})

	var notecoba entitiesNote.Note
	database.DB.First(&notecoba)
	fmt.Println("=====================================ini adalah path==========================================")
	fmt.Println(path)
	fmt.Printf("ini adalah data note \n %s", notecoba)
	fmt.Println("=====================================ini adalah path==========================================")

	t.Run("Get Note by ID", func(t *testing.T) {
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

			// database.DB.Exec("DELETE FROM notes")
		})
	})

	t.Run("Edit Note by Id", func(t *testing.T) {
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
	})

	t.Run("Delete Note by Id", func(t *testing.T) {
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
			var note entitiesNote.Note
			database.DB.First(&note)

			req := httptest.NewRequest(
				http.MethodDelete,
				path,
				nil,
			)

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
		database.DB.Exec("DELETE FROM users")
	})
}
