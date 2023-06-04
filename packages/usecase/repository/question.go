package repository

type QuestionRepository interface {
	GetIdsOfUnansweredQuestions(fingerprint string) ([]int, error)
	GetQuestion(id int) (string, []string, error)
	InsertAnswer(fingerprint, questionID, answer string) error
}
