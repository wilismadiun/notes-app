package app

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"notes-app/src/commons/database"
	"notes-app/src/modul/user/Domains/entities"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	database.ConnectPostgresql(".env")

	gin.SetMode(gin.TestMode)

	code := m.Run()

	os.Exit(code)
}

func TestUserModule(t *testing.T) {

	router := gin.New()

	RegisterRouter(router, database.DB)

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
			Email:    "jaya@gmail.com",
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
