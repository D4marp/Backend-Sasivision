package repositories

import (
	"database/sql"
	"time"

	"github.com/sasivision/backend/internal/models"
)

type QuizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) GetActiveCategories() ([]models.QuizCategory, error) {
	rows, err := r.db.Query(`
		SELECT id, name, slug, description, display_order, is_active, created_at
		FROM quiz_categories WHERE is_active = 1 ORDER BY display_order ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.QuizCategory
	for rows.Next() {
		var category models.QuizCategory
		var isActive bool
		if err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.Description,
			&category.DisplayOrder, &isActive, &category.CreatedAt,
		); err != nil {
			return nil, err
		}
		category.IsActive = isActive
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *QuizRepository) GetCategoryByNameOrSlug(category string) (*models.QuizCategory, error) {
	var cat models.QuizCategory
	var isActive bool
	err := r.db.QueryRow(`
		SELECT id, name, slug, description, display_order, is_active, created_at
		FROM quiz_categories WHERE name = ? OR slug = ?`, category, category).Scan(
		&cat.ID, &cat.Name, &cat.Slug, &cat.Description,
		&cat.DisplayOrder, &isActive, &cat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	cat.IsActive = isActive
	return &cat, nil
}

func (r *QuizRepository) GetQuestionsByCategory(category string) ([]models.Quiz, error) {
	cat, err := r.GetCategoryByNameOrSlug(category)
	if err != nil {
		return nil, err
	}
	return r.GetQuestionsByCategoryID(cat.ID)
}

func (r *QuizRepository) GetQuestionsByCategoryID(categoryID int) ([]models.Quiz, error) {
	rows, err := r.db.Query(`
		SELECT id, category_id, type, question, image_url, sequence_order, created_at
		FROM quizzes WHERE category_id = ? ORDER BY sequence_order ASC`, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		var imageURL sql.NullString
		if err := rows.Scan(
			&quiz.ID, &quiz.CategoryID, &quiz.Type, &quiz.Question,
			&imageURL, &quiz.SequenceOrder, &quiz.CreatedAt,
		); err != nil {
			return nil, err
		}
		if imageURL.Valid {
			quiz.ImageURL = &imageURL.String
		}

		answers, err := r.getAnswersForQuiz(quiz.ID)
		if err != nil {
			return nil, err
		}
		quiz.Answers = answers
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func (r *QuizRepository) GetQuestionByID(id int) (*models.Quiz, error) {
	var quiz models.Quiz
	var imageURL sql.NullString
	err := r.db.QueryRow(`
		SELECT id, category_id, type, question, image_url, sequence_order, created_at
		FROM quizzes WHERE id = ?`, id).Scan(
		&quiz.ID, &quiz.CategoryID, &quiz.Type, &quiz.Question,
		&imageURL, &quiz.SequenceOrder, &quiz.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if imageURL.Valid {
		quiz.ImageURL = &imageURL.String
	}
	answers, err := r.getAnswersForQuiz(quiz.ID)
	if err != nil {
		return nil, err
	}
	quiz.Answers = answers
	return &quiz, nil
}

func (r *QuizRepository) GetAttemptOwnerEmail(attemptID int) (string, error) {
	var email string
	err := r.db.QueryRow(`
		SELECT u.email FROM quiz_attempts qa
		JOIN users u ON u.id = qa.user_id
		WHERE qa.id = ?`, attemptID).Scan(&email)
	return email, err
}

func (r *QuizRepository) getAnswersForQuiz(quizID int) ([]models.QuizAnswer, error) {
	rows, err := r.db.Query(`
		SELECT id, quiz_id, answer_key, answer_text, is_correct, created_at
		FROM quiz_answers WHERE quiz_id = ? ORDER BY answer_key ASC`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.QuizAnswer
	for rows.Next() {
		var answer models.QuizAnswer
		var answerKey sql.NullString
		var isCorrect bool
		if err := rows.Scan(
			&answer.ID, &answer.QuizID, &answerKey, &answer.AnswerText,
			&isCorrect, &answer.CreatedAt,
		); err != nil {
			return nil, err
		}
		if answerKey.Valid {
			answer.AnswerKey = &answerKey.String
		}
		answer.IsCorrect = isCorrect
		answers = append(answers, answer)
	}
	return answers, rows.Err()
}

func (r *QuizRepository) GetUserIDByEmail(email string) (int, error) {
	var userID int
	err := r.db.QueryRow(`SELECT id FROM users WHERE email = ?`, email).Scan(&userID)
	return userID, err
}

func (r *QuizRepository) CreateAttempt(attempt models.QuizAttempt) (int64, error) {
	userID, err := r.GetUserIDByEmail(attempt.Email)
	if err != nil {
		return 0, err
	}

	result, err := r.db.Exec(`
		INSERT INTO quiz_attempts (user_id, category_id, correct_count, total_count, score, start_time, end_time, finish_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		userID, attempt.CategoryID, attempt.CorrectCount, attempt.TotalCount,
		attempt.Score, attempt.StartTime, attempt.EndTime, attempt.FinishDate,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *QuizRepository) CreateAttemptDetails(attemptID int64, details []models.AttemptDetail) error {
	for _, detail := range details {
		_, err := r.db.Exec(`
			INSERT INTO quiz_attempt_details (quiz_attempt_id, quiz_id, type, user_answer)
			VALUES (?, ?, ?, ?)`,
			attemptID, detail.QuizID, detail.Type, detail.UserAnswer,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *QuizRepository) GetAttemptsByEmail(email string) ([]models.QuizAttempt, error) {
	rows, err := r.db.Query(`
		SELECT qa.id, u.email, qa.category_id, qc.name, qa.correct_count, qa.total_count,
		       qa.score, qa.start_time, qa.end_time, qa.finish_date, qa.created_at
		FROM quiz_attempts qa
		JOIN users u ON u.id = qa.user_id
		JOIN quiz_categories qc ON qc.id = qa.category_id
		WHERE u.email = ?
		ORDER BY qa.end_time DESC`, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []models.QuizAttempt
	for rows.Next() {
		var attempt models.QuizAttempt
		if err := rows.Scan(
			&attempt.ID, &attempt.Email, &attempt.CategoryID, &attempt.CategoryName,
			&attempt.CorrectCount, &attempt.TotalCount, &attempt.Score,
			&attempt.StartTime, &attempt.EndTime, &attempt.FinishDate, &attempt.CreatedAt,
		); err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	return attempts, rows.Err()
}

func (r *QuizRepository) GetAttemptDetails(attemptID int) ([]models.AttemptDetail, error) {
	rows, err := r.db.Query(`
		SELECT id, quiz_attempt_id, quiz_id, type, user_answer, created_at
		FROM quiz_attempt_details WHERE quiz_attempt_id = ?`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []models.AttemptDetail
	for rows.Next() {
		var detail models.AttemptDetail
		if err := rows.Scan(
			&detail.ID, &detail.QuizAttemptID, &detail.QuizID,
			&detail.Type, &detail.UserAnswer, &detail.CreatedAt,
		); err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, rows.Err()
}

func (r *QuizRepository) CountQuestions() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM quizzes`).Scan(&count)
	return count, err
}

func (r *QuizRepository) CountAttempts() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM quiz_attempts`).Scan(&count)
	return count, err
}

func (r *QuizRepository) IsCorrectAnswer(quizID int, userAnswer string) (bool, error) {
	var isCorrect bool
	err := r.db.QueryRow(`
		SELECT is_correct FROM quiz_answers
		WHERE quiz_id = ? AND (answer_key = ? OR answer_text = ?)
		LIMIT 1`, quizID, userAnswer, userAnswer).Scan(&isCorrect)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return isCorrect, err
}

func ParseFinishDate(value string) (time.Time, error) {
	if value == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", value)
}

// --- Quiz category CRUD ---

func (r *QuizRepository) GetAllCategories() ([]models.QuizCategory, error) {
	rows, err := r.db.Query(`
		SELECT id, name, slug, description, display_order, is_active, created_at
		FROM quiz_categories ORDER BY display_order ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.QuizCategory
	for rows.Next() {
		var category models.QuizCategory
		var isActive bool
		if err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &category.Description,
			&category.DisplayOrder, &isActive, &category.CreatedAt,
		); err != nil {
			return nil, err
		}
		category.IsActive = isActive
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *QuizRepository) CreateCategory(req models.QuizCategoryRequest) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO quiz_categories (name, slug, description, display_order, is_active)
		VALUES (?, ?, ?, ?, ?)`,
		req.Name, req.Slug, req.Description, req.DisplayOrder, req.IsActive,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *QuizRepository) UpdateCategory(id int, req models.QuizCategoryRequest) error {
	_, err := r.db.Exec(`
		UPDATE quiz_categories SET name = ?, slug = ?, description = ?, display_order = ?, is_active = ?
		WHERE id = ?`,
		req.Name, req.Slug, req.Description, req.DisplayOrder, req.IsActive, id,
	)
	return err
}

func (r *QuizRepository) DeleteCategory(id int) error {
	_, err := r.db.Exec(`DELETE FROM quiz_categories WHERE id = ?`, id)
	return err
}

// --- Quiz question CRUD ---

func (r *QuizRepository) CreateQuestion(req models.QuizQuestionRequest) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO quizzes (category_id, type, question, image_url, sequence_order)
		VALUES (?, ?, ?, ?, ?)`,
		req.CategoryID, req.Type, req.Question, req.ImageURL, req.SequenceOrder,
	)
	if err != nil {
		return 0, err
	}
	quizID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := r.replaceAnswers(int(quizID), req.Answers); err != nil {
		return 0, err
	}
	return quizID, nil
}

func (r *QuizRepository) UpdateQuestion(id int, req models.QuizQuestionRequest) error {
	_, err := r.db.Exec(`
		UPDATE quizzes SET category_id = ?, type = ?, question = ?, image_url = ?, sequence_order = ?
		WHERE id = ?`,
		req.CategoryID, req.Type, req.Question, req.ImageURL, req.SequenceOrder, id,
	)
	if err != nil {
		return err
	}
	return r.replaceAnswers(id, req.Answers)
}

func (r *QuizRepository) replaceAnswers(quizID int, answers []models.QuizAnswerRequest) error {
	if _, err := r.db.Exec(`DELETE FROM quiz_answers WHERE quiz_id = ?`, quizID); err != nil {
		return err
	}
	for _, a := range answers {
		var key *string
		if a.AnswerKey != "" {
			k := a.AnswerKey
			key = &k
		}
		if _, err := r.db.Exec(`
			INSERT INTO quiz_answers (quiz_id, answer_key, answer_text, is_correct)
			VALUES (?, ?, ?, ?)`, quizID, key, a.AnswerText, a.IsCorrect); err != nil {
			return err
		}
	}
	return nil
}

func (r *QuizRepository) DeleteQuestion(id int) error {
	_, err := r.db.Exec(`DELETE FROM quizzes WHERE id = ?`, id)
	return err
}

// --- Analytics helpers ---

func (r *QuizRepository) CountAttemptsSince(t time.Time) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM quiz_attempts WHERE created_at >= ?`, t).Scan(&count)
	return count, err
}

func (r *QuizRepository) AverageScore() (float64, error) {
	var avg sql.NullFloat64
	err := r.db.QueryRow(`SELECT AVG(score) FROM quiz_attempts`).Scan(&avg)
	if err != nil {
		return 0, err
	}
	if avg.Valid {
		return avg.Float64, nil
	}
	return 0, nil
}

// AttemptsByCategory returns attempt counts grouped by category name.
func (r *QuizRepository) AttemptsByCategory() ([]models.CategoryCount, error) {
	rows, err := r.db.Query(`
		SELECT qc.name, COUNT(qa.id) AS total
		FROM quiz_categories qc
		LEFT JOIN quiz_attempts qa ON qa.category_id = qc.id
		GROUP BY qc.id, qc.name
		ORDER BY total DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.CategoryCount
	for rows.Next() {
		var item models.CategoryCount
		if err := rows.Scan(&item.Label, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
