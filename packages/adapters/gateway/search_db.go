package gateway

import (
	"context"
	"encoding/json"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"

	"github.com/redis/go-redis/v9"
)

type searchRepository struct {
	rdb *redis.Client
}

func NewDBRepository(rdb *redis.Client) repository.SearchDBRepository {
	return &searchRepository{rdb}
}

func (sr *searchRepository) LoadCarCardsData(cards *[]models.CarCard) error {

	carCardsJSON, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := sr.rdb.Set(ctx, "carCards", string(carCardsJSON), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (sr *searchRepository) GetCarCardData() (*[]models.CarCard, error) {

	ctx := context.Background()
	carCardsJSON, err := sr.rdb.Get(ctx, "carCards").Result()
	if err != nil {
		return nil, err
	}

	var carCards *[]models.CarCard
	if err := json.Unmarshal([]byte(carCardsJSON), &carCards); err != nil {
		return nil, err
	}
	return carCards, nil
}
