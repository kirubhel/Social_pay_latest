   ALTER TABLE public.transactions
     DROP CONSTRAINT IF EXISTS fk_merchant;
   
   ALTER TABLE public.transactions
     DROP COLUMN IF EXISTS merchant_id;

   DROP INDEX IF EXISTS idx_transactions_merchant_id;