package registry

import (
	"vehicles/packages/adapters/controller"
	"vehicles/packages/adapters/repository"
	"vehicles/packages/usecase/usecase"

	"github.com/gin-gonic/gin"
)

func NewUserController(c *gin.Context) controller.User {
	u := usecase.NewUserUseCase(
		repository.NewUserRepository(c),
	)
	return controller.NewUserController(u)
}
