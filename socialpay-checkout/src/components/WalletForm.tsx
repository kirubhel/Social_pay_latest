'use client'

import { useState } from 'react'
import Image from 'next/image'
import PhoneNumberInput from './PhoneNumberInput'
import PaymentButtons from './PaymentButtons'
import TipSection from './TipSection'
import { Gateway } from '@/services/gateway.service'
import { usePayment } from '@/contexts/PaymentContext'
import { useMerchantPayment } from '@/contexts/MerchantPaymentContext'
import { useCheckoutPayment } from '@/contexts/CheckoutPaymentContext'

export type WalletType = 'TELEBIRR' | 'CBE' | 'EBIRR' | 'MPESA'

interface WalletFormProps {
  gateways: Gateway[];
  merchantId?: string;
  checkoutId?: string;
  initialPhoneNumber?: string;
  isTipEnabled?: boolean;
}

export default function WalletForm({ gateways, merchantId, checkoutId, initialPhoneNumber, isTipEnabled }: WalletFormProps) {
  console.log('initialPhoneNumber', initialPhoneNumber)
  const [selectedWallet, setSelectedWallet] = useState<string>(gateways[0]?.key || 'TELEBIRR')
  const [phoneNumber, setPhoneNumber] = useState(initialPhoneNumber || '')
  const [isValid, setIsValid] = useState(false)
  
  // Use the appropriate context based on what props are provided
  const regularPayment = usePayment();
  const merchantPayment = useMerchantPayment();
  const checkoutPayment = useCheckoutPayment();
  
  const paymentContext = checkoutId ? checkoutPayment : merchantId ? merchantPayment : regularPayment;
  const { amount, processPayment, isLoading } = paymentContext;
  
  // Get tip data if it's a checkout payment
  const tipAmount = checkoutId && 'tipAmount' in paymentContext ? paymentContext.tipAmount : '';
  const tipeePhone = checkoutId && 'tipeePhone' in paymentContext ? paymentContext.tipeePhone : '';
  const tipMedium = checkoutId && 'tipMedium' in paymentContext ? paymentContext.tipMedium : '';

  const handleWalletSelection = (walletKey: string) => {
    setSelectedWallet(walletKey)
  }

  const handlePhoneNumberChange = (number: string, isValid: boolean) => {
    setPhoneNumber(number)
    setIsValid(isValid)
  }

  const handlePayNow = async () => {
    if (!isValid || !selectedWallet) return;

    if (checkoutId) {
      // Checkout payment - cast to the correct type
      const checkoutProcessPayment = processPayment as (medium: string, checkoutId: string, phone: string, tipData?: { amount: number, phone: string, medium: string }) => Promise<void>;
      const tipData = parseFloat(tipAmount) > 0 
        ? {
            amount: parseFloat(tipAmount),
            phone: tipeePhone,
            medium: tipMedium
          }
        : undefined;
      await checkoutProcessPayment(selectedWallet, checkoutId, phoneNumber, tipData);
    } else if (merchantId) {
      // Merchant payment - cast to the correct type
      const merchantProcessPayment = processPayment as (medium: string, merchantId: string, phone: string) => Promise<void>;
      await merchantProcessPayment(selectedWallet, merchantId, phoneNumber);
    } else {
      // Regular payment - cast to the correct type
      const regularProcessPayment = processPayment as (medium: string, details: { phone: string }) => Promise<void>;
      await regularProcessPayment(selectedWallet, {
        phone: phoneNumber
      });
    }
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
    console.log('key', key)
    const gateway = wallets.find(g => g.id === key);
    return gateway?.logo || null;
  }

  return (
    <div className="w-full max-w-[466px] pt-4">
      <div className="px-7 mb-3">
        {/* Wallet Grid */}
        <div className="grid grid-cols-3 gap-4">
          {gateways.map((gateway) => (
            <div 
              key={gateway.key}
              onClick={() => handleWalletSelection(gateway.key)}
              className={`
                relative flex items-center justify-center
                w-full h-[63px] rounded-[7.2px] cursor-pointer
                ${selectedWallet === gateway.key 
                  ? 'bg-[#eaf3fa] border-[1.35px] border-[#2B3A67] shadow-md' 
                  : 'bg-white border-[1.35px] border-[#EEEEEE] hover:border-[#f5a414]'}
              `}
            >
              {/* Checkmark for selected wallet */}
              {selectedWallet === gateway.key && (
                <div className="absolute -top-2.7 -right-2.7 w-5.4 h-5.4 bg-[#2B3A67] rounded-full flex items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" className="w-3.6 h-3.6">
                    <path fillRule="evenodd" d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z" clipRule="evenodd" />
                  </svg>
                </div>
              )}
              
              {/* Wallet logo and name */}
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
          initialPhoneNumber={initialPhoneNumber}
          walletType={selectedWallet as WalletType} 
          onChange={handlePhoneNumberChange}
          logoPath={getWalletLogo(selectedWallet) || undefined}
        />
      </div>
      
      {/* Tip Section for checkout payments */}
      {checkoutId && isTipEnabled && (
        <TipSection 
          gateways={gateways}
          isTipEnabled={isTipEnabled}
          contextType="checkout"
        />
      )}
      
      <div className="mt-6">
        <PaymentButtons 
          onPay={handlePayNow}
          disabled={!isValid}
          totalPrice={checkoutId && parseFloat(tipAmount) > 0 
            ? (parseFloat(amount || '0') + parseFloat(tipAmount)).toString()
            : amount
          }
          isLoading={isLoading}
        />
      </div>
    </div>
  )
} 