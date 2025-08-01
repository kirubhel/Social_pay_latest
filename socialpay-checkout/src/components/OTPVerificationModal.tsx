'use client'

import { useState, useEffect, useRef, KeyboardEvent, ClipboardEvent } from 'react'
import Image from 'next/image'

interface OTPVerificationModalProps {
  isOpen: boolean
  onClose: () => void
  onVerify: (otp: string) => Promise<boolean>
  phoneNumber?: string // Make phoneNumber optional since it's not used yet
}

export default function OTPVerificationModal({ 
  isOpen, 
  onClose, 
  onVerify
}: OTPVerificationModalProps) {
  const [otp, setOtp] = useState<string[]>(Array(6).fill(''))
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isError, setIsError] = useState(false)
  const [isSuccess, setIsSuccess] = useState(false)
  const [isResending, setIsResending] = useState(false)
  
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])
  
  // Reset state when modal opens
  useEffect(() => {
    if (isOpen) {
      setOtp(Array(6).fill(''))
      setIsSubmitting(false)
      setIsError(false)
      setIsSuccess(false)
      setIsResending(false)
      
      // Focus first input when modal opens
      setTimeout(() => {
        if (inputRefs.current[0]) {
          inputRefs.current[0]?.focus()
        }
      }, 100)
    }
  }, [isOpen])
  
  // Handle input change
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>, index: number) => {
    const value = e.target.value
    
    // Only allow digits
    if (!/^\d*$/.test(value)) return
    
    // Update OTP array
    const newOtp = [...otp]
    newOtp[index] = value.substring(0, 1)
    setOtp(newOtp)
    
    // Auto-focus next input if value is entered
    if (value && index < 5 && inputRefs.current[index + 1]) {
      inputRefs.current[index + 1]?.focus()
    }
    
    // Check if all digits are filled
    if (!newOtp.includes('') && newOtp.every(digit => digit.length === 1)) {
      handleSubmit(newOtp.join(''))
    }
  }
  
  // Handle key press
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>, index: number) => {
    // Move to previous input on backspace if current input is empty
    if (e.key === 'Backspace' && !otp[index] && index > 0 && inputRefs.current[index - 1]) {
      inputRefs.current[index - 1]?.focus()
    }
  }
  
  // Handle paste
  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault()
    const pastedData = e.clipboardData.getData('text/plain').trim()
    
    // Check if pasted content is a 6-digit number
    if (/^\d{6}$/.test(pastedData)) {
      const newOtp = pastedData.split('')
      setOtp(newOtp)
      
      // Focus last input
      if (inputRefs.current[5]) {
        inputRefs.current[5]?.focus()
      }
      
      // Submit OTP
      handleSubmit(pastedData)
    }
  }
  
  // Submit OTP
  const handleSubmit = async (otpValue: string) => {
    if (isSubmitting || isSuccess) return
    
    setIsSubmitting(true)
    setIsError(false)
    
    try {
      // Call onVerify but ignore the result for testing purposes
      await onVerify(otpValue);
      
      // For testing: randomly alternate between success and error
      const randomSuccess = Math.random() > 0.5
      
      if (randomSuccess) {
        setIsSuccess(true)
        // Don't close modal automatically, wait for user to click confirm
      } else {
        setIsError(true)
      }
    } catch {
      // Set error state if verification fails
      setIsError(true)
    } finally {
      setIsSubmitting(false)
    }
  }
  
  // Resend OTP
  const handleResend = () => {
    setIsResending(true)
    
    // Reset OTP fields
    setOtp(Array(6).fill(''))
    setIsError(false)
    
    // Simulate resend (replace with actual resend logic)
    setTimeout(() => {
      setIsResending(false)
      // Focus first input
      if (inputRefs.current[0]) {
        inputRefs.current[0]?.focus()
      }
    }, 1500)
  }
  
  if (!isOpen) return null
  
  return (
    <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
      <div className="bg-white rounded-[16px] w-[380px] relative shadow-[0px_4px_20px_rgba(0,0,0,0.25)] p-8">
        {/* X Close Button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-400 hover:text-gray-700 text-2xl font-bold focus:outline-none"
          aria-label="Close"
        >
          &times;
        </button>
        {/* OTP Icon */}
        <div className="flex justify-center mb-8">
          <div className="relative w-[60px] h-[60px]">
            <Image 
              src={isError ? "/submition-error.svg" : isSuccess ? "/submition-success.svg" : "/otp-icon.svg"} 
              alt="Verification" 
              fill
              style={{ objectFit: 'contain' }}
            />
          </div>
        </div>
        
        {/* OTP Input */}
        <div className="flex justify-center gap-2 mb-6">
          {otp.map((digit, index) => (
            <input
              key={index}
              type="text"
              value={digit}
              onChange={(e) => handleChange(e, index)}
              onKeyDown={(e) => handleKeyDown(e, index)}
              onPaste={index === 0 ? handlePaste : undefined}
              ref={(el) => {
                inputRefs.current[index] = el
              }}
              className={`
                w-[40px] h-[40px] rounded-md text-center text-2xl font-bold
                bg-[#f2f2f2]
                border ${isError ? 'border-red-500' : 'border-gray-200'}
                focus:outline-none focus:border-[#30BB54] focus:ring-1 focus:ring-[#30BB54]
              `}
              maxLength={1}
              disabled={isSubmitting || isSuccess}
            />
          ))}
        </div>
        {/* Title */}
        <h2 className="text-center text-[16px] leading-[16px] font-medium mb-6">Verification Code</h2>
        {/* Description - Moved below OTP input */}
        <p className="px-[0px] text-center text-[#505050] mb-10 px-4 text-[12px] leading-[12px]">
          To complete your request, a 6-digit verification code has been sent
          to your mobile number. Please enter the code to confirm.
        </p>
        
        {/* Error message - Only show when there's an error */}
        {isError && (
          <p className="text-center text-red-500 font-medium mb-4 text-[10px]">
            The numbers you entered are incorrect!
          </p>
        )}
        
        {/* Buttons */}
        <div className="flex justify-center">
          {isSuccess ? (
            /* Confirm Button - Only show if verification was successful */
            <button
              onClick={onClose}
              className="bg-[#30BB54] text-white py-3 px-10 rounded-md font-medium w-[200px] text-[14.12px] leading-[14.12px]"
            >
              Confirm
            </button>
          ) : (
            /* Resend Button - Show when not successful yet */
            <button
              onClick={handleResend}
              disabled={isResending || isSubmitting}
              className={`
                py-3 px-10 rounded-md font-medium text-white w-[200px]
                ${isError ? 'bg-red-500' : 'bg-black'}
                disabled:opacity-50
                text-[14.12px] leading-[14.12px]
              `}
            >
              {isResending ? 'Sending...' : 'Resend Code'}
            </button>
          )}
        </div>
      </div>
    </div>
  )
} 