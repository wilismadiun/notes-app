package security

type AuthToken interface {
	GenerateToken(id string) (string, error)
	ValidateToken(tokenString string) (string, error)
}
