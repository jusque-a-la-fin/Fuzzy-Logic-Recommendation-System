package gateway

import (
	"database/sql"
	"fmt"
	"vehicles/packages/usecase/repository"
)

type questionRepository struct {
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) repository.QuestionRepository {
	return &questionRepository{db}
}

func (qr *questionRepository) GetIdsOfUnansweredQuestions(fingerprint string) ([]int, error) {

	// формируем запрос
	query := `
        SELECT id
        FROM questions
        WHERE id NOT IN (
          SELECT question_id
          FROM user_responses
          WHERE user_id = (
            SELECT id
            FROM users
            WHERE fingerprint = $1
          )
        );
    `

	// отправляем запрос к БД
	rows, err := qr.db.Query(query, fingerprint)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// обрабатываем результаты запроса
	var questionIDs []int
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		questionIDs = append(questionIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questionIDs, nil
}

func (qr *questionRepository) GetQuestion(questionID int) (string, []string, error) {

	query := `
        SELECT question
        FROM questions
        WHERE id = $1;
    `

	// отправляем запрос к БД
	var questionText string
	err := qr.db.QueryRow(query, questionID).Scan(&questionText)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Запись не найдена")
		} else {
			return "", nil, nil
		}
	}

	// выводим результаты запроса на экран
	fmt.Printf("Вопрос с id %d: %s\n", questionID, questionText)

	// формируем запрос
	query = `
        SELECT possible_answer
        FROM possible_answers
        WHERE question_id = $1;
    `

	// отправляем запрос к БД
	rows, err := qr.db.Query(query, questionID)
	if err != nil {
		return "", nil, nil
	}
	defer rows.Close()

	// обрабатываем результаты запроса
	var possibleAnswers []string
	var possibleAnswer string
	for rows.Next() {
		err := rows.Scan(&possibleAnswer)
		if err != nil {
			return "", nil, nil
		}
		possibleAnswers = append(possibleAnswers, possibleAnswer)
	}
	if err := rows.Err(); err != nil {
		return "", nil, nil
	}

	return questionText, possibleAnswers, nil
}

func (qr *questionRepository) InsertAnswer(fingerprint, questionID, answer string) error {

	var exists bool
	err := qr.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE fingerprint = $1)", fingerprint).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {

		sqlQuery := `
    WITH new_user AS (
        INSERT INTO users (fingerprint)
        SELECT CAST($1 AS VARCHAR(40))
        WHERE NOT EXISTS (
            SELECT 1 FROM users WHERE fingerprint = $1
        )
        RETURNING id
    ),
    upsert_response AS (
        INSERT INTO user_responses (user_id, question_id, answer)
        SELECT new_user.id, $2, $3
        FROM new_user
        RETURNING *
    )
    SELECT * FROM upsert_response
`

		_, err := qr.db.Exec(sqlQuery, fingerprint, questionID, answer)
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else {

		sqlQuery := `
  WITH existing_user AS (
	  SELECT id FROM users WHERE fingerprint = $1
  ),
  upsert_response AS (
	  INSERT INTO user_responses (user_id, question_id, answer)
	  SELECT existing_user.id, $2, $3
	  FROM existing_user
	  RETURNING *
  )
  SELECT * FROM upsert_response
`

		_, err := qr.db.Exec(sqlQuery, fingerprint, questionID, answer)
		if err != nil {
			panic(err)
		}

	}
	return nil
}
