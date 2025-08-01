'use client'

import { useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import Link from 'next/link'
import Image from 'next/image'
import { CheckCircle, XCircle, Loader2 } from 'lucide-react'
import { v2Client, V2ClientError } from '@/services/v2client.service'
import LanguageSelector from '@/components/LanguageSelector'

interface TransactionData {
  id: string
  amount: number
  currency: string
  status: string
  reference: string
  reference_number: string
  medium: string
  created_at: string
  updated_at: string
  description: string
  success_url: string
  failed_url: string
}

export default function PaymentResultClient() {
 const searchParams = useSearchParams()
  const transactionId = searchParams.get('transactionId')

  const [transaction, setTransaction] = useState<TransactionData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!transactionId) {
      setError('No transaction ID found in the URL.')
      setLoading(false)
      return
    }

    async function fetchTransaction() {
      try {
        const data = await v2Client.get<TransactionData>(`/payment/transaction/${transactionId}`)
        setTransaction(data)
      } catch (err) {
        if (err instanceof V2ClientError) {
          setError(err.message)
        } else {
          setError('Unknown error occurred while fetching transaction.')
        }
      } finally {
        setLoading(false)
      }
    }

    fetchTransaction()
  }, [transactionId])

  const isSuccess = transaction?.status?.toUpperCase() === 'SUCCESS'
  const merchantUrl = isSuccess ? transaction?.success_url : transaction?.failed_url

  return (
    <>
      {/* Header */}
      <header className="w-full px-5 py-4 border-b border-gray-200 bg-white shadow-sm">
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
      </header>

      {/* Content */}
      <main className="flex items-center justify-center min-h-[80vh] bg-gradient-to-br from-gray-50 to-gray-100 px-4 py-8">
        {loading ? (
          <div className="flex flex-col items-center text-gray-600 animate-fadeIn">
            <Loader2 className="w-12 h-12 animate-spin mb-4" />
            <p className="text-xl font-semibold">Verifying your payment...</p>
          </div>
        ) : error ? (
          <div className="max-w-lg w-full bg-white/80 backdrop-blur-md shadow-xl rounded-3xl p-10 text-center border border-red-200 transition-transform duration-300 hover:scale-[1.01] animate-fadeIn">
            <XCircle className="w-20 h-20 text-red-600 mx-auto mb-6" />
            <h1 className="text-3xl font-bold mb-3 text-red-700">Payment Verification Failed</h1>
            <p className="text-gray-600 mb-8">{error}</p>
            <Link
              href="/checkout"
              className="inline-block px-8 py-3 bg-red-600 text-white rounded-xl hover:bg-red-700 hover:shadow-lg transition transform hover:scale-105"
            >
              Try Again
            </Link>
          </div>
        ) : (
          <div className="max-w-lg w-full bg-white/80 backdrop-blur-md shadow-xl rounded-3xl p-10 text-center border transition-transform duration-300 hover:scale-[1.01] animate-fadeIn">
            {isSuccess ? (
              <>
                <CheckCircle className="w-20 h-20 text-green-600 mx-auto mb-6" />
                <h1 className="text-4xl font-bold mb-3 text-green-700">Payment Successful</h1>
                <p className="text-gray-600 mb-8">
                  Thank you! Your payment was processed successfully.
                </p>
                <div className="text-left mb-8 space-y-3 text-gray-700">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Reference:</span>
                    <span className="font-semibold">{transaction?.reference}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Method:</span>
                    <span className="font-semibold">{transaction?.medium}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Amount:</span>
                    <span className="font-semibold">
                      {transaction?.amount} {transaction?.currency}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Date:</span>
                    <span className="font-semibold">
                      {new Date(transaction?.created_at || '').toLocaleString()}
                    </span>
                  </div>
                </div>
                <Link
                  href={merchantUrl || '/'}
                  className="inline-block px-8 py-3 bg-green-600 text-white rounded-xl hover:bg-green-700 hover:shadow-lg transition transform hover:scale-105"
                >
                  Go back to merchant
                </Link>
              </>
            ) : (
              <>
                <XCircle className="w-20 h-20 text-red-600 mx-auto mb-6" />
                <h1 className="text-4xl font-bold mb-3 text-red-700">Payment Failed</h1>
                <p className="text-gray-600 mb-8">
                  Unfortunately, your payment could not be completed.
                </p>
                <div className="text-left mb-8 space-y-3 text-gray-700">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Reference:</span>
                    <span className="font-semibold">{transaction?.reference}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Status:</span>
                    <span className="font-semibold">{transaction?.status}</span>
                  </div>
                </div>
                <Link
                  href={merchantUrl || '/'}
                  className="inline-block px-8 py-3 bg-red-600 text-white rounded-xl hover:bg-red-700 hover:shadow-lg transition transform hover:scale-105"
                >
                  Go back to merchant
                </Link>
              </>
            )}
          </div>
        )}
      </main>
    </>
  )
}