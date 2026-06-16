-- Add role-based access control to users
-- Roles: user (default), editor (manage content), admin (full access)

ALTER TABLE users
  ADD COLUMN role ENUM('user', 'editor', 'admin') NOT NULL DEFAULT 'user' AFTER full_name;

CREATE INDEX idx_users_role ON users (role);

-- Promote existing demo user to admin for testing
UPDATE users SET role = 'admin' WHERE email = 'demo@sasivision.com';

-- Seed dedicated admin + editor accounts (password: Sasivision123)
INSERT INTO users (email, password_hash, full_name, role) VALUES
('admin@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Administrator', 'admin'),
('editor@sasivision.com', '$2a$10$LQqWB/xuSpIA.CzD/vmYr.j4YWROwhUuLB3ziMkQQ8.6EXtqOdIAe', 'Content Editor', 'editor')
ON DUPLICATE KEY UPDATE role = VALUES(role);
