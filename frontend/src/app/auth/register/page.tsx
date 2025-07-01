'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { EyeIcon, EyeSlashIcon, ArrowRightIcon, UserPlusIcon } from '@heroicons/react/24/outline'
import { CheckCircleIcon, ShieldCheckIcon, CreditCardIcon, StarIcon } from '@heroicons/react/24/solid'
import { authAPI } from '@/lib/api'

export default function RegisterPage() {
  const [step, setStep] = useState<'register' | 'verify-otp'>('register')
  const [otpToken, setOtpToken] = useState('')
  const [phoneData, setPhoneData] = useState({ prefix: '', number: '' })
  
  const [formData, setFormData] = useState({
    title: 'Mr',
    firstName: '',
    lastName: '',
    phoneNumber: '',
    businessName: '',
    password: '',
    confirmPassword: '',
    passwordHint: '',
    agreedToTerms: false
  })
  
  const [otpCode, setOtpCode] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const router = useRouter()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target
    const checked = (e.target as HTMLInputElement).checked
    
    let processedValue = value
    
    // Handle phone number input - allow only digits and process leading zero
    if (name === 'phoneNumber') {
      // Remove any non-digit characters
      processedValue = value.replace(/\D/g, '')
      
      // Limit based on starting digits
      if (processedValue.startsWith('07') || processedValue.startsWith('09')) {
        // For 07/09 formats, limit to 10 characters
        if (processedValue.length > 10) {
          processedValue = processedValue.slice(0, 10)
        }
      } else if (processedValue.startsWith('7') || processedValue.startsWith('9')) {
        // For 7/9 formats, limit to 9 characters
        if (processedValue.length > 9) {
          processedValue = processedValue.slice(0, 9)
        }
      } else {
        // For other formats, limit to 10 characters to allow typing
        if (processedValue.length > 10) {
          processedValue = processedValue.slice(0, 10)
        }
      }
    }
    
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : processedValue
    }))
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }))
    }
  }

  // Function to normalize phone number (remove leading 0 if present)
  const normalizePhoneNumber = (phoneNumber: string) => {
    if ((phoneNumber.startsWith('07') || phoneNumber.startsWith('09')) && phoneNumber.length === 10) {
      return phoneNumber.slice(1) // Remove the leading 0
    }
    return phoneNumber
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.title.trim()) newErrors.title = 'Title is required'
    if (!formData.firstName.trim()) newErrors.firstName = 'First name is required'
    if (!formData.lastName.trim()) newErrors.lastName = 'Last name is required'
    if (!formData.phoneNumber.trim()) newErrors.phoneNumber = 'Phone number is required'
    if (!formData.businessName.trim()) newErrors.businessName = 'Business name is required'
    if (!formData.password) newErrors.password = 'Password is required'
    if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = 'Passwords do not match'
    if (!formData.agreedToTerms) newErrors.agreedToTerms = 'You must agree to the terms'

    // Validate Ethiopian phone number format
    if (formData.phoneNumber) {
      const normalizedPhone = normalizePhoneNumber(formData.phoneNumber)
      
      // Accept 07xxxxxxxx, 09xxxxxxxx (10 digits) or 7xxxxxxxx, 9xxxxxxxx (9 digits)
      const isValid10Digit = /^(07|09)\d{8}$/.test(formData.phoneNumber)
      const isValid9Digit = /^[79]\d{8}$/.test(formData.phoneNumber)
      const isNormalizedValid = /^[79]\d{8}$/.test(normalizedPhone)
      
      if (!(isValid10Digit || isValid9Digit) || !isNormalizedValid) {
        newErrors.phoneNumber = 'Phone number must start with 07, 09 (10 digits) or 7, 9 (9 digits)'
      }
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return
    
    setIsLoading(true)

    try {
      // Normalize phone number before sending to backend
      const normalizedPhone = normalizePhoneNumber(formData.phoneNumber)
      
      const response = await authAPI.signUp({
        title: formData.title,
        first_name: formData.firstName,
        last_name: formData.lastName,
        phone_prefix: '+251',
        phone_number: normalizedPhone,
        password: formData.password,
        password_hint: formData.passwordHint,
        confirm_password: formData.confirmPassword
      })

      if (response.success) {
        setOtpToken(response.data.auth.token)
        setPhoneData({ prefix: '+251', number: normalizedPhone })
        setStep('verify-otp')
      } else {
        setErrors({ general: response.error?.message || 'Registration failed' })
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'Registration failed. Please try again.'
      setErrors({ general: errorMessage })
    } finally {
      setIsLoading(false)
    }
  }

  const handleOTPSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!otpCode.trim()) {
      setErrors({ otp: 'OTP code is required' })
      return
    }
    
    setIsLoading(true)

    try {
      const response = await authAPI.verifyOTP(otpToken, otpCode)

      if (response.success) {
        // Store auth tokens
        localStorage.setItem('authToken', response.data.token.active)
        localStorage.setItem('refreshToken', response.data.token.refresh)
        
        // Redirect to dashboard
        router.push('/dashboard?message=Registration successful! Welcome to Social Pay.')
      } else {
        setErrors({ otp: response.error?.message || 'Invalid OTP code' })
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error?.message || 'OTP verification failed. Please try again.'
      setErrors({ otp: errorMessage })
    } finally {
      setIsLoading(false)
    }
  }

  const resendOTP = async () => {
    setIsLoading(true)
    try {
      // Normalize phone number before sending to backend
      const normalizedPhone = normalizePhoneNumber(formData.phoneNumber)
      
      // Re-register to get new OTP
      await authAPI.signUp({
        title: formData.title,
        first_name: formData.firstName,
        last_name: formData.lastName,
        phone_prefix: '+251',
        phone_number: normalizedPhone,
        password: formData.password,
        password_hint: formData.passwordHint,
        confirm_password: formData.confirmPassword
      })
      
      // Clear OTP field
      setOtpCode('')
      setErrors({})
    } catch (err) {
      setErrors({ otp: 'Failed to resend OTP. Please try again.' })
    } finally {
      setIsLoading(false)
    }
  }

  if (step === 'verify-otp') {
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
                  Verify Your Phone Number
                </h2>
                <p className="text-gray-600 text-sm">
                  We've sent a verification code to +251{phoneData.number}
                </p>
              </div>

              {errors.otp && (
                <div className="mb-4 p-3 bg-red-50/80 backdrop-blur-sm border border-red-200/60 rounded-xl">
                  <p className="text-red-600 text-sm font-medium text-center">{errors.otp}</p>
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
                    <span>Verify Code</span>
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
                  onClick={() => setStep('register')}
                  className="text-sm font-semibold text-gray-600 hover:text-gray-500 transition-colors"
                >
                  ← Back to Registration
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

      {/* Centered Register Container */}
      <div className="relative z-10 flex min-h-screen items-center justify-center p-4">
        <div className="w-full max-w-2xl">
          {/* Header Section */}
          <div className="text-center mb-8">
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
                Create Your Account
              </h2>
              <p className="text-gray-600 text-sm">
                Get started with your payment gateway today
              </p>
            </div>

            {errors.general && (
              <div className="mb-4 p-3 bg-red-50/80 backdrop-blur-sm border border-red-200/60 rounded-xl">
                <p className="text-red-600 text-sm font-medium text-center">{errors.general}</p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              {/* Title and First Name */}
              <div className="grid grid-cols-3 gap-3">
                <div className="group">
                  <label htmlFor="title" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Title
                  </label>
                  <select
                    id="title"
                    name="title"
                    required
                    value={formData.title}
                    onChange={handleChange}
                    className="w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 text-sm"
                  >
                    <option value="Mr">Mr</option>
                    <option value="Mrs">Mrs</option>
                    <option value="Ms">Ms</option>
                    <option value="Dr">Dr</option>
                  </select>
                  {errors.title && <p className="mt-1 text-xs text-red-600 font-medium">{errors.title}</p>}
                </div>
                <div className="group col-span-2">
                  <label htmlFor="firstName" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    First Name
                  </label>
                  <input
                    id="firstName"
                    name="firstName"
                    type="text"
                    required
                    value={formData.firstName}
                    onChange={handleChange}
                    className={`w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.firstName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                    placeholder="First name"
                  />
                  {errors.firstName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.firstName}</p>}
                </div>
              </div>

              {/* Last Name */}
              <div className="group">
                <label htmlFor="lastName" className="block text-sm font-semibold text-gray-700 mb-1.5">
                  Last Name
                </label>
                <input
                  id="lastName"
                  name="lastName"
                  type="text"
                  required
                  value={formData.lastName}
                  onChange={handleChange}
                  className={`w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.lastName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                  placeholder="Last name"
                />
                {errors.lastName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.lastName}</p>}
              </div>

              {/* Phone and Business Name */}
              <div className="grid grid-cols-2 gap-3">
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
                      value={formData.phoneNumber}
                      onChange={handleChange}
                      className={`flex-1 px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-r-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.phoneNumber ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                      placeholder="0911123456"
                    />
                  </div>
                  {errors.phoneNumber && <p className="mt-1 text-xs text-red-600 font-medium">{errors.phoneNumber}</p>}
                </div>
                <div className="group">
                  <label htmlFor="businessName" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Business Name
                  </label>
                  <input
                    id="businessName"
                    name="businessName"
                    type="text"
                    required
                    value={formData.businessName}
                    onChange={handleChange}
                    className={`w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.businessName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                    placeholder="Your business"
                  />
                  {errors.businessName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.businessName}</p>}
                </div>
              </div>

              {/* Password and Confirm Password */}
              <div className="grid grid-cols-2 gap-3">
                <div className="group">
                  <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Password
                  </label>
                  <div className="relative">
                    <input
                      id="password"
                      name="password"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="new-password"
                      required
                      value={formData.password}
                      onChange={handleChange}
                      className={`w-full px-3 py-2.5 pr-10 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.password ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                      placeholder="Password"
                    />
                    <button
                      type="button"
                      className="absolute inset-y-0 right-0 pr-3 flex items-center hover:scale-110 transition-transform duration-200"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {showPassword ? (
                        <EyeSlashIcon className="h-4 w-4 text-gray-400 hover:text-brand-green-500 transition-colors" />
                      ) : (
                        <EyeIcon className="h-4 w-4 text-gray-400 hover:text-brand-green-500 transition-colors" />
                      )}
                    </button>
                  </div>
                  {errors.password && <p className="mt-1 text-xs text-red-600 font-medium">{errors.password}</p>}
                </div>
                <div className="group">
                  <label htmlFor="confirmPassword" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Confirm Password
                  </label>
                  <div className="relative">
                    <input
                      id="confirmPassword"
                      name="confirmPassword"
                      type={showConfirmPassword ? 'text' : 'password'}
                      autoComplete="new-password"
                      required
                      value={formData.confirmPassword}
                      onChange={handleChange}
                      className={`w-full px-3 py-2.5 pr-10 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.confirmPassword ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                      placeholder="Confirm password"
                    />
                    <button
                      type="button"
                      className="absolute inset-y-0 right-0 pr-3 flex items-center hover:scale-110 transition-transform duration-200"
                      onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    >
                      {showConfirmPassword ? (
                        <EyeSlashIcon className="h-4 w-4 text-gray-400 hover:text-brand-green-500 transition-colors" />
                      ) : (
                        <EyeIcon className="h-4 w-4 text-gray-400 hover:text-brand-green-500 transition-colors" />
                      )}
                    </button>
                  </div>
                  {errors.confirmPassword && <p className="mt-1 text-xs text-red-600 font-medium">{errors.confirmPassword}</p>}
                </div>
              </div>

              {/* Password Hint (Optional) */}
              <div className="group">
                <label htmlFor="passwordHint" className="block text-sm font-semibold text-gray-700 mb-1.5">
                  Password Hint <span className="text-gray-400">(Optional)</span>
                </label>
                <input
                  id="passwordHint"
                  name="passwordHint"
                  type="text"
                  value={formData.passwordHint}
                  onChange={handleChange}
                  className="w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm"
                  placeholder="Something to help you remember your password"
                />
              </div>

              {/* Terms Agreement */}
              <div className="flex items-start space-x-3 pt-2">
                <input
                  id="agreedToTerms"
                  name="agreedToTerms"
                  type="checkbox"
                  checked={formData.agreedToTerms}
                  onChange={handleChange}
                  required
                  className="mt-1 h-4 w-4 text-brand-green-600 focus:ring-brand-green-500 border-gray-300 rounded transition-colors"
                />
                <label htmlFor="agreedToTerms" className="text-sm text-gray-600 leading-tight">
                  I agree to the{' '}
                  <Link href="/terms" className="text-brand-green-600 hover:text-brand-green-500 font-semibold hover:underline transition-colors">
                    Terms of Service
                  </Link>{' '}
                  and{' '}
                  <Link href="/privacy" className="text-brand-green-600 hover:text-brand-green-500 font-semibold hover:underline transition-colors">
                    Privacy Policy
                  </Link>
                </label>
              </div>
              {errors.agreedToTerms && <p className="mt-1 text-xs text-red-600 font-medium">{errors.agreedToTerms}</p>}

              <button
                type="submit"
                disabled={isLoading}
                className="group w-full bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white font-semibold py-2.5 px-6 rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    <span>Creating account...</span>
                  </div>
                ) : (
                  <div className="flex items-center justify-center">
                    <span>Create Account</span>
                    <ArrowRightIcon className="w-4 h-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                  </div>
                )}
              </button>

              

          
            </form>

            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600">
                Already have an account?{' '}
                <Link 
                  href="/auth/login" 
                  className="font-semibold text-brand-green-600 hover:text-brand-green-500 transition-colors hover:underline"
                >
                  Sign in
                </Link>
              </p>
            </div>
          </div>

          
        </div>
      </div>
    </div>
  )
} 