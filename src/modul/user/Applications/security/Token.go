package security

type Token interface {
	GenerateToken(id string) (string, string, error)
}
