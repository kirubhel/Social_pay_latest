'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { useAuthStore } from '@/stores/auth'
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const router = useRouter()
  const login = useAuthStore((state) => state.login)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError('')

    try {
      // Mock login - replace with actual API call
      if (email === 'admin@socialpay.com' && password === 'password') {
        const mockUser = {
          id: '1',
          email: 'admin@socialpay.com',
          name: 'John Doe',
          role: 'admin' as const,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        }
        
        login(mockUser, 'mock-token-123')
        router.push('/dashboard')
      } else {
        setError('Invalid email or password')
      }
    } catch (err) {
      setError('Login failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex">
      {/* Left side - Login Form */}
      <div className="flex-1 flex items-center justify-center p-8 bg-white">
        <div className="w-full max-w-md">
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

          {/* Sign In Form */}
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2 text-center bg-gradient-to-r from-brand-green-600 to-brand-gold-600 bg-clip-text text-transparent">
              Sign in
            </h2>
            <p className="text-gray-600 text-sm mb-8 text-center">Welcome back! Please enter your details</p>

            {error && (
              <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg shadow-sm">
                <p className="text-red-600 text-sm font-medium">{error}</p>
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="group">
                <label htmlFor="email" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  Email
                </label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="socialpay-input shadow-sm hover:shadow-md transition-shadow duration-200"
                  placeholder="Enter your email"
                />
              </div>

              <div className="group">
                <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-2 group-focus-within:text-brand-green-600 transition-colors">
                  Password
                </label>
                <div className="relative">
                  <input
                    id="password"
                    name="password"
                    type={showPassword ? 'text' : 'password'}
                    autoComplete="current-password"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="socialpay-input pr-12 shadow-sm hover:shadow-md transition-shadow duration-200"
                    placeholder="Enter your password"
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
              </div>

              <div className="flex items-center justify-end">
                <Link href="/auth/forgot-password" className="text-sm text-brand-green-600 hover:text-brand-green-500 font-semibold hover:underline transition-all duration-200">
                  Forgot your password?
                </Link>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full socialpay-button-primary disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl"
              >
                {isLoading ? (
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white mr-2"></div>
                    <span className="font-medium">Signing in...</span>
                  </div>
                ) : (
                  <span className="font-semibold">Sign in</span>
                )}
              </button>

              <div className="relative my-6">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-300" />
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-4 bg-white text-gray-500 font-medium">or</span>
                </div>
              </div>

              <button
                type="button"
                className="w-full flex items-center justify-center px-4 py-3 border-2 border-gray-200 rounded-lg shadow-sm bg-white text-gray-700 hover:bg-gray-50 hover:border-brand-green-300 font-semibold transition-all duration-200 hover:shadow-md"
              >
                <svg className="w-5 h-5 mr-3" viewBox="0 0 24 24">
                  <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                Sign in with Google
              </button>
            </form>

            <div className="mt-8 text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{' '}
                <Link href="/auth/register" className="font-semibold text-brand-green-600 hover:text-brand-green-500 hover:underline transition-all duration-200">
                  Sign up
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Welcome Section */}
      <div className="flex-1 bg-brand-gradient flex items-center justify-center p-8 relative overflow-hidden">
        {/* Background Decorations */}
        <div className="absolute top-20 left-20 w-32 h-32 bg-white/10 rounded-full blur-xl animate-pulse"></div>
        <div className="absolute bottom-20 right-20 w-40 h-40 bg-brand-gold-400/20 rounded-full blur-2xl"></div>
        <div className="absolute top-1/2 left-10 w-20 h-20 bg-white/5 rounded-full blur-lg"></div>
        
        <div className="text-center text-white max-w-lg z-10">
          <div className="floating-animation">
            <h1 className="text-5xl font-bold mb-4 drop-shadow-lg">Welcome Back!</h1>
            <p className="text-xl mb-8 text-white/90 font-medium">
              Please sign in to your SocialPay Dashboard
            </p>
            <p className="text-lg text-white/80 mb-8 leading-relaxed">
              Your all-in-one payment management platform.
            </p>
          </div>

          {/* Dashboard Preview */}
          <div className="glass-effect rounded-2xl p-6 mx-auto max-w-sm hover:scale-105 transition-transform duration-300">
            <div className="bg-white rounded-xl p-4 mb-4 shadow-lg">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-gray-800 font-semibold text-sm">Sales by Products</h3>
                <div className="w-8 h-8 bg-red-500 rounded-lg flex items-center justify-center text-white text-xs font-bold animate-bounce">
                  N
                </div>
              </div>
              <div className="text-red-600 text-sm font-medium">1 Issue ×</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
} 