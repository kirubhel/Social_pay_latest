'use client'

import { useState } from 'react'
import Image from 'next/image'
// import { useTranslations } from 'next-intl'
import PhoneNumberInput from './PhoneNumberInput'
import PaymentButtons from './PaymentButtons'
import { Gateway } from '@/services/gateway.service'
import Divider from './Divider'
import { usePayment } from '@/contexts/PaymentContext'
import { useQRPayment } from '@/contexts/QRPaymentContext'

export type BankType = 'ahadu' | 'abyssinia' | 'awash' | 'buna'

interface BankFormProps {
  gateways: Gateway[];
  qrLinkId?: string;
}

export default function BankForm({ gateways, qrLinkId }: BankFormProps) {
  const [selectedBank, setSelectedBank] = useState<string>(gateways[0]?.key || '')
  const [phoneNumber, setPhoneNumber] = useState('')
  const [isValid, setIsValid] = useState(false)
  
  // Use the appropriate context based on whether qrLinkId is provided
  const regularPayment = usePayment();
  const qrPayment = useQRPayment();
  
  const paymentContext = qrLinkId ? qrPayment : regularPayment;
  const { amount, processPayment, isLoading } = paymentContext;
  // const t = useTranslations('checkout');

  const handleBankSelection = (bankKey: string) => {
    setSelectedBank(bankKey)
  }

  const handlePhoneNumberChange = (number: string, isValid: boolean) => {
    setPhoneNumber(number)
    setIsValid(isValid)
  }

  const handlePayNow = async () => {
    if (!isValid || !selectedBank) return;

    if (qrLinkId) {
      // QR payment - cast to the correct type
      const qrProcessPayment = processPayment as (medium: string, qrLinkId: string, phone: string) => Promise<void>;
      await qrProcessPayment(selectedBank, qrLinkId, phoneNumber);
    } else {
      // Regular payment - cast to the correct type
      const regularProcessPayment = processPayment as (medium: string, details: { phone: string }) => Promise<void>;
      await regularProcessPayment(selectedBank, {
        phone: phoneNumber
      });
    }
  }



  // Bank data array
const banks = [
  {
    id: 'ahadu',
    name: 'Ahadu Bank',
    icon: '/bank/ahadu.svg'
  },
  {
    id: 'DASH',
    name: 'Dashen Bank',
    icon: '/bank/dashen.svg'
  },
  {
    id: 'abyssinia',
    name: 'Abyssinia Bank',
    icon: '/bank/boa.svg' // Using boa.svg for Abyssinia as per available icons
  },
  {
    id: 'AWINETAA',
    name: 'Awash Bank',
    icon: '/bank/awash.svg'
  },
  {
    id: 'BUNAETAA',
    name: 'Buna Bank',
    icon: '/bank/buna.png'
  }
]

  const getBankLogo = (key: string) => {
    const gateway = banks.find(g => g.id === key);
    return gateway?.icon || null;
  }

  return (
    <div className="w-full max-w-[479px] pt-4 relative">
      <div className="relative w-full mb-[38px]">
        <div className="grid grid-cols-3 gap-4 mb-[34px]">
          {gateways.map((gateway, index) => (
            <button 
              key={gateway.key}
              onClick={() => handleBankSelection(gateway.key)}
              className={`
                relative flex items-center h-[50px] rounded-[5px] px-3
                ${selectedBank === gateway.key 
                  ? 'bg-[#eaf3fa] border border-[#2B3A67] shadow-md' 
                  : 'bg-white border border-[#EEEEEE] hover:border-[#f5a414]'}
                ${index > 2 && 'mt-4'}
              `}
            >
              <div className="w-[40px] h-[40px] relative mr-1">
                <Image 
                  src={getBankLogo(gateway.key) || ''} 
                  alt={gateway.name} 
                  fill
                  style={{ objectFit: 'contain' }}
                />
              </div>
              <span className="text-[#2B3A67] text-[13.7px] font-medium leading-[120%] ml-1">
                {gateway.short_name || gateway.name}
              </span>
              {selectedBank === gateway.key && (
                <div className="absolute -top-2 -right-2 w-[17.5px] h-[17.5px] bg-[#2B3A67] rounded-full flex items-center justify-center">
                  <svg width="8.4" height="6.3" viewBox="0 0 12 9" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M4.00001 6.586L1.70701 4.293L0.292908 5.707L4.00001 9.414L12 1.414L10.586 0L4.00001 6.586Z" fill="white"/>
                  </svg>
                </div>
              )}
            </button>
          ))}
          
          {/* Add empty slots to maintain grid layout if needed */}
          {gateways.length % 3 !== 0 && Array(3 - (gateways.length % 3)).fill(0).map((_, i) => (
            <div key={`empty-${i}`} className="h-[50px]"></div>
          ))}
        </div>
        <Divider />
      </div>

      {/* Phone Number Input */}
      <div className="mt-4">
        <PhoneNumberInput 
          walletType={selectedBank as BankType} 
          onChange={handlePhoneNumberChange}
          logoPath={getBankLogo(selectedBank) || undefined}
        />
      </div>
      
      <div className="mt-[42px]">
        <PaymentButtons 
          onPay={handlePayNow}
          disabled={!isValid}
          totalPrice={amount}
          isLoading={isLoading}
        />
      </div>
    </div>
  )
} 