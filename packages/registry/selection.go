package registry

import (
	"database/sql"
	"vehicles/packages/adapters/controller"
	"vehicles/packages/adapters/gateway"
	"vehicles/packages/adapters/presenter"
	"vehicles/packages/adapters/repository"
	"vehicles/packages/usecase/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewSelectionController(preferencesDB *sql.DB, vehiclesDB *sql.DB, rdb *redis.Client, c *gin.Context) controller.Selection {

	sl := usecase.NewSelectionUseCase(
		gateway.NewSelectionRepository(preferencesDB, vehiclesDB, rdb),
		repository.NewSelectionCoockiesRepository(c),
		usecase.NewUserUseCase(repository.NewUserRepository(c)),
		presenter.NewSelectionPresenter(c),
	)
	return controller.NewSelectionController(sl)
}
