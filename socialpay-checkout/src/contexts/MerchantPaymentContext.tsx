"use client"

import React, { createContext, useContext, useState } from 'react';
import { initiatePayment as initiateMerchantPayment, getTransactionStatus } from '../services/merchant.service';
import { PaymentInitiationRequest, TransactionStatusResponse } from '../types/merchant.dto';
import { V2ClientError } from '../services/merchant.service';
import LoadingSpinner from '../components/LoadingSpinner';
import PaymentIframe from '../components/PaymentIframe';
import { FaCheckCircle, FaTimesCircle } from 'react-icons/fa';

interface MerchantPaymentContextType {
  amount: string;
  phoneNumber: string;
  setAmount: (amount: string) => void;
  setPhoneNumber: (phone: string) => void;
  processPayment: (medium: string, merchantId: string, phone: string) => Promise<void>;
  isLoading: boolean;
}

const MerchantPaymentContext = createContext<MerchantPaymentContextType>({
  amount: '',
  phoneNumber: '',
  setAmount: () => {},
  setPhoneNumber: () => {},
  processPayment: async () => {},
  isLoading: false,
});

export const useMerchantPayment = () => {
  const context = useContext(MerchantPaymentContext);
  if (!context) {
    throw new Error('useMerchantPayment must be used within a MerchantPaymentProvider');
  }
  return context;
};

interface PaymentPopupProps {
  isOpen: boolean;
  onClose: () => void;
  transactionId: string;
  merchantName: string;
  amount: number;
}

const PaymentPopup: React.FC<PaymentPopupProps> = ({ 
  isOpen, 
  onClose, 
  transactionId, 
  merchantName, 
  amount 
}) => {
  const [status, setStatus] = useState<string>('PENDING');
  const [transaction, setTransaction] = useState<TransactionStatusResponse | null>(null);

  React.useEffect(() => {
    if (!isOpen || !transactionId) return;

    const checkStatus = async () => {
      try {
        const response = await getTransactionStatus(transactionId);
        setTransaction(response);
        setStatus(response.status);
        
        if (response.status !== 'PENDING') {
          // Stop polling when status is not pending
          return;
        }
      } catch (error) {
        console.error('Error checking transaction status:', error);
      }
    };

    // Initial check
    checkStatus();

    // Poll every 3 seconds while status is PENDING
    const interval = setInterval(() => {
      if (status === 'PENDING') {
        checkStatus();
      } else {
        clearInterval(interval);
      }
    }, 3000);

    return () => clearInterval(interval);
  }, [isOpen, transactionId, status]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-4">Payment Status</h3>
          
          <div className="mb-4">
            <p className="text-sm text-gray-600">Merchant: {merchantName}</p>
            <p className="text-lg font-medium">{amount.toFixed(2)} ETB</p>
          </div>

          <div className="mb-6">
            {status === 'PENDING' && (
              <div className="flex flex-col items-center justify-center">
                <div className="mb-4">
                  <LoadingSpinner size="lg" />
                </div>
                <div className="text-center">
                  <p className="text-lg font-medium text-gray-800">Processing Payment</p>
                  <p className="text-sm text-gray-500 mt-1">Please wait while we confirm your transaction...</p>
                  <div className="flex items-center justify-center gap-1 mt-2">
                    <span className="w-2 h-2 bg-[#30BB54] rounded-full animate-bounce" style={{animationDelay: '0ms'}}></span>
                    <span className="w-2 h-2 bg-[#30BB54] rounded-full animate-bounce" style={{animationDelay: '150ms'}}></span>
                    <span className="w-2 h-2 bg-[#30BB54] rounded-full animate-bounce" style={{animationDelay: '300ms'}}></span>
                  </div>
                </div>
              </div>
            )}
            
            {status === 'SUCCESS' && (
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <FaCheckCircle size={64} className="text-[#30BB54]" />
                </div>
                <h3 className="text-xl font-semibold text-gray-800 mb-2">Payment Successful!</h3>
                <p className="text-gray-600 mb-4">Your payment has been processed successfully</p>
                {transaction && (
                  <div className="bg-gray-50 rounded-lg p-4 mb-4">
                    <p className="text-sm text-gray-600">Transaction Reference</p>
                    <p className="font-mono text-sm font-medium text-gray-800">{transaction.reference}</p>
                  </div>
                )}
                <button
                  onClick={() => window.open(`/receipt/${transactionId}`, '_blank')}
                  className="w-full bg-gray-100 text-gray-700 py-2 px-4 rounded-lg hover:bg-gray-200 transition-colors mb-3 flex items-center justify-center gap-2"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
                    <path fillRule="evenodd" d="M5.625 1.5c-1.036 0-1.875.84-1.875 1.875v17.25c0 1.035.84 1.875 1.875 1.875h12.75c1.035 0 1.875-.84 1.875-1.875V12.75A3.75 3.75 0 0016.5 9h-1.875a1.875 1.875 0 01-1.875-1.875V5.25A3.75 3.75 0 009 1.5H5.625zM7.5 15a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5A.75.75 0 017.5 15zm.75 2.25a.75.75 0 000 1.5H12a.75.75 0 000-1.5H8.25z" clipRule="evenodd" />
                    <path d="M12.971 1.816A5.23 5.23 0 0114.25 5.25v1.875c0 .207.168.375.375.375H16.5a5.23 5.23 0 013.434 1.279 9.768 9.768 0 00-6.963-6.963z" />
                  </svg>
                  View Receipt
                </button>
              </div>
            )}
            
            {status === 'FAILED' && (
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <FaTimesCircle size={64} className="text-red-500" />
                </div>
                <h3 className="text-xl font-semibold text-red-600 mb-2">Payment Failed</h3>
                <p className="text-gray-600 mb-4">Unfortunately, your payment could not be processed</p>
              </div>
            )}
          </div>

          <button
            onClick={onClose}
            className="w-full bg-[#30BB54] text-white py-2 px-4 rounded-lg hover:bg-[#28a745]"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export const MerchantPaymentProvider: React.FC<{ 
  children: React.ReactNode;
  merchantId: string;
  merchantName: string;
}> = ({ children, merchantId, merchantName }) => {
  const [amount, setAmount] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showPaymentPopup, setShowPaymentPopup] = useState(false);
  const [currentTransactionId, setCurrentTransactionId] = useState<string>('');
  const [showPaymentIframe, setShowPaymentIframe] = useState(false);
  const [paymentUrl, setPaymentUrl] = useState<string>('');
 console.log("merchantId", merchantId)
  const handleSetAmount = (newAmount: string) => {
    console.log('Setting amount to:', newAmount);
    setAmount(newAmount);
  };

  const handlePaymentProcess = async (medium: string, merchantId: string, phone: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const payload: PaymentInitiationRequest = {
        medium,
        phone_number:  phone.startsWith('251') ? phone : `251${phone}`,
        amount: parseFloat(amount),
        merchant_id: merchantId,
      };

      const response = await initiateMerchantPayment(payload);

      if (response.success) {
        // Check if payment_url exists and open in iframe
        if (response.payment_url && response.payment_url.trim() !== '') {
          setPaymentUrl(response.payment_url);
          setShowPaymentIframe(true);
        }
        
        setCurrentTransactionId(response.socialpay_transaction_id);
        setShowPaymentPopup(true);
      } else {
        setError('Payment initiation failed');
      }
    } catch (error) {
      if (error instanceof V2ClientError) {
        setError(error.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <MerchantPaymentContext.Provider value={{
      amount,
      phoneNumber,
      setAmount: handleSetAmount,
      setPhoneNumber,
      processPayment: handlePaymentProcess,
      isLoading,
    }}>
      {children}
      
      <PaymentPopup
        isOpen={showPaymentPopup}
        onClose={() => setShowPaymentPopup(false)}
        transactionId={currentTransactionId}
        merchantName={merchantName}
        amount={parseFloat(amount) || 0}
      />
      
      <PaymentIframe
        isOpen={showPaymentIframe}
        onClose={() => setShowPaymentIframe(false)}
        paymentUrl={paymentUrl}
      />
      
      {error && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
            <div className="text-center">
              <h3 className="text-lg font-semibold mb-4 text-red-600">Error</h3>
              <p className="mb-4">{error}</p>
              <button
                onClick={() => setError(null)}
                className="w-full bg-[#30BB54] text-white py-2 px-4 rounded-lg hover:bg-[#28a745]"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </MerchantPaymentContext.Provider>
  );
}; 