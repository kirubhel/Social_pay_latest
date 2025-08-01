'use client'

import { useTranslations } from 'next-intl'

interface PaymentButtonsProps {
  onPay?: (e: React.FormEvent) => void;
  disabled?: boolean; 
  totalPrice: string;
  isLoading?: boolean;
}

export default function PaymentButtons({ onPay, disabled, totalPrice, isLoading }: PaymentButtonsProps) {
  const t = useTranslations('checkout');

  return (
    <div>

    <div className="gap-7 hidden md:flex">
      <button 
        className={`h-[48px] w-[160px] bg-[#f5a414] text-white rounded-lg text-[16px] font-semibold shadow hover:bg-[#ffb940] active:bg-[#d98c00] transition-all duration-150 ${(disabled || !totalPrice || isLoading) ? 'opacity-70' : ''}`}
        onClick={onPay}
        disabled={disabled || !totalPrice || isLoading}
      >
        {isLoading ? (
          <div className="w-5 h-5 border-t-2 border-white rounded-full animate-spin mx-auto" />
        ) : (
          <> {t("pay")} <span className="font-semibold">{Number(totalPrice).toFixed(2)}</span></>
        )}
      </button>
    </div>
    <div className="fixed bottom-0 left-0 w-full md:hidden bg-white shadow-[0px_4px_20px_rgba(0,0,0,0.15)] h-[70px]" style={{ zIndex: '9999' }}>
      <div className="flex justify-end items-center px-6 h-full">
        <button 
          className={`h-[49px] w-full bg-[#f5a414] text-white rounded-xl text-[16px] font-bold shadow-md hover:bg-[#ffb940] active:bg-[#d98c00] transition-all duration-150 ${(disabled || !totalPrice || isLoading) ? 'opacity-70' : ''}`}
          onClick={onPay}
          disabled={disabled || !totalPrice || isLoading}
        >
          {isLoading ? (
            <div className="w-5 h-5 border-t-2 border-white rounded-full animate-spin mx-auto" />
          ) : (
            <> {t("pay")} <span className="font-extrabold">{Number(totalPrice).toFixed(2)}</span></>
          )}
        </button>
      </div>
    </div>
    </div>
  )
} 