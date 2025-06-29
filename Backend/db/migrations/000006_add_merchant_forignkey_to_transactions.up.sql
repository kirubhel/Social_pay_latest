ALTER TABLE public.transactions
ADD COLUMN merchant_id uuid,
ADD CONSTRAINT fk_merchant FOREIGN KEY (merchant_id) 
    REFERENCES merchants.merchants(id) ON DELETE SET NULL;

CREATE INDEX idx_transactions_merchant_id ON public.transactions(merchant_id);
