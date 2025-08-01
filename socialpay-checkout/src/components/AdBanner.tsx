'use client'

import Image from 'next/image'
import { useEffect, useState } from 'react'
import { checkoutService, AdResponse } from '@/services/checkout.service'

export default function AdBanner() {
  // Use null as initial state to avoid hydration mismatch
  const [adData, setAdData] = useState<AdResponse | null>(null)
  // Start with true for loading to ensure consistent rendering
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)
  // Add a mounted state to prevent hydration mismatch
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
    
    const fetchAd = async () => {
      try {
        setLoading(true)
        console.log('Fetching ad data...')
        const data = await checkoutService.fetchAds('checkout-banner')
        console.log('Ad data received:', data)
        setAdData(data)
      } catch (error) {
        console.error('Failed to fetch ad:', error)
        setError(error instanceof Error ? error : new Error('Unknown error'))
      } finally {
        setLoading(false)
      }
    }

    fetchAd()
  }, [])

  // Return a placeholder during server rendering to ensure consistency
  if (!mounted) {
    return (
      <div className="w-full relative overflow-hidden rounded-2xl">
        <div className="w-full h-24 bg-gray-200" />
      </div>
    )
  }

  if (error) {
    console.error('AdBanner error:', error)
    return (
      <div className="w-full relative overflow-hidden rounded-2xl bg-gray-100 flex items-center justify-center">
        <div className="w-full h-24 flex items-center justify-center text-gray-500 text-sm">
          Advertisement unavailable
        </div>
      </div>
    )
  }

  return (
    <div className="w-full relative overflow-hidden rounded-2xl">
      {loading ? (
        // Shimmer effect - minimal height for loading state
        <div className="w-full h-24 animate-pulse bg-gradient-to-r from-gray-200 via-gray-300 to-gray-200 background-animate" />
      ) : adData ? (
        <a 
          href={adData.linkUrl || '#'} 
          className="block relative w-full"
          target={adData.linkUrl && adData.linkUrl !== '#' ? '_blank' : '_self'}
          rel={adData.linkUrl && adData.linkUrl !== '#' ? 'noopener noreferrer' : undefined}
        >
          <div className="relative w-full">
          <Image
            src={adData.imagePath}
            alt={adData.altText || 'Advertisement'}
              width={1200}
              height={0}
              className="w-full h-auto"
              sizes="100vw"
            priority
              onError={() => {
                console.error('Image failed to load:', adData.imagePath)
                setError(new Error('Failed to load banner image'))
              }}
              onLoad={() => {
                console.log('Image loaded successfully:', adData.imagePath)
              }}
          />
          </div>
        </a>
      ) : (
        <div className="w-full h-24 bg-gray-100 flex items-center justify-center text-gray-500 text-sm">
          No advertisement available
        </div>
      )}
    </div>
  )
} 