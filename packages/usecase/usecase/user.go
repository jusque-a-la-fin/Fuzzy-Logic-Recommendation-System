package usecase

import (
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"
)

type UserInput interface {
	SetFingerprint(user *models.User) error
	GetFingerprint(coockieName string) (string, error)
}

type userUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) UserInput {
	return &userUseCase{r}
}

func (uu *userUseCase) SetFingerprint(user *models.User) error {
	err := uu.userRepository.SetFingerprintAsCookie(user.Fingerprint)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUseCase) GetFingerprint(cookieName string) (string, error) {

	fingerprint, err := uu.userRepository.GetFingerprintFromCookie(cookieName)
	if err != nil {
		return "", err
	}

	return fingerprint, nil
}
