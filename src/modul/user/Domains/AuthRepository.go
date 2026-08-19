package domains

type AuthenticationRepository interface {
	SaveRefreshToken(refreshToken string) error
}
