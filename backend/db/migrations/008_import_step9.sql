CREATE TABLE IF NOT EXISTS import_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  import_type TEXT NOT NULL CHECK (import_type IN ('members', 'wallet')),
  status TEXT NOT NULL CHECK (status IN ('previewed', 'committed', 'failed')),
  source TEXT NOT NULL DEFAULT 'spreadsheet_import',
  file_name TEXT NOT NULL,
  file_hash TEXT NOT NULL,
  created_by UUID NOT NULL REFERENCES users(id),
  duplicate_strategy TEXT CHECK (duplicate_strategy IN ('skip', 'fail')),
  total_rows INTEGER NOT NULL DEFAULT 0 CHECK (total_rows >= 0),
  valid_rows INTEGER NOT NULL DEFAULT 0 CHECK (valid_rows >= 0),
  invalid_rows INTEGER NOT NULL DEFAULT 0 CHECK (invalid_rows >= 0),
  committed_rows INTEGER NOT NULL DEFAULT 0 CHECK (committed_rows >= 0),
  skipped_rows INTEGER NOT NULL DEFAULT 0 CHECK (skipped_rows >= 0),
  failed_rows INTEGER NOT NULL DEFAULT 0 CHECK (failed_rows >= 0),
  preview_payload JSONB NOT NULL DEFAULT '{}'::JSONB,
  commit_summary JSONB,
  committed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_import_jobs_created_by
  ON import_jobs(created_by, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_import_jobs_type_status
  ON import_jobs(import_type, status, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS ux_import_jobs_committed_hash
  ON import_jobs(file_hash)
  WHERE status = 'committed'
    AND import_type = 'wallet';

ALTER TABLE wallet_transactions
  ADD COLUMN IF NOT EXISTS transaction_date DATE,
  ADD COLUMN IF NOT EXISTS proof_url TEXT,
  ADD COLUMN IF NOT EXISTS source TEXT NOT NULL DEFAULT 'manual',
  ADD COLUMN IF NOT EXISTS import_job_id UUID REFERENCES import_jobs(id);

UPDATE wallet_transactions
SET
  transaction_date = COALESCE(transaction_date, DATE(created_at)),
  source = COALESCE(NULLIF(source, ''), 'manual')
WHERE transaction_date IS NULL
   OR source IS NULL
   OR source = '';

ALTER TABLE wallet_transactions
  ALTER COLUMN transaction_date SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'wallet_transactions_source_check'
  ) THEN
    ALTER TABLE wallet_transactions
      ADD CONSTRAINT wallet_transactions_source_check
      CHECK (source IN ('manual', 'spreadsheet_import'));
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_transaction_date
  ON wallet_transactions(transaction_date DESC, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_import_job_id
  ON wallet_transactions(import_job_id);
