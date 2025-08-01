
'use client'

import { useState, useEffect } from 'react'
import { useTranslations } from 'next-intl'
import Image from 'next/image'
import { CreditCard, Wallet, Banknote, ShieldCheck } from 'lucide-react';

type PaymentMethod = 'socialpay' | 'wallet' | 'bank' | 'card'

interface PaymentMethodSelectorProps {
  onSelect?: (method: PaymentMethod) => void
  defaultMethod?: PaymentMethod
  availableTypes?: string[]
}

export default function PaymentMethodSelector({
  onSelect,
  defaultMethod = 'wallet',
  availableTypes = ['wallet', 'bank', 'card', 'socialpay'],
}: PaymentMethodSelectorProps) {
  const [selectedMethod, setSelectedMethod] = useState<PaymentMethod>(defaultMethod)
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const t = useTranslations('checkout');

  useEffect(() => {
    if (availableTypes.length > 0 && !availableTypes.includes(defaultMethod)) {
      const firstAvailable = availableTypes[0] as PaymentMethod;
      setSelectedMethod(firstAvailable);
      onSelect?.(firstAvailable);
    }
  }, [availableTypes, defaultMethod, onSelect]);

  const handleSelectMethod = (method: PaymentMethod) => {
    setSelectedMethod(method)
    setSelectedProvider(null); // Reset provider when category changes
    if (onSelect) {
      onSelect(method)
    }
  }

  const handleSelectProvider = (providerKey: string) => {
    setSelectedProvider(providerKey);
    // You might want to pass the selected provider up to the parent component here
    // if (onSelect) { onSelect(selectedMethod, providerKey); }
  };

  const methodIcons = {
    wallet: <Wallet className="w-6 h-6 mr-2" />,
    card: <CreditCard className="w-6 h-6 mr-2" />,
    bank: <Banknote className="w-6 h-6 mr-2" />,
    socialpay: <ShieldCheck className="w-6 h-6 mr-2" />,
  };

  const paymentProviders: Record<PaymentMethod, { key: string; name: string; logo: string }[]> = {
    wallet: [
      { key: 'telebirr', name: 'Telebirr', logo: '/telebirr.png' },
      { key: 'cbebirr', name: 'CBE Birr', logo: '/cbebirr.png' },
      { key: 'mpesa', name: 'M-Pesa', logo: '/mpesa.png' },
      { key: 'ebirr', name: 'eBirr', logo: '/ebirr.png' },
    ],
    card: [
      { key: 'visa', name: 'Visa', logo: '/visa.png' },
      { key: 'mastercard', name: 'Mastercard', logo: '/mastercard.png' },
      { key: 'amex', name: 'Amex', logo: '/amex.png' },
      { key: 'ethswitch', name: 'EthSwitch', logo: '/ethswitch.png' },
    ],
    bank: [
      { key: 'cbe', name: 'CBE', logo: '/cbe.png' },
      { key: 'awash', name: 'Awash Bank', logo: '/awash.png' },
      { key: 'dashen', name: 'Dashen Bank', logo: '/dashen.png' },
      { key: 'boa', name: 'Bank of Abyssinia', logo: '/boa.png' },
    ],
    socialpay: [
      { key: 'paypal', name: 'PayPal', logo: '/paypal.png' },
      { key: 'global_cards', name: 'Global Cards', logo: '/global_cards.png' },
    ],
  };

  const renderMethodButton = (method: PaymentMethod, translationKey: string) => {
    if (!availableTypes.includes(method)) return null;
    const isSelected = selectedMethod === method;
    return (
      <button
        key={method}
        className={`flex items-center w-full px-6 py-5 mb-1 rounded-2xl text-lg font-bold transition-all duration-200
          border-2 shadow-sm
          ${isSelected ? 'bg-gradient-to-r from-[#eaf3fa] to-[#f5a414]/20 border-[#f5a414] text-[#2B3A67] shadow-lg scale-105' : 'bg-white border-transparent text-[#8A8A8A] hover:bg-[#f4f8fb] hover:shadow-md hover:text-[#2B3A67]'}
        `}
        onClick={() => handleSelectMethod(method)}
        style={{ outline: 'none' }}
      >
        {methodIcons[method]}
        <span className="ml-2">{t(translationKey)}</span>
      </button>
    );
  };

  return (
    <div className="w-full max-w-[400px] bg-white rounded-xl border border-gray-200 shadow-sm flex flex-col gap-2 py-6 px-4 mb-6">
      <div className="text-gray-700 text-sm font-medium mb-2 text-center">{t("select_your_payment_method_here")}</div>
      <div className="border-t border-gray-200 mb-4"></div>
      {renderMethodButton('wallet', 'wallets')}
      {renderMethodButton('card', 'local_cards')}
      {renderMethodButton('bank', 'bank_transfers')}
      {renderMethodButton('socialpay', 'international_payments')}

      {selectedMethod && paymentProviders[selectedMethod] && (
        <div className="mt-4">
          <h3 className="text-gray-700 text-sm font-medium mb-2 text-center">Select Provider</h3>
          <div className="grid grid-cols-2 gap-2">
            {paymentProviders[selectedMethod].map((provider) => (
              <button
                key={provider.key}
                className={`flex flex-col items-center justify-center p-3 rounded-lg border-2 transition-all duration-200
                  ${selectedProvider === provider.key ? 'border-[#f5a414] bg-[#f5a414]/10' : 'border-gray-200 hover:border-gray-300'}
                `}
                onClick={() => handleSelectProvider(provider.key)}
              >
                <Image src={provider.logo} alt={provider.name} width={40} height={40} className="mb-1" />
                <span className="text-xs font-semibold text-gray-700">{provider.name}</span>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
} 


