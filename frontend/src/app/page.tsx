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
  StarIcon,
  PhoneIcon,
  BuildingStorefrontIcon,
  CurrencyDollarIcon,
  ClockIcon
} from '@heroicons/react/24/outline'
import {
  LockClosedIcon,
  SparklesIcon,
  CheckBadgeIcon
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
      <nav className="relative z-10 px-6 py-4 bg-white/80 backdrop-blur-sm border-b border-gray-200/50">
        <div className="max-w-7xl mx-auto flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Image
              src="/logo.png"
              alt="Social Pay"
              width={160}
              height={40}
              className="rounded-lg"
            />
          </div>
          
          <div className="hidden md:flex items-center space-x-8">
            <a href="#features" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Solutions</a>
            <a href="#pricing" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Pricing</a>
            <a href="#developers" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Developers</a>
            <a href="#about" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">About</a>
            <a href="#contact" className="text-gray-600 hover:text-brand-green-600 font-medium transition-colors">Contact</a>
          </div>

          <div className="flex items-center space-x-4">
            <Link 
              href="/auth/login"
              className="text-brand-green-600 hover:text-brand-green-500 font-semibold transition-colors"
            >
              Login
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
      <section className="relative z-10 px-6 py-16">
        <div className="max-w-7xl mx-auto">
          <div className="text-center">
            <div className="inline-flex items-center px-4 py-2 bg-white/80 backdrop-blur-sm rounded-full border border-brand-green-200/60 mb-8">
              <CheckBadgeIcon className="w-4 h-4 text-brand-gold-500 mr-2" />
              <span className="text-sm font-medium text-gray-700">Trusted by 1000+ Ethiopian businesses</span>
            </div>
            
            <h1 className="text-4xl md:text-6xl lg:text-7xl font-bold text-gray-900 mb-6 leading-tight">
              Ethiopia's Leading
              <br />
              <span className="bg-gradient-to-r from-brand-green-600 via-brand-green-500 to-brand-gold-500 bg-clip-text text-transparent">
                Payment Gateway
              </span>
            </h1>
            
            <p className="text-xl md:text-2xl text-gray-600 mb-12 max-w-4xl mx-auto leading-relaxed">
              Accept payments seamlessly with Telebirr, CBE Birr, M-Birr, AmolePay, and international cards. 
              Built for Ethiopian businesses, designed for growth.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
              <Link 
                href="/auth/register"
                className="group bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white px-8 py-4 rounded-xl font-semibold text-lg shadow-xl hover:shadow-2xl transform hover:scale-105 transition-all duration-200 flex items-center"
              >
                Start Accepting Payments
                <ArrowRightIcon className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
              </Link>
              
              <Link 
                href="#demo"
                className="group bg-white/80 backdrop-blur-sm hover:bg-white text-gray-700 px-8 py-4 rounded-xl font-semibold text-lg border border-gray-200/60 hover:border-brand-green-300 shadow-lg hover:shadow-xl transition-all duration-200 flex items-center"
              >
                <DevicePhoneMobileIcon className="w-5 h-5 mr-2" />
                View Live Demo
              </Link>
            </div>

            {/* Trust Indicators */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6 text-center">
              {[
                { icon: CheckCircleIcon, text: '99.9% Uptime', color: 'text-brand-green-500' },
                { icon: ShieldCheckIcon, text: 'PCI Compliant', color: 'text-brand-green-500' },
                { icon: ClockIcon, text: '24/7 Support', color: 'text-brand-gold-500' },
                { icon: StarIcon, text: '4.9/5 Rating', color: 'text-yellow-400' }
              ].map((item, index) => (
                <div key={index} className="flex flex-col items-center">
                  <item.icon className={`w-6 h-6 ${item.color} mb-2`} />
                  <span className="text-sm font-medium text-gray-600">{item.text}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Payment Methods Section */}
      <section className="relative z-10 px-6 py-16 bg-white/60 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              All Ethiopian Payment Methods
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Support every customer with comprehensive payment options
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-6">
            {[
              { name: 'Telebirr', logo: '📱', desc: 'Ethiopia\'s #1 mobile money' },
              { name: 'CBE Birr', logo: '🏦', desc: 'Commercial Bank of Ethiopia' },
              { name: 'M-Birr', logo: '💳', desc: 'Mobile payments platform' },
              { name: 'AmolePay', logo: '🔔', desc: 'Digital wallet solution' },
              { name: 'Visa/Mastercard', logo: '💳', desc: 'International cards' },
              { name: 'Bank Transfer', logo: '🏛️', desc: 'Direct bank payments' }
            ].map((method, index) => (
              <div key={index} className="bg-white/80 backdrop-blur-sm p-6 rounded-2xl border border-white/20 shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200 text-center">
                <div className="text-3xl mb-3">{method.logo}</div>
                <h3 className="font-bold text-gray-900 mb-1">{method.name}</h3>
                <p className="text-xs text-gray-600">{method.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="relative z-10 px-6 py-20">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Built for Ethiopian Business Success
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Everything you need to grow your business in Ethiopia's digital economy
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {[
              {
                icon: BoltIcon,
                title: "Instant Settlement",
                description: "Get paid instantly with real-time settlement to your Ethiopian bank account"
              },
              {
                icon: LockClosedIcon,
                title: "Enterprise Security",
                description: "Bank-grade security with PCI DSS compliance and fraud protection"
              },
              {
                icon: ChartBarIcon,
                title: "Smart Analytics",
                description: "Real-time insights, transaction reports, and business intelligence"
              },
              {
                icon: PhoneIcon,
                title: "Mobile-First Design",
                description: "Optimized checkout experience for Ethiopian mobile users"
              },
              {
                icon: BuildingStorefrontIcon,
                title: "Multi-Business Support",
                description: "Manage multiple stores, brands, or business units from one dashboard"
              },
              {
                icon: GlobeAltIcon,
                title: "Developer APIs",
                description: "RESTful APIs, webhooks, and SDKs for seamless integration"
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

      {/* Success Stories Section */}
      <section className="relative z-10 px-6 py-20 bg-white/60 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Powering Ethiopian Businesses
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              From startups to enterprises, businesses trust Social Pay
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                category: "E-commerce",
                business: "Online Retail Store",
                growth: "+300% sales growth",
                quote: "Social Pay simplified our payment process and helped us reach more customers across Ethiopia."
              },
              {
                category: "Restaurant",
                business: "Food Delivery Service", 
                growth: "+250% order volume",
                quote: "Instant payments and multiple payment options increased our customer satisfaction significantly."
              },
              {
                category: "Service",
                business: "Digital Marketing Agency",
                growth: "+180% client retention",
                quote: "Professional payment experience helped us build trust with enterprise clients."
              }
            ].map((story, index) => (
              <div key={index} className="bg-white/80 backdrop-blur-sm p-8 rounded-2xl border border-white/20 shadow-lg hover:shadow-xl transition-all duration-200">
                <div className="flex items-center mb-4">
                  <span className="bg-brand-green-100 text-brand-green-700 px-3 py-1 rounded-full text-sm font-medium">
                    {story.category}
                  </span>
                </div>
                <h3 className="text-xl font-bold text-gray-900 mb-2">{story.business}</h3>
                <div className="text-2xl font-bold text-brand-green-600 mb-4">{story.growth}</div>
                <p className="text-gray-600 italic">"{story.quote}"</p>
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
                Growing Ethiopia's Digital Economy
              </h2>
              <p className="text-xl opacity-90">
                Trusted numbers that speak for themselves
              </p>
            </div>
            
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
              {[
                { number: "1,000+", label: "Active Merchants" },
                { number: "₿50M+", label: "Monthly Volume" },
                { number: "99.9%", label: "Uptime" },
                { number: "<2s", label: "Transaction Speed" }
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

      {/* Pricing Section */}
      <section id="pricing" className="relative z-10 px-6 py-20 bg-white/60 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
              Simple, Transparent Pricing
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Pay only for what you use. No hidden fees, no setup costs.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {[
              {
                name: "Starter",
                price: "2.5%",
                description: "Perfect for small businesses",
                features: ["All payment methods", "Basic analytics", "Email support", "Standard checkout"]
              },
              {
                name: "Business",
                price: "2.0%",
                description: "Best for growing businesses",
                features: ["All payment methods", "Advanced analytics", "Priority support", "Custom checkout", "Multi-store management"],
                popular: true
              },
              {
                name: "Enterprise",
                price: "Custom",
                description: "For large organizations",
                features: ["All payment methods", "Custom analytics", "24/7 dedicated support", "White-label solution", "Custom integration"]
              }
            ].map((plan, index) => (
              <div key={index} className={`bg-white/80 backdrop-blur-sm p-8 rounded-2xl border-2 shadow-lg hover:shadow-xl transition-all duration-200 ${plan.popular ? 'border-brand-green-500 transform scale-105' : 'border-white/20'}`}>
                {plan.popular && (
                  <div className="bg-brand-green-500 text-white px-4 py-1 rounded-full text-sm font-medium mb-4 text-center">
                    Most Popular
                  </div>
                )}
                <h3 className="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
                <div className="text-4xl font-bold text-brand-green-600 mb-2">{plan.price}</div>
                <p className="text-gray-600 mb-6">{plan.description}</p>
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature, fIndex) => (
                    <li key={fIndex} className="flex items-center">
                      <CheckCircleIcon className="w-5 h-5 text-brand-green-500 mr-3" />
                      <span className="text-gray-700">{feature}</span>
                    </li>
                  ))}
                </ul>
                <Link 
                  href="/auth/register"
                  className={`w-full inline-flex items-center justify-center px-6 py-3 rounded-xl font-semibold transition-all duration-200 ${plan.popular ? 'bg-brand-green-500 hover:bg-brand-green-600 text-white shadow-lg hover:shadow-xl' : 'bg-gray-100 hover:bg-gray-200 text-gray-700'}`}
                >
                  Get Started
                </Link>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="relative z-10 px-6 py-20">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Ready to Transform Your Business?
          </h2>
          <p className="text-xl text-gray-600 mb-12">
            Join thousands of Ethiopian businesses already growing with Social Pay
          </p>
          
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
            <Link 
              href="/auth/register"
              className="group bg-gradient-to-r from-brand-green-500 to-brand-gold-400 hover:from-brand-green-600 hover:to-brand-gold-500 text-white px-8 py-4 rounded-xl font-semibold text-lg shadow-xl hover:shadow-2xl transform hover:scale-105 transition-all duration-200 flex items-center"
            >
              Start Your Free Account
              <ArrowRightIcon className="w-5 h-5 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
            </Link>
            
            <div className="text-center">
              <p className="text-sm text-gray-500">Free setup • No monthly fees • Cancel anytime</p>
            </div>
          </div>

          <div className="bg-white/80 backdrop-blur-sm p-6 rounded-2xl border border-white/20 shadow-lg">
            <p className="text-lg font-semibold text-gray-900 mb-2">Need help getting started?</p>
            <p className="text-gray-600 mb-4">Our team is ready to help you integrate and optimize your payment flow</p>
            <div className="flex items-center justify-center space-x-6 text-sm">
              <a href="tel:+251911000000" className="flex items-center text-brand-green-600 hover:text-brand-green-500 font-medium">
                <PhoneIcon className="w-4 h-4 mr-2" />
                +251 911 000 000
              </a>
              <a href="mailto:support@socialpay.et" className="flex items-center text-brand-green-600 hover:text-brand-green-500 font-medium">
                <span className="mr-2">✉️</span>
                support@socialpay.et
              </a>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="relative z-10 px-6 py-12 bg-gray-900 text-white">
        <div className="max-w-7xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
            <div className="col-span-1 md:col-span-2">
              <div className="flex items-center space-x-3 mb-4">
                <Image
                  src="/logo.png"
                  alt="Social Pay"
                  width={120}
                  height={30}
                  className="rounded-lg"
                />
              </div>
              <p className="text-gray-400 mb-4 max-w-md">
                Ethiopia's leading payment gateway, empowering businesses with seamless digital payment solutions.
              </p>
              <div className="flex space-x-4">
                <a href="#" className="text-gray-400 hover:text-white transition-colors">Twitter</a>
                <a href="#" className="text-gray-400 hover:text-white transition-colors">LinkedIn</a>
                <a href="#" className="text-gray-400 hover:text-white transition-colors">Telegram</a>
              </div>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Solutions</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">Online Payments</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Mobile Payments</a></li>
                <li><a href="#" className="hover:text-white transition-colors">In-Store Payments</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Subscription Billing</a></li>
              </ul>
            </div>
            
            <div>
              <h3 className="text-lg font-semibold mb-4">Company</h3>
              <ul className="space-y-2 text-gray-400">
                <li><a href="#" className="hover:text-white transition-colors">About Us</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Careers</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Press</a></li>
                <li><a href="#" className="hover:text-white transition-colors">Contact</a></li>
              </ul>
            </div>
          </div>
          
          <div className="border-t border-gray-800 pt-8 flex flex-col md:flex-row items-center justify-between">
            <p className="text-gray-400 mb-4 md:mb-0">© 2024 Social Pay. All rights reserved.</p>
            <div className="flex space-x-6 text-sm">
              <a href="#" className="text-gray-400 hover:text-white transition-colors">Privacy Policy</a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">Terms of Service</a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">Security</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
