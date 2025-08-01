'use client'

import { use, useState, useEffect } from 'react'
import Image from 'next/image'
import { getReceipt, V2ClientError } from '@/services/merchant.service'
import { ReceiptResponse } from '@/types/merchant.dto'

export default function ReceiptPage(props: { params: Promise<{ id: string }> }) {
  const params = use(props.params)
  const [receipt, setReceipt] = useState<ReceiptResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchReceipt = async () => {
      try {
        const response = await getReceipt(params.id)
        setReceipt(response)
      } catch (err) {
        if (err instanceof V2ClientError) {
          setError(err.message)
        } else {
          setError('Failed to fetch receipt')
        }
      } finally {
        setIsLoading(false)
      }
    }

    fetchReceipt()
  }, [params.id])

  const handlePrint = () => {
    window.print()
  }

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(amount)
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-[#30BB54]"></div>
      </div>
    )
  }

  if (error || !receipt) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-800 mb-4">Receipt Not Found</h1>
          <p className="text-gray-600">{error || 'The requested receipt could not be found.'}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4">
        {/* Print Button - Hidden in print */}
        <div className="mb-6 print:hidden">
          <button
            onClick={handlePrint}
            className="bg-[#30BB54] text-white px-6 py-2 rounded-lg hover:bg-[#28a745] transition-colors flex items-center gap-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-5 h-5">
              <path fillRule="evenodd" d="M7.875 1.5C6.839 1.5 6 2.34 6 3.375v2.25C6 6.66 6.84 7.5 7.875 7.5h8.25C17.16 7.5 18 6.66 18 5.625V3.375C18 2.339 17.16 1.5 16.125 1.5h-8.25zM3.375 9C2.339 9 1.5 9.84 1.5 10.875V18a1.5 1.5 0 001.5 1.5h1.5v-3.375c0-.621.504-1.125 1.125-1.125h12.75c.621 0 1.125.504 1.125 1.125V19.5h1.5a1.5 1.5 0 001.5-1.5v-7.125C22.5 9.839 21.66 9 20.625 9H3.375z" clipRule="evenodd" />
              <path d="M5.25 18.75h13.5v-7.5H5.25v7.5z" />
            </svg>
            Print Receipt
          </button>
        </div>

        {/* Receipt Container */}
        <div className="bg-white rounded-lg shadow-lg overflow-hidden print:shadow-none print:rounded-none relative">
          {/* Header */}
          <div className="bg-gradient-to-r from-[#30BB54] to-[#28a745] text-white p-8 text-center">
            <div className="flex items-center justify-center mb-4">
              <Image
                src="/socialpay.webp"
                alt="SocialPay"
                width={120}
                height={35}
                className="brightness-0 invert"
              />
            </div>
            <h1 className="text-2xl font-bold mb-2">Payment Receipt</h1>
            <p className="text-green-100">Transaction Confirmation</p>
          </div>

          {/* Receipt Content */}
          <div className="p-8 relative">
            {/* Status Badge */}
            <div className="flex justify-center mb-6">
              <span className={`px-4 py-2 rounded-full text-sm font-medium ${
                receipt.status === 'SUCCESS' 
                  ? 'bg-green-100 text-green-800' 
                  : receipt.status === 'PENDING'
                  ? 'bg-yellow-100 text-yellow-800'
                  : 'bg-red-100 text-red-800'
              }`}>
                {receipt.status}
              </span>
            </div>

            {/* Transaction Details */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
              <div>
                <h3 className="text-lg font-semibold text-gray-800 mb-4">Transaction Details</h3>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm text-gray-600">Transaction ID</p>
                    <p className="font-mono text-sm font-medium">{receipt.id}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Transaction Type</p>
                    <p className="font-medium">{receipt.type}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Reference</p>
                    <p className="font-mono text-sm font-medium">{receipt.reference}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Payment Method</p>
                    <p className="font-medium">{receipt.medium}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Date & Time</p>
                    <p className="font-medium">{formatDate(receipt.created_at)}</p>
                  </div>
                </div>
              </div>

              <div>
                <h3 className="text-lg font-semibold text-gray-800 mb-4">Merchant Information</h3>
                <div className="space-y-3">
                  <div>
                    <p className="text-sm text-gray-600">Business Name</p>
                    <p className="font-medium">{receipt.merchant.trading_name || receipt.merchant.legal_name}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Legal Name</p>
                    <p className="font-medium">{receipt.merchant.legal_name}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Business Type</p>
                    <p className="font-medium capitalize">{receipt.merchant.business_type}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600">Registration Number</p>
                    <p className="font-medium">{receipt.merchant.business_registration_number}</p>
                  </div>
                </div>
              </div>
            </div>

            {/* Payment Breakdown */}
            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Payment Breakdown</h3>
              <div className="bg-gray-50 rounded-lg p-6">
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Amount</span>
                    <span className="font-medium">{formatCurrency(receipt.amount)} {receipt.currency}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Transaction Fee</span>
                    <span className="font-medium">{formatCurrency(receipt.fee_amount)} {receipt.currency}</span>
                  </div>
                  {receipt.vat_amount > 0 && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">VAT</span>
                      <span className="font-medium">{formatCurrency(receipt.vat_amount)} {receipt.currency}</span>
                    </div>
                  )}
                  <div className="border-t border-gray-300 pt-3">
                    <div className="flex justify-between text-lg font-semibold">
                      <span>Total Amount</span>
                      <span className="text-[#30BB54]">{formatCurrency(receipt.total_amount)} {receipt.currency}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Description */}
            {receipt.description && (
              <div className="border-t border-gray-200 pt-6 mt-6">
                <h3 className="text-lg font-semibold text-gray-800 mb-2">Description</h3>
                <p className="text-gray-600">{receipt.description}</p>
              </div>
            )}

            {/* Footer */}
            <div className="border-t border-gray-200 pt-6 mt-8 text-center text-sm text-gray-500">
              <p>This is an electronic receipt generated by SocialPay</p>
              <p className="mt-1">For support, please contact us at support@socialpay.co</p>
            </div>

            {/* Stamp - positioned in lower right corner */}
            <div className="absolute bottom-20 right-16">
              <Image
                src="/stamp-new.png"
                alt="Official Stamp"
                width={150}
                height={150}
                className="opacity-80"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Print Styles */}
      <style jsx global>{`
        @media print {
          body {
            background: white !important;
          }
          .print\\:hidden {
            display: none !important;
          }
          .print\\:shadow-none {
            box-shadow: none !important;
          }
          .print\\:rounded-none {
            border-radius: 0 !important;
          }
        }
      `}</style>
    </div>
  )
} 