-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create user_sessions table for auto-login
CREATE TABLE IF NOT EXISTS user_sessions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    token VARCHAR(512) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id_expires_at (user_id, expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create quiz_categories table
CREATE TABLE IF NOT EXISTS quiz_categories (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    display_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_is_active_display_order (is_active, display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create quizzes table
CREATE TABLE IF NOT EXISTS quizzes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    category_id INT NOT NULL,
    type ENUM('multiple_choice', 'essay') NOT NULL,
    question TEXT NOT NULL,
    image_url VARCHAR(255),
    sequence_order INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES quiz_categories(id) ON DELETE CASCADE,
    INDEX idx_category_id_sequence_order (category_id, sequence_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create quiz_answers table
CREATE TABLE IF NOT EXISTS quiz_answers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    quiz_id INT NOT NULL,
    answer_key VARCHAR(1),
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
    INDEX idx_quiz_id_is_correct (quiz_id, is_correct)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create quiz_attempts table
CREATE TABLE IF NOT EXISTS quiz_attempts (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    category_id INT NOT NULL,
    correct_count INT DEFAULT 0,
    total_count INT NOT NULL,
    score INT DEFAULT 0,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    finish_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES quiz_categories(id),
    INDEX idx_user_id_finish_date (user_id, finish_date),
    INDEX idx_user_id_category_id (user_id, category_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create quiz_attempt_details table
CREATE TABLE IF NOT EXISTS quiz_attempt_details (
    id INT PRIMARY KEY AUTO_INCREMENT,
    quiz_attempt_id INT NOT NULL,
    quiz_id INT NOT NULL,
    type ENUM('multiple_choice', 'essay') NOT NULL,
    user_answer TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (quiz_attempt_id) REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
    INDEX idx_quiz_attempt_id (quiz_attempt_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create markers table
CREATE TABLE IF NOT EXISTS markers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    image_file VARCHAR(255) NOT NULL,
    audio_file VARCHAR(255) NOT NULL,
    model_path VARCHAR(255),
    sentences JSON,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_display_order (display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create videos table
CREATE TABLE IF NOT EXISTS videos (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL UNIQUE,
    slug VARCHAR(255),
    description TEXT,
    source VARCHAR(100),
    video_url VARCHAR(255) NOT NULL,
    thumbnail VARCHAR(255) NOT NULL,
    discussion_form_url VARCHAR(255),
    view_count INT DEFAULT 0,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_display_order_created_at (display_order, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create feature_switches table
CREATE TABLE IF NOT EXISTS feature_switches (
    id INT PRIMARY KEY AUTO_INCREMENT,
    feature_name VARCHAR(100) NOT NULL UNIQUE,
    status ENUM('active', 'inactive') DEFAULT 'inactive',
    description TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_feature_name (feature_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create feature_logs table
CREATE TABLE IF NOT EXISTS feature_logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    feature_switch_id INT NOT NULL,
    action ENUM('activated', 'deactivated') NOT NULL,
    changed_by_user_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (feature_switch_id) REFERENCES feature_switches(id) ON DELETE CASCADE,
    INDEX idx_feature_switch_id_created_at (feature_switch_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- Seed data for SasiVision development
-- Password for demo user: Sasivision123 (bcrypt hash below)

INSERT INTO users (email, password_hash, full_name) VALUES
('demo@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Demo User')
ON DUPLICATE KEY UPDATE email = email;

INSERT INTO feature_switches (feature_name, status, description) VALUES
('AR Sasirangan', 'active', '3D AR visualization of Sasirangan motifs'),
('Quizzes', 'active', 'Interactive quiz questions'),
('Vocabulary Quiz', 'inactive', 'Vocabulary learning mini-quiz')
ON DUPLICATE KEY UPDATE feature_name = feature_name;

INSERT INTO markers (title, slug, description, image_file, audio_file, model_path, sentences, display_order) VALUES
(
  'Bintang Bahambur',
  'bintang-bahambur',
  'The Bintang Bahambur motif in Sasirangan fabric is a visual representation of stars scattered across the sky, symbolizing how the universe is full of small elements that radiate beauty in harmony.',
  'markers/bintang_bahambur.png',
  'audio/bintang_bahambur.mp3',
  'models/stars.glb',
  '["The Bintang Bahambur motif is a visual representation of stars scattered across the sky.", "The philosophy teaches that human life should be peaceful and simple.", "It symbolizes that the universe is full of small elements radiating beauty in harmony."]',
  1
),
(
  'Naga Balimbur',
  'naga-balimbur',
  'The Naga Balimbur motif depicts a dragon in Banjar mythology. The word Balimbur comes from limbur, meaning to wash or bathe.',
  'markers/naga_balimbur.png',
  'audio/naga_balimbur.mp3',
  'models/dragon_draco_2.glb',
  '["The Naga Balimbur motif depicts a dragon in Banjar mythology.", "The word Balimbur comes from limbur, meaning to wash or bathe."]',
  2
),
(
  'Kulat Karikit',
  'kulat-karikit',
  'The Kulat Karikit motif imitates the shape of mushrooms. It also shows the relationship between humans and the environment.',
  'markers/kulat_karikit.png',
  'audio/kulat_karikit.mp3',
  'models/mushroom_clump.glb',
  '["The Kulat Karikit motif imitates the shape of mushrooms.", "It also shows the relationship between humans and the environment."]',
  3
)
ON DUPLICATE KEY UPDATE title = title;

INSERT INTO videos (title, slug, description, source, video_url, thumbnail, discussion_form_url, view_count, display_order) VALUES
(
  'History of Sasirangan',
  'history-of-sasirangan',
  'Learn about the origins and cultural significance of Sasirangan batik from South Kalimantan.',
  'SasiVision Team',
  'videos/history_sasirangan.mp4',
  'videos/thumbnails/history.jpg',
  'https://forms.gle/example-history',
  1250,
  1
),
(
  'The Art of Sasirangan Dyeing',
  'art-of-sasirangan-dyeing',
  'Discover the traditional needle-resist dyeing technique used to create Sasirangan patterns.',
  'SasiVision Team',
  'videos/art_of_dyeing.mp4',
  'videos/thumbnails/dyeing.jpg',
  'https://forms.gle/example-dyeing',
  890,
  2
),
(
  'Motif Meanings in Sasirangan',
  'motif-meanings',
  'Explore the philosophy and symbolism behind popular Sasirangan motifs.',
  'SasiVision Team',
  'videos/motif_meanings.mp4',
  'videos/thumbnails/motif.jpg',
  NULL,
  654,
  3
)
ON DUPLICATE KEY UPDATE title = title;

-- Quiz categories & questions: see 005_quiz_seed_sasirangan.sql (idempotent)
-- Add role-based access control to users (idempotent)
-- Roles: user (default), editor (manage content), admin (full access)

SET @role_col_exists = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND COLUMN_NAME = 'role'
);

SET @add_role_sql = IF(
  @role_col_exists = 0,
  'ALTER TABLE users ADD COLUMN role ENUM(''user'', ''editor'', ''admin'') NOT NULL DEFAULT ''user'' AFTER full_name',
  'SELECT 1'
);
PREPARE add_role_stmt FROM @add_role_sql;
EXECUTE add_role_stmt;
DEALLOCATE PREPARE add_role_stmt;

SET @role_idx_exists = (
  SELECT COUNT(*)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'users'
    AND INDEX_NAME = 'idx_users_role'
);

SET @add_idx_sql = IF(
  @role_idx_exists = 0,
  'CREATE INDEX idx_users_role ON users (role)',
  'SELECT 1'
);
PREPARE add_idx_stmt FROM @add_idx_sql;
EXECUTE add_idx_stmt;
DEALLOCATE PREPARE add_idx_stmt;

UPDATE users SET role = 'user' WHERE email = 'demo@sasivision.com';

INSERT INTO users (email, password_hash, full_name, role) VALUES
('admin@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Administrator', 'admin'),
('editor@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Content Editor', 'editor')
ON DUPLICATE KEY UPDATE role = VALUES(role), full_name = VALUES(full_name);
-- Analytics event stream for the super-app dashboard
CREATE TABLE IF NOT EXISTS analytics_events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    event_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50),
    entity_id INT,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_event_type_created (event_type, created_at),
    INDEX idx_entity (entity_type, entity_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- Idempotent quiz seed for Post-Test (can be re-run safely)
-- Password for all demo accounts: Sasivision123

INSERT INTO users (email, password_hash, full_name) VALUES
('demo@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Demo Mahasiswa')
ON DUPLICATE KEY UPDATE full_name = VALUES(full_name);

INSERT INTO users (email, password_hash, full_name, role) VALUES
('admin@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Administrator', 'admin'),
('editor@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Content Editor', 'editor')
ON DUPLICATE KEY UPDATE role = VALUES(role), full_name = VALUES(full_name);

UPDATE users SET role = 'user' WHERE email = 'demo@sasivision.com' AND role = 'admin';

INSERT INTO feature_switches (feature_name, status, description) VALUES
('AR Sasirangan', 'active', 'Visualisasi AR motif Sasirangan 3D'),
('Quizzes', 'active', 'Soal quiz interaktif'),
('Vocabulary Quiz', 'inactive', 'Quiz kosakata — segera hadir')
ON DUPLICATE KEY UPDATE status = VALUES(status), description = VALUES(description);

INSERT INTO quiz_categories (name, slug, description, display_order, is_active) VALUES
('Post-Test', 'post-test', 'Evaluasi pemahaman motif dan budaya Sasirangan', 1, 1),
('Basics', 'basics', 'Pertanyaan pengantar Sasirangan', 2, 0)
ON DUPLICATE KEY UPDATE
  description = VALUES(description),
  display_order = VALUES(display_order),
  is_active = VALUES(is_active);

DELETE qa FROM quiz_answers qa
INNER JOIN quizzes q ON qa.quiz_id = q.id
INNER JOIN quiz_categories qc ON q.category_id = qc.id
WHERE qc.slug = 'post-test';

DELETE q FROM quizzes q
INNER JOIN quiz_categories qc ON q.category_id = qc.id
WHERE qc.slug = 'post-test';

INSERT INTO quizzes (category_id, type, question, image_url, sequence_order) VALUES
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice',
 'Sasirangan berasal dari daerah mana?', NULL, 1),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice',
 'Teknik apa yang digunakan untuk membuat motif Sasirangan?', NULL, 2),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice',
 'Motif Bintang Bahambur melambangkan:', NULL, 3),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice',
 'Kata "Balimbur" pada Naga Balimbur berasal dari kata Banjar yang berarti:', NULL, 4),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice',
 'Motif Kulat Karikit menyerupai bentuk:', NULL, 5),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'essay',
 'Jelaskan makna budaya kain Sasirangan dalam masyarakat Banjar.', NULL, 6);

INSERT INTO quiz_answers (quiz_id, answer_key, answer_text, is_correct) VALUES
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'A', 'Kalimantan Selatan', 1),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'B', 'Jawa Tengah', 0),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'C', 'Bali', 0),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'D', 'Sumatera Barat', 0),

((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'A', 'Batik tulis', 0),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'B', 'Canting / jarum resist', 1),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'C', 'Tenun ikat', 0),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'D', 'Cap batik', 0),

((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'A', 'Perang dan konflik', 0),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'B', 'Bintang dan harmoni alam semesta', 1),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'C', 'Panen pertanian', 0),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'D', 'Kekuasaan kerajaan', 0),

((SELECT id FROM quizzes WHERE sequence_order = 4 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'A', 'Terbang', 0),
((SELECT id FROM quizzes WHERE sequence_order = 4 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'B', 'Mandi / mencuci', 1),
((SELECT id FROM quizzes WHERE sequence_order = 4 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'C', 'Berkembang biak', 0),
((SELECT id FROM quizzes WHERE sequence_order = 4 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'D', 'Bertempur', 0),

((SELECT id FROM quizzes WHERE sequence_order = 5 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'A', 'Bunga', 0),
((SELECT id FROM quizzes WHERE sequence_order = 5 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'B', 'Jamur', 1),
((SELECT id FROM quizzes WHERE sequence_order = 5 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'C', 'Ikan', 0),
((SELECT id FROM quizzes WHERE sequence_order = 5 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')),
 'D', 'Burung', 0);
-- JWT tokens exceed 255 characters; widen session storage (idempotent)
SET @token_len = (
  SELECT CHARACTER_MAXIMUM_LENGTH
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'user_sessions'
    AND COLUMN_NAME = 'token'
);

SET @widen_sql = IF(
  @token_len IS NULL OR @token_len < 512,
  'ALTER TABLE user_sessions MODIFY COLUMN token VARCHAR(512) NOT NULL',
  'SELECT 1'
);
PREPARE widen_stmt FROM @widen_sql;
EXECUTE widen_stmt;
DEALLOCATE PREPARE widen_stmt;
