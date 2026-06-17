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
