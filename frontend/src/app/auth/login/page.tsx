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
  const [step, setStep] = useState<'phone' | 'otp'>('phone')
  const [phoneNumber, setPhoneNumber] = useState('')
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

  const handlePhoneSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

    if (!validatePhoneNumber(phoneNumber)) {
      setError('Phone number must start with 07, 09 (10 digits) or 7, 9 (9 digits)')
      return
    }

    setIsLoading(true)

    try {
      // Normalize phone number before proceeding
      const normalizedPhone = normalizePhoneNumber(phoneNumber).replace(/\s/g, '') // Remove any spaces
      console.log('Original phone:', phoneNumber)
      console.log('Normalized phone:', normalizedPhone)
      
      // Initialize pre-session for login with phone number
      const preSessionResponse = await authAPI.initPreSessionForLogin({
        prefix: '', // Will be overridden by API to +251
        number: normalizedPhone
      })
      
      if (preSessionResponse.success) {
        setOtpToken(preSessionResponse.data?.token || 'temp-token')
        setPhoneNumber(normalizedPhone) // Update to use normalized phone number
        setStep('otp')
        setError('')
      } else {
        setError('Failed to initialize login. Please try again.')
      }
    } catch (err: any) {
      console.log('Login error:', err)
      const errorMessage = err.response?.data?.error?.message || 'Login initialization failed. Please try again.'
      setError(errorMessage)
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
      const response = await authAPI.signIn(otpToken, otpCode, {
        prefix: '', // Will be overridden by API to +251
        number: phoneNumber
      })

      if (response.success) {
        // Store auth tokens
        localStorage.setItem('authToken', response.data.token.active)
        localStorage.setItem('refreshToken', response.data.token.refresh)
        
        // Update auth store
        const user = {
          id: response.data.user.id,
          name: `${response.data.user.first_name} ${response.data.user.last_name}`,
          email: '', // Phone-based auth might not have email
          role: response.data.user.user_type as 'admin' | 'merchant' | 'user',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
        
        login(user, response.data.token.active)
        router.push('/dashboard')
      } else {
        setError(response.error?.message || 'Invalid verification code')
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'Login failed. Please try again.'
      setError(errorMessage)
    } finally {
      setIsLoading(false)
    }
  }

  const resendOTP = async () => {
    setIsLoading(true)
    try {
      // Re-initialize with phone number to get new OTP
      const preSessionResponse = await authAPI.initPreSessionForLogin({
        prefix: '', // Will be overridden by API to +251
        number: phoneNumber
      })
      if (preSessionResponse.success) {
        setOtpToken(preSessionResponse.data?.token || 'temp-token')
        setOtpCode('')
        setError('')
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'Failed to resend code. Please try again.'
      setError(errorMessage)
    } finally {
      setIsLoading(false)
    }
  }

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

              <form onSubmit={handleOTPSubmit} className="space-y-4">
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
                  disabled={isLoading}
                  className="group w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
                >
                  {isLoading ? (
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      <span>Verifying...</span>
                    </div>
                  ) : (
                    <div className="flex items-center justify-center">
                      <span>Sign In</span>
                      <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                    </div>
                  )}
                </button>

                <div className="text-center">
                  <p className="text-sm text-gray-600 mb-2">Didn't receive the code?</p>
                  <button
                    type="button"
                    onClick={resendOTP}
                    disabled={isLoading}
                    className="text-sm font-semibold text-brand-green-600 hover:text-brand-green-500 transition-colors hover:underline"
                  >
                    Resend Code
                  </button>
                </div>
              </form>

              <div className="mt-6 text-center">
                <button
                  onClick={() => setStep('phone')}
                  className="text-sm font-semibold text-gray-600 hover:text-gray-500 transition-colors"
                >
                  ← Back to Phone Number
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

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
                Sign in with your phone number
              </p>
            </div>

            {error && (
              <div className="mb-4 p-3 bg-red-50/80 backdrop-blur-sm border border-red-200/60 rounded-xl">
                <p className="text-red-600 text-sm font-medium text-center">{error}</p>
              </div>
            )}

            <form onSubmit={handlePhoneSubmit} className="space-y-4">
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
                        
                        // Limit based on starting digits
                        if (value.startsWith('07') || value.startsWith('09')) {
                          // For 07/09 formats, limit to 10 characters
                          if (value.length > 10) {
                            value = value.slice(0, 10)
                          }
                        } else if (value.startsWith('7') || value.startsWith('9')) {
                          // For 7/9 formats, limit to 9 characters
                          if (value.length > 9) {
                            value = value.slice(0, 9)
                          }
                        } else {
                          // For other formats, limit to 10 characters to allow typing
                          if (value.length > 10) {
                            value = value.slice(0, 10)
                          }
                        }
                        
                        setPhoneNumber(value)
                      }}
                      className="flex-1 px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-r-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm"
                      placeholder="0911123456"
                    />
                  </div>
                </div>
              </div>

              <button
                type="submit"
                disabled={isLoading || !phoneNumber}
                className="group w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    <span>Sending code...</span>
                  </div>
                ) : (
                  <div className="flex items-center justify-center">
                    <span>Send Verification Code</span>
                    <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                  </div>
                )}
              </button>

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-200/60" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-3 bg-white/80 text-gray-500 font-medium">or continue with</span>
                </div>
              </div>

              <button
                type="button"
                className="w-full flex items-center justify-center px-4 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl hover:bg-white/80 hover:border-gray-300/60 font-semibold text-gray-700 transition-all duration-200 group text-sm"
              >
                <svg className="w-4 h-4 mr-2" viewBox="0 0 24 24">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                <span>Continue with Google</span>
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