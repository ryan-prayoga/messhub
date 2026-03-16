ALTER TABLE wallet_transactions
  DROP COLUMN IF EXISTS transaction_date,
  DROP COLUMN IF EXISTS proof_url,
  DROP COLUMN IF EXISTS updated_by;

ALTER TABLE wallet_transactions
  ALTER COLUMN amount TYPE INTEGER USING amount::integer;

ALTER TABLE wallet_transactions
  ALTER COLUMN type SET NOT NULL,
  ALTER COLUMN category SET NOT NULL,
  ALTER COLUMN amount SET NOT NULL,
  ALTER COLUMN description SET NOT NULL,
  ALTER COLUMN created_by SET NOT NULL,
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

DROP INDEX IF EXISTS idx_wallet_transactions_date;

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_created_at
  ON wallet_transactions(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_wallet_transactions_created_by
  ON wallet_transactions(created_by);
