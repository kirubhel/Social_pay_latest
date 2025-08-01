'use client'

import Image from 'next/image';
import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { Wallet, CreditCard, Banknote, Globe, Shield, CheckCircle, Clock, X, ChevronDown, Pointer } from 'lucide-react';
export default function UniqueSocialPayCheckout() {
  const [paymentMethod, setPaymentMethod] = useState<'wallet' | 'card' | 'bank' | 'international'>('wallet')
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)
  const [showPaymentModal, setShowPaymentModal] = useState(false)
  const [timeLeft, setTimeLeft] = useState(60)
  const [tipAmount, setTipAmount] = useState(0)
  const [showTipOptions, setShowTipOptions] = useState(false)
  const router = useRouter()
  const [phoneNumber, setPhoneNumber] = useState('');
  const [isValidPhone, setIsValidPhone] = useState(false);
  const [showDigitAlert, setShowDigitAlert] = useState(false);
  const [isDetailsExpanded, setIsDetailsExpanded] = useState(false);
  const [hasInteractedWithDetails, setHasInteractedWithDetails] = useState(false);
  

  useEffect(() => {
    setSelectedProvider(null)
  }, [paymentMethod])

  useEffect(() => {
    if (!showPaymentModal) return

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timer)
          setTimeout(() => {
            setShowPaymentModal(false)
            router.push('/')
          }, 1000)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => clearInterval(timer)
  }, [showPaymentModal, router])
  
  useEffect(() => {
  if (!showPaymentModal) return;

  const timer = setInterval(() => {
    setTimeLeft((prev) => {
      if (prev <= 1) {
        clearInterval(timer);
        setTimeout(() => {
          setShowPaymentModal(false);
          setIsProcessing(false);
          router.push('/');
        }, 1000);
        return 0;
      }
      return prev - 1;
    });
  }, 1000);

  return () => clearInterval(timer);
}, [showPaymentModal, router]);

  const handlePaymentMethodSelect = (method: 'wallet' | 'card' | 'bank' | 'international') => {
    setPaymentMethod(method)
  }

  const handlePaymentSubmit = () => {
    if (!selectedProvider) return
    setIsProcessing(true)
    setShowPaymentModal(true)
    setTimeLeft(60)
  }

  useEffect(() => {
    if (selectedProvider === 'M-Pesa' || selectedProvider === 'mpesa') {
      setIsValidPhone(/^7\d{8}$/.test(phoneNumber));
    } else {
      setIsValidPhone(/^9\d{8}$/.test(phoneNumber));
    }
  }, [phoneNumber, selectedProvider]);

  const getMethodDescription = () => {
    switch (paymentMethod) {
      case 'wallet': return 'Pay instantly with your mobile wallet'
      case 'card': return 'Secure card payments with instant processing'
      case 'bank': return 'Direct bank transfer - safe and reliable'
      case 'international': return 'Global payment options for international customers'
      default: return ''
    }
  }

  const getPaymentIcons = () => {
    switch (paymentMethod) {
      case 'wallet':
        return [
          { src: '/wallet/telebirr.svg', alt: 'TeleBirr' },
          { src: '/wallet/cbe_birr.png', alt: 'CBE Birr' },
          { src: '/wallet/mpesa.svg', alt: 'M-Pesa' },
          { src: '/wallet/ebirr.svg', alt: 'eBirr' },
          { src: '/wallet/kacha_logo.jpg', alt: 'Kacha' }
        ]
      case 'card':
        return [
          { src: '/card/visa.svg', alt: 'Visa' },
          { src: '/card/mastercard.svg', alt: 'Mastercard' },
          { src: '/card/amex.svg', alt: 'American Express' },
          { src: '/card/ethswitch.webp', alt: 'EthSwitch' }
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
      default: return []
    }
  }

  const handleTipSelect = (amount: number) => {
    setTipAmount(amount)
    setShowTipOptions(false)
  }

  const totalAmount = 1000 + tipAmount

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#0f172a] via-[#1e293b] to-[#334155] flex items-center justify-center px-4 py-8">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0" style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.1'%3E%3Ccircle cx='30' cy='30' r='2'%3E%3C/circle%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
        }} />
      </div>

      <div className="w-full max-w-[900px] relative">
        {/* Main Container */}
        <div className="bg-white/10 backdrop-blur-xl rounded-2xl overflow-hidden shadow-2xl border border-white/20">
          <div className="flex flex-col lg:flex-row">

            {/* Left Sidebar - Payment Methods */}
            <div className="lg:w-[380px] bg-gradient-to-b from-[#1e293b] to-[#334155] p-8 relative">
              <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-[#f59e0b]/20 to-transparent rounded-full blur-3xl" />
              <div className="absolute bottom-0 left-0 w-24 h-24 bg-gradient-to-tr from-[#10b981]/20 to-transparent rounded-full blur-2xl" />

              <div className="relative z-10">
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
                        onClick={() => handlePaymentMethodSelect(method.key as any)}
                        className={`w-full p-4 rounded-xl text-left transition-all duration-300 group relative overflow-hidden ${isSelected
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
              <button
                onClick={() => router.back()}
                className="absolute top-4 right-4 text-gray-400 hover:text-gray-600 text-2xl font-bold transition-colors"
                aria-label="Close"
              >
                ×
              </button>

              <div className="flex flex-col h-full max-w-md mx-auto">
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

                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 mb-4">
                    <p className="text-blue-800 text-sm font-medium">{getMethodDescription()}</p>
                  </div>
                </div>

                <div className="mb-8">
                  <p className="text-sm font-medium text-gray-700 mb-3">Select your payment provider:</p>
                  <div className="flex gap-3 flex-wrap">
                    {getPaymentIcons().map((icon, index) => {
                      const isSelected = selectedProvider === icon.alt;
                      return (
                        <button
                          key={index}
                          onClick={() => setSelectedProvider(icon.alt)}
                          className={`
            w-20 h-14 rounded-lg border-2 transition-all duration-300 
            flex items-center justify-center p-1 relative
            ${isSelected
                              ? 'bg-orange-100 border-orange-400 shadow-lg scale-[1.05]'
                              : 'bg-white border-gray-200 hover:border-orange-300 hover:shadow-md'
                            }
            group
          `}
                        >
                          {/* Glow effect for selected state */}
                          {isSelected && (
                            <div className="absolute inset-0 rounded-lg bg-orange-100/30 blur-[6px] -z-10" />
                          )}

                          {/* Payment provider logo */}
                          <div className={`
            w-full h-full flex items-center justify-center
            transition-transform duration-300
            ${isSelected ? 'scale-110' : 'group-hover:scale-105'}
          `}>
                            <Image
                              src={icon.src}
                              alt={icon.alt}
                              width={56}
                              height={40}
                              className={`
                w-full h-full object-contain
                ${!isSelected && 'opacity-80 group-hover:opacity-100'}
                transition-opacity duration-200
              `}
                            />
                          </div>

                          {/* Selection indicator */}
                          {isSelected && (
                            <div className="absolute inset-0 overflow-hidden rounded-lg">
                              <div className="
      absolute -inset-y-full -inset-x-8 bg-gradient-to-r 
      from-transparent via-white/30 to-transparent 
      transform rotate-12 animate-shimmer
    " />
                            </div>
                          )}

                          {/* Hover overlay */}
                          <div className={`
            absolute inset-0 rounded-lg bg-gradient-to-br from-white/50 to-transparent
            opacity-0 group-hover:opacity-100 transition-opacity duration-200
            ${isSelected && 'hidden'}
          `} />
                        </button>
                      );
                    })}
                  </div>
                </div>
                <div className="relative mb-8">
                  <label className="block text-gray-700 text-sm font-semibold mb-3">
                    Phone Number
                  </label>
                  <div className="flex items-stretch h-14">
                    {/* Country Code Selector */}
                    <div className="relative flex-1 max-w-[120px] h-full">
                      <select
                        className="w-full h-full border-2 border-gray-200 rounded-l-xl px-4 focus:border-[#f59e0b] focus:outline-none appearance-none pl-12 pr-8"
                        value="+251"
                        onChange={() => { }}
                        disabled={!selectedProvider}
                      >
                        <option value="+251">+251</option>
                      </select>
                      <div className="absolute left-3 top-1/2 transform -translate-y-1/2 flex items-center">
                        <Image
                          src="/ethiopia-flag.svg"
                          alt="Ethiopia"
                          width={24}
                          height={16}
                          className="w-6 h-4 mr-2"
                        />
                      </div>
                      <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-500" />
                    </div>

                    {/* Phone Number Input */}
                    <input
                      type="tel"
                      placeholder={
                        selectedProvider
                          ? (selectedProvider === 'M-Pesa' || selectedProvider === 'mpesa'
                            ? '710032610'
                            : '965822675')
                          : 'Select payment method first'
                      }
                      value={phoneNumber}
                      className={`flex-1 h-full px-4 border-2 border-l-0 border-gray-200 rounded-r-xl text-gray-900 placeholder-gray-400 focus:border-[#f59e0b] focus:outline-none text-lg ${!selectedProvider ? 'bg-gray-100 cursor-not-allowed' : ''
                        }`}
                      maxLength={9}
                      disabled={!selectedProvider}
                      onChange={(e) => {
                        if (!selectedProvider) return;

                        const value = e.target.value;
                        // Only allow numbers
                        if (!/^\d*$/.test(value)) return;

                        if (value.length > 0) {
                          const firstDigit = value[0];
                          const isMpesa = selectedProvider === 'M-Pesa' || selectedProvider === 'mpesa';
                          const isValidFirstDigit = isMpesa ? firstDigit === '7' : firstDigit === '9';

                          if (!isValidFirstDigit) {
                            setShowDigitAlert(true);
                            setTimeout(() => setShowDigitAlert(false), 3000);
                            return;
                          } else {
                            setShowDigitAlert(false);
                          }
                        }

                        // Prevent repeating digits
                        if (/^(\d)\1{8}$/.test(value)) return;

                        setPhoneNumber(value);
                      }}
                    />
                  </div>
                  {!selectedProvider && (
                    <p className="text-gray-500 text-xs mt-1">
                      Please select a payment method first
                    </p>
                  )}
                  {showDigitAlert && (
                    <div className="absolute top-full left-0 mt-1 bg-red-100 border border-red-400 text-red-700 px-3 py-1 rounded text-xs animate-bounce">
                      {selectedProvider === 'M-Pesa' || selectedProvider === 'mpesa'
                        ? 'M-Pesa numbers must start with 7'
                        : 'Phone numbers must start with 9'}
                    </div>
                  )}
                  {phoneNumber.length > 0 && !isValidPhone && selectedProvider && (
                    <p className="text-red-500 text-xs mt-1">
                      {selectedProvider === 'M-Pesa' || selectedProvider === 'mpesa'
                        ? 'Please enter a valid M-Pesa number starting with 7'
                        : 'Please enter a valid number starting with 9'}
                    </p>
                  )}
                </div>
                <div className="bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-4 mb-4 border border-gray-200">
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-600">Item Subtotal</span>
                    <span className="font-medium text-gray-800">ETB 1,000.00</span>
                  </div>
                  {tipAmount > 0 && (
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-gray-600">Tip Amount</span>
                      <span className="font-medium text-orange-600">+ETB {tipAmount.toFixed(2)}</span>
                    </div>
                  )}
                  <div className="border-t border-gray-200 mt-2 pt-2 flex justify-between">
                    <span className="text-base font-semibold text-gray-800">Total</span>
                    <span className="text-base font-bold text-orange-600">ETB {totalAmount.toFixed(2)}</span>
                  </div>
                </div>

                <div className="mb-6">
                  <button
                    onClick={() => setShowTipOptions(!showTipOptions)}
                    className="text-sm text-orange-600 hover:text-orange-700 font-medium flex items-center gap-1"
                  >
                    {tipAmount > 0 ? 'Change tip amount' : 'Add a tip?'}
                    {tipAmount > 0 && <CheckCircle className="w-4 h-4" />}
                  </button>

                  {showTipOptions && (
                    <div className="mt-3 flex gap-2">
                      {[50, 100, 200].map((amount) => (
                        <button
                          key={amount}
                          onClick={() => handleTipSelect(amount)}
                          className={`px-3 py-2 rounded-lg text-sm font-medium ${tipAmount === amount
                            ? 'bg-orange-100 text-orange-700 border border-orange-300'
                            : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                            }`}
                        >
                          ETB {amount}
                        </button>
                      ))}
                      <button
                        onClick={() => handleTipSelect(0)}
                        className={`px-3 py-2 rounded-lg text-sm font-medium ${tipAmount === 0
                          ? 'bg-orange-100 text-orange-700 border border-orange-300'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                          }`}
                      >
                        No tip
                      </button>
                    </div>
                  )}
                </div>

                <button
                  onClick={handlePaymentSubmit}
                  className={`w-full bg-gradient-to-r from-[#f59e0b] to-[#d97706] text-white font-bold py-4 rounded-xl transition-all duration-200 shadow-lg hover:shadow-xl transform hover:scale-[1.02] mb-8 text-lg ${!selectedProvider || !isValidPhone ? 'opacity-50 cursor-not-allowed' : 'hover:from-[#d97706] hover:to-[#b45309]'
                    }`}
                  disabled={!selectedProvider || !isValidPhone || isProcessing}
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

        <div className="mt-6 bg-white/5 backdrop-blur-sm rounded-xl p-4 border border-white/10">
          <div className="flex items-center justify-center gap-8 text-gray-300 text-xs">
            <div className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              <span>High-level security</span>
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

      {/* Payment Processing Modal */}
 {showPaymentModal && (
  <div className="fixed inset-0 bg-black/80 backdrop-blur-xl z-50 flex items-center justify-center p-4">
    {/* Particle background canvas */}
    <div className="absolute inset-0 opacity-20">
      {[...Array(8)].map((_, i) => (
        <div
          key={i}
          className="absolute text-yellow-300 opacity-70"
          style={{
            fontSize: `${Math.random() * 20 + 10}px`,
            top: `${Math.random() * 100}%`,
            left: `${Math.random() * 100}%`,
            animation: `coinFloat ${Math.random() * 15 + 10}s linear infinite`,
            animationDelay: `${Math.random() * 5}s`,
            transform: `rotate(${Math.random() * 360}deg)`
          }}
        >
          {['₦','€','$','¥','£','₹','₩','₱'][i % 8]}
        </div>
      ))}
    </div>
    
    {/* Main modal container */}
    <div className="relative bg-white rounded-3xl shadow-2xl w-full max-w-md overflow-hidden border-2 border-orange-200/50">
      {/* Animated gradient header */}
      <div className="relative bg-gradient-to-br from-[#f59e0b] via-[#e67e22] to-[#d97706] p-6 overflow-hidden">
        {/* Holographic stripe effect */}
        <div className="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-transparent via-white/70 to-transparent animate-holographic"></div>
        
        <div className="relative z-10 flex justify-between items-center">
          <div className="flex items-center gap-3">
            {/* 3D logo with shine */}
            <div className="relative">
              <div className="bg-white p-2 rounded-lg shadow-lg">
                <div className="bg-gradient-to-br from-white to-gray-50 p-1 rounded overflow-hidden relative">
                  <Image 
                    src="/socialpay.webp"
                    alt="SocialPay"
                    width={80}
                    height={24}
                    className="h-8 w-auto relative z-10"
                  />
                </div>
              </div>
            </div>
            <div>
              <h3 className="text-xl font-bold text-white drop-shadow-lg">Payment Processing</h3>
              <p className="text-orange-100/90 text-sm font-medium">Securing your transaction</p>
            </div>
          </div>
          
          {/* Close button */}
         <button 
  onClick={() => {
    setShowPaymentModal(false);
    setIsProcessing(false);
  }}
  className="p-1 rounded-full bg-white/10 hover:bg-white/20 transition-all"
>
  <X className="w-5 h-5 text-white" />
</button>
        </div>
      </div>

      {/* Main content */}
      <div className="p-6 relative overflow-hidden">
        {/* Confetti burst effect */}
        {timeLeft < 5 && (
          <div className="absolute inset-0 overflow-hidden pointer-events-none">
            {[...Array(30)].map((_, i) => (
              <div
                key={i}
                className="absolute text-yellow-400 text-xl"
                style={{
                  animation: `confettiFall ${Math.random() * 3 + 2}s linear forwards`,
                  left: `${Math.random() * 100}%`,
                  top: `-10%`,
                  transform: `rotate(${Math.random() * 360}deg)`,
                  opacity: 0,
                  animationDelay: `${(timeLeft / 60) * 2}s`
                }}
              >
                {['✦','✧','✺','✤','❈'][i % 5]}
              </div>
            ))}
          </div>
        )}
        
        {/* Payment verification animation */}
        <div className="flex flex-col items-center justify-center mb-6">
          <div className="relative w-24 h-24 mb-4">
            {/* Animated verification circle */}
            <svg className="w-full h-full" viewBox="0 0 100 100">
              <circle
                cx="50"
                cy="50"
                r="45"
                fill="none"
                stroke="#f3f4f6"
                strokeWidth="8"
              />
              <circle
                cx="50"
                cy="50"
                r="45"
                fill="none"
                stroke="#f59e0b"
                strokeWidth="8"
                strokeDasharray="283"
                strokeDashoffset={(283 * (1 - (60-timeLeft)/60))}
                transform="rotate(-90 50 50)"
                className="transition-all duration-300 ease-out"
                strokeLinecap="round"
              />
            </svg>
            
            {/* Checkmark center */}
            <div className="absolute inset-0 flex items-center justify-center">
              <CheckCircle className="w-10 h-10 text-orange-500" />
            </div>
          </div>
          
          <h4 className="text-xl font-bold text-gray-800 mb-1 text-center">
            {timeLeft > 5 ? 'Processing Payment' : 'Payment Successful!'}
          </h4>
          <p className="text-gray-600 text-center">
            {timeLeft > 5 
              ? `Completing your ETB ${totalAmount.toFixed(2)} transaction`
              : 'Your payment has been processed successfully'}
          </p>
        </div>

        {/* Progress tracker */}
        <div className="mb-6">
          <div className="flex justify-between items-center mb-2">
            <span className="text-sm text-gray-600">Progress</span>
            <span className="text-sm font-medium text-orange-600">
              {Math.floor((60-timeLeft)/60*100)}% Complete
            </span>
          </div>
          <div className="relative h-2 bg-gray-200 rounded-full overflow-hidden">
            <div 
              className="absolute top-0 left-0 h-full bg-gradient-to-r from-orange-400 to-orange-500 rounded-full transition-all duration-300 ease-out"
              style={{ width: `${(60-timeLeft)/60*100}%` }}
            />
          </div>
          <p className="text-right text-xs text-gray-500 mt-1">
            About {timeLeft} seconds remaining
          </p>
        </div>

        {/* Collapsible Payment Details with enhanced visual cues */}
  <div className="mb-6 relative">
  {/* Animated hand pointer that disappears after first interaction */}
{!hasInteractedWithDetails && (
  <div className="absolute -left-2 top-0 z-10 animate-bounce pointer-events-none">
    <div className="relative">
      <Pointer className="w-6 h-6 text-orange-500 rotate-[135deg]" />
      <div className="absolute inset-0 bg-orange-500 rounded-full animate-ping opacity-30"></div>
    </div>
  </div>
)}
  
  <button 
    onClick={() => {
      setIsDetailsExpanded(!isDetailsExpanded);
      setHasInteractedWithDetails(true);
    }}
    className="w-full flex justify-between items-center p-3 bg-gray-50 rounded-lg border border-gray-200 hover:bg-gray-100 transition-colors relative group"
  >
            <div className="flex items-center gap-2">
              <div className="relative">
                <CreditCard className="w-4 h-4 text-gray-600" />
                {/* Pulse animation on icon for attention */}
                {!hasInteractedWithDetails && (
                  <div className="absolute inset-0 rounded-full bg-orange-400/20 animate-ping"></div>
                )}
              </div>
              <span className="font-medium text-gray-700">Payment Details</span>
            </div>
            <div className="flex items-center">
              {/* Bouncing arrow animation */}
              <ChevronDown className={`w-4 h-4 text-gray-500 transition-transform ${isDetailsExpanded ? 'rotate-180' : 'group-hover:translate-y-0.5'}`} /> 
            </div>
          </button>

          {isDetailsExpanded && (
            <div className="mt-2 bg-white rounded-lg p-4 border border-gray-200 animate-fadeIn">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-gray-500">Amount</p>
                  <p className="font-medium">ETB {totalAmount.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Method</p>
                  <p className="font-medium">{selectedProvider}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Recipient</p>
                  <p className="font-medium">SocialPay</p>
                </div>
                <div>
                  <p className="text-xs text-gray-500">Status</p>
                  <p className="font-medium flex items-center gap-1">
                    <span className={`w-2 h-2 rounded-full ${
                      timeLeft > 5 ? 'bg-yellow-500' : 'bg-green-500'
                    }`}></span>
                    {timeLeft > 5 ? 'Processing' : 'Completed'}
                  </p>
                </div>
                <div className="col-span-2">
                  <p className="text-xs text-gray-500">Transaction ID</p>
                  <p className="font-mono text-sm">SP-{Math.random().toString(36).substring(2, 10).toUpperCase()}</p>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Security reassurance */}
        <div className="flex items-center justify-center gap-2 text-sm text-blue-600 bg-blue-50/50 p-3 rounded border border-blue-200">
          <Shield className="w-4 h-4" />
          <span>Payment secured with 256-bit encryption</span>
        </div>

        {/* Completion message */}
        {timeLeft < 5 && (
          <div className="mt-4 flex justify-center">
            <div className="flex items-center gap-2 text-sm bg-green-50 text-green-700 px-3 py-1.5 rounded-full border border-green-200">
              <CheckCircle className="w-4 h-4" />
              <span>Almost done! Redirecting shortly...</span>
            </div>
          </div>
        )}
      </div>
    </div>
  </div>
)}
    </div>
  )
}