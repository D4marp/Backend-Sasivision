package repositories

import (
	"database/sql"
	"time"

	"github.com/sasivision/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, string, error) {
	var user models.User
	var passwordHash string
	err := r.db.QueryRow(
		`SELECT id, email, full_name, role, password_hash, created_at, updated_at
		 FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &passwordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

func (r *UserRepository) Create(email, passwordHash, fullName string) (*models.User, error) {
	result, err := r.db.Exec(
		`INSERT INTO users (email, password_hash, full_name) VALUES (?, ?, ?)`,
		email, passwordHash, fullName,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.FindByID(int(id))
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, email, full_name, role, created_at, updated_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *UserRepository) List() ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, full_name, role, created_at, updated_at
		 FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.FullName, &user.Role,
			&user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) UpdateRole(id int, role string) error {
	_, err := r.db.Exec(`UPDATE users SET role = ? WHERE id = ?`, role, id)
	return err
}

func (r *UserRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (r *UserRepository) CountByRole() (map[string]int, error) {
	rows, err := r.db.Query(`SELECT role, COUNT(*) FROM users GROUP BY role`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var role string
		var count int
		if err := rows.Scan(&role, &count); err != nil {
			return nil, err
		}
		result[role] = count
	}
	return result, rows.Err()
}

func (r *UserRepository) CountCreatedSince(t time.Time) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE created_at >= ?`, t).Scan(&count)
	return count, err
}

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(userID int, token string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO user_sessions (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *SessionRepository) FindValid(token string) (*models.UserSession, error) {
	var session models.UserSession
	err := r.db.QueryRow(
		`SELECT id, user_id, token, created_at, expires_at
		 FROM user_sessions WHERE token = ? AND expires_at > NOW()`,
		token,
	).Scan(&session.ID, &session.UserID, &session.Token, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM user_sessions WHERE token = ?`, token)
	return err
}
