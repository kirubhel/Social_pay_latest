'use client'

import Image from 'next/image'

export default function TransactionSuccessPage() {
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
            {/* Success Icon */}
            <div className="flex justify-center mb-6">
              <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center">
                <svg 
                  className="w-12 h-12 text-green-500" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M5 13l4 4L19 7" 
                  />
                </svg>
              </div>
            </div>

            {/* Success Message */}
            <h1 className="text-2xl font-semibold text-gray-900 mb-4">
              Transaction Successful
            </h1>
            <p className="text-gray-600 mb-8">
              Your transaction has been completed successfully.
            </p>

            {/* Close Button */}
            <button
              onClick={() => window.close()}
              className="w-full bg-[#30BB54] text-white py-3 px-6 rounded-lg font-medium hover:bg-[#28a745] transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  )
} 