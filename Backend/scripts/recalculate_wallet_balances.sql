-- Standalone script to recalculate correct wallet balances
-- This script recalculates admin and merchant wallet balances based on transaction history
-- Run this script: psql -d your_database -f scripts/recalculate_wallet_balances.sql

DO $$
DECLARE
    wallet_record RECORD;
    calculated_balance DECIMAL(20,2);
    current_balance DECIMAL(20,2);
    affected_rows INTEGER := 0;
    deposit_sum DECIMAL(20,2);
    withdrawal_sum DECIMAL(20,2);
BEGIN
    RAISE NOTICE 'Starting wallet balance recalculation...';
    
    -- Recalculate merchant wallet balances
    FOR wallet_record IN 
        SELECT w.id, w.user_id, w.merchant_id, w.balance as current_balance
        FROM merchant.wallet w 
        WHERE w.type = 'MERCHANT'
    LOOP
        -- Calculate deposits (positive merchant_net) from successful transactions
        SELECT COALESCE(SUM(merchant_net), 0)
        INTO deposit_sum
        FROM public.transactions 
        WHERE merchant_id = wallet_record.merchant_id 
        AND status = 'SUCCESS'
        AND type IN ('deposit', 'payment')
        AND merchant_net IS NOT NULL
        AND merchant_net > 0;
        
        -- Calculate withdrawals (positive merchant_net) from successful transactions
        SELECT COALESCE(SUM(ABS(merchant_net)), 0)
        INTO withdrawal_sum
        FROM public.transactions 
        WHERE merchant_id = wallet_record.merchant_id 
        AND status = 'SUCCESS'
        AND type = 'withdrawal'
        AND merchant_net IS NOT NULL;
        
        -- Balance = DEPOSITS - WITHDRAWALS
        calculated_balance := deposit_sum - withdrawal_sum;
        
        -- Update if balance is different
        IF wallet_record.current_balance != calculated_balance THEN
            UPDATE merchant.wallet 
            SET balance = calculated_balance,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = wallet_record.id;
            
            affected_rows := affected_rows + 1;
            
            RAISE NOTICE 'Updated merchant wallet % from % to % (Deposits: %, Withdrawals: %)', 
                wallet_record.merchant_id, wallet_record.current_balance, calculated_balance, deposit_sum, withdrawal_sum;
        END IF;
    END LOOP;
    
    -- Recalculate admin wallet balances (single admin wallet)
    FOR wallet_record IN 
        SELECT w.id, w.user_id, w.balance as current_balance
        FROM merchant.wallet w 
        WHERE w.type = 'ADMIN'
        LIMIT 1
    LOOP
        -- Calculate total admin net (commission) from successful transactions
        SELECT COALESCE(SUM(admin_net), 0)
        INTO calculated_balance
        FROM public.transactions 
        WHERE status = 'SUCCESS'
        AND admin_net IS NOT NULL;
        
        -- Update if balance is different
        IF wallet_record.current_balance != calculated_balance THEN
            UPDATE merchant.wallet 
            SET balance = calculated_balance,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = wallet_record.id;
            
            affected_rows := affected_rows + 1;
            
            RAISE NOTICE 'Updated single admin wallet % from % to %', 
                wallet_record.user_id, wallet_record.current_balance, calculated_balance;
        END IF;
    END LOOP;
    
    RAISE NOTICE 'Wallet balance recalculation completed. Updated % wallets.', affected_rows;
END $$; 