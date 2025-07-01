'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/stores/auth'
import Link from 'next/link'
import Image from 'next/image'
import { 
  ArrowRightIcon, 
  CheckCircleIcon, 
  CreditCardIcon, 
  ShieldCheckIcon,
  GlobeAltIcon,
  ChartBarIcon,
  BoltIcon,
  DevicePhoneMobileIcon,
  StarIcon
} from '@heroicons/react/24/outline'
import {
  CurrencyDollarIcon,
  LockClosedIcon,
  SparklesIcon
} from '@heroicons/react/24/solid'

export default function LandingPage() {
  const { isAuthenticated, isHydrated } = useAuthStore()
  const router = useRouter()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (mounted && isHydrated && isAuthenticated) {
      router.push('/dashboard')
    }
  }, [isAuthenticated, isHydrated, mounted, router])

  // Don't render anything until component is mounted and store is hydrated
  if (!mounted || !isHydrated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-brand-green-50 via-white to-brand-gold-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-brand-green-500"></div>
      </div>
    )
  }

  if (isAuthenticated) {
    return null // Will redirect to dashboard
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-brand-green-50 via-white to-brand-gold-50 relative overflow-hidden">
      {/* Animated Background Elements */}
      <div className="absolute inset-0">
        <div className="absolute top-20 left-20 w-64 h-64 bg-brand-green-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
        <div className="absolute top-32 right-20 w-64 h-64 bg-brand-gold-200/20 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
        <div className="absolute -bottom-8 left-1/2 w-64 h-64 bg-brand-green-300/10 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
      </div>

      {/* Navigation */}
      <nav className="relative z-10 px-6 py-4">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Image
              src="/logo.png"
              alt="Social Pay"
              width={160}
              height={40}
              className="rounded-lg"
            />
            {/* <span className="text-2xl font-bold bg-gradient-to-r from-brand-green-600 to-brand-gold-500 bg-clip-text text-transparent">
              Social Pay
            </span> */}
          </div>
          
          <div className="hidden md:flex items-center space-x-8">
            <a href="#features" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Features</a>
            <a href="#pricing" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Pricing</a>
            <a href="#about" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">About</a>
          </div>

          <div className="flex items-center space-x-4">
            <Link 
              href="/auth/login"
              className="text-brand-green-600 hover:text-brand-green-500 font-semibold transition-colors"
            >
              Sign In
            </Link>
            <Link 
              href="/auth/register"
              className="bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white px-6 py-2 rounded-xl font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200"
            >
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative z-10 px-6 py-20">
        <div className="max-w-7xl mx-auto">
          <div className="text-center">
            <div className="inline-flex items-center px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full border border-brand-green-200/60 mb-8">
              <SparklesIcon className="w-4 h-4 text-brand-gold-500 mr-2" />
              <span className="text-sm font-medium text-gray-700">Trusted by 10,000+ businesses in Ethiopia</span>
            </div>
            
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold text-gray-900 mb-6 leading-tight">
              Accept Payments
              <br />
              <span className="bg-gradient-to-r from-brand-green-600 via-brand-green-500 to-brand-gold-500 bg-clip-text text-transparent">
                Effortlessly
              </span>
            </h1>
            
            <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-3xl mx-auto leading-relaxed">
              The modern payment gateway for Ethiopian businesses. Accept mobile money, 
              bank transfers, and cards with our simple API and beautiful checkout experience.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
              <Link 
                href="/auth/register"
                className="group bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white px-8 py-4 rounded-xl font-semibold text-lg shadow-xl hover:shadow-2xl transform hover:scale-105 transition-all duration-200 flex items-center"
              >
                Start Free Trial
                <ArrowRightIcon className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
              </Link>
              
              <Link 
                href="/auth/login"
                className="group bg-white/80 backdrop-blur-sm hover:bg-white text-gray-700 px-8 py-4 rounded-xl font-semibold text-lg border border-gray-200/60 hover:border-brand-green-300 shadow-lg hover:shadow-xl transition-all duration-200 flex items-center"
              >
                <DevicePhoneMobileIcon className="w-5 h-5 mr-2" />
                View Demo
              </Link>
            </div>

            {/* Social Proof */}
            <div className="flex items-center justify-center space-x-8 text-sm text-gray-500">
              <div className="flex items-center">
                <StarIcon className="w-4 h-4 text-yellow-400 mr-1" />
                <span className="font-medium">4.9/5 Rating</span>
              </div>
              <div className="flex items-center">
                <CheckCircleIcon className="w-4 h-4 text-brand-green-500 mr-1" />
                <span className="font-medium">99.9% Uptime</span>
              </div>
              <div className="flex items-center">
                <ShieldCheckIcon className="w-4 h-4 text-brand-green-500 mr-1" />
                <span className="font-medium">Bank-level Security</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="relative z-10 px-6 py-20 bg-white/60 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Everything you need to accept payments
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Powerful features designed to help your business grow faster
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: CreditCardIcon,
                title: "Multiple Payment Methods",
                description: "Accept Telebirr, CBE Birr, AmolePay, bank transfers, and international cards"
              },
              {
                icon: BoltIcon,
                title: "Lightning Fast",
                description: "Process payments in seconds with our optimized infrastructure"
              },
              {
                icon: LockClosedIcon,
                title: "Bank-level Security",
                description: "PCI DSS compliant with end-to-end encryption for all transactions"
              },
              {
                icon: ChartBarIcon,
                title: "Real-time Analytics",
                description: "Track payments, revenue, and customer insights with detailed reports"
              },
              {
                icon: DevicePhoneMobileIcon,
                title: "Mobile Optimized",
                description: "Beautiful checkout experience across all devices and screen sizes"
              },
              {
                icon: GlobeAltIcon,
                title: "API Integration",
                description: "Simple REST API and webhooks for seamless integration"
              }
            ].map((feature, index) => (
              <div key={index} className="group bg-white/80 backdrop-blur-sm p-8 rounded-2xl border border-white/20 shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200">
                <div className="w-12 h-12 bg-gradient-to-r from-brand-green-500 to-brand-gold-400 rounded-xl flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-200">
                  <feature.icon className="w-6 h-6 text-white" />
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-3">{feature.title}</h3>
                <p className="text-gray-600 leading-relaxed">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="relative z-10 px-6 py-20">
        <div className="max-w-7xl mx-auto">
          <div className="bg-gradient-to-r from-brand-green-500 to-brand-gold-400 rounded-3xl p-12 text-white">
            <div className="text-center mb-12">
              <h2 className="text-3xl md:text-4xl font-bold mb-4">
                Trusted by thousands of businesses
              </h2>
              <p className="text-xl opacity-90">
                Join the growing community of successful merchants
              </p>
            </div>
            
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
              {[
                { number: "10,000+", label: "Active Merchants" },
                { number: "₹50M+", label: "Processed Monthly" },
                { number: "99.9%", label: "Uptime" },
                { number: "24/7", label: "Support" }
              ].map((stat, index) => (
                <div key={index}>
                  <div className="text-3xl md:text-4xl font-bold mb-2">{stat.number}</div>
                  <div className="text-lg opacity-90">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="relative z-10 px-6 py-20 bg-white/60 backdrop-blur-sm">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Ready to get started?
          </h2>
          <p className="text-xl text-gray-600 mb-12">
            Join thousands of businesses already using Social Pay to grow their revenue
          </p>
          
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <Link 
              href="/auth/register"
              className="group bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white px-8 py-4 rounded-xl font-semibold text-lg shadow-xl hover:shadow-2xl transform hover:scale-105 transition-all duration-200 flex items-center"
            >
              Create Free Account
              <ArrowRightIcon className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
            </Link>
            
            <Link 
              href="/auth/login"
              className="group bg-white/80 backdrop-blur-sm hover:bg-white text-gray-700 px-8 py-4 rounded-xl font-semibold text-lg border border-gray-200/60 hover:border-brand-green-300 shadow-lg hover:shadow-xl transition-all duration-200 flex items-center"
            >
              <CurrencyDollarIcon className="w-5 h-5 mr-2" />
              Already have an account?
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 px-6 py-12 bg-gray-900 text-white">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col md:flex-row items-center justify-between">
            <div className="flex items-center space-x-3 mb-4 md:mb-0">
                             <Image
                 src="/logo.png"
                 alt="Social Pay"
                 width={32}
                 height={32}
                 className="rounded-lg"
               />
              <span className="text-xl font-bold">Social Pay</span>
            </div>
            
            <div className="text-center md:text-right">
              <p className="text-gray-400 mb-2">© 2024 Social Pay. All rights reserved.</p>
              <p className="text-sm text-gray-500">Secure payment processing for Ethiopian businesses</p>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
