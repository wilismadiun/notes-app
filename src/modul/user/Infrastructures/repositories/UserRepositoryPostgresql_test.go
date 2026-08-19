package repositories

import (
	"os"
	"testing"

	"notes-app/src/commons/database"
	"notes-app/src/modul/user/Domains/entities"

	"github.com/stretchr/testify/assert"
)

var repo *UserRepository

func TestMain(m *testing.M) {
	database.ConnectPostgresql(".test.env")

	repo = &UserRepository{
		DB: database.DB,
	}

	code := m.Run()

	os.Exit(code)
}

func TestVerifyUsername_NotExists(t *testing.T) {
	user := entities.User{
		Username: "username_yang_tidak_ada",
	}

	err := repo.VerifyUsername(&user)

	assert.NoError(t, err)
}

func TestVerifyUsername_AlreadyExists(t *testing.T) {

	user := entities.User{
		Username: "john",
		Password: "password",
	}

	err := repo.Createuser(&user)
	assert.NoError(t, err)

	err = repo.VerifyUsername(&user)

	assert.EqualError(t, err, "username sudah digunakan")
}

func TestAddUser_Success(t *testing.T) {
	user := entities.User{
		Username: "joko",
		Password: "12345678",
	}

	err := repo.Createuser(&user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	database.DB.Exec("DELETE FROM users")
}

func Test_FindUserByUsername(t *testing.T) {
	user := entities.User{
		ID:       "user-12345",
		Username: "Jaya123",
		Password: "12345678",
	}

	err := repo.Createuser(&user)
	assert.NoError(t, err)

	exisistUser := entities.User{
		Username: "Jaya123",
	}

	err = repo.FindUserByUsername(&exisistUser)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, exisistUser.ID)
	assert.Equal(t, user.Password, exisistUser.Password)
}
