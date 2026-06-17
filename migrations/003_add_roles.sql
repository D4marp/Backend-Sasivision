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
