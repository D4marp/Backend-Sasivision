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

INSERT INTO quiz_categories (name, slug, description, display_order, is_active) VALUES
('Post-Test', 'post-test', 'Test your knowledge about Sasirangan culture', 1, 1),
('Basics', 'basics', 'Introductory questions about Sasirangan', 2, 0)
ON DUPLICATE KEY UPDATE name = name;

INSERT INTO quizzes (category_id, type, question, image_url, sequence_order) VALUES
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice', 'What region is Sasirangan batik originally from?', NULL, 1),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice', 'What technique is used to create Sasirangan patterns?', NULL, 2),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'multiple_choice', 'The Bintang Bahambur motif symbolizes:', NULL, 3),
((SELECT id FROM quiz_categories WHERE slug = 'post-test'), 'essay', 'Explain the cultural significance of Sasirangan batik in Banjar society.', NULL, 4);

INSERT INTO quiz_answers (quiz_id, answer_key, answer_text, is_correct) VALUES
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'A', 'South Kalimantan', 1),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'B', 'Central Java', 0),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'C', 'Bali', 0),
((SELECT id FROM quizzes WHERE sequence_order = 1 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'D', 'West Sumatra', 0),

((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'A', 'Batik tulis', 0),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'B', 'Needle-resist dyeing', 1),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'C', 'Ikat weaving', 0),
((SELECT id FROM quizzes WHERE sequence_order = 2 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'D', 'Block printing', 0),

((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'A', 'War and conflict', 0),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'B', 'Stars and harmony in the universe', 1),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'C', 'Agricultural harvest', 0),
((SELECT id FROM quizzes WHERE sequence_order = 3 AND category_id = (SELECT id FROM quiz_categories WHERE slug = 'post-test')), 'D', 'Royal authority', 0);
