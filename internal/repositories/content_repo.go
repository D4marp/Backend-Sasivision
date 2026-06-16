package repositories

import (
	"database/sql"
	"encoding/json"

	"github.com/sasivision/backend/internal/models"
)

type ContentRepository struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) GetMarkers() ([]models.Marker, error) {
	rows, err := r.db.Query(`
		SELECT id, title, slug, description, image_file, audio_file, model_path, sentences, display_order, created_at
		FROM markers ORDER BY display_order ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var markers []models.Marker
	for rows.Next() {
		var marker models.Marker
		var sentences sql.NullString
		var modelPath sql.NullString
		if err := rows.Scan(
			&marker.ID, &marker.Title, &marker.Slug, &marker.Description,
			&marker.ImageFile, &marker.AudioFile, &modelPath, &sentences,
			&marker.DisplayOrder, &marker.CreatedAt,
		); err != nil {
			return nil, err
		}
		if modelPath.Valid {
			marker.ModelPath = &modelPath.String
		}
		if sentences.Valid {
			marker.Sentences = &sentences.String
		}
		markers = append(markers, marker)
	}
	return markers, rows.Err()
}

func (r *ContentRepository) GetMarkerByID(id int) (*models.Marker, error) {
	var marker models.Marker
	var sentences sql.NullString
	var modelPath sql.NullString
	err := r.db.QueryRow(`
		SELECT id, title, slug, description, image_file, audio_file, model_path, sentences, display_order, created_at
		FROM markers WHERE id = ?`, id).Scan(
		&marker.ID, &marker.Title, &marker.Slug, &marker.Description,
		&marker.ImageFile, &marker.AudioFile, &modelPath, &sentences,
		&marker.DisplayOrder, &marker.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if modelPath.Valid {
		marker.ModelPath = &modelPath.String
	}
	if sentences.Valid {
		marker.Sentences = &sentences.String
	}
	return &marker, nil
}

func (r *ContentRepository) GetVideos() ([]models.Video, error) {
	rows, err := r.db.Query(`
		SELECT id, title, slug, description, source, video_url, thumbnail, discussion_form_url, view_count, display_order, created_at
		FROM videos ORDER BY display_order ASC, created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var video models.Video
		var discussion sql.NullString
		if err := rows.Scan(
			&video.ID, &video.Title, &video.Slug, &video.Description, &video.Source,
			&video.VideoURL, &video.Thumbnail, &discussion, &video.ViewCount,
			&video.DisplayOrder, &video.CreatedAt,
		); err != nil {
			return nil, err
		}
		if discussion.Valid {
			video.DiscussionFormURL = &discussion.String
		}
		videos = append(videos, video)
	}
	return videos, rows.Err()
}

func (r *ContentRepository) GetVideoByID(id int) (*models.Video, error) {
	var video models.Video
	var discussion sql.NullString
	err := r.db.QueryRow(`
		SELECT id, title, slug, description, source, video_url, thumbnail, discussion_form_url, view_count, display_order, created_at
		FROM videos WHERE id = ?`, id).Scan(
		&video.ID, &video.Title, &video.Slug, &video.Description, &video.Source,
		&video.VideoURL, &video.Thumbnail, &discussion, &video.ViewCount,
		&video.DisplayOrder, &video.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if discussion.Valid {
		video.DiscussionFormURL = &discussion.String
	}
	return &video, nil
}

func ParseSentences(raw *string) []string {
	if raw == nil || *raw == "" {
		return nil
	}
	var sentences []string
	_ = json.Unmarshal([]byte(*raw), &sentences)
	return sentences
}

// --- Video CRUD ---

func (r *ContentRepository) CreateVideo(v models.VideoRequest) (int64, error) {
	result, err := r.db.Exec(`
		INSERT INTO videos (title, slug, description, source, video_url, thumbnail, discussion_form_url, display_order)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		v.Title, v.Slug, v.Description, v.Source, v.VideoURL, v.Thumbnail, v.DiscussionFormURL, v.DisplayOrder,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ContentRepository) UpdateVideo(id int, v models.VideoRequest) error {
	_, err := r.db.Exec(`
		UPDATE videos SET title = ?, slug = ?, description = ?, source = ?,
		video_url = ?, thumbnail = ?, discussion_form_url = ?, display_order = ?
		WHERE id = ?`,
		v.Title, v.Slug, v.Description, v.Source, v.VideoURL, v.Thumbnail,
		v.DiscussionFormURL, v.DisplayOrder, id,
	)
	return err
}

func (r *ContentRepository) DeleteVideo(id int) error {
	_, err := r.db.Exec(`DELETE FROM videos WHERE id = ?`, id)
	return err
}

func (r *ContentRepository) IncrementViewCount(id int) error {
	_, err := r.db.Exec(`UPDATE videos SET view_count = view_count + 1 WHERE id = ?`, id)
	return err
}

// --- Marker CRUD ---

func (r *ContentRepository) CreateMarker(m models.MarkerRequest) (int64, error) {
	var sentences *string
	if len(m.Sentences) > 0 {
		raw, _ := json.Marshal(m.Sentences)
		s := string(raw)
		sentences = &s
	}
	result, err := r.db.Exec(`
		INSERT INTO markers (title, slug, description, image_file, audio_file, model_path, sentences, display_order)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Title, m.Slug, m.Description, m.ImageFile, m.AudioFile, m.ModelPath, sentences, m.DisplayOrder,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *ContentRepository) UpdateMarker(id int, m models.MarkerRequest) error {
	var sentences *string
	if len(m.Sentences) > 0 {
		raw, _ := json.Marshal(m.Sentences)
		s := string(raw)
		sentences = &s
	}
	_, err := r.db.Exec(`
		UPDATE markers SET title = ?, slug = ?, description = ?, image_file = ?,
		audio_file = ?, model_path = ?, sentences = ?, display_order = ?
		WHERE id = ?`,
		m.Title, m.Slug, m.Description, m.ImageFile, m.AudioFile,
		m.ModelPath, sentences, m.DisplayOrder, id,
	)
	return err
}

func (r *ContentRepository) DeleteMarker(id int) error {
	_, err := r.db.Exec(`DELETE FROM markers WHERE id = ?`, id)
	return err
}
