package repository

import "vehicles/packages/domain/models"

type SearchRepository interface {
	GetCarsUsingScraping(search *models.Search) *[]models.CarCard
}
