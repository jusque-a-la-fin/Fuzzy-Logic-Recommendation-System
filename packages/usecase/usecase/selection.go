package usecase

import (
	"fmt"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"
)

type SelectionInput interface {
	ChoosePriorities() error
	SetPriorities(priorities *[]string) error
	ChoosePrice() error
	SetPrice(minPrice, maxPrice, deviation string) error
	ChooseManufacturers() error
	SetManufacturers(manufacturers *[]string) error
	GetSelection() error
	ShowCarCard(id int) error
}

type SelectionOutput interface {
	ChoosePriorities()
	ChoosePrices()
	ChooseManufacturers()
	ShowResultOfFuzzyAlgorithm(*[]models.Car)
	ShowCarCard(cars *[]models.Car, id int)
}

type selectionUseCase struct {
	selectionRepository        repository.SelectionRepository
	selectionCookiesRepository repository.SelectionCookiesRepository
	userUseCase                UserInput
	output                     SelectionOutput
}

func NewSelectionUseCase(r repository.SelectionRepository, c repository.SelectionCookiesRepository, u UserInput, o SelectionOutput) SelectionInput {
	return &selectionUseCase{r, c, u, o}
}

func (su *selectionUseCase) ChoosePriorities() error {
	su.output.ChoosePriorities()
	return nil
}

func (su *selectionUseCase) SetPriorities(priorities *[]string) error {

	cookieName := "fingerprint"
	fingerprint, err := su.userUseCase.GetFingerprint(cookieName)

	if err != nil {
		return err
	}
	cookieName = "preferencesID"

	preferencesID, _ := su.selectionCookiesRepository.GetUserPreferencesID(cookieName)
	if preferencesID == "" {

		preferencesID, err := su.selectionRepository.InsertPriorities(fingerprint, priorities)
		if err != nil {
			return err
		}

		err = su.selectionCookiesRepository.SetUserPreferencesID(preferencesID)
		if err != nil {
			return err
		}
	} else {

		fmt.Println(preferencesID)
		err = su.selectionRepository.UpdatePriorities(preferencesID, fingerprint, priorities)
		if err != nil {
			return err
		}
	}

	return nil
}

func (su *selectionUseCase) ChoosePrice() error {
	su.output.ChoosePrices()
	return nil
}

func (su *selectionUseCase) SetPrice(minPrice, maxPrice, deviation string) error {
	cookieName := "fingerprint"
	fingerprint, err := su.userUseCase.GetFingerprint(cookieName)
	if err != nil {
		return err
	}
	cookieName = "preferencesID"
	preferencesID, _ := su.selectionCookiesRepository.GetUserPreferencesID(cookieName)

	err = su.selectionRepository.SetPrice(preferencesID, fingerprint, minPrice, maxPrice, deviation)
	if err != nil {
		return err
	}
	return nil
}

func (su *selectionUseCase) ChooseManufacturers() error {
	su.output.ChooseManufacturers()
	return nil
}

func (su *selectionUseCase) SetManufacturers(manufacturers *[]string) error {
	cookieName := "fingerprint"
	fingerprint, err := su.userUseCase.GetFingerprint(cookieName)

	if err != nil {
		return err
	}
	cookieName = "preferencesID"
	preferencesID, _ := su.selectionCookiesRepository.GetUserPreferencesID(cookieName)

	err = su.selectionRepository.SetManufacturers(preferencesID, fingerprint, manufacturers)
	if err != nil {
		return err
	}
	return nil
}

func (su *selectionUseCase) GetSelection() error {
	cookieName := "preferencesID"
	preferencesID, _ := su.selectionCookiesRepository.GetUserPreferencesID(cookieName)

	selection, err := su.selectionRepository.GetSelection(preferencesID)
	if err != nil {
		return err
	}

	cars, err := su.selectionRepository.SelectCars(selection)

	if err != nil {
		return err
	}

	cars = generateResultOfFuzzyAlgorithm(cars, selection.Priorities)

	su.selectionRepository.LoadCarsData(*cars)

	su.output.ShowResultOfFuzzyAlgorithm(cars)

	return nil
}

func (su *selectionUseCase) ShowCarCard(id int) error {
	cars, _ := su.selectionRepository.GetCarsData()
	su.output.ShowCarCard(&cars, id)
	return nil
}
