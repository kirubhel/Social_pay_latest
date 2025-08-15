-- Create transaction status enum type
  DO $$
  BEGIN
      IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_status') THEN
          CREATE TYPE transaction_status AS ENUM (
              'INITIATED',
              'PENDING',
              'SUCCESS', 
              'FAILED',
              'REFUNDED',
              'EXPIRED',
              'CANCELED'
          );
      END IF;
  END
  $$;

-- Create transactions table
CREATE TABLE public.transactions (
    id UUID PRIMARY KEY,
    phone_number VARCHAR(50),
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    medium VARCHAR(50) NOT NULL,
    reference VARCHAR(100),
    comment TEXT,
    verified BOOLEAN DEFAULT false,
    ttl BIGINT,
    details JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirm_timestamp TIMESTAMP WITH TIME ZONE,
    reference_number VARCHAR(100),
    test BOOLEAN DEFAULT false,
    status transaction_status NOT NULL DEFAULT 'INITIATED',
    description TEXT,
    token VARCHAR(255),
    amount DECIMAL(20,2) NOT NULL,
    has_challenge BOOLEAN DEFAULT false,
    webhook_received BOOLEAN DEFAULT false,
    fee_amount DECIMAL(20,2),
    admin_net DECIMAL(20,2),
    vat_amount DECIMAL(20,2),
    merchant_net DECIMAL(20,2),
    total_amount DECIMAL(20,2),
    currency VARCHAR(3) DEFAULT 'ETB',
    callback_url TEXT,
    success_url TEXT,
    failed_url TEXT,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES auth.users(id)
);

-- Create indexes
CREATE INDEX idx_transactions_user_id ON public.transactions(user_id);
CREATE INDEX idx_transactions_status ON public.transactions(status);
CREATE INDEX idx_transactions_type ON public.transactions(type);
CREATE INDEX idx_transactions_created_at ON public.transactions(created_at);
CREATE INDEX idx_transactions_reference ON public.transactions(reference);
CREATE INDEX idx_transactions_reference_number ON public.transactions(reference_number);

-- Create trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON public.transactions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();