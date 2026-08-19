package security

type HashPassword interface {
	HashingPassword(password string) (string, error)
	CompareHashPassword(password, hashPassword string) error
}
