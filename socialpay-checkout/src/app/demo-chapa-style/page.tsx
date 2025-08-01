'use client'

import Image from 'next/image'
import { useState } from 'react'
import { useRouter } from 'next/navigation';
import { Wallet, CreditCard, Banknote, Globe } from 'lucide-react';

export default function ChapaStyleCheckoutPage() {
  const [paymentMethod, setPaymentMethod] = useState<'wallet' | 'card' | 'bank' | 'international'>('wallet')
  const router = useRouter();

  const handlePaymentMethodSelect = (method: 'wallet' | 'card' | 'bank' | 'international') => {
    setPaymentMethod(method)
  }

  const getPlaceholder = () => {
    switch (paymentMethod) {
      case 'wallet':
        return '09********'
      case 'card':
        return '(09|07)********'
      case 'bank':
        return '07********'
      case 'international':
        return '(09|07)********'
      default:
        return '09********'
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#1a2332] to-[#2a3441] flex items-center justify-center px-4">
      <div className="w-full max-w-[800px] flex rounded-xl overflow-hidden shadow-2xl">
        {/* Left Sidebar - Payment Methods */}
        <div className="w-[350px] bg-[#2a3441] p-8 flex flex-col">
          <div className="mb-8">
            <div className="flex items-center gap-3 mb-2">
              <Image
                src="/socialpay.webp"
                alt="SocialPay"
                width={120}
                height={36}
                className="h-8 w-auto"
              />
            </div>
            <p className="text-gray-300 text-sm">Select your payment method here</p>
          </div>
          
          <div className="border-t border-gray-600 mb-6"></div>
          
          <div className="space-y-2">
            <button
              onClick={() => handlePaymentMethodSelect('wallet')}
              className={`w-full flex items-center gap-3 px-4 py-4 rounded-lg text-left transition-all duration-200 ${
                paymentMethod === 'wallet' 
                  ? 'bg-[#4CAF50] text-white' 
                  : 'text-gray-300 hover:bg-gray-700'
              }`}
            >
              <Wallet className="w-5 h-5" />
              <span className="font-medium">Wallets</span>
              {paymentMethod !== 'wallet' && <span className="ml-auto">›</span>}
            </button>
            
            <button
              onClick={() => handlePaymentMethodSelect('card')}
              className={`w-full flex items-center gap-3 px-4 py-4 rounded-lg text-left transition-all duration-200 ${
                paymentMethod === 'card' 
                  ? 'bg-[#4CAF50] text-white' 
                  : 'text-gray-300 hover:bg-gray-700'
              }`}
            >
              <CreditCard className="w-5 h-5" />
              <span className="font-medium">Local Cards</span>
              {paymentMethod !== 'card' && <span className="ml-auto">›</span>}
            </button>
            
            <button
              onClick={() => handlePaymentMethodSelect('bank')}
              className={`w-full flex items-center gap-3 px-4 py-4 rounded-lg text-left transition-all duration-200 ${
                paymentMethod === 'bank' 
                  ? 'bg-[#4CAF50] text-white' 
                  : 'text-gray-300 hover:bg-gray-700'
              }`}
            >
              <Banknote className="w-5 h-5" />
              <span className="font-medium">Bank Transfers</span>
              {paymentMethod !== 'bank' && <span className="ml-auto">›</span>}
            </button>
            
            <button
              onClick={() => handlePaymentMethodSelect('international')}
              className={`w-full flex items-center gap-3 px-4 py-4 rounded-lg text-left transition-all duration-200 ${
                paymentMethod === 'international' 
                  ? 'bg-[#4CAF50] text-white' 
                  : 'text-gray-300 hover:bg-gray-700'
              }`}
            >
              <Globe className="w-5 h-5" />
              <span className="font-medium">International Payments</span>
              {paymentMethod !== 'international' && <span className="ml-auto">›</span>}
            </button>
          </div>
        </div>

        {/* Right Panel - Payment Form */}
        <div className="flex-1 bg-[#3a4a5c] p-8 relative">
          {/* Close Button */}
          <button
            onClick={() => router.back()}
            className="absolute top-4 right-4 text-gray-400 hover:text-white text-2xl font-bold"
            aria-label="Close"
          >
            ×
          </button>

          <div className="flex flex-col h-full">
            {/* Header */}
            <div className="mb-8">
              <div className="flex items-center gap-3 mb-4">
                <Image
                  src="/socialpay.webp"
                  alt="SocialPay"
                  width={120}
                  height={36}
                  className="h-8 w-auto"
                />
                <span className="text-white font-bold text-lg">Checkout</span>
              </div>
              <p className="text-gray-300 text-sm">
                Helping you serve international and local customers by providing access to secure digital payment solutions.
              </p>
            </div>

            {/* Payment Method Icons */}
            <div className="flex gap-3 mb-8">
              <div className="w-16 h-12 bg-white rounded flex items-center justify-center p-1">
                <Image
                  src="/wallet/telebirr.svg"
                  alt="TeleBirr"
                  width={48}
                  height={32}
                  className="w-full h-full object-contain"
                />
              </div>
              <div className="w-16 h-12 bg-white rounded flex items-center justify-center p-1">
                <Image
                  src="/wallet/cbe_birr.png"
                  alt="CBE Birr"
                  width={48}
                  height={32}
                  className="w-full h-full object-contain"
                />
              </div>
              <div className="w-16 h-12 bg-white rounded flex items-center justify-center p-1">
                <Image
                  src="/wallet/mpesa.svg"
                  alt="M-Pesa"
                  width={48}
                  height={32}
                  className="w-full h-full object-contain"
                />
              </div>
              <div className="w-16 h-12 bg-white rounded flex items-center justify-center p-1">
                <Image
                  src="/wallet/ebirr.svg"
                  alt="eBirr"
                  width={48}
                  height={32}
                  className="w-full h-full object-contain"
                />
              </div>
            </div>

            {/* Phone Number Input */}
            <div className="mb-8">
              <label className="block text-white text-sm font-medium mb-3">
                Phone Number
              </label>
              <input
                type="tel"
                placeholder={getPlaceholder()}
                className="w-full px-4 py-3 bg-[#2a3441] border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:border-[#4CAF50] focus:outline-none"
              />
            </div>

            {/* Pay Button */}
            <button className="w-full bg-[#8BC34A] hover:bg-[#7CB342] text-white font-bold py-4 rounded-lg transition-colors duration-200 mb-8">
              Pay ETB 1.00
            </button>

            {/* Security Footer */}
            <div className="mt-auto text-center">
              <div className="flex items-center justify-center gap-2 text-gray-300 mb-2">
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                </svg>
                <span className="font-medium">Secured By SocialPay</span>
              </div>
              <p className="text-gray-400 text-xs">Your payment is encrypted and protected</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

