ALTER TABLE users
  ADD COLUMN IF NOT EXISTS auth_version INTEGER NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS archived_at TIMESTAMPTZ;

UPDATE users
SET
  auth_version = COALESCE(auth_version, 1),
  is_active = CASE
    WHEN archived_at IS NOT NULL THEN FALSE
    ELSE is_active
  END,
  left_at = CASE
    WHEN archived_at IS NOT NULL THEN COALESCE(left_at, archived_at, NOW())
    ELSE left_at
  END;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'users_archived_requires_inactive'
  ) THEN
    ALTER TABLE users
      ADD CONSTRAINT users_archived_requires_inactive
      CHECK (archived_at IS NULL OR is_active = FALSE);
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_archived_at
  ON users(archived_at DESC NULLS LAST);

CREATE INDEX IF NOT EXISTS idx_users_active_archived
  ON users(is_active, archived_at, role);

ALTER TABLE shared_expenses
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE shared_expenses
SET updated_at = COALESCE(updated_at, created_at, NOW());

CREATE OR REPLACE FUNCTION set_shared_expenses_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_shared_expenses_set_updated_at ON shared_expenses;

CREATE TRIGGER trg_shared_expenses_set_updated_at
BEFORE UPDATE ON shared_expenses
FOR EACH ROW
EXECUTE FUNCTION set_shared_expenses_updated_at();

CREATE OR REPLACE FUNCTION set_proposals_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_proposals_set_updated_at ON proposals;

CREATE TRIGGER trg_proposals_set_updated_at
BEFORE UPDATE ON proposals
FOR EACH ROW
EXECUTE FUNCTION set_proposals_updated_at();

CREATE INDEX IF NOT EXISTS idx_shared_expenses_status_date
  ON shared_expenses(status, expense_date DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_shared_expenses_paid_by
  ON shared_expenses(paid_by_user_id, expense_date DESC);

CREATE INDEX IF NOT EXISTS idx_proposals_status_updated_at
  ON proposals(status, updated_at DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_proposal_votes_proposal_id
  ON proposal_votes(proposal_id, created_at ASC);
