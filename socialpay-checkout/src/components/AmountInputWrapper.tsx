'use client'

import { usePayment } from '@/contexts/PaymentContext'
import { useMerchantPayment } from '@/contexts/MerchantPaymentContext'
import { useCheckoutPayment } from '@/contexts/CheckoutPaymentContext'
import { useQRPayment } from '@/contexts/QRPaymentContext'
import AmountInput from './AmountInput'
import { useEffect } from 'react'

interface AmountInputWrapperProps {
  merchantId?: string;
  checkoutId?: string;
  qrLinkId?: string;
}

export default function AmountInputWrapper({ merchantId, checkoutId, qrLinkId }: AmountInputWrapperProps = {}) {
  // Use the appropriate context based on what props are provided
  const regularPayment = usePayment();
  const merchantPayment = useMerchantPayment();
  const checkoutPayment = useCheckoutPayment();
  const qrPayment = useQRPayment();
  
  const paymentContext = qrLinkId ? qrPayment : checkoutId ? checkoutPayment : merchantId ? merchantPayment : regularPayment;
  const { amount, setAmount } = paymentContext;

  useEffect(() => {
    console.log('Amount in AmountInputWrapper:', amount);
  }, [amount]);

  const handleAmountChange = (value: string, isValid: boolean) => {
    if (isValid) {
      console.log('AmountInputWrapper handling change:', value);
      setAmount(value);
    }
  };

  return (
    <AmountInput
      value={amount}
      onChange={handleAmountChange}
    />
  );
} 