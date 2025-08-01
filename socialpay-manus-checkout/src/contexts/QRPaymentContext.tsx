"use client"

import React, { createContext, useContext, useState } from 'react';
import { processQRPayment } from '../services/qr.service';
import { getTransactionStatus } from '../services/merchant.service';
import { QRPaymentRequest } from '../types/qr.dto';
import { TransactionStatusResponse } from '../types/merchant.dto';
import { V2ClientError } from '../services/qr.service';
import LoadingSpinner from '../components/LoadingSpinner';
import PaymentIframe from '../components/PaymentIframe';
import { FaCheckCircle, FaTimesCircle } from 'react-icons/fa';

interface QRPaymentContextType {
  amount: string;
  phoneNumber: string;
  tipAmount: string;
  tipeePhone: string;
  tipMedium: string;
  setAmount: (amount: string) => void;
  setPhoneNumber: (phone: string) => void;
  setTipAmount: (amount: string) => void;
  setTipeePhone: (phone: string) => void;
  setTipMedium: (medium: string) => void;
  processPayment: (medium: string, qrLinkId: string, phone: string, amount?: number, tipData?: { amount: number, phone: string, medium: string }) => Promise<void>;
  isLoading: boolean;
}

const QRPaymentContext = createContext<QRPaymentContextType>({
  amount: '',
  phoneNumber: '',
  tipAmount: '',
  tipeePhone: '',
  tipMedium: '',
  setAmount: () => {},
  setPhoneNumber: () => {},
  setTipAmount: () => {},
  setTipeePhone: () => {},
  setTipMedium: () => {},
  processPayment: async () => {},
  isLoading: false,
});

export const useQRPayment = () => {
  const context = useContext(QRPaymentContext);
  if (!context) {
    throw new Error('useQRPayment must be used within a QRPaymentProvider');
  }
  return context;
};

interface PaymentPopupProps {
  isOpen: boolean;
  onClose: () => void;
  transactionId: string;
  title: string;
  amount: number;
  tipAmount?: number;
}

const PaymentPopup: React.FC<PaymentPopupProps> = ({ 
  isOpen, 
  onClose, 
  transactionId, 
  title, 
  amount,
  tipAmount
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
            <p className="text-sm text-gray-600">{title}</p>
            <p className="text-lg font-medium">{amount.toFixed(2)} ETB</p>
            {tipAmount && tipAmount > 0 && (
              <p className="text-sm text-gray-500">+ {tipAmount.toFixed(2)} ETB tip</p>
            )}
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
                  className="bg-[#30BB54] text-white px-6 py-2 rounded-lg hover:bg-[#28a049] transition-colors mr-2"
                >
                  View Receipt
                </button>
              </div>
            )}
            
            {status === 'FAILED' && (
              <div className="text-center">
                <div className="flex justify-center mb-4">
                  <FaTimesCircle size={64} className="text-red-500" />
                </div>
                <h3 className="text-xl font-semibold text-gray-800 mb-2">Payment Failed</h3>
                <p className="text-gray-600 mb-4">Unfortunately, your payment could not be processed</p>
                {transaction && (
                  <div className="bg-red-50 rounded-lg p-4 mb-4">
                    <p className="text-sm text-red-600">Error Details</p>
                    <p className="text-sm text-red-800">{transaction.details || 'Payment processing failed'}</p>
                  </div>
                )}
              </div>
            )}
          </div>

          <div className="flex gap-2 justify-center">
            <button
              onClick={onClose}
              className="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

interface QRPaymentProviderProps {
  children: React.ReactNode;
  qrLinkId: string;
  title: string;
  initialAmount?: number;
}

export const QRPaymentProvider: React.FC<QRPaymentProviderProps> = ({ 
  children, 
  qrLinkId, 
  title, 
  initialAmount 
}) => {
  const [amount, setAmount] = useState<string>(initialAmount?.toString() || '');
  const [phoneNumber, setPhoneNumber] = useState<string>('');
  const [tipAmount, setTipAmount] = useState<string>('');
  const [tipeePhone, setTipeePhone] = useState<string>('');
  const [tipMedium, setTipMedium] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [showPaymentPopup, setShowPaymentPopup] = useState<boolean>(false);
  const [currentTransactionId, setCurrentTransactionId] = useState<string>('');
  const [showPaymentIframe, setShowPaymentIframe] = useState<boolean>(false);
  const [paymentUrl, setPaymentUrl] = useState<string>('');

  const handleSetAmount = (newAmount: string) => {
    setAmount(newAmount);
  };

  const handlePaymentProcess = async (
    medium: string, 
    qrId: string, 
    phone: string,
    paymentAmount?: number,
    tipData?: { amount: number, phone: string, medium: string }
  ) => {
    setIsLoading(true);
    try {
      // Format phone number to include country code
      const formattedPhone = phone.startsWith('251') ? phone : `251${phone}`;
      
      const request: QRPaymentRequest = {
        medium,
        phone_number: formattedPhone,
      };

      // Add amount for dynamic QR or use paymentAmount parameter
      if (paymentAmount !== undefined) {
        request.amount = paymentAmount;
      } else if (amount) {
        request.amount = parseFloat(amount);
      }

      // Add tip data if provided
      if (tipData && tipData.amount > 0) {
        request.tip_amount = tipData.amount;
        request.tipee_phone = tipData.phone.startsWith('251') ? tipData.phone : `251${tipData.phone}`;
        request.tip_medium = tipData.medium;
      }

      // Use the qrLinkId from the provider context for the actual API call
      const response = await processQRPayment(qrLinkId, request);
      
      if (response.success) {

        // Checking if redirect for ethswitch 
        if (response.payment_url && medium=="ETHSWITCH") {

           const payment_url = new URL(response.payment_url)
           payment_url.searchParams.set("cb",Date.now().toString());
           window.location.href=payment_url.toString()
           return;
        }

        // Check if payment_url exists and open in iframe
        if (response.payment_url && response.payment_url.trim() !== '') {
          setPaymentUrl(response.payment_url);
          setShowPaymentIframe(true);
        }
        
        setCurrentTransactionId(response.socialpay_transaction_id);
        setShowPaymentPopup(true);
      } else {
        throw new Error(response.message || 'Payment failed');
      }
    } catch (error) {
      console.error('Payment processing error:', error);
      let errorMessage = 'Payment failed. Please try again.';
      
      if (error instanceof V2ClientError) {
        errorMessage = error.message;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }
      
      alert(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClosePopup = () => {
    setShowPaymentPopup(false);
    setCurrentTransactionId('');
  };

  return (
    <QRPaymentContext.Provider
      value={{
        amount,
        phoneNumber,
        tipAmount,
        tipeePhone,
        tipMedium,
        setAmount: handleSetAmount,
        setPhoneNumber,
        setTipAmount,
        setTipeePhone,
        setTipMedium,
        processPayment: handlePaymentProcess,
        isLoading,
      }}
    >
      {children}
      <PaymentPopup
        isOpen={showPaymentPopup}
        onClose={handleClosePopup}
        transactionId={currentTransactionId}
        title={title}
        amount={parseFloat(amount) || 0}
        tipAmount={parseFloat(tipAmount) || 0}
      />
      <PaymentIframe
        isOpen={showPaymentIframe}
        onClose={() => setShowPaymentIframe(false)}
        paymentUrl={paymentUrl}
      />
    </QRPaymentContext.Provider>
  );
}; 