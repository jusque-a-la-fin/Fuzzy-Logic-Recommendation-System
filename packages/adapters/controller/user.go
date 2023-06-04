package controller

import (
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/usecase"
)

type userController struct {
	userUseCase usecase.UserInput
}

type User interface {
	SetFingerprint(ctx adapters.Context) error
}

func NewUserController(ur usecase.UserInput) User {
	return &userController{ur}
}

func (uc *userController) SetFingerprint(ctx adapters.Context) error {

	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return nil
	}

	err := uc.userUseCase.SetFingerprint(&user)
	if err != nil {
		return err
	}

	return nil
}
