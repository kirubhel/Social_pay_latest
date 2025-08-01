'use client'

import React, { useState, useEffect } from 'react'
import Image from 'next/image'
import { WalletType } from './WalletForm'
import { BankType } from './BankForm'
import { useTranslations } from 'next-intl'

interface PhoneNumberInputProps {
  walletType: WalletType | BankType | 'socialpay'
  onChange?: (phoneNumber: string, isValid: boolean) => void
  logoPath?: string,
  initialPhoneNumber?: string
}

export default function PhoneNumberInput({ walletType, onChange, logoPath, initialPhoneNumber }: PhoneNumberInputProps) {
  const [phoneNumber, setPhoneNumber] = useState(initialPhoneNumber || '')
  const [isValid, setIsValid] = useState<boolean | null>(null)
  const [isFocused, setIsFocused] = useState(false)
  const t = useTranslations('checkout')

  // Validate phone number based on wallet type
  const validatePhoneNumber = React.useCallback((number: string) => {

    if (!number || number.length === 0) {
      console.log("Empty number, returning null");
      return null;
    }
    
    // Must be 9 digits
    if (number.length !== 9) {
      return false;
    }
    
    // For Telebirr, CBE, EBIRR: must start with 9
    if (walletType != "MPESA") {
      const isValid = number.startsWith("9");
      return isValid;
    }
    
    // For MPesa: must start with 7
    if (walletType === "MPESA") {
      const isValid = number.startsWith("7");
      return isValid;
    }

    if (walletType === "socialpay") {
      const isValid = number.startsWith("7") || number.startsWith("9");
      return isValid;
    }
    
    return false;
  }, [walletType]);

  // Handle phone number change
  const handlePhoneNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    
    // Only allow digits
    const digitsOnly = value.replace(/\D/g, '')
    
    // Limit to 9 digits (plus any leading zeros)
    if (digitsOnly.length <= 10) { // Increased to 10 to allow for a leading zero
      setPhoneNumber(digitsOnly)
      
      // For validation, remove leading zeros but keep original input
      const withoutLeadingZeros = digitsOnly.replace(/^0+/, '')
      
      const validationResult = validatePhoneNumber(withoutLeadingZeros)
      
      setIsValid(validationResult)
      
      if (onChange) {
        // Pass the original input with leading zeros, but also pass the cleaned version for validation
        onChange(digitsOnly, validationResult === true)
      }
    }
  }

  // Update validation when wallet type changes
  useEffect(() => {
    if (phoneNumber) {
      const withoutLeadingZeros = phoneNumber.replace(/^0+/, '')
      const validationResult = validatePhoneNumber(withoutLeadingZeros)
      setIsValid(validationResult)
      
      if (onChange) {
        onChange(phoneNumber, validationResult === true)
      }
    }
  }, [walletType, phoneNumber, onChange, validatePhoneNumber])

  return (
    <div className="w-full">
      <label className="block text-[#505050] text-[14px] font-medium mb-2">
        {t('phoneNumber')}
      </label>
      
      <div className="relative">
        <div 
          className={`
            w-full h-[52px] bg-white rounded-xl px-5
            flex items-center
            border transition-all duration-200
            ${isValid === false ? 'border-red-500' : isFocused ? 'border-[#2B3A67] shadow-md' : 'border-[#e3e8ee]'}
          `}
        >
          {/* Flag and country code */}
          <div className="flex items-center gap-2 mr-4">
            <div className="w-[15px] h-[14px] relative">
              <Image 
                src="/ethiopia-flag.svg" 
                alt="Ethiopia" 
                fill
                style={{ objectFit: 'contain' }}
              />
            </div>
            <span className="text-[#1a3129] text-[15px] font-semibold">+251</span>
          </div>
          {/* Vertical divider */}
          <div className="h-[18px] w-[1px] bg-[#D1D1D1] mr-4"></div>
          {/* Phone number input */}
          <input
            type="text"
            value={phoneNumber}
            onChange={handlePhoneNumberChange}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            className="bg-transparent outline-none border-none flex-1 text-[16px] font-semibold text-[#1a3129] focus:border-none focus:outline-none active:border-none active:outline-none placeholder-[#8A8A8A]"
            placeholder="900000000"
            style={{ border: 'none', outline: 'none' }}
          />
          
          {/* Wallet logo */}
          {logoPath && (
            <div className={`relative transition-none ${walletType === "TELEBIRR" ? "w-[40px] h-[20px]" : walletType === "socialpay" ? "w-[70px] h-[40px]" : "w-[60px] h-[40px]"}`}>
              <Image 
                src={logoPath} 
                alt={`${walletType} logo`} 
                fill
                style={{ objectFit: 'contain' }}
              />
            </div>
          )}
        </div>
        
        {/* Error message */}
        {isValid === false && (
          <div className="flex items-center gap-2 mt-2 text-red-500 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
              <path fillRule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clipRule="evenodd" />
            </svg>
            <span>{t('phoneNumberInvalid')}</span>
          </div>
        )}
      </div>
      
      {/* Notification message */}
      <div className="flex items-center gap-2 mt-3 text-[#505050] text-[14px]">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4 text-gray-400">
          <path fillRule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clipRule="evenodd" />
        </svg>
        <span className="text-[12px]">{t('phoneNumberRequired')}</span>
      </div>
    </div>
  )
} 