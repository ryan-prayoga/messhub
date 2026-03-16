CREATE TABLE IF NOT EXISTS audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id),
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id UUID,
  old_value JSONB,
  new_value JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE wifi_bills
  ALTER COLUMN nominal_per_person TYPE INTEGER USING nominal_per_person::integer,
  ALTER COLUMN nominal_per_person SET NOT NULL,
  ALTER COLUMN deadline_date SET NOT NULL,
  ALTER COLUMN status SET DEFAULT 'active';

ALTER TABLE wifi_bills
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE wifi_bills
SET updated_at = COALESCE(updated_at, created_at, NOW());

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'wifi_bills_status_check'
  ) THEN
    ALTER TABLE wifi_bills
      ADD CONSTRAINT wifi_bills_status_check
      CHECK (status IN ('draft', 'active', 'closed'));
  END IF;
END $$;

ALTER TABLE wifi_bill_members
  ALTER COLUMN amount TYPE INTEGER USING amount::integer;

ALTER TABLE wifi_bill_members
  ADD COLUMN IF NOT EXISTS note TEXT,
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE wifi_bill_members
SET
  created_at = COALESCE(created_at, submitted_at, verified_at, NOW()),
  updated_at = COALESCE(updated_at, submitted_at, verified_at, NOW());

CREATE INDEX IF NOT EXISTS idx_wifi_bills_year_month
  ON wifi_bills(year DESC, month DESC);

CREATE INDEX IF NOT EXISTS idx_wifi_bills_status
  ON wifi_bills(status);

CREATE INDEX IF NOT EXISTS idx_wifi_bill_members_bill_status
  ON wifi_bill_members(wifi_bill_id, payment_status);

CREATE INDEX IF NOT EXISTS idx_wifi_bill_members_user_id
  ON wifi_bill_members(user_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at
  ON audit_logs(created_at DESC);
