package controller

import (
	"database/sql"
	"strconv"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/usecase"

	"github.com/redis/go-redis/v9"
)

type searchController struct {
	searchUseCase usecase.SearchInput
}

type Search interface {
	GetCars(ctx adapters.Context) error
	PassCarsData(ctx adapters.Context, rdb *redis.Client, db *sql.DB) error
	ShowCarCard(ctx adapters.Context, rdb *redis.Client) error
	ShowMainPage(ctx adapters.Context, htmlFileName string) error
}

func NewSearchController(sr usecase.SearchInput) Search {
	return &searchController{sr}
}

func (sc *searchController) GetCars(ctx adapters.Context) error {

	var search models.Search
	if err := ctx.Bind(&search); err != nil {
		return err
	}

	err := sc.searchUseCase.GetCars(&search)

	if err != nil {
		return err
	}

	return nil
}

func (sc *searchController) PassCarsData(ctx adapters.Context, rdb *redis.Client, db *sql.DB) error {

	err := sc.searchUseCase.PassCarsData()
	if err != nil {
		return err
	}

	return nil
}

func (sc *searchController) ShowCarCard(ctx adapters.Context, rdb *redis.Client) error {
	_id := ctx.Param("id")
	id, _ := strconv.Atoi(_id)

	err := sc.searchUseCase.ShowCarCard(id)
	if err != nil {
		return err
	}
	return nil
}

func (sc *searchController) ShowMainPage(ctx adapters.Context, htmlFileName string) error {
	err := sc.searchUseCase.ShowMainPage(htmlFileName)
	return err
}
