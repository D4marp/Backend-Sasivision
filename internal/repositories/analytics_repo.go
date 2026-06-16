package repositories

import (
	"database/sql"
	"encoding/json"

	"github.com/sasivision/backend/internal/models"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) RecordEvent(userID *int, eventType, entityType string, entityID int, metadata map[string]interface{}) error {
	var meta *string
	if len(metadata) > 0 {
		raw, _ := json.Marshal(metadata)
		s := string(raw)
		meta = &s
	}
	var entType *string
	if entityType != "" {
		entType = &entityType
	}
	var entID *int
	if entityID > 0 {
		entID = &entityID
	}
	_, err := r.db.Exec(`
		INSERT INTO analytics_events (user_id, event_type, entity_type, entity_id, metadata)
		VALUES (?, ?, ?, ?, ?)`,
		userID, eventType, entType, entID, meta,
	)
	return err
}

// CountByEventType returns total events grouped by event_type.
func (r *AnalyticsRepository) CountByEventType() ([]models.CategoryCount, error) {
	rows, err := r.db.Query(`
		SELECT event_type, COUNT(*) AS total
		FROM analytics_events
		GROUP BY event_type
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

// TopEntities returns the most frequent entity_ids for a given event type,
// joined to a label from the given table/column.
func (r *AnalyticsRepository) TopEntities(eventType, table, labelColumn string, limit int) ([]models.CategoryCount, error) {
	query := `
		SELECT COALESCE(t.` + labelColumn + `, CONCAT('#', ae.entity_id)) AS label, COUNT(*) AS total
		FROM analytics_events ae
		LEFT JOIN ` + table + ` t ON t.id = ae.entity_id
		WHERE ae.event_type = ?
		GROUP BY ae.entity_id, label
		ORDER BY total DESC
		LIMIT ?`
	rows, err := r.db.Query(query, eventType, limit)
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

// EventsDaily returns a daily count of events for the last n days.
func (r *AnalyticsRepository) EventsDaily(days int) ([]models.TimePoint, error) {
	rows, err := r.db.Query(`
		SELECT DATE(created_at) AS day, COUNT(*) AS total
		FROM analytics_events
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY day
		ORDER BY day ASC`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TimePoint
	for rows.Next() {
		var item models.TimePoint
		if err := rows.Scan(&item.Date, &item.Count); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}
