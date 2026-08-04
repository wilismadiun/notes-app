package domains

import "notes-app/src/modul/user/Domains/entities"

type UserRepository interface {
	Createuser(user *entities.User) error
	VerifyUsername(user *entities.User) error
}
