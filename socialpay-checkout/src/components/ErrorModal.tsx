'use client'

import { useTranslations } from 'next-intl';
import { WifiOff, AlertCircle, ServerCrash, X } from 'lucide-react';

interface ErrorModalProps {
  isOpen: boolean;
  onClose: () => void;
  message: string;
  type?: 'NETWORK_ERROR' | 'BAD_REQUEST' | 'USER_ERROR' | 'NOT_FOUND';
}

export default function ErrorModal({ isOpen, onClose, message, type = 'USER_ERROR' }: ErrorModalProps) {
  const t = useTranslations('checkout');

  if (!isOpen) return null;

  const getIcon = () => {
    switch (type) {
      case 'NETWORK_ERROR':
        return <WifiOff size={48} className="text-red-600" />;
      case 'NOT_FOUND':
        return <AlertCircle size={48} className="text-red-600" />;
      case 'USER_ERROR':
      case 'BAD_REQUEST':
      default:
        return <ServerCrash size={48} className="text-red-600" />;
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[10000]">
      <div className="bg-white rounded-lg p-8 w-[400px] relative">
        {/* Close button */}
        <button 
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 transition-colors"
        >
          <X size={24} />
        </button>

        {/* Icon and content */}
        <div className="flex flex-col items-center">
          <div className="mb-6">
            {getIcon()}
          </div>
          
          <h2 className="text-xl font-semibold mb-4 text-gray-900">
            {t('error')}
          </h2>
          
          <p className="text-gray-600 text-center mb-6">
            {message}
          </p>

          {/* Action buttons */}
          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="px-6 py-2.5 border border-gray-200 text-gray-700 rounded-md hover:bg-gray-50 transition-colors"
            >
              {t('cancel')}
            </button>
            <button
              onClick={() => window.location.reload()}
              className="px-6 py-2.5 bg-[#30BB54] text-white rounded-md hover:bg-[#28a349] transition-colors"
            >
              {t('tryAgain')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
} 