package presenter

import (
	"fmt"
	"net/http"
	"regexp"
	"vehicles/packages/adapters"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/usecase"

	"github.com/gin-gonic/gin"
)

type searchPresenter struct {
	ctx adapters.Context
}

type Search interface {
	ShowCarsWithSurvey(htmlFileName string, cards *[]models.CarCard, question string, possible_answers []string)
	ShowCars(htmlFileName string, cards *[]models.CarCard)
	ShowMainPage(htmlFileName string)
}

func NewSearchPresenter(ctx adapters.Context) usecase.SearchOutput {
	return &searchPresenter{ctx}
}

func (s *searchPresenter) ShowCarsWithSurvey(htmlFileName string, cards *[]models.CarCard, question string, possibleAnswers []string) {

	formattedQuestion := formatQuestion(question)
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"Cars": cards, "Quantity": len(*cards), "NotAnswered": true, "QuestionPart1": formattedQuestion["Part1"], "QuestionPart2": formattedQuestion["Part2"], "QuestionPart3": formattedQuestion["Part3"], "PossibleAnswers": possibleAnswers})
}

func (s *searchPresenter) ShowCars(htmlFileName string, cards *[]models.CarCard) {
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"Cars": cards, "Quantity": len(*cards)})
}

func (s *searchPresenter) ShowCarCard(htmlFileName string, cards *[]models.CarCard, id int) {
	part_of_link := "http://localhost:8080/search/card/"
	s.ctx.HTML(http.StatusOK, htmlFileName, gin.H{"Name": (*cards)[id-1].Name, "Price": (*cards)[id-1].Price, "Images": (*cards)[id-1].Images[:6], "Characteristics": (*cards)[id-1].Characteristics, "Part_of_link": part_of_link})
}
func (s *searchPresenter) SendResponseCarData(cars []models.CarCard) {
	s.ctx.JSON(http.StatusOK, cars)
}

func (s *searchPresenter) ShowMainPage(htmlFileName string) {

	s.ctx.HTML(http.StatusOK, htmlFileName, nil)
}

func formatQuestion(question string) map[string]string {
	// Объявляем значения для регулярных выражений
	reg1 := `\d+\s*л/\d+\s*км`
	reg2 := `\d+\s*секунд[ы]?\s*от\s*0\s*до\s*\d+\s*км\/ч`

	reg1Compile := regexp.MustCompile(reg1)
	reg1Match := reg1Compile.FindString(question)

	result := make(map[string]string)
	if reg1Match == "" {
		reg2Compile := regexp.MustCompile(reg2)
		reg2Match := reg2Compile.FindString(question)
		fmt.Println("A", reg2Match)
		result["Part1"] = "Как Вы думаете, время разгона "
		result["Part2"] = reg2Match
		result["Part3"] = " — ,  это:"
	} else {
		result["Part1"] = "Как Вы думаете, расход топлива в смешанном цикле "
		result["Part2"] = reg1Match
		result["Part3"] = " — ,  это:"
	}
	return result
}
