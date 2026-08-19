package app

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"notes-app/src/commons/database"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_Server(t *testing.T) {
	router := gin.New()

	Router(router, database.DB)

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
	log.Println(w.Body.String())

	type LoginResponse struct {
		Data struct {
			Token    string `json:"token"`
			Username string `json:"username"`
		} `json:"data"`
		Message string `json:"message"`
	}

	var response LoginResponse

	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	token := response.Data.Token

	log.Println("TOKEN:", token)
	// assert.Equal(t, "hai", "halo")
}
