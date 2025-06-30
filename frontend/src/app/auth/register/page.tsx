'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { EyeIcon, EyeSlashIcon, ArrowRightIcon, UserPlusIcon } from '@heroicons/react/24/outline'
import { CheckCircleIcon, ShieldCheckIcon, CreditCardIcon, StarIcon } from '@heroicons/react/24/solid'

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phone: '+251',
    businessName: '',
    password: '',
    confirmPassword: '',
    agreedToTerms: false
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  const router = useRouter()

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }))
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }))
    }
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.firstName.trim()) newErrors.firstName = 'First name is required'
    if (!formData.lastName.trim()) newErrors.lastName = 'Last name is required'
    if (!formData.email.trim()) newErrors.email = 'Email is required'
    if (!formData.phone.trim() || formData.phone === '+251') newErrors.phone = 'Valid phone number is required'
    if (!formData.businessName.trim()) newErrors.businessName = 'Business name is required'
    if (!formData.password) newErrors.password = 'Password is required'
    if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = 'Passwords do not match'
    if (!formData.agreedToTerms) newErrors.agreedToTerms = 'You must agree to the terms'

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!validateForm()) return
    
    setIsLoading(true)

    try {
      // Mock registration - replace with actual API call
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      // Redirect to login with success message
      router.push('/auth/login?message=Registration successful! Please log in.')
    } catch (err) {
      setErrors({ general: 'Registration failed. Please try again.' })
    } finally {
      setIsLoading(false)
    }
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
              <Image
                src="/logo.png"
                alt="Social Pay Logo"
                width={200}
                height={10}
                className="object-contain mx-auto mb-4"
                priority
              />
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
              {/* First Name and Last Name */}
              <div className="grid grid-cols-2 gap-3">
                <div className="group">
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
              </div>

              {/* Email */}
              <div className="group">
                <label htmlFor="email" className="block text-sm font-semibold text-gray-700 mb-1.5">
                  Email Address
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  className={`w-full px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.email ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                  placeholder="your@email.com"
                />
                {errors.email && <p className="mt-1 text-xs text-red-600 font-medium">{errors.email}</p>}
              </div>

              {/* Phone and Business Name */}
              <div className="grid grid-cols-2 gap-3">
                <div className="group">
                  <label htmlFor="phone" className="block text-sm font-semibold text-gray-700 mb-1.5">
                    Phone Number
                  </label>
                  <div className="flex">
                    <div className="flex items-center px-3 py-2.5 bg-gray-100 border border-r-0 border-gray-200/60 rounded-l-xl">
                      <span className="text-sm mr-1">🇪🇹</span>
                      <span className="text-xs text-gray-600 font-medium">+251</span>
                    </div>
                    <input
                      id="phone"
                      name="phone"
                      type="tel"
                      required
                      value={formData.phone.replace('+251', '')}
                      onChange={(e) => handleChange({
                        ...e,
                        target: { ...e.target, name: 'phone', value: '+251' + e.target.value }
                      })}
                      className={`flex-1 px-3 py-2.5 bg-white/60 backdrop-blur-sm border border-gray-200/60 rounded-r-xl focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all duration-200 placeholder-gray-400 text-sm ${errors.phone ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                      placeholder="911123456"
                    />
                  </div>
                  {errors.phone && <p className="mt-1 text-xs text-red-600 font-medium">{errors.phone}</p>}
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

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-200/60" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-3 bg-white/80 text-gray-500 font-medium">or sign up with</span>
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

          {/* Features & Trust Section */}
          <div className="mt-6 space-y-4">
            {/* Quick Features */}
            <div className="flex justify-center space-x-6 text-center">
              {[
                { icon: ShieldCheckIcon, text: 'Secure' },
                { icon: CreditCardIcon, text: 'Fast Setup' },
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

            {/* Trust Indicators */}
            <div className="text-center">
              <p className="text-xs text-gray-500 mb-2">Join 50,000+ successful businesses</p>
              <div className="flex items-center justify-center space-x-3 opacity-60">
                <div className="w-5 h-5 bg-gray-300 rounded"></div>
                <div className="w-5 h-5 bg-gray-300 rounded"></div>
                <div className="w-5 h-5 bg-gray-300 rounded"></div>
                <div className="w-5 h-5 bg-gray-300 rounded"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 