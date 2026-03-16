ALTER TABLE users
  ADD COLUMN IF NOT EXISTS phone TEXT,
  ADD COLUMN IF NOT EXISTS avatar_url TEXT;

CREATE TABLE IF NOT EXISTS mess_settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  singleton BOOLEAN NOT NULL DEFAULT TRUE,
  mess_name TEXT NOT NULL,
  wifi_price INTEGER NOT NULL CHECK (wifi_price > 0),
  wifi_deadline_day SMALLINT NOT NULL CHECK (wifi_deadline_day BETWEEN 1 AND 31),
  bank_account_name TEXT NOT NULL,
  bank_account_number TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (singleton)
);

UPDATE mess_settings
SET updated_at = COALESCE(updated_at, created_at, NOW());

INSERT INTO mess_settings (
  mess_name,
  wifi_price,
  wifi_deadline_day,
  bank_account_name,
  bank_account_number
)
SELECT
  'MessHub',
  20000,
  10,
  'Ryan Prayoga',
  '104987106615'
WHERE NOT EXISTS (
  SELECT 1
  FROM mess_settings
);

CREATE OR REPLACE FUNCTION set_mess_settings_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_mess_settings_set_updated_at ON mess_settings;

CREATE TRIGGER trg_mess_settings_set_updated_at
BEFORE UPDATE ON mess_settings
FOR EACH ROW
EXECUTE FUNCTION set_mess_settings_updated_at();

CREATE INDEX IF NOT EXISTS idx_mess_settings_updated_at
  ON mess_settings(updated_at DESC);
