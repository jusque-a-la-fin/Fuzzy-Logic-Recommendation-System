package repository

import (
	"vehicles/packages/adapters"
	"vehicles/packages/usecase/repository"
)

type selectionCookiesRepository struct {
	ctx adapters.Context
}

func NewSelectionCoockiesRepository(ctx adapters.Context) repository.SelectionCookiesRepository {
	return &selectionCookiesRepository{ctx}
}

func (s selectionCookiesRepository) SetUserPreferencesID(preferencesID string) error {

	s.ctx.SetCookie("preferencesID", preferencesID, 0, "preferencesID", "localhost", false, true)
	return nil
}

func (s selectionCookiesRepository) GetUserPreferencesID(cookieName string) (string, error) {

	preferencesID, err := s.ctx.Cookie(cookieName)
	if err != nil {

		return "", err
	}
	return preferencesID, nil
}
