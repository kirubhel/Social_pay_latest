'use client'

import { useTranslations } from 'next-intl'
import { WifiOff, AlertCircle, ServerCrash } from 'lucide-react'

interface ErrorPageProps {
  title?: string;
  message: string;
  type: 'NETWORK_ERROR' | 'BAD_REQUEST' | 'USER_ERROR' | 'NOT_FOUND';
  onRetry?: () => void;
}

export default function ErrorPage({ title, message, type, onRetry }: ErrorPageProps) {
  const t = useTranslations('checkout');

  const getIcon = () => {
    switch (type) {
      case 'NETWORK_ERROR':
        return <WifiOff size={64} className="text-gray-600" />;
      case 'NOT_FOUND':
        return <AlertCircle size={64} className="text-gray-600" />;
      case 'USER_ERROR':
      case 'BAD_REQUEST':
      default:
        return <ServerCrash size={64} className="text-gray-600" />;
    }
  };

  return (
    <div className="min-h-screen bg-[rgb(248,248,248)] flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-md text-center">
        <div className="mb-8 flex justify-center">
          {getIcon()}
        </div>
        
        <h1 className="text-2xl font-medium text-gray-900 mb-4">
          {title || t('error')}
        </h1>
        
        <p className="text-gray-600 mb-8">
          {message}
        </p>
        <p className="text-gray-600 text-sm mb-8">
          Error code: {type}
        </p>

        {onRetry && (
          <button
            onClick={onRetry}
            className="bg-[#30BB54] text-white px-6 py-3 rounded-lg font-medium hover:bg-[#28a349] transition-colors"
          >
            {t('tryAgain')}
          </button>
        )}
      </div>
    </div>
  );
} 