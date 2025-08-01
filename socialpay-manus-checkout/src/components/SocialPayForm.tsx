'use client'

import { useState } from "react"
import PaymentButtons from "./PaymentButtons"
import PhoneNumberInput from "./PhoneNumberInput"
import OTPVerificationModal from "./OTPVerificationModal"
import { usePayment } from "@/contexts/PaymentContext"
import { useQRPayment } from "@/contexts/QRPaymentContext"

interface SocialPayFormProps {
  qrLinkId?: string;
}

export default function SocialPayForm({ qrLinkId }: SocialPayFormProps) {
  const [phoneNumber, setPhoneNumber] = useState('')
  const [isValid, setIsValid] = useState(false)
  const [isOTPModalOpen, setIsOTPModalOpen] = useState(false)
  const [isPaymentProcessing, setIsPaymentProcessing] = useState(false)
  const [isPaymentComplete, setIsPaymentComplete] = useState(false)
  
  // Use the appropriate context based on whether qrLinkId is provided
  const regularPayment = usePayment();
  const qrPayment = useQRPayment();
  
  const paymentContext = qrLinkId ? qrPayment : regularPayment;
  const { amount } = paymentContext;

  function handlePhoneNumberChange(phoneNumber: string, isValid: boolean): void {
    setPhoneNumber(phoneNumber)
    setIsValid(isValid)
  }

  function handlePayNow(): void {
    // Open OTP verification modal
    setIsOTPModalOpen(true)
  }



  // Handle OTP verification
  async function handleVerifyOTP(otp: string): Promise<boolean> {
    // Simulate API call to verify OTP
    setIsPaymentProcessing(true)
    
    return new Promise((resolve) => {
      setTimeout(() => {
        // For demo purposes, consider any OTP valid except "000000"
        const isValid = otp !== "000000"
        
        if (isValid) {
          setIsPaymentComplete(true)
        }
        
        setIsPaymentProcessing(false)
        resolve(isValid)
      }, 1500)
    })
  }

  // Handle OTP modal close
  function handleOTPModalClose(): void {
    setIsOTPModalOpen(false)
    
    // If payment is complete, you might want to redirect or show a success message
    if (isPaymentComplete) {
      console.log("Payment completed successfully!")
      // You could redirect here or show a success message
    }
  }

  return (
    <div className="w-full max-w-[466px]">
      <div className="mt-2">
        <PhoneNumberInput 
          walletType='socialpay' 
          onChange={handlePhoneNumberChange}
          logoPath={'/socialpay.svg'}
        />
      </div>
      <div className="mt-6">
        <PaymentButtons 
          onPay={handlePayNow}
          disabled={!isValid || isPaymentProcessing}
          totalPrice={amount}
        />
      </div>
      
      {/* OTP Verification Modal */}
      <OTPVerificationModal
        isOpen={isOTPModalOpen}
        onClose={handleOTPModalClose}
        onVerify={handleVerifyOTP}
        phoneNumber={phoneNumber}
      />
    </div>
  )
} 