'use client'

import { useState } from 'react'
import PaymentMethodSelector from '@/components/PaymentMethodSelector'
import { useRouter } from 'next/navigation';

export default function Checkout() {
  const [paymentMethod, setPaymentMethod] = useState<'socialpay' | 'wallet' | 'bank' | 'card'>('wallet')
  const router = useRouter();

  const handlePaymentMethodSelect = (
    method: 'wallet' | 'bank' | 'card' | 'socialpay',
  ) => {
    setPaymentMethod(method)
  }

  const renderPaymentForm = () => {
    switch (paymentMethod) {
      case 'wallet':
        return (
          <div className="space-y-4">
            <div className="flex gap-4 justify-center">
              <div className="w-16 h-12 bg-blue-100 rounded flex items-center justify-center">
                <span className="text-xs font-bold text-blue-600">TeleBirr</span>
              </div>
              <div className="w-16 h-12 bg-orange-100 rounded flex items-center justify-center">
                <span className="text-xs font-bold text-orange-600">CBE Birr</span>
              </div>
              <div className="w-16 h-12 bg-green-100 rounded flex items-center justify-center">
                <span className="text-xs font-bold text-green-600">M-Pesa</span>
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Phone Number</label>
              <input 
                type="tel" 
                placeholder="09********" 
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#f5a414] focus:border-transparent"
              />
            </div>
            <button className="w-full bg-[#f5a414] text-white py-3 rounded-lg font-semibold hover:bg-[#e6940f] transition-colors">
              Pay ETB 1.00
            </button>
          </div>
        )
      case 'card':
        return (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Card Number</label>
              <input 
                type="text" 
                placeholder="1234 5678 9012 3456" 
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#f5a414] focus:border-transparent"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">Expiry Date</label>
                <input 
                  type="text" 
                  placeholder="MM/YY" 
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#f5a414] focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">CVV</label>
                <input 
                  type="text" 
                  placeholder="123" 
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#f5a414] focus:border-transparent"
                />
              </div>
            </div>
            <button className="w-full bg-[#f5a414] text-white py-3 rounded-lg font-semibold hover:bg-[#e6940f] transition-colors">
              Pay ETB 1.00
            </button>
          </div>
        )
      default:
        return (
          <div className="space-y-4">
            <div className="text-center text-gray-500">
              {paymentMethod} payment form would appear here
            </div>
            <button className="w-full bg-[#f5a414] text-white py-3 rounded-lg font-semibold hover:bg-[#e6940f] transition-colors">
              Pay ETB 1.00
            </button>
          </div>
        )
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#eaf3fa] to-[#f4f8fb] flex flex-col pb-[70px] md:pb-0">
      {/* Main Content */}
      <div className="flex-1 flex items-center justify-center px-4 py-8">
        <div className="w-full max-w-[700px] flex flex-col gap-0 rounded-xl shadow-lg overflow-hidden border border-[#e3e8ee]">
          <main className="flex-1 bg-white flex flex-col gap-8 p-6 md:p-8">
            <div className="flex-1 flex flex-col gap-2 justify-center relative">
              {/* X Close Button for the whole card */}
              <button
                onClick={() => router.back()}
                className="absolute -top-7 -right-7 bg-white shadow rounded-full text-gray-400 hover:text-gray-700 text-2xl font-semibold focus:outline-none z-20 w-10 h-10 flex items-center justify-center border border-[#e3e8ee]"
                aria-label="Close"
              >
                &times;
              </button>
              {/* Logo and Slogan above payment method selector */}
              <div className="w-full flex flex-col items-center justify-center mb-4">
                <div className="w-20 h-6 bg-[#f5a414] rounded flex items-center justify-center mb-2">
                  <span className="text-white font-bold text-sm">SocialPay</span>
                </div>
                <div className="mt-2 text-base font-semibold tracking-wide text-center">
                  <span style={{ color: '#2B3A67' }}>Connect.</span>
                  <span className="mx-1" style={{ color: '#f5a414' }}>Accept.</span>
                  <span style={{ color: '#1a3129' }}>Grow.</span>
                </div>
              </div>
              <h1 className="text-xl font-semibold mb-4 text-[#1a3129] tracking-wide text-center">How would you like to pay today?</h1>
              <PaymentMethodSelector
                onSelect={handlePaymentMethodSelect}
                defaultMethod={paymentMethod}
                availableTypes={['wallet', 'card', 'bank', 'socialpay']}
              />
              {renderPaymentForm()}
            </div>
          </main>
        </div>
      </div>
      {/* Trust/Security Footer */}
      <div className="max-w-[700px] mx-auto pt-8 flex flex-col items-center gap-2">
        <div className="flex items-center gap-2 text-[#2B3A67] text-lg font-semibold">
          <svg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='currentColor' className='w-6 h-6'><path strokeLinecap='round' strokeLinejoin='round' strokeWidth={2} d='M12 11c.5304 0 1.0391.2107 1.4142.5858C13.7893 11.9609 14 12.4696 14 13v2a2 2 0 11-4 0v-2c0-.5304.2107-1.0391.5858-1.4142C10.9609 11.2107 11.4696 11 12 11zm0 0V7a4 4 0 10-8 0v4m8 0a4 4 0 018 0v4m-8 0v2a2 2 0 104 0v-2' /></svg>
          Secured by SocialPay
        </div>
        <span className="text-xs text-[#8A8A8A]">Your payment is encrypted and protected</span>
      </div>
    </div>
  )
}

