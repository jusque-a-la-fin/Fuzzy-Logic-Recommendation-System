package repository

import "vehicles/packages/domain/models"

type SearchDBRepository interface {
	LoadCarCardsData(cards *[]models.CarCard) error
	GetCarCardData() (*[]models.CarCard, error)
}
