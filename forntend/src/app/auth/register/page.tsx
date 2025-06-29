'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'

export default function RegisterPage() {
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
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
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-green-50/30 to-yellow-50/30 flex items-center justify-center p-4">
      <div className="w-full max-w-lg">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="flex flex-col items-center justify-center mb-6">
            <div className="mb-4 p-2">
              <Image
                src="/logo-removebg-preview.png"
                alt="SocialPay Logo"
                width={160}
                height={60}
                className="object-contain max-w-full h-auto"
                priority
                unoptimized
              />
            </div>
            <p className="text-sm font-medium text-gray-600">
              SocialPay — <span className="text-brand-green-600 font-semibold">Connect.</span>{' '}
              <span className="text-brand-gold-600 font-semibold">Accept.</span>{' '}
              <span className="text-brand-green-600 font-semibold">Grow.</span>
            </p>
          </div>
        </div>
        
        {/* Registration Form */}
        <div className="socialpay-card p-8 shadow-xl hover:shadow-2xl transition-shadow duration-300">
          <div className="text-center mb-8">
            <h2 className="text-3xl font-bold text-gray-900 mb-2 bg-gradient-to-r from-brand-green-600 to-brand-gold-600 bg-clip-text text-transparent">
              Create an account
            </h2>
            <p className="text-gray-600 text-sm">Please enter your details</p>
          </div>

          {errors.general && (
            <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg shadow-sm">
              <p className="text-red-600 text-sm font-medium">{errors.general}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            {/* First Name and Last Name */}
            <div className="grid grid-cols-2 gap-4">
              <div className="group">
                <label htmlFor="firstName" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  First Name
                </label>
                <input
                  id="firstName"
                  name="firstName"
                  type="text"
                  required
                  value={formData.firstName}
                  onChange={handleChange}
                  className={`w-full px-4 py-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 shadow-sm hover:shadow-md ${errors.firstName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                  placeholder="First name"
                />
                {errors.firstName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.firstName}</p>}
              </div>
              <div className="group">
                <label htmlFor="lastName" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  Last Name
                </label>
                <input
                  id="lastName"
                  name="lastName"
                  type="text"
                  required
                  value={formData.lastName}
                  onChange={handleChange}
                  className={`w-full px-4 py-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 shadow-sm hover:shadow-md ${errors.lastName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                  placeholder="Last name"
                />
                {errors.lastName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.lastName}</p>}
              </div>
            </div>

            {/* Phone Number and Business Name */}
            <div className="grid grid-cols-2 gap-4">
              <div className="group">
                <label htmlFor="phone" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  Phone Number
                </label>
                <div className="flex h-12">
                  <div className="flex items-center px-4 py-4 bg-gradient-to-r from-brand-green-50 to-brand-gold-50 border border-r-0 border-gray-300 rounded-l-lg">
                    <span className="text-lg mr-2">🇪🇹</span>
                    <span className="text-sm text-gray-600 font-medium">+251</span>
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
                    className={`flex-1 px-4 py-4 border border-gray-300 rounded-r-lg focus:outline-none focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 ${errors.phone ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                    placeholder="911 123 456"
                  />
                </div>
                {errors.phone && <p className="mt-1 text-xs text-red-600 font-medium">{errors.phone}</p>}
              </div>
              <div className="group">
                <label htmlFor="businessName" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  Business Name
                </label>
                <input
                  id="businessName"
                  name="businessName"
                  type="text"
                  required
                  value={formData.businessName}
                  onChange={handleChange}
                  className={`w-full px-4 py-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 shadow-sm hover:shadow-md h-12 ${errors.businessName ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                  placeholder="Business name"
                />
                {errors.businessName && <p className="mt-1 text-xs text-red-600 font-medium">{errors.businessName}</p>}
              </div>
            </div>

            {/* Password and Confirm Password */}
            <div className="grid grid-cols-2 gap-4">
              <div className="group">
                <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
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
                    className={`w-full px-4 py-4 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 shadow-sm hover:shadow-md ${errors.password ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center hover:scale-110 transition-transform duration-200"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400 hover:text-brand-green-500 transition-colors" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400 hover:text-brand-green-500 transition-colors" />
                    )}
                  </button>
                </div>
                {errors.password && <p className="mt-1 text-xs text-red-600 font-medium">{errors.password}</p>}
              </div>
              <div className="group">
                <label htmlFor="confirmPassword" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
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
                    className={`w-full px-4 py-4 pr-12 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 bg-white placeholder-gray-400 shadow-sm hover:shadow-md ${errors.confirmPassword ? 'border-red-300 focus:border-red-500 focus:ring-red-500' : ''}`}
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-3 flex items-center hover:scale-110 transition-transform duration-200"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? (
                      <EyeSlashIcon className="h-5 w-5 text-gray-400 hover:text-brand-green-500 transition-colors" />
                    ) : (
                      <EyeIcon className="h-5 w-5 text-gray-400 hover:text-brand-green-500 transition-colors" />
                    )}
                  </button>
                </div>
                {errors.confirmPassword && <p className="mt-1 text-xs text-red-600 font-medium">{errors.confirmPassword}</p>}
              </div>
            </div>

            {/* Terms Agreement */}
            <div className="flex items-start space-x-3 pt-4">
              <input
                id="agreedToTerms"
                name="agreedToTerms"
                type="checkbox"
                checked={formData.agreedToTerms}
                onChange={handleChange}
                required
                className="mt-1 h-4 w-4 text-brand-green-600 focus:ring-brand-green-500 border-gray-300 rounded hover:scale-110 transition-transform duration-200"
              />
              <label htmlFor="agreedToTerms" className="text-sm text-gray-600 leading-tight">
                By creating an account, you agree to the{' '}
                <Link href="/terms" className="text-brand-green-600 hover:text-brand-green-500 font-semibold hover:underline transition-all duration-200">
                  Merchant Service Agreement
                </Link>{' '}
                and{' '}
                <Link href="/privacy" className="text-brand-green-600 hover:text-brand-green-500 font-semibold hover:underline transition-all duration-200">
                  Privacy Policy
                </Link>
                .
              </label>
              {errors.agreedToTerms && <p className="mt-1 text-xs text-red-600 font-medium">{errors.agreedToTerms}</p>}
            </div>

            <button
              type="submit"
              disabled={isLoading}
              className="w-full socialpay-button-primary disabled:opacity-50 disabled:cursor-not-allowed mt-6 shadow-lg hover:shadow-xl"
            >
              {isLoading ? (
                <div className="flex items-center justify-center">
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                  <span className="font-medium">Creating account...</span>
                </div>
              ) : (
                <span className="font-semibold">Sign Up</span>
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm text-gray-600">
              Already have an account?{' '}
              <Link href="/auth/login" className="font-semibold text-brand-green-600 hover:text-brand-green-500 hover:underline transition-all duration-200">
                Sign In
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
} 