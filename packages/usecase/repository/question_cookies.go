package repository

type QuestionCookiesRepository interface {
	SetQuestionID(questionID string) error
	GetQuestionID(cookieName string) (string, error)
}
