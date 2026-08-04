package security

type HashPassword interface {
	HashingPassword(password string) (string, error)
}
