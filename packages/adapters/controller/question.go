package controller

import (
	"database/sql"
	"vehicles/packages/adapters"
	"vehicles/packages/usecase/usecase"

	"github.com/redis/go-redis/v9"
)

type questionController struct {
	questionUseCase usecase.QuestionInput
}

type Question interface {
	GetAnswer(ctx adapters.Context, rdb *redis.Client, db *sql.DB) error
}

func NewQuestionController(qn usecase.QuestionInput) Question {
	return &questionController{qn}
}

func (qc *questionController) GetAnswer(ctx adapters.Context, rdb *redis.Client, db *sql.DB) error {

	answer := ctx.PostForm("radio")

	err := qc.questionUseCase.GetAnswer(answer)
	if err != nil {
		return err
	}
	return nil
}
