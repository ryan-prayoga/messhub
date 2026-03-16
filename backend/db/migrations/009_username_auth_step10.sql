ALTER TABLE users
  ADD COLUMN IF NOT EXISTS username TEXT;

WITH source_rows AS (
  SELECT
    id,
    created_at,
    COALESCE(
      NULLIF(BTRIM(LOWER(REGEXP_REPLACE(name, '[^a-zA-Z0-9]+', '-', 'g')), '-'), ''),
      NULLIF(BTRIM(LOWER(REGEXP_REPLACE(SPLIT_PART(email, '@', 1), '[^a-zA-Z0-9]+', '-', 'g')), '-'), ''),
      'user'
    ) AS base_username
  FROM users
  WHERE username IS NULL OR BTRIM(username) = ''
),
ranked_rows AS (
  SELECT
    id,
    base_username,
    ROW_NUMBER() OVER (PARTITION BY base_username ORDER BY created_at, id) AS occurrence
  FROM source_rows
),
resolved_rows AS (
  SELECT
    id,
    CASE
      WHEN occurrence = 1 THEN LEFT(base_username, 32)
      ELSE LEFT(base_username, GREATEST(1, 31 - LENGTH(occurrence::TEXT))) || '-' || occurrence::TEXT
    END AS resolved_username
  FROM ranked_rows
)
UPDATE users AS target
SET username = resolved_rows.resolved_username
FROM resolved_rows
WHERE target.id = resolved_rows.id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower
  ON users (LOWER(username));

ALTER TABLE users
  ALTER COLUMN username SET NOT NULL;
