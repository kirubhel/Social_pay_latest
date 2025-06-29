DROP INDEX IF EXISTS merchants.idx_audit_logs_created_at;
DROP INDEX IF EXISTS merchants.idx_audit_logs_merchant_id;
DROP TABLE IF EXISTS merchants.audit_logs;

DROP TABLE IF EXISTS merchants.settings;

DROP INDEX IF EXISTS merchants.idx_api_keys_merchant_id;
DROP TABLE IF EXISTS merchants.api_keys;

DROP INDEX IF EXISTS merchants.idx_bank_accounts_merchant_id;
DROP TABLE IF EXISTS merchants.bank_accounts;

DROP INDEX IF EXISTS merchants.idx_documents_merchant_id;
DROP TABLE IF EXISTS merchants.documents;

DROP INDEX IF EXISTS merchants.idx_contacts_email_unique;
DROP INDEX IF EXISTS merchants.idx_contacts_merchant_id;
DROP TABLE IF EXISTS merchants.contacts;

DROP INDEX IF EXISTS merchants.idx_addresses_merchant_id;
DROP TABLE IF EXISTS merchants.addresses;

DROP TABLE IF EXISTS merchants.merchants;
