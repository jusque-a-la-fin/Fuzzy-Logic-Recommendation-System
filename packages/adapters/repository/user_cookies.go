package repository

import (
	"vehicles/packages/adapters"
	"vehicles/packages/usecase/repository"
)

type userRepository struct {
	ctx adapters.Context
}

func NewUserRepository(ctx adapters.Context) repository.UserRepository {
	return &userRepository{ctx}
}

func (u userRepository) SetFingerprintAsCookie(fingerprint string) error {
	u.ctx.SetCookie("fingerprint", fingerprint, 0, "fingerprint", "localhost", false, true)
	return nil
}

func (u userRepository) GetFingerprintFromCookie(cookieName string) (string, error) {

	fingerprint, err := u.ctx.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return fingerprint, nil
}
