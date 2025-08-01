'use client'

import { useState } from 'react'
import Image from 'next/image'
import PhoneNumberInput from './PhoneNumberInput'
import { Gateway } from '@/services/gateway.service'
import { useQRPayment } from '@/contexts/QRPaymentContext'
import { useCheckoutPayment } from '@/contexts/CheckoutPaymentContext'

export type WalletType = 'TELEBIRR' | 'CBE' | 'EBIRR' | 'MPESA'

interface TipSectionProps {
  gateways: Gateway[];
  isTipEnabled: boolean;
  contextType?: 'qr' | 'checkout';
}

export default function TipSection({ gateways, isTipEnabled, contextType = 'qr' }: TipSectionProps) {
  const [showTipFields, setShowTipFields] = useState(false)
  const [selectedTipMedium, setSelectedTipMedium] = useState<string>(gateways[0]?.key || 'TELEBIRR')
  
  const qrPayment = useQRPayment();
  const checkoutPayment = useCheckoutPayment();
  
  // Use the appropriate context based on contextType
  const { tipAmount, setTipAmount, setTipeePhone, setTipMedium } = contextType === 'checkout' ? checkoutPayment : qrPayment;

  const handleTipToggle = () => {
    setShowTipFields(!showTipFields)
    if (!showTipFields) {
      setTipAmount('')
      setTipeePhone('')
      setTipMedium('TELEBIRR')
    } else {
      // Set initial tip medium when enabling tips
      setTipMedium(selectedTipMedium)
    }
  }

  const handleTipMediumSelection = (walletKey: string) => {
    setSelectedTipMedium(walletKey)
    setTipMedium(walletKey) // Update context
  }

  const handleTipeePhoneChange = (number: string) => {
    setTipeePhone(number) // Update context
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

  if (!isTipEnabled) return null

  return (
    <div className="mt-4">
      <div className="border border-gray-200 rounded-lg p-4">
        <div className="flex items-center justify-between mb-3">
          <h3 className="text-sm font-medium text-gray-700">Add Tip</h3>
          <label className="flex items-center">
            <input
              type="checkbox"
              checked={showTipFields}
              onChange={handleTipToggle}
              className="w-4 h-4 text-[#30BB54] border-gray-300 rounded focus:ring-[#30BB54]"
            />
            <span className="ml-2 text-sm text-gray-600">Enable tip</span>
          </label>
        </div>

        {showTipFields && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tip Amount (ETB)
              </label>
              <input
                type="number"
                value={tipAmount}
                onChange={(e) => setTipAmount(e.target.value)}
                placeholder="Enter tip amount"
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-[#30BB54] focus:border-transparent"
                min="0"
                step="0.01"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tip Payment Method
              </label>
              <div className="grid grid-cols-3 gap-2">
                {gateways.map((gateway) => (
                  <div 
                    key={`tip-${gateway.key}`}
                    onClick={() => handleTipMediumSelection(gateway.key)}
                    className={`
                      relative flex items-center justify-center
                      w-full h-[50px] rounded-[7.2px] cursor-pointer
                      ${selectedTipMedium === gateway.key 
                        ? 'bg-[#F2F2F2] border-[1.35px] border-[#30BB54] shadow-md' 
                        : 'bg-white border-[1.35px] border-[#EEEEEE]'
                      }
                    `}
                  >
                    {selectedTipMedium === gateway.key && (
                      <div className="absolute -top-1.5 -right-1.5 w-4 h-4 bg-[#30BB54] rounded-full flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="white" className="w-2.5 h-2.5">
                          <path fillRule="evenodd" d="M19.916 4.626a.75.75 0 01.208 1.04l-9 13.5a.75.75 0 01-1.154.114l-6-6a.75.75 0 011.06-1.06l5.353 5.353 8.493-12.739a.75.75 0 011.04-.208z" clipRule="evenodd" />
                        </svg>
                      </div>
                    )}
                    
                    <div className="flex items-center justify-center w-full h-full relative">
                      <Image 
                        src={getWalletLogo(gateway.key) || ''} 
                        alt={`Tip via ${gateway.key}`} 
                        fill
                        className="object-contain p-1" 
                      />
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tipee Phone Number
              </label>
              <PhoneNumberInput 
                walletType={selectedTipMedium as WalletType} 
                logoPath={getWalletLogo(selectedTipMedium) || undefined}
                onChange={handleTipeePhoneChange}
              />
            </div>
          </div>
        )}
      </div>
    </div>
  )
} 