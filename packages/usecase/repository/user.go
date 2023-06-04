package repository

type UserRepository interface {
	SetFingerprintAsCookie(coockiename string) error
	GetFingerprintFromCookie(coockieName string) (string, error)
}
