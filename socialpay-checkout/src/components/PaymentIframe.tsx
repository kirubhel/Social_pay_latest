import Image from 'next/image';
import React from 'react';
import { FaTimes } from 'react-icons/fa';

interface PaymentIframeProps {
  isOpen: boolean;
  onClose: () => void;
  paymentUrl: string;
}

const PaymentIframe: React.FC<PaymentIframeProps> = ({ 
  isOpen, 
  onClose, 
  paymentUrl 
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-white z-[9999] flex flex-col">
      {/* Close button at the top */}
      <div className="bg-[#30BB54] p-4 flex justify-between items-center">
        <div className="flex items-center gap-2">
          <Image src="/socialpay-black.svg" alt="SocialPay" width={100} height={32} className="h-8 w-auto" />
          <span className="text-white font-medium">Secure Payment</span>
        </div>
        <button
          onClick={onClose}
          className="text-white hover:text-gray-200 transition-colors p-2 rounded-full hover:bg-black hover:bg-opacity-20"
          aria-label="Close payment window"
        >
          <FaTimes size={20} />
        </button>
      </div>
      
      {/* Fullscreen iframe */}
      <iframe
        src={paymentUrl}
        className="flex-1 w-full h-full border-0"
        title="Payment Gateway"
        allow="payment; microphone; camera; geolocation"
        sandbox="allow-same-origin allow-scripts allow-forms allow-popups allow-top-navigation allow-storage-access-by-user-activation"
      />
    </div>
  );
};

export default PaymentIframe; 