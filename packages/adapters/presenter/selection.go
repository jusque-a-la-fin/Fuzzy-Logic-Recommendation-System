package presenter

import (
	"fmt"
	"net/http"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/usecase"

	"github.com/gin-gonic/gin"
)

type selectionPresenter struct {
	ctx adapters.Context
}

type Selection interface {
	ChoosePriorities()
	ChoosePrices()
	ChooseManufacturers()
	ShowResultOfFuzzyAlgorithm(cars *[]models.Car)
	ShowCarCard(cars *[]models.Car, id int)
}

func NewSelectionPresenter(ctx adapters.Context) usecase.SelectionOutput {
	return &selectionPresenter{ctx}
}

func (s *selectionPresenter) ChoosePriorities() {
	fmt.Println("приколитесь")
	fmt.Println(s)
	fmt.Println(s.ctx)
	s.ctx.HTML(http.StatusOK, "priorities.html", nil)
}

func (s *selectionPresenter) ChoosePrices() {
	s.ctx.HTML(http.StatusOK, "prices.html", nil)
}

func (s *selectionPresenter) ChooseManufacturers() {
	s.ctx.HTML(http.StatusOK, "manufacturers.html", nil)
}

func (s *selectionPresenter) ShowResultOfFuzzyAlgorithm(cars *[]models.Car) {
	s.ctx.HTML(http.StatusOK, "offer_for_selection.html", gin.H{"Cars": cars, "Quantity": len(*cars)})
}

func (s *selectionPresenter) ShowCarCard(cars *[]models.Car, id int) {

	s.ctx.HTML(http.StatusOK, "car_card_selection.html", gin.H{"Make": (*cars)[id-1].Make, "Model": (*cars)[id-1].Model, "Price": (*cars)[id-1].Offering.Price, "Images": (*cars)[id-1].Offering.PhotoURLs[:2]})
}
