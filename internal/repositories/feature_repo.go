package repositories

import (
	"database/sql"

	"github.com/sasivision/backend/internal/models"
)

type FeatureRepository struct {
	db *sql.DB
}

func NewFeatureRepository(db *sql.DB) *FeatureRepository {
	return &FeatureRepository{db: db}
}

func (r *FeatureRepository) GetByName(featureName string) (*models.FeatureSwitch, error) {
	var feature models.FeatureSwitch
	err := r.db.QueryRow(`
		SELECT id, feature_name, status, description, updated_at
		FROM feature_switches WHERE feature_name = ?`, featureName).Scan(
		&feature.ID, &feature.FeatureName, &feature.Status, &feature.Description, &feature.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &feature, nil
}

func (r *FeatureRepository) GetAll() ([]models.FeatureSwitch, error) {
	rows, err := r.db.Query(`
		SELECT id, feature_name, status, description, updated_at
		FROM feature_switches ORDER BY feature_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []models.FeatureSwitch
	for rows.Next() {
		var feature models.FeatureSwitch
		if err := rows.Scan(
			&feature.ID, &feature.FeatureName, &feature.Status,
			&feature.Description, &feature.UpdatedAt,
		); err != nil {
			return nil, err
		}
		features = append(features, feature)
	}
	return features, rows.Err()
}

func (r *FeatureRepository) SetStatus(featureName, status string) error {
	_, err := r.db.Exec(
		`UPDATE feature_switches SET status = ? WHERE feature_name = ?`,
		status, featureName,
	)
	return err
}

func (r *FeatureRepository) LogChange(featureName, action string, userID *int) error {
	var featureID int
	err := r.db.QueryRow(
		`SELECT id FROM feature_switches WHERE feature_name = ?`, featureName,
	).Scan(&featureID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO feature_logs (feature_switch_id, action, changed_by_user_id)
		VALUES (?, ?, ?)`, featureID, action, userID)
	return err
}
