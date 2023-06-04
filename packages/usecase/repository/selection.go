package repository

import "vehicles/packages/domain/models"

type SelectionRepository interface {
	InsertPriorities(fingerprint string, priorities *[]string) (string, error)
	UpdatePriorities(preferencesID, fingerprint string, priorities *[]string) error
	SetPrice(preferencesID, fingerprint, minPrice, maxPrice, deviation string) error
	SetManufacturers(preferencesID, fingerprint string, manufacturers *[]string) error
	GetSelection(preferencesID string) (*models.Selection, error)
	SelectCars(*models.Selection) (*[]models.Car, error)
	LoadCarsData(cars []models.Car) error
	GetCarsData() ([]models.Car, error)
}
