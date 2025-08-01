'use client'

import { useState } from 'react'
import Image from 'next/image'
import PhoneNumberInput from './PhoneNumberInput'
import PaymentButtons from './PaymentButtons'
import TipSection from './TipSection'
import { Gateway } from '@/services/gateway.service'
import { useQRPayment } from '@/contexts/QRPaymentContext'

export type WalletType = 'TELEBIRR' | 'CBE' | 'EBIRR' | 'MPESA'

interface QRWalletFormProps {
  gateways: Gateway[];
  qrLinkId: string;
  isTipEnabled: boolean;
  isStaticAmount: boolean;
  staticAmount?: number;
}

export default function QRWalletForm({ 
  gateways, 
  qrLinkId, 
  isTipEnabled, 
  isStaticAmount, 
  staticAmount 
}: QRWalletFormProps) {
  const [selectedWallet, setSelectedWallet] = useState<string>(gateways[0]?.key || 'TELEBIRR')
  const [phoneNumber, setPhoneNumber] = useState('')
  const [isValid, setIsValid] = useState(false)
  
  const { 
    amount, 
    tipAmount,
    tipeePhone,
    tipMedium,
    processPayment, 
    isLoading 
  } = useQRPayment();

  const handleWalletSelection = (walletKey: string) => {
    setSelectedWallet(walletKey)
  }

  const handlePhoneNumberChange = (number: string, isValid: boolean) => {
    setPhoneNumber(number)
    setIsValid(isValid)
  }

  const handlePayNow = async () => {
    if (!isValid || !selectedWallet) return;

    const paymentAmount = isStaticAmount ? staticAmount : parseFloat(amount);
    const tipData = parseFloat(tipAmount) > 0 
      ? {
          amount: parseFloat(tipAmount),
          phone: tipeePhone,
          medium: tipMedium
        }
      : undefined;

    await processPayment(selectedWallet, qrLinkId, phoneNumber, paymentAmount, tipData);
  }



  const wallets = [
    {
      id: 'TELEBIRR',
      name: 'Telebirr',
      logo: '/wallet/telebirr.svg',
    },
    {
      id: 'CBE',
      name: 'CBE',
      logo: '/wallet/cbe.svg',
    },
    {
      id: 'EBIRR',
      name: 'E Birr',
      logo: '/wallet/ebirr.svg',
    },
    {
      id: 'MPESA',
      name: 'M-Pessa',
      logo: '/wallet/mpesa.svg',
    },
  ]

  const getWalletLogo = (key: string) => {
    const gateway = wallets.find(g => g.id === key);
    return gateway?.logo || null;
  }

  const totalAmount = isStaticAmount 
    ? (staticAmount || 0) + (parseFloat(tipAmount) > 0 ? parseFloat(tipAmount) : 0)
    : parseFloat(amount || '0') + (parseFloat(tipAmount) > 0 ? parseFloat(tipAmount) : 0);

  return (
    <div className="w-full max-w-[466px] pt-4">
      {/* Payment Method Selection */}
      <div className="mb-3">
        <h3 className="text-sm font-medium text-gray-700 mb-2">Choose Payment Method</h3>
        <div className="grid grid-cols-3 gap-4">
          {gateways.map((gateway) => (
            <div 
              key={gateway.key}
              onClick={() => handleWalletSelection(gateway.key)}
              className={`
                relative flex items-center justify-center
                w-full h-[63px] rounded-[7.2px] cursor-pointer
                ${selectedWallet === gateway.key 
                  ? 'bg-[#F2F2F2] border-[1.35px] border-[#30BB54] shadow-md' 
                  : 'bg-white border-[1.35px] border-[#EEEEEE]'
                }
              `}
            >
              {selectedWallet === gateway.key && (
                <div className="absolute -top-2.7 -right-2.7 w-5.4 h-5.4 bg-[#30BB54] rounded-full flex items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" className="w-3.6 h-3.6">
                    <path fillRule="evenodd" d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z" clipRule="evenodd" />
                  </svg>
                </div>
              )}
              
              <div className="flex items-center justify-center w-full h-full relative">
                <Image 
                  src={getWalletLogo(gateway.key) || ''} 
                  alt={gateway.key} 
                  fill
                  className="object-contain p-2" 
                />
              </div>
            </div>
          ))}
        </div>
      </div>
      
      {/* Phone Number Input */}
      <div className="mt-2">
        <PhoneNumberInput 
          walletType={selectedWallet as WalletType} 
          onChange={handlePhoneNumberChange}
          logoPath={getWalletLogo(selectedWallet) || undefined}
        />
      </div>

      {/* Tip Section */}
      <TipSection 
        gateways={gateways}
        isTipEnabled={isTipEnabled}
      />
      
      <div className="mt-6">
        <PaymentButtons 
          onPay={handlePayNow}
          disabled={!isValid || (isStaticAmount ? false : !amount)}
          totalPrice={totalAmount.toString()}
          isLoading={isLoading}
        />
      </div>
    </div>
  )
} 