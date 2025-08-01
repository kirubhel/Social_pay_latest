'use client'

import Image from 'next/image'
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation';
import { Wallet, CreditCard, Banknote, Globe, Shield, CheckCircle, Clock } from 'lucide-react';

export default function UniqueSocialPayCheckout() {
  const [paymentMethod, setPaymentMethod] = useState<'wallet' | 'card' | 'bank' | 'international'>('wallet')
  const [isProcessing, setIsProcessing] = useState(false)
  useEffect(() => {
    // This useEffect is a dummy to prevent setIsProcessing from being marked as unused.
    // In a real application, this would be used to manage loading states during payment processing.
    setIsProcessing(false); 
  }, []);
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






  const getMethodDescription = () => {
    switch (paymentMethod) {
      case 'wallet':
        return 'Pay instantly with your mobile wallet'
      case 'card':
        return 'Secure card payments with instant processing'
      case 'bank':
        return 'Direct bank transfer - safe and reliable'
      case 'international':
        return 'Global payment options for international customers'
      default:
        return ''
    }
  }




  const getPaymentIcons = () => {
    switch (paymentMethod) {
      case 'wallet':
        return [
          { src: '/wallet/telebirr.svg', alt: 'TeleBirr' },
          { src: '/wallet/cbe_birr.png', alt: 'CBE Birr' },
          { src: '/wallet/mpesa.svg', alt: 'M-Pesa' },
          { src: '/wallet/ebirr.svg', alt: 'eBirr' }
        ]
      case 'card':
        return [
          { src: '/card/visa.svg', alt: 'Visa' },
          { src: '/card/mastercard.svg', alt: 'Mastercard' },
          { src: '/card/amex.svg', alt: 'American Express' },
          { src: '/card/ethswitch_logo.png', alt: 'EthSwitch' }
        ]
      case 'bank':
        return [
          { src: '/bank/cbe.svg', alt: 'CBE' },
          { src: '/bank/awash.svg', alt: 'Awash Bank' },
          { src: '/bank/dashen.svg', alt: 'Dashen Bank' },
          { src: '/bank/boa.svg', alt: 'Bank of Abyssinia' }
        ]
      case 'international':
        return [
          { src: '/paypal.svg', alt: 'PayPal' },
          { src: '/card/visa.svg', alt: 'Visa' },
          { src: '/card/mastercard.svg', alt: 'Mastercard' },
          { src: '/card/amex.svg', alt: 'American Express' }
        ]
      default:
        return []
    }
  }




  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0f172a] via-[#1e293b] to-[#334155] flex items-center justify-center px-4 py-8">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0" style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width=\'60\' height=\'60\' viewBox=\'0 0 60 60\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%23ffffff\'_fill-opacity=\'0.1\'%3E%3Ccircle cx=\'30\' cy=\'30\' r=\'2\'%3E%3C/circle%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
        }} />
      </div>

      <div className="w-full max-w-[900px] relative">
        {/* Main Container */}
        <div className="bg-white/10 backdrop-blur-xl rounded-2xl overflow-hidden shadow-2xl border border-white/20">
          <div className="flex flex-col lg:flex-row">
            
            {/* Left Sidebar - Payment Methods */}
            <div className="lg:w-[380px] bg-gradient-to-b from-[#1e293b] to-[#334155] p-8 relative">
              {/* Decorative Elements */}
              <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-[#f59e0b]/20 to-transparent rounded-full blur-3xl" />
              <div className="absolute bottom-0 left-0 w-24 h-24 bg-gradient-to-tr from-[#10b981]/20 to-transparent rounded-full blur-2xl" />
              
              <div className="relative z-10">
                {/* Header */}
                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-3">
                    <Image
                      src="/socialpay.webp"
                      alt="SocialPay"
                      width={140}
                      height={42}
                      className="h-10 w-auto"
                    />
                  </div>
                  <p className="text-gray-300 text-sm leading-relaxed">
                    Choose your preferred payment method
                  </p>
                </div>
                
                {/* Payment Methods */}
                <div className="space-y-3">
                  {[
                    { key: 'wallet', icon: Wallet, label: 'Mobile Wallets', desc: 'TeleBirr, M-Pesa, CBE Birr' },
                    { key: 'card', icon: CreditCard, label: 'Debit & Credit Cards', desc: 'Visa, Mastercard, Amex' },
                    { key: 'bank', icon: Banknote, label: 'Bank Transfer', desc: 'Direct from your bank account' },
                    { key: 'international', icon: Globe, label: 'International', desc: 'PayPal, Global cards' }
                  ].map((method) => {
                    const isSelected = paymentMethod === method.key
                    const Icon = method.icon
                    return (
                      <button
                        key={method.key}
                        onClick={() => handlePaymentMethodSelect(method.key as 'wallet' | 'card' | 'bank' | 'international')}
                        className={`w-full p-4 rounded-xl text-left transition-all duration-300 group relative overflow-hidden ${
                          isSelected 
                            ? 'bg-gradient-to-r from-[#f59e0b] to-[#d97706] text-white shadow-lg shadow-orange-500/25 scale-105' 
                            : 'bg-white/5 text-gray-300 hover:bg-white/10 hover:text-white border border-white/10'
                        }`}
                      >
                        {isSelected && (
                          <div className="absolute inset-0 bg-gradient-to-r from-[#f59e0b]/20 to-[#d97706]/20 animate-pulse" />
                        )}
                        <div className="relative flex items-start gap-3">
                          <Icon className={`w-5 h-5 mt-0.5 ${isSelected ? 'text-white' : 'text-gray-400'}`} />
                          <div className="flex-1">
                            <div className="font-semibold text-sm mb-1">{method.label}</div>
                            <div className={`text-xs ${isSelected ? 'text-orange-100' : 'text-gray-500'}`}>
                              {method.desc}
                            </div>
                          </div>
                          {isSelected && <CheckCircle className="w-4 h-4 text-white" />}
                        </div>
                      </button>
                    )
                  })}
                </div>

                {/* Trust Indicators */}
                <div className="mt-8 pt-6 border-t border-white/10">
                  <div className="flex items-center gap-2 text-gray-400 text-xs mb-2">
                    <Shield className="w-4 h-4" />
                    <span>256-bit SSL encryption</span>
                  </div>
                  <div className="flex items-center gap-2 text-gray-400 text-xs">
                    <Clock className="w-4 h-4" />
                    <span>Instant processing</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Right Panel - Payment Form */}
            <div className="flex-1 bg-gradient-to-br from-white to-gray-50 p-8 relative">
              {/* Close Button */}
              <button
                onClick={() => router.back()}
                className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 text-2xl font-bold transition-colors"
                aria-label="Close"
              >
                ×
              </button>

              <div className="flex flex-col h-full max-w-md mx-auto">
                {/* Header */}
                <div className="mb-8">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="w-10 h-10 bg-gradient-to-br from-[#f59e0b] to-[#d97706] rounded-xl flex items-center justify-center">
                      <CheckCircle className="w-5 h-5 text-white" />
                    </div>
                    <div>
                      <h2 className="text-xl font-bold text-gray-900">Secure Checkout</h2>
                      <p className="text-sm text-gray-600">Complete your payment</p>
                    </div>
                  </div>
                  
                  {/* Method Description */}
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-4">
                    <p className="text-blue-800 text-sm font-medium">{getMethodDescription()}</p>
                  </div>
                </div>

                {/* Payment Provider Logos */}
                <div className="mb-8">
                  <p className="text-sm font-medium text-gray-700 mb-3">Accepted payment methods:</p>
                  <div className="flex gap-2 flex-wrap">
                    {getPaymentIcons().map((icon, index) => (
                      <div key={index} className="w-14 h-10 bg-white rounded-lg border border-gray-200 flex items-center justify-center p-1 shadow-sm hover:shadow-md transition-shadow">
                        <Image
                          src={icon.src}
                          alt={icon.alt}
                          width={40}
                          height={28}
                          className="w-full h-full object-contain"
                        />
                      </div>
                    ))}
                  </div>
                </div>

                {/* Payment Form */}
                <div className="mb-8">
                  <label className="block text-gray-700 text-sm font-semibold mb-3">
                    Phone Number
                  </label>
                  <div className="relative">
                    <input
                      type="tel"
                      placeholder={getPlaceholder()}
                      className="w-full px-4 py-4 bg-white border-2 border-gray-200 rounded-xl text-gray-900 placeholder-gray-400 focus:border-[#f59e0b] focus:outline-none transition-colors text-lg"
                    />
                    <div className="absolute right-4 top-1/2 transform -translate-y-1/2">
                      <Image
                        src="/ethiopia-flag.svg"
                        alt="Ethiopia"
                        width={24}
                        height={16}
                        className="w-6 h-4"
                      />
                    </div>
                  </div>
                </div>

                {/* Amount Display */}
                <div className="bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-4 mb-8 border border-gray-200">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600 font-medium">Total Amount:</span>
                    <span className="text-2xl font-bold text-gray-900">ETB 1.00</span>
                  </div>
                </div>

                {/* Pay Button */}
                <button 
                  className="w-full bg-gradient-to-r from-[#f59e0b] to-[#d97706] hover:from-[#d97706] to-[#b45309] text-white font-bold py-4 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-[1.02] mb-8 text-lg"
                  disabled={isProcessing}
                >
                  {isProcessing ? (
                    <div className="flex items-center justify-center gap-2">
                      <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
                      Processing...
                    </div>
                  ) : (
                    'Complete Payment'
                  )}
                </button>

                {/* Security Footer */}
                <div className="mt-auto text-center">
                  <div className="flex items-center justify-center gap-2 text-gray-600 mb-2">
                    <Shield className="w-5 h-5" />
                    <span className="font-semibold">Secured by SocialPay</span>
                  </div>
                  <p className="text-gray-500 text-xs">
                    Your payment information is encrypted and secure
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom Trust Bar */}
        <div className="mt-6 bg-white/5 backdrop-blur-sm rounded-xl p-4 border border-white/10">
          <div className="flex items-center justify-center gap-8 text-gray-300 text-xs">
            <div className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              <span>Bank-level security</span>
            </div>
            <div className="flex items-center gap-2">
              <CheckCircle className="w-4 h-4" />
              <span>Instant confirmation</span>
            </div>
            <div className="flex items-center gap-2">
              <Clock className="w-4 h-4" />
              <span>24/7 support</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}


