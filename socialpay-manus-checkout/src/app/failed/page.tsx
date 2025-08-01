'use client'

import Image from 'next/image'

export default function TransactionFailedPage() {
  return (
    <div className="min-h-screen bg-[rgb(248,248,248)] flex flex-col">
      {/* Header */}
      <div className="w-full px-5 py-3 mt-5">
        <div className="max-w-[1240px] mx-auto flex justify-center items-center">
          <Image
            src="/socialpay.webp"
            alt="SocialPay logo"
            width={148}
            height={43}
            priority
          />
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 flex items-center justify-center px-4">
        <div className="w-full max-w-md text-center">
          <div className="bg-white rounded-lg shadow-lg p-8">
            {/* Error Icon */}
            <div className="flex justify-center mb-6">
              <div className="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center">
                <svg 
                  className="w-12 h-12 text-red-500" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M6 18L18 6M6 6l12 12" 
                  />
                </svg>
              </div>
            </div>

            {/* Error Message */}
            <h1 className="text-2xl font-semibold text-gray-900 mb-4">
              Transaction Failed
            </h1>
            <p className="text-gray-600 mb-8">
              Your transaction could not be completed. Please try again.
            </p>

            {/* Close Button */}
            <button
              onClick={() => window.close()}
              className="w-full bg-red-500 text-white py-3 px-6 rounded-lg font-medium hover:bg-red-600 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  )
} 