"use client"

import React, { createContext, useContext, useState } from 'react';
import { initiatePayment, processPayment, PaymentDetails, PaymentResponse, PaymentError } from '../services/payment.service';
import PaymentConfirmationModal from '../components/PaymentConfirmationModal';
import PaymentIframe from '../components/PaymentIframe';
import ErrorModal from '../components/ErrorModal';

interface PaymentContextType {
  amount: string;
  phoneNumber: string;
  setAmount: (amount: string) => void;
  setPhoneNumber: (phone: string) => void;
  processPayment: (medium: string, details: PaymentDetails) => Promise<void>;
  isLoading: boolean;
}

const PaymentContext = createContext<PaymentContextType>({
  amount: '',
  phoneNumber: '',
  setAmount: () => {},
  setPhoneNumber: () => {},
  processPayment: async () => {},
  isLoading: false,
});

export const usePayment = () => {
  const context = useContext(PaymentContext);
  if (!context) {
    throw new Error('usePayment must be used within a PaymentProvider');
  }
  return context;
};

export const PaymentProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [amount, setAmount] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [paymentData, setPaymentData] = useState<PaymentResponse['data'] | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showPaymentIframe, setShowPaymentIframe] = useState(false);
  const [paymentUrl, setPaymentUrl] = useState<string>('');

  // Add debugging


  const handleSetAmount = (newAmount: string) => {
    console.log('Setting amount to:', newAmount);
    setAmount(newAmount);
  };

  const handlePaymentProcess = async (medium: string, details: PaymentDetails) => {
    setIsLoading(true);
    setError(null);
   const user_id = window.location.pathname.split('/checkout/')[1];
   console.log("user_id", user_id)
    try {
      const payload = {
        to: user_id,
        medium,
        amount: parseFloat(amount),
        details,
        redirects: {
          success: window.location.origin + "/success",
          cancel: window.location.origin + "/cancel",
          declined: window.location.origin + "/declined"
        }
      };

      const response = await initiatePayment(payload);

      if (response.success) {
        if (response.data.status.next === 'TXN_PROCESS') {
          setPaymentData(response.data);
          setShowConfirmation(true);
        } else {
          setError('Unknown payment step encountered');
        }
      } else {
        setError('Payment initiation failed');
      }
    } catch (error) {
      if (error instanceof PaymentError) {
        setError(error.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleConfirmPayment = async () => {
    setIsLoading(true);
    setError(null);
    console.log("paymentData", paymentData)
    try {
      if (!paymentData?.transaction?.id) return;

      const response = await processPayment(paymentData.transaction.id);
      console.log("payment process", response)
      if (response.success) {
        if (response.data.status.next === 'TXN_CHECKOUT') {
          if (response.data.checkout?.data) {
            setPaymentUrl(response.data.checkout.data);
            setShowPaymentIframe(true);
            setShowConfirmation(false);
          } else {
            setError('Payment URL not provided');
          }
        } else {
          setError('Unknown payment step encountered');
        }
      } else {
        setError('Payment processing failed');
      }
    } catch (error) {
      if (error instanceof PaymentError) {
        setError(error.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <PaymentContext.Provider value={{
      amount,
      phoneNumber,
      setAmount: handleSetAmount,
      setPhoneNumber,
      processPayment: handlePaymentProcess,
      isLoading,
    }}>
      {children}
      <PaymentConfirmationModal
        isOpen={showConfirmation}
        onConfirm={handleConfirmPayment}
        onCancel={() => setShowConfirmation(false)}
        amount={paymentData?.transaction?.pricing?.amount || 0}
        fees={paymentData?.transaction?.pricing?.fees?.[0]?.transaction || 0}
        totalAmount={paymentData?.transaction?.pricing?.total_amount || 0}
        isLoading={isLoading}
      />
      <PaymentIframe
        isOpen={showPaymentIframe}
        onClose={() => setShowPaymentIframe(false)}
        paymentUrl={paymentUrl}
      />
      <ErrorModal
        isOpen={!!error}
        onClose={() => setError(null)}
        message={error || ''}
      />
    </PaymentContext.Provider>
  );
}; 