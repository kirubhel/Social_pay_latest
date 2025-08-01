'use client'

import Image from 'next/image'
import { use, useState, useEffect } from 'react'
import { useTranslations } from 'next-intl'
import AdBanner from '@/components/AdBanner'
import PaymentMethodSelector from '@/components/PaymentMethodSelector'
import Divider from '@/components/Divider'
import QRWalletForm from '@/components/QRWalletForm'
import CardForm from '@/components/CardForm'
import BankForm from '@/components/BankForm'
import SocialPayForm from '@/components/SocialPayForm'
import AmountInput from '@/components/AmountInput'
import LanguageSelector from '@/components/LanguageSelector'
import LoadingSpinner from '@/components/LoadingSpinner'
import { QRPaymentProvider, useQRPayment } from '@/contexts/QRPaymentContext'
import { Gateway, getGateways } from '@/services/gateway.service'
import ErrorPage from '@/components/ErrorPage'
import { GatewayError } from '@/services/gateway.service'
import { getQRLink, V2ClientError } from '@/services/qr.service'
import { QRLinkResponse } from '@/types/qr.dto'

export default function QRPaymentPage(props: {
  params: Promise<{ id: string }>
}) {
  const params = use(props.params)
  const [paymentMethod, setPaymentMethod] = useState<'socialpay' | 'wallet' | 'bank' | 'card'>('wallet')
  const [isLoading, setIsLoading] = useState(true)
  const [gateways, setGateways] = useState<Gateway[]>([])
  const [error, setError] = useState<GatewayError | null>(null)
  const [qrLink, setQRLink] = useState<QRLinkResponse | null>(null)
  const t = useTranslations('checkout')

  const PaymentDetails = () => {
    const { amount, tipAmount } = useQRPayment()

    const displayAmount = qrLink?.type === 'STATIC' 
      ? qrLink.amount || 0 
      : parseFloat(amount) || 0;
    
    const displayTipAmount = parseFloat(tipAmount) || 0;
    const totalAmount = displayAmount + displayTipAmount;

    return (
      <div className="w-[440px] flex-1 hidden md:flex flex-col bg-[#F1F2F3] rounded-[9.6px] ml-[47px]">
        <div className="p-5">
          <h2 className="text-lg font-medium mb-3">{t('payment')}</h2>
          <div className="text-gray-600 mb-5">{t('paymentDetails')}</div>
          
          {qrLink && (
            <>
              <div className="mb-4 p-3 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-gray-50 rounded-full">
                    {qrLink.image_url ? (
                      <Image 
                        src={qrLink.image_url} 
                        alt={qrLink.title}
                        width={40}
                        height={40}
                        className="rounded-full object-cover"
                      />
                    ) : (
                      <Image 
                        src="/merchant-icon.png" 
                        alt="QR Payment"
                        width={40}
                        height={40}
                        className="opacity-100"
                      />
                    )}
                  </div>
                  <div>
                    <span className="text-sm text-gray-500">{qrLink.tag}</span>
                    <div className="font-medium text-gray-900">{qrLink.title}</div>
                  </div>
                </div>
              </div>
              
              {qrLink.description && (
                <div className="mb-4 p-3 bg-white rounded-lg">
                  <span className="text-sm text-gray-500">Description</span>
                  <div className="text-gray-900">{qrLink.description}</div>
                </div>
              )}
              
              <Divider />
            </>
          )}

          <div className="space-y-3">
            <div className="flex justify-between items-center">
              <span className="text-gray-600 text-sm">Type</span>
              <span className="text-gray-900 font-medium">
                {qrLink?.type === 'STATIC' ? 'Fixed Amount' : 'Variable Amount'}
              </span>
            </div>
            
            {qrLink?.type === 'STATIC' && (
              <div className="flex justify-between items-center">
                <span className="text-gray-600 text-sm">Amount</span>
                <span className="text-gray-900 font-medium">
                  {qrLink.amount?.toFixed(2)} ETB
                </span>
              </div>
            )}
          </div>
          
          <Divider />
          
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <span className="text-gray-600">Payment Amount</span>
              <span className="text-gray-900">{displayAmount.toFixed(2)} ETB</span>
            </div>
            
            {displayTipAmount > 0 && (
              <div className="flex justify-between items-center">
                <span className="text-gray-600">Tip Amount</span>
                <span className="text-gray-900">+{displayTipAmount.toFixed(2)} ETB</span>
              </div>
            )}
            
            <div className="flex justify-between pt-3 border-t">
              <span className="text-gray-600 font-medium">Total</span>
              <span className="text-green-500 font-medium">{totalAmount.toFixed(2)} ETB</span>
            </div>
          </div>
        </div>
      </div>
    )
  }

  const AmountInputSection = () => {
    const { amount, setAmount } = useQRPayment()

    const handleAmountChange = (value: string, isValid: boolean) => {
      if (isValid) {
        setAmount(value)
      }
    }

    // Only show for dynamic QR codes
    if (qrLink?.type === 'STATIC') {
      return null
    }

    return (
      <div className="mb-4">
        <AmountInput
          value={amount}
          onChange={handleAmountChange}
        />
      </div>
    )
  }

  useEffect(() => {
    const fetchQRLink = async () => {
      try {
        const response = await getQRLink(params.id)
        setQRLink(response)
      } catch (err) {
        if (err instanceof V2ClientError) {
          setError(new GatewayError(err.message, err.status, 'USER_ERROR'))
        } else {
          setError(new GatewayError('Failed to fetch QR link details', 500, 'USER_ERROR'))
        }
      }
    }

    const fetchGateways = async () => {
      setIsLoading(true)
      setError(null)
      try {
        const data = await getGateways()
        setGateways(data.filter((gateway) => gateway.can_process))
      } catch (err) {
        if (err instanceof GatewayError) {
          setError(err)
        } else {
          setError(
            new GatewayError(
              'An unexpected error occurred',
              500,
              'USER_ERROR',
            ),
          )
        }
      } finally {
        setIsLoading(false)
      }
    }

    fetchQRLink()
    fetchGateways()
  }, [params.id])

  const handlePaymentMethodSelect = (
    method: 'wallet' | 'bank' | 'card' | 'socialpay',
  ) => {
    setPaymentMethod(method)
  }

  const getAvailableGateways = (type?: string) => {
    const targetType = type || paymentMethod
    const typeGateways = gateways.filter((gateway) => gateway.type === targetType.toUpperCase())
    
    // Filter gateways by QR link supported methods
    if (qrLink?.supported_methods) {
      return typeGateways.filter(gateway => 
        qrLink.supported_methods.includes(gateway.key)
      )
    }
    return typeGateways
  }

  const renderPaymentForm = () => {
    const availableGateways = getAvailableGateways()
    
    switch (paymentMethod) {
      case 'card':
        return <CardForm 
          gateways={availableGateways} 
          qrLinkId={params.id} 
          isStaticAmount={qrLink?.type === 'STATIC'}
          staticAmount={qrLink?.amount}
          initialPhoneNumber=''
        />
      case 'wallet':
        return (
          <QRWalletForm
            gateways={availableGateways}
            qrLinkId={params.id}
            isTipEnabled={qrLink?.is_tip_enabled || false}
            isStaticAmount={qrLink?.type === 'STATIC'}
            staticAmount={qrLink?.amount}
          />
        )
      case 'bank':
        return <BankForm gateways={availableGateways} qrLinkId={params.id} />
      case 'socialpay':
        return <SocialPayForm qrLinkId={params.id} />
      default:
        return (
          <QRWalletForm
            gateways={availableGateways}
            qrLinkId={params.id}
            isTipEnabled={qrLink?.is_tip_enabled || false}
            isStaticAmount={qrLink?.type === 'STATIC'}
            staticAmount={qrLink?.amount}
          />
        )
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-[rgb(248,248,248)] flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (error) {
    return (
      <ErrorPage
        message={error.message}
        type={error.type}
        onRetry={() => window.location.reload()}
      />
    )
  }

  if (!qrLink) {
    return (
      <ErrorPage
        message="QR link not found"
        type="NOT_FOUND"
        onRetry={() => window.location.reload()}
      />
    )
  }

  // Check if QR link is active
  if (!qrLink.is_active) {
    return (
      <ErrorPage
        message="This QR link is no longer active"
        type="BAD_REQUEST"
        onRetry={() => window.location.reload()}
      />
    )
  }

  const availableGateways = getAvailableGateways()

  if (availableGateways.length === 0 && paymentMethod === 'wallet') {
    return (
      <ErrorPage
        message="No payment methods available for this QR link"
        type="BAD_REQUEST"
        onRetry={() => window.location.reload()}
      />
    )
  }

  return (
    <QRPaymentProvider 
      qrLinkId={params.id} 
      title={qrLink.title}
      initialAmount={qrLink.type === 'STATIC' ? qrLink.amount : undefined}
    >
      <div className="min-h-screen bg-[rgb(248,248,248)] flex flex-col pb-[70px] md:pb-0">
        {/* Header */}
        <div className="w-full px-5 py-3 mt-5">
          <div className="max-w-[1240px] mx-auto flex justify-between items-center">
            <Image
              src="/socialpay.webp"
              alt="SocialPay logo"
              width={148}
              height={43}
              priority
            />
            <LanguageSelector />
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 flex items-start justify-center px-4 py-2 mb-[30px]">
          <div className="w-full max-w-[893px]">
            <div className="flex flex-col md:flex-row">
              {/* Left Column - Payment Form */}
              <div className="flex-1">
                <h1 className="text-[18.5px] font-medium mb-3">
                  {qrLink.title}
                </h1>
                
                {qrLink.description && (
                  <p className="text-gray-600 mb-4">{qrLink.description}</p>
                )}

                {/* Payment Method Selector */}
                <div className="mb-4">
                  <PaymentMethodSelector
                    onSelect={handlePaymentMethodSelect}
                    defaultMethod={
                      paymentMethod as 'wallet' | 'bank' | 'card' | 'socialpay'
                    }
                    availableTypes={[
                      ...new Set(gateways.map((g) => g.type.toLowerCase())),
                    ]}
                  />
                </div>

                {/* Amount Input - Shared across all payment methods for dynamic QR */}
                <AmountInputSection />

                {/* Accepted Cards */}
                {paymentMethod === 'card' && (
                  <div className="flex flex-wrap items-center gap-4 mb-4">
                    <span className="text-sm text-gray-600">We accept:</span>
                    <div className="flex items-center gap-4">
                      <Image
                        src="/card/visa.svg"
                        alt="Visa"
                        width={40}
                        height={25}
                        className="object-contain"
                      />
                      <Image
                        src="/card/mastercard.svg"
                        alt="Mastercard"
                        width={40}
                        height={25}
                        className="object-contain"
                      />
                      <Image
                        src="/card/boa.svg"
                        alt="Bank of Abyssinia"
                        width={40}
                        height={25}
                        className="object-contain"
                      />
                      <div className="bg-[#112349] p-[5px]">
                        <Image
                          src="/card/ethswitch_logo.png"
                          alt="ETSwitch"
                          width={74}
                          height={47}
                          className="object-contain"
                        />
                      </div>
                    </div>
                  </div>
                )}

                {/* Payment Form based on selected method */}
                {renderPaymentForm()}
              </div>

              {/* Right Column - Payment Details */}
              <PaymentDetails />
            </div>
            
            {/* Ad Banner */}
            <div className="max-w-[893px] pt-6">
              <AdBanner />
            </div>
          </div>
        </div>
      </div>
    </QRPaymentProvider>
  )
} 