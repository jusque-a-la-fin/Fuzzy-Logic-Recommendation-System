package usecase

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"vehicles/packages/usecase/repository"
)

type QuestionInput interface {
	ChooseQuestion(fingerprint string) (string, []string, error)
	SetQuestionID(questionID string) error
	GetQuestionID(cookieName string) (string, error)
	GetAnswer(answer string) error
}

type questionUseCase struct {
	questionRepository        repository.QuestionRepository
	questionCookiesRepository repository.QuestionCookiesRepository
	dbRepository              repository.SearchDBRepository
	userUseCase               UserInput
	output                    SearchOutput
}

func NewQuestionUseCase(r repository.QuestionRepository, c repository.QuestionCookiesRepository, d repository.SearchDBRepository, u UserInput, o SearchOutput) QuestionInput {
	return &questionUseCase{r, c, d, u, o}
}

func (qu *questionUseCase) ChooseQuestion(fingerprint string) (string, []string, error) {

	questionIDs, err := qu.questionRepository.GetIdsOfUnansweredQuestions(fingerprint)
	if err != nil {
		return "", nil, err
	}

	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(questionIDs))
	questionID := questionIDs[randIndex]
	fmt.Println(questionID)

	question, possibleAnswers, err := qu.questionRepository.GetQuestion(questionID)
	if err != nil {
		return "", nil, err
	}

	_questionID := strconv.Itoa(questionID)

	err = qu.questionCookiesRepository.SetQuestionID(_questionID)
	if err != nil {
		return "", nil, err
	}

	return question, possibleAnswers, nil
}

func (qu *questionUseCase) GetAnswer(answer string) error {

	cookieName := "fingerprint"
	fingerprint, err := qu.userUseCase.GetFingerprint(cookieName)

	if err != nil {
		return err
	}

	cookieName = "questionID"
	questionID, err := qu.questionCookiesRepository.GetQuestionID(cookieName)
	if err != nil {
		return err
	}

	err = qu.questionRepository.InsertAnswer(fingerprint, questionID, answer)
	if err != nil {
		return err
	}

	cards, err := qu.dbRepository.GetCarCardData()
	if err != nil {
		return err
	}
	htmlFileName := "offer_for_search.html"
	qu.output.ShowCars(htmlFileName, cards)

	return nil
}

func (qu *questionUseCase) SetQuestionID(questionID string) error {
	err := qu.questionCookiesRepository.SetQuestionID(questionID)
	if err != nil {
		return err
	}
	return nil
}

func (qu *questionUseCase) GetQuestionID(cookieName string) (string, error) {
	questionID, err := qu.questionCookiesRepository.GetQuestionID(cookieName)
	if err != nil {
		return "", err
	}
	return questionID, nil
}
