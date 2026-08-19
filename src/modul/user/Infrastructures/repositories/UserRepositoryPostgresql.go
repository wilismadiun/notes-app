package repositories

import (
	"errors"
	"notes-app/src/modul/user/Domains/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (h *UserRepository) Createuser(user *entities.User) error {
	user.ID = uuid.New().String()
	err := h.DB.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (h *UserRepository) VerifyUsername(user *entities.User) error {
	var existingUser entities.User

	err := h.DB.
		Where("username = ?", user.Username).
		First(&existingUser).Error

	if err == nil {
		// Username ditemukan
		return errors.New("username sudah digunakan")
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Username belum ada, boleh lanjut
		return nil
	}

	// Error database
	return err
}

func (h *UserRepository) FindUserByUsername(user *entities.User) error {
	err := h.DB.Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		return err
	}

	return nil
}
