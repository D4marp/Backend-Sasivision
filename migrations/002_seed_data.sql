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
