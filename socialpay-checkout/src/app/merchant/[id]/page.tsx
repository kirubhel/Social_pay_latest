'use client'

import Image from 'next/image'
import { use, useState, useEffect } from 'react'
import { useTranslations } from 'next-intl'
import AdBanner from '@/components/AdBanner'
import PaymentMethodSelector from '@/components/PaymentMethodSelector'
import Divider from '@/components/Divider'
import CardForm from '@/components/CardForm'
import WalletForm from '@/components/WalletForm'
import BankForm from '@/components/BankForm'
import SocialPayForm from '@/components/SocialPayForm'
import LanguageSelector from '@/components/LanguageSelector'
import LoadingSpinner from '@/components/LoadingSpinner'
import { MerchantPaymentProvider, useMerchantPayment } from '@/contexts/MerchantPaymentContext'
import { Gateway, getGateways } from '@/services/gateway.service'
import ErrorPage from '@/components/ErrorPage'
import { GatewayError } from '@/services/gateway.service'
import AmountInputWrapper from '@/components/AmountInputWrapper'
import { getMerchant, V2ClientError } from '@/services/merchant.service'
import { Merchant } from '@/types/merchant.dto'

export default function Checkout(props: {
  params: Promise<{ id: string }>
}) {
  const params = use(props.params)
  console.log(params)
  const [paymentMethod, setPaymentMethod] = useState<'socialpay' | 'wallet' | 'bank' | 'card'>('wallet')
  const [isLoading, setIsLoading] = useState(true)
  const [gateways, setGateways] = useState<Gateway[]>([])
  const [error, setError] = useState<GatewayError | null>(null)
  const [merchant, setMerchant] = useState<Merchant | null>(null)
  const t = useTranslations('checkout')

  const PaymentDetails = () => {
    const { amount } = useMerchantPayment()

    return (
      <div className="w-[440px] flex-1 hidden md:flex flex-col bg-[#F1F2F3] rounded-[9.6px] ml-[47px]">
        <div className="p-5">
          <h2 className="text-lg font-medium mb-3">{t('payment')}</h2>
          <div className="text-gray-600 mb-5">{t('paymentDetails')}</div>
          {merchant && (
            <>
              <div className="mb-4 p-3 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-gray-50 rounded-full">
                    <Image 
                      src="/matrix.png" 
                      alt="Merchant"
                      width={40}
                      height={40}
                      className="opacity-100"
                    />
                  </div>
                  <div>
                    <span className="text-sm text-gray-500">Merchant</span>
                    <div className="font-medium text-gray-900">{merchant.trading_name || merchant.legal_name}</div>
                  </div>
                </div>
              </div>
              <Divider />
            </>
          )}
          <div className="space-y-3">
            <div className="flex justify-between">
              <span className="text-gray-600">{t('amount')}</span>
              <span>{amount || '0.00'} ETB</span>
            </div>
            <div className="flex justify-between">
              <span className="text-gray-600">{t('tax')}</span>
              <span>0.00 ETB</span>
            </div>
          </div>
          <Divider />
          <div className="flex justify-between pt-3">
            <span className="text-gray-600">{t('total')}</span>
            <span className="text-green-500">{amount || '0.00'} ETB</span>
          </div>
          <Divider />
        </div>
      </div>
    )
  }

  useEffect(() => {
    const fetchMerchant = async () => {
      try {
        const response = await getMerchant(params.id)
        setMerchant(response.data)
      } catch (err) {
        if (err instanceof V2ClientError) {
          setError(new GatewayError(err.message, err.status, 'USER_ERROR'))
        } else {
          setError(new GatewayError('Failed to fetch merchant', 500, 'USER_ERROR'))
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
    
    fetchMerchant()
    fetchGateways()
  }, [params.id])

  const handlePaymentMethodSelect = (
    method: 'wallet' | 'bank' | 'card' | 'socialpay',
  ) => {
    setPaymentMethod(method)
  }

  const getAvailableGateways = (type: string) => {
    return gateways.filter((gateway) => gateway.type === type.toUpperCase())
  }

  const renderPaymentForm = () => {
    const availableGateways = getAvailableGateways(paymentMethod)

    switch (paymentMethod) {
      case 'card':
        return <CardForm gateways={availableGateways}  checkoutId={params.id}/>
      case 'wallet':
        return <WalletForm gateways={availableGateways} merchantId={params.id} />
      case 'bank':
        return <BankForm gateways={availableGateways} />
      case 'socialpay':
        return <SocialPayForm />
      default:
        return <WalletForm gateways={availableGateways} merchantId={params.id} />
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

  return (
    <MerchantPaymentProvider 
      merchantId={params.id} 
      merchantName={merchant?.trading_name || merchant?.legal_name || 'Unknown Merchant'}
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
                  {t('choosePayment')}
                </h1>

                {/* Payment Method Selector */}
                <div className="mb-4">
                  <PaymentMethodSelector
                    onSelect={handlePaymentMethodSelect}
                    defaultMethod={
                      paymentMethod as 'wallet' | 'bank' | 'card' | 'socialpay'
                    }
                    availableTypes={[
                      ...new Set(gateways.map((g) => g.type.toLowerCase()))
                    ]}
                  />
                </div>
                {/* Accepted Cards */}
                {paymentMethod == 'card' && <div className="flex flex-wrap items-center gap-4 mb-4">
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
                }

                <AmountInputWrapper merchantId={params.id} />
                {/* Payment Form based on selected method */}
                {renderPaymentForm()}
              </div>

              {/* Right Column - Now using the PaymentDetails component */}
              <PaymentDetails />
            </div>
            {/* Ad Banner */}
            <div className="max-w-[893px] pt-6">
              <AdBanner />
            </div>
          </div>
        </div>
      </div>
    </MerchantPaymentProvider>
  )
}
