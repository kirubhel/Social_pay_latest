ALTER TABLE api_keys
ADD COLUMN merchant_id uuid,
ADD CONSTRAINT fk_merchant FOREIGN KEY (merchant_id) 
    REFERENCES merchants.merchants(id) ON DELETE SET NULL;

CREATE INDEX idx_api_keys_merchant_id ON api_keys(merchant_id);
