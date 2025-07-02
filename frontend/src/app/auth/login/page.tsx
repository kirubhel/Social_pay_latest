'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { useAuthStore } from '@/stores/auth'
import { EyeIcon, EyeSlashIcon, ArrowRightIcon, SparklesIcon, LockClosedIcon } from '@heroicons/react/24/outline'
import { CheckCircleIcon, ShieldCheckIcon, CreditCardIcon, StarIcon } from '@heroicons/react/24/solid'
import { authAPI } from '@/lib/api'

export default function LoginPage() {
  const [step, setStep] = useState<'login' | 'otp'>('login')
  const [phoneNumber, setPhoneNumber] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [otpCode, setOtpCode] = useState('')
  const [otpToken, setOtpToken] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const router = useRouter()
  const login = useAuthStore((state) => state.login)

  // Function to normalize phone number (remove leading 0 if present)
  const normalizePhoneNumber = (phoneNumber: string) => {
    if ((phoneNumber.startsWith('07') || phoneNumber.startsWith('09')) && phoneNumber.length === 10) {
      return phoneNumber.slice(1) // Remove the leading 0
    }
    return phoneNumber
  }

  const validatePhoneNumber = (phone: string) => {
    const normalizedPhone = normalizePhoneNumber(phone)
    // Accept 07xxxxxxxx, 09xxxxxxxx (10 digits) or 7xxxxxxxx, 9xxxxxxxx (9 digits)
    const isValid10Digit = /^(07|09)\d{8}$/.test(phone)
    const isValid9Digit = /^[79]\d{8}$/.test(phone)
    const isNormalizedValid = /^[79]\d{8}$/.test(normalizedPhone)
    return (isValid10Digit || isValid9Digit) && isNormalizedValid
  }

  const handleLoginSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (!validatePhoneNumber(phoneNumber)) {
      setError('Phone number must start with 07, 09 (10 digits) or 7, 9 (9 digits)')
      return
    }
    if (!password) {
      setError('Password is required')
      return
    }
    setIsLoading(true)
    try {
      const normalizedPhone = normalizePhoneNumber(phoneNumber).replace(/\s/g, '')
      // Call backend password check API
      const response = await authAPI.login({
        prefix: '',
        number: normalizedPhone,
        password,
      })
     
      if (response.success) {
         if (response.data?.next_step === 'OTP_REQUIRED' || response.data?.next_step === 'CHECK_OTP') {
         
          setOtpToken(response.data.token)
          setStep('otp')
        } else {
          await login(response.data.user, response.data.token)
          router.push('/dashboard')
        }
      } else {
        setError(response.error?.message || 'Login failed. Please try again.')
      }
    } catch (err: any) {
      console.error('Login error:', err)
      setError('An unexpected error occurred. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleOTPSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (!otpCode.trim()) {
      setError('Please enter the verification code')
      return
    }
    setIsLoading(true)
    try {
      const response = await authAPI.verifyOTP(otpToken, otpCode)
      if (response.success) {
        await login(response.data.user, response.data.token.active)
        router.push('/dashboard')
      } else {
        setError(response.error?.message || 'Invalid verification code')
      }
    } catch (err: any) {
      console.error('OTP verification error:', err)
      setError('An unexpected error occurred. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  // UI rendering
  if (step === 'otp') {
    return (
      <div className="min-h-screen relative overflow-hidden">
        {/* Animated Background */}
        <div className="absolute inset-0 bg-gradient-to-br from-brand-green-50 via-white to-brand-gold-50">
          <div className="absolute top-0 left-0 w-full h-full">
            <div className="absolute top-20 left-20 w-64 h-64 bg-brand-green-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
            <div className="absolute top-32 right-20 w-64 h-64 bg-brand-gold-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
            <div className="absolute -bottom-8 left-1/2 w-64 h-64 bg-brand-green-300/10 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
          </div>
        </div>
        {/* OTP Verification Container */}
        <div className="relative z-10 flex min-h-screen items-center justify-center p-4">
          <div className="w-full max-w-md">
            <div className="bg-white/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 p-6">
              <div className="text-center mb-6">
                <Link href="/" className="inline-block">
                  <Image
                    src="/logo.png"
                    alt="Social Pay Logo"
                    width={200}
                    height={10}
                    className="object-contain mx-auto mb-4 cursor-pointer hover:opacity-80 transition-opacity"
                    priority
                  />
                </Link>
                <h2 className="text-xl font-bold text-gray-900 mb-1">
                  Enter Verification Code
                </h2>
                <p className="text-gray-600 text-sm">
                  We've sent a verification code to +251{phoneNumber}
                </p>
              </div>
              {error && (
                <div className="mb-4 p-3 bg-red-50/80 backdrop-blur-sm border border-red-200/60 rounded-xl">
                  <p className="text-red-600 text-sm font-medium text-center">{error}</p>
                </div>
              )}
              <form onSubmit={handleOTPSubmit} className="space-y-4" noValidate>
                <div className="group">
                  <label htmlFor="otpCode" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Verification Code
                  </label>
                  <input
                    id="otpCode"
                    name="otpCode"
                    type="text"
                    maxLength={6}
                    required
                    value={otpCode}
                    onChange={(e) => setOtpCode(e.target.value.replace(/\D/g, ''))}
                    className="w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-center text-lg tracking-widest"
                    placeholder="123456"
                  />
                </div>
                <button
                  type="submit"
                  disabled={isLoading || !otpCode.trim()}
                  className="group w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                >
                  {isLoading ? (
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      <span>Verifying...</span>
                    </div>
                  ) : (
                    <div className="flex items-center justify-center">
                      <span>Verify</span>
                      <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                    </div>
                  )}
                </button>
              </form>
              <div className="mt-6 text-center">
                <button
                  onClick={() => setStep('login')}
                  className="text-sm font-semibold text-gray-600 hover:text-gray-500 transition-colors"
                >
                  ← Back to Login
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Login form UI
  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Animated Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-brand-green-50 via-white to-brand-gold-50">
        <div className="absolute top-0 left-0 w-full h-full">
          <div className="absolute top-20 left-20 w-64 h-64 bg-brand-green-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
          <div className="absolute top-32 right-20 w-64 h-64 bg-brand-gold-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
          <div className="absolute -bottom-8 left-1/2 w-64 h-64 bg-brand-green-300/10 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
        </div>
      </div>
      {/* Centered Login Container */}
      <div className="relative z-10 flex min-h-screen items-center justify-center p-4">
        <div className="w-full max-w-md">
          {/* Logo */}
          <div className="text-center mb-8">
            <h1 className="text-2xl font-bold text-gray-900 mb-2">
              Welcome to{' '}
              <span className="bg-gradient-to-r from-brand-green-600 to-brand-gold-500 bg-clip-text text-transparent">
                Social Pay
              </span>
            </h1>
            <p className="text-gray-600 text-sm">
              Secure payment solutions for modern businesses
            </p>
          </div>
          {/* Form Container */}
          <div className="bg-white/80 backdrop-blur-lg rounded-2xl shadow-2xl border border-white/20 p-6">
            <div className="text-center mb-6">
              <Link href="/" className="inline-block">
                <Image
                  src="/logo.png"
                  alt="Social Pay Logo"
                  width={200}
                  height={10}
                  className="object-contain mx-auto mb-4 cursor-pointer hover:opacity-80 transition-opacity"
                  priority
                />
              </Link>
              <h2 className="text-xl font-bold text-gray-900 mb-1">
                Welcome Back
              </h2>
              <p className="text-gray-600 text-sm">
                Sign in with your phone number and password
              </p>
            </div>
            {error && (
              <div className="mb-4 p-3 bg-red-50/80 backdrop-blur-sm border border-red-200/60 rounded-xl">
                <p className="text-red-600 text-sm font-medium text-center">{error}</p>
              </div>
            )}
            <form onSubmit={handleLoginSubmit} className="space-y-4" noValidate>
              <div className="space-y-3">
                <div className="group">
                  <label htmlFor="phoneNumber" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Phone Number
                  </label>
                  <div className="flex">
                    <div className="flex items-center px-3 py-2.5 bg-gray-100 border border-r-0 border-gray-200/60 rounded-l-xl">
                      <span className="text-sm mr-1">🇪🇹</span>
                      <span className="text-xs text-gray-600 font-medium">+251</span>
                    </div>
                    <input
                      id="phoneNumber"
                      name="phoneNumber"
                      type="tel"
                      required
                      value={phoneNumber}
                      onChange={(e) => {
                        let value = e.target.value.replace(/\D/g, '')
                        if (value.startsWith('07') || value.startsWith('09')) {
                          if (value.length > 10) value = value.slice(0, 10)
                        } else if (value.startsWith('7') || value.startsWith('9')) {
                          if (value.length > 9) value = value.slice(0, 9)
                        } else {
                          if (value.length > 10) value = value.slice(0, 10)
                        }
                        setPhoneNumber(value)
                      }}
                      className="flex-1 px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-r-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm"
                      placeholder="0911123456"
                    />
                  </div>
                </div>
                <div className="group">
                  <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Password
                  </label>
                  <div className="relative">
                    <input
                      id="password"
                      name="password"
                      type={showPassword ? "text" : "password"}
                      required
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm pr-10"
                      placeholder="Your password"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none"
                    >
                      {showPassword ? (
                        <EyeSlashIcon className="h-5 w-5" />
                      ) : (
                        <EyeIcon className="h-5 w-5" />
                      )}
                    </button>
                  </div>
                </div>
              </div>
              <button
                type="submit"
                disabled={isLoading || !phoneNumber || !password}
                className="group w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    <span>Signing in...</span>
                  </div>
                ) : (
                  <div className="flex items-center justify-center">
                    <span>Sign In</span>
                    <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                  </div>
                )}
              </button>
            </form>
            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{' '}
                <Link 
                  href="/auth/register" 
                  className="font-semibold text-brand-green-600 hover:text-brand-green-500 transition-colors hover:underline"
                >
                  Sign up
                </Link>
              </p>
            </div>
          </div>
          {/* Features Section */}
          <div className="mt-6 space-y-4">
            <div className="flex justify-center space-x-6 text-center">
              {[
                { icon: ShieldCheckIcon, text: 'Secure' },
                { icon: CreditCardIcon, text: 'Fast' },
                { icon: CheckCircleIcon, text: 'Trusted' }
              ].map((feature, index) => (
                <div key={index} className="flex flex-col items-center">
                  <div className="w-8 h-8 bg-gradient-to-r from-brand-green-500 to-brand-gold-400 rounded-lg flex items-center justify-center mb-1">
                    <feature.icon className="w-4 h-4 text-white" />
                  </div>
                  <span className="text-xs text-gray-600 font-medium">{feature.text}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 