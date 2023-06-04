package usecase

import (
	"fmt"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"
)

type SearchInput interface {
	GetCars(search *models.Search) error
	PassCarsData() error
	ShowMainPage(htmlFileName string) error
	ShowCarCard(id int) error
}

type SearchOutput interface {
	ShowCarsWithSurvey(htmlFileName string, cards *[]models.CarCard, question string, possible_answers []string)
	ShowCars(htmlFileName string, cards *[]models.CarCard)
	ShowCarCard(htmlFileName string, cards *[]models.CarCard, id int)
	ShowMainPage(htmlFileName string)
}

type searchUseCase struct {
	searchRepository repository.SearchRepository
	dbRepository     repository.SearchDBRepository
	userUseCase      UserInput
	questionUseCase  QuestionInput
	output           SearchOutput
}

func NewSearchUseCase(r repository.SearchRepository, d repository.SearchDBRepository, u UserInput, q QuestionInput, o SearchOutput) SearchInput {
	return &searchUseCase{r, d, u, q, o}
}

func (su *searchUseCase) GetCars(search *models.Search) error {
	fmt.Println(su)
	fmt.Println(search)
	cars := su.searchRepository.GetCarsUsingScraping(search)

	err := su.dbRepository.LoadCarCardsData(cars)

	if err != nil {
		return err
	}

	return nil
}

func (su *searchUseCase) PassCarsData() error {

	cards, err := su.dbRepository.GetCarCardData()
	if err != nil {
		return err
	}

	cookieName := "fingerprint"

	fingerprint, err := su.userUseCase.GetFingerprint(cookieName)

	if err != nil {
		return err
	}

	question, possibleAnswers, err := su.questionUseCase.ChooseQuestion(fingerprint)
	if err != nil {
		return err
	}

	htmlFileName := "offer_for_search.html"

	su.output.ShowCarsWithSurvey(htmlFileName, cards, question, possibleAnswers)

	return nil
}

func (su *searchUseCase) ShowCarCard(id int) error {
	cards, err := su.dbRepository.GetCarCardData()
	if err != nil {
		return err
	}
	htmlFileName := "car_card.html"
	su.output.ShowCarCard(htmlFileName, cards, id)
	return nil
}

func (su *searchUseCase) ShowMainPage(htmlFileName string) error {
	su.output.ShowMainPage(htmlFileName)
	return nil
}
