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

func NewQuestionController(pdb *sql.DB, rdb *redis.Client, c *gin.Context) controller.Question {
	qr := gateway.NewQuestionRepository(pdb)
	u := usecase.NewUserUseCase(repository.NewUserRepository(c))
	p := presenter.NewSearchPresenter(c)
	d := gateway.NewDBRepository(rdb)
	qc := repository.NewQuestionCookiesRepository(c)
	q := usecase.NewQuestionUseCase(qr, qc, d, u, p)
	return controller.NewQuestionController(q)
}
