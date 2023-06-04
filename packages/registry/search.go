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

func NewSearchController(pdb *sql.DB, rdb *redis.Client, c *gin.Context) controller.Search {

	u := usecase.NewUserUseCase(repository.NewUserRepository(c))
	p := presenter.NewSearchPresenter(c)
	d := gateway.NewDBRepository(rdb)
	cr := repository.NewQuestionCookiesRepository(c)
	s := usecase.NewSearchUseCase(
		repository.NewSearchRepository(),
		d,
		u,
		usecase.NewQuestionUseCase(
			gateway.NewQuestionRepository(pdb), cr, d, u, p,
		),
		p,
	)
	return controller.NewSearchController(s)
}
