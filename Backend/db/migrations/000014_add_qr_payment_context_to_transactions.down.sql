-- Rollback QR payment context and tip processing fields

-- Drop indexes first
DROP INDEX IF EXISTS idx_transactions_qr_tag;
DROP INDEX IF EXISTS idx_transactions_source;
DROP INDEX IF EXISTS idx_transactions_tip_processing;
DROP INDEX IF EXISTS idx_transactions_hosted_checkout_id;
DROP INDEX IF EXISTS idx_transactions_qr_link_id;

-- Drop foreign key constraints
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_hosted_checkout;
ALTER TABLE transactions DROP CONSTRAINT IF EXISTS fk_transactions_qr_link;

-- Drop columns
ALTER TABLE transactions DROP COLUMN IF EXISTS tip_processed;
ALTER TABLE transactions DROP COLUMN IF EXISTS tip_transaction_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS tip_medium;
ALTER TABLE transactions DROP COLUMN IF EXISTS tipee_phone;
ALTER TABLE transactions DROP COLUMN IF EXISTS tip_amount;
ALTER TABLE transactions DROP COLUMN IF EXISTS has_tip;
ALTER TABLE transactions DROP COLUMN IF EXISTS qr_tag;
ALTER TABLE transactions DROP COLUMN IF EXISTS hosted_checkout_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS qr_link_id;
ALTER TABLE transactions DROP COLUMN IF EXISTS transaction_source;

-- Drop enum type
DROP TYPE IF EXISTS transaction_source; 