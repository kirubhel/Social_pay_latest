'use client'

import { useTranslations } from 'next-intl';

interface PaymentConfirmationModalProps {
  isOpen: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  amount: number;
  fees: number;
  totalAmount: number;
  isLoading: boolean;
}

export default function PaymentConfirmationModal({
  isOpen,
  onConfirm,
  onCancel,
  amount,
  fees,
  totalAmount,
  isLoading
}: PaymentConfirmationModalProps) {
  const t = useTranslations('checkout');

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[10000]">
      <div className="bg-white rounded-lg p-6 w-[400px]">
        <h2 className="text-xl font-semibold mb-4">{t('confirmPayment')}</h2>
        <div className="space-y-3">
          <div className="flex justify-between">
            <span>{t('amount')}</span>
            <span>{amount.toFixed(2)} ETB</span>
          </div>
          <div className="flex justify-between">
            <span>{t('fees')}</span>
            <span>{fees.toFixed(2)} ETB</span>
          </div>
          <div className="flex justify-between font-semibold">
            <span>{t('total')}</span>
            <span>{totalAmount.toFixed(2)} ETB</span>
          </div>
        </div>
        <div className="flex justify-end gap-3 mt-6">
          <button
            onClick={onCancel}
            className="px-4 py-2 border rounded-md"
            disabled={isLoading}
          >
            {t('cancel')}
          </button>
          <button
            onClick={onConfirm}
            className="px-4 py-2 bg-[#30BB54] text-white rounded-md flex items-center justify-center min-w-[100px]"
            disabled={isLoading}
          >
            {isLoading ? (
              <div className="w-5 h-5 border-t-2 border-white rounded-full animate-spin" />
            ) : (
              t('confirm')
            )}
          </button>
        </div>
      </div>
    </div>
  );
} 