package repository

import (
	"vehicles/packages/adapters"
	"vehicles/packages/usecase/repository"
)

type questionCookiesRepository struct {
	ctx adapters.Context
}

func NewQuestionCookiesRepository(ctx adapters.Context) repository.QuestionCookiesRepository {
	return &questionCookiesRepository{ctx}
}

func (q questionCookiesRepository) SetQuestionID(questionID string) error {
	q.ctx.SetCookie("questionID", questionID, 0, "questionID", "localhost", false, true)
	return nil
}

func (q questionCookiesRepository) GetQuestionID(cookieName string) (string, error) {

	questionID, err := q.ctx.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	return questionID, nil
}
