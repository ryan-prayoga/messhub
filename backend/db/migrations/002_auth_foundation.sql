ALTER TABLE users
  ALTER COLUMN joined_at SET DEFAULT NOW();

UPDATE users
SET
  joined_at = COALESCE(joined_at, created_at, NOW()),
  updated_at = COALESCE(updated_at, NOW());

ALTER TABLE users
  ALTER COLUMN joined_at SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'users_left_at_after_joined_at'
  ) THEN
    ALTER TABLE users
      ADD CONSTRAINT users_left_at_after_joined_at
      CHECK (left_at IS NULL OR left_at >= joined_at);
  END IF;
END $$;

CREATE OR REPLACE FUNCTION set_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_users_set_updated_at ON users;

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_users_updated_at();
