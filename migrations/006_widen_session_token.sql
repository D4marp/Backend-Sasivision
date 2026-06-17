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
