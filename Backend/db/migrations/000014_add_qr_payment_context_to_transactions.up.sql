-- Add QR payment context and tip processing fields to transactions table

-- Create transaction source enum
CREATE TYPE transaction_source AS ENUM (
    'DIRECT',           -- Direct API calls
    'HOSTED_CHECKOUT',  -- Hosted checkout page
    'QR_PAYMENT',       -- QR payment link
    'WITHDRAWAL'        -- Withdrawal/tip transactions
);

-- Add new columns to transactions table
ALTER TABLE transactions ADD COLUMN transaction_source transaction_source DEFAULT 'DIRECT';
ALTER TABLE transactions ADD COLUMN qr_link_id UUID NULL;
ALTER TABLE transactions ADD COLUMN hosted_checkout_id UUID NULL;
ALTER TABLE transactions ADD COLUMN qr_tag VARCHAR(30) NULL;

-- Tip tracking fields
ALTER TABLE transactions ADD COLUMN has_tip BOOLEAN DEFAULT FALSE;
ALTER TABLE transactions ADD COLUMN tip_amount DECIMAL(20,2) NULL;
ALTER TABLE transactions ADD COLUMN tipee_phone VARCHAR(20) NULL;
ALTER TABLE transactions ADD COLUMN tip_medium VARCHAR(20) NULL;
ALTER TABLE transactions ADD COLUMN tip_transaction_id UUID NULL;
ALTER TABLE transactions ADD COLUMN tip_processed BOOLEAN DEFAULT FALSE;

-- Add foreign key constraints
ALTER TABLE transactions 
ADD CONSTRAINT fk_transactions_qr_link 
FOREIGN KEY (qr_link_id) REFERENCES qr_links(id);

ALTER TABLE transactions 
ADD CONSTRAINT fk_transactions_hosted_checkout 
FOREIGN KEY (hosted_checkout_id) REFERENCES hosted_payments(id);

-- Add indexes for performance
CREATE INDEX idx_transactions_qr_link_id ON transactions(qr_link_id);
CREATE INDEX idx_transactions_hosted_checkout_id ON transactions(hosted_checkout_id);
CREATE INDEX idx_transactions_tip_processing ON transactions(has_tip, tip_processed, status);
CREATE INDEX idx_transactions_source ON transactions(transaction_source);
CREATE INDEX idx_transactions_qr_tag ON transactions(qr_tag); 