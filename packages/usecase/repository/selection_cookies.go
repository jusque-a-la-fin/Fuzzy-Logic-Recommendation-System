package repository

type SelectionCookiesRepository interface {
	SetUserPreferencesID(preferencesId string) error
	GetUserPreferencesID(cookieName string) (string, error)
}
