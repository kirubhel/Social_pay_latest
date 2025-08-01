'use client'

import Image from 'next/image'

export default function TransactionCompletePage() {
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
            {/* Processing Icon */}
            <div className="flex justify-center mb-6">
              <div className="w-20 h-20 bg-blue-100 rounded-full flex items-center justify-center">
                <svg 
                  className="w-12 h-12 text-blue-500" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" 
                  />
                </svg>
              </div>
            </div>

            {/* Complete Message */}
            <h1 className="text-2xl font-semibold text-gray-900 mb-4">
              Transaction Complete
            </h1>
            <p className="text-gray-600 mb-6">
              Your transaction has been submitted and is being processed.
            </p>
            
            {/* Status Check Information */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-8">
              <div className="flex items-start gap-3">
                <svg 
                  className="w-5 h-5 text-blue-500 mt-0.5 flex-shrink-0" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path 
                    strokeLinecap="round" 
                    strokeLinejoin="round" 
                    strokeWidth={2} 
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" 
                  />
                </svg>
                <div className="text-left">
                  <p className="text-blue-800 font-medium text-sm mb-1">
                    Check Transaction Status
                  </p>
                  <p className="text-blue-700 text-sm">
                    Please return to the checkout page to view your transaction status and details.
                  </p>
                </div>
              </div>
            </div>

            {/* Close Button */}
            <button
              onClick={() => window.close()}
              className="w-full bg-blue-500 text-white py-3 px-6 rounded-lg font-medium hover:bg-blue-600 transition-colors"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  )
} 