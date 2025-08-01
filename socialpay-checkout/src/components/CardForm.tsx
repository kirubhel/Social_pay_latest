/* eslint-disable @typescript-eslint/no-unused-vars */
'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import { useTranslations } from 'next-intl'
import PaymentButtons from './PaymentButtons';
import TipSection, { WalletType } from './TipSection';
import PhoneNumberInput from './PhoneNumberInput';
import { usePayment } from '@/contexts/PaymentContext';
import { useMerchantPayment } from '@/contexts/MerchantPaymentContext';
import { useCheckoutPayment } from '@/contexts/CheckoutPaymentContext';
import { useQRPayment } from '@/contexts/QRPaymentContext';
import { Gateway } from '@/services/gateway.service';

// Card types and their regex patterns
const cardTypes = [
  // Mastercard: Starts with 51-55 or 2221-2720, length 16
  { type: 'mastercard', pattern: /^(5[1-5][0-9]{14}|2(22[1-9][0-9]{12}|2[3-9][0-9]{13}|[3-6][0-9]{14}|7[0-1][0-9]{13}|720[0-9]{12}))$/ },
  // Visa: Starts with 4, length 13, 16, or 19
  { type: 'visa', pattern: /^4[0-9]{12}(?:[0-9]{3})?(?:[0-9]{3})?$/ },
  // American Express: Starts with 34 or 37, length 15
  { type: 'amex', pattern: /^3[47][0-9]{13}$/ },
  // Discover: Starts with 6011, 622126-622925, 644-649, 65, length 16-19
  { type: 'discover', pattern: /^6(?:011|5[0-9]{2}|4[4-9][0-9]|22(?:1(?:2[6-9]|[3-9][0-9])|[2-8][0-9]{2}|9(?:[01][0-9]|2[0-5])))[0-9]{10,13}$/ },
  // Diners Club: Starts with 300-305, 36, 38-39, length 14-19
  { type: 'diners', pattern: /^3(?:0[0-5]|[68][0-9])[0-9]{11,16}$/ },
  // JCB: Starts with 2131, 1800, or 35, length 16-19
  { type: 'jcb', pattern: /^(?:2131|1800|35[0-9]{3})[0-9]{11,16}$/ }
];

// Early detection patterns (for partial numbers)
const earlyDetectionPatterns = [
  { type: 'visa', pattern: /^4/ },
  { type: 'mastercard', pattern: /^(5[1-5]|2[2-7])/ },
  { type: 'amex', pattern: /^3[47]/ },
  { type: 'discover', pattern: /^(6011|65|64[4-9]|622)/ },
  { type: 'diners', pattern: /^(30[0-5]|36|38|39)/ },
  { type: 'jcb', pattern: /^(2131|1800|35)/ }
];

interface CardFormProps {
  gateways: Gateway[];
  merchantId?: string;
  checkoutId?: string;
  qrLinkId?: string;
  isStaticAmount?: boolean;
  staticAmount?: number;
  initialPhoneNumber?: string;
  isTipEnabled?: boolean;
}

export default function CardForm({ gateways, merchantId, checkoutId, qrLinkId, isStaticAmount, staticAmount, initialPhoneNumber, isTipEnabled }: CardFormProps) {
  const [cardNumber, setCardNumber] = useState('')
  const [expiryDate, setExpiryDate] = useState('')
  const [cvv, setCvv] = useState('')
  const [name, setName] = useState('')
  const [saveCard, setSaveCard] = useState(true)
  const [phoneNumber, setPhoneNumber] = useState(initialPhoneNumber || '')
  const [isPhoneValid, setIsPhoneValid] = useState(false)
  const [selectedWallet] = useState<string>('CYBERSOURCE') // Default for card payments
  const t = useTranslations('checkout');
  const [selectedGateway, setSelectedGateway] = useState<Gateway | null>(null)
  // Validation states
  const [cardType, setCardType] = useState<string | null>(null)
  const [isCardNumberValid, setIsCardNumberValid] = useState<boolean | null>(null)
  const [isExpiryDateValid, setIsExpiryDateValid] = useState<boolean | null>(null)
  const [isCvvValid, setIsCvvValid] = useState<boolean | null>(null)
  const [isNameValid, setIsNameValid] = useState<boolean | null>(null)
  const [isFocused, setIsFocused] = useState<string | null>(null)

  // Use the appropriate context based on what props are provided
  const regularPayment = usePayment();
  const merchantPayment = useMerchantPayment();
  const checkoutPayment = useCheckoutPayment();
  const qrPayment = useQRPayment();
  
  const paymentContext = qrLinkId ? qrPayment : checkoutId ? checkoutPayment : merchantId ? merchantPayment : regularPayment;
  const { amount, processPayment, isLoading } = paymentContext;
  
  // Only get tip data if it's a QR payment context
  const tipAmount = qrLinkId && 'tipAmount' in paymentContext ? paymentContext.tipAmount : '';
  const tipeePhone = qrLinkId && 'tipeePhone' in paymentContext ? paymentContext.tipeePhone : '';
  const tipMedium = qrLinkId && 'tipMedium' in paymentContext ? paymentContext.tipMedium : '';

  // Format card number with spaces
  const formatCardNumber = (value: string) => {
    const numbers = value.replace(/\D/g, '')
    
    // Format in groups of 4 as the user types
    let formatted = '';
    for (let i = 0; i < numbers.length; i += 4) {
      const chunk = numbers.slice(i, i + 4);
      if (formatted) formatted += ' ';
      formatted += chunk;
    }
    
    return formatted;
  }

  // Format expiry date as MM/YY
  const formatExpiryDate = (value: string) => {
    const numbers = value.replace(/\D/g, '')
    if (numbers.length <= 2) return numbers;
    return `${numbers.slice(0, 2)}/${numbers.slice(2, 4)}`;
  }

  // Detect card type based on number
  const detectCardType = (number: string) => {
    const cleanNumber = number.replace(/\D/g, '');
    
    // Early detection for partial numbers
    for (const card of earlyDetectionPatterns) {
      if (card.pattern.test(cleanNumber)) {
        return card.type;
      }
    }
    
    // Full validation for completed numbers
    for (const card of cardTypes) {
      if (card.pattern.test(cleanNumber)) {
        return card.type;
      }
    }
    
    return null;
  }

  // Validate card number using Luhn algorithm
  const validateCardNumber = (number: string) => {
    const cleanNumber = number.replace(/\D/g, '');
    if (cleanNumber.length < 13) return false;
    
    let sum = 0;
    let shouldDouble = false;
    
    // Loop through values starting from the rightmost digit
    for (let i = cleanNumber.length - 1; i >= 0; i--) {
      let digit = parseInt(cleanNumber.charAt(i));
      
      if (shouldDouble) {
        digit *= 2;
        if (digit > 9) digit -= 9;
      }
      
      sum += digit;
      shouldDouble = !shouldDouble;
    }
    
    return sum % 10 === 0;
  }

  // Validate expiry date
  const validateExpiryDate = (date: string) => {
    const [monthStr, yearStr] = date.split('/');
    if (!monthStr || !yearStr) return false;
    
    const month = parseInt(monthStr, 10);
    const year = parseInt(`20${yearStr}`, 10);
    
    // Check if month is valid
    if (month < 1 || month > 12) return false;
    
    const now = new Date();
    const currentYear = now.getFullYear();
    const currentMonth = now.getMonth() + 1; // JavaScript months are 0-indexed
    
    // Check if the card is expired
    if (year < currentYear) return false;
    if (year === currentYear && month < currentMonth) return false;
    
    return true;
  }

  // Validate CVV
  const validateCvv = (cvvValue: string, cardTypeValue: string | null) => {
    const cleanCvv = cvvValue.replace(/\D/g, '');
    
    // Amex CVV is 4 digits, others are 3
    if (cardTypeValue === 'amex') {
      return cleanCvv.length === 4;
    }
    
    return cleanCvv.length === 3;
  }

  // Validate name
  const validateName = (nameValue: string) => {
    return nameValue.trim().length >= 3 && /^[a-zA-Z\s]+$/.test(nameValue);
  }

  // Handle card number change
  const handleCardNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const cleanValue = value.replace(/\D/g, '');
    
    // Limit length based on card type
    let maxLength = 16;
    if (detectCardType(cleanValue) === 'amex') maxLength = 15;
    
    if (cleanValue.length <= maxLength) {
      const formattedValue = formatCardNumber(cleanValue);
      setCardNumber(formattedValue);
      
      // Detect card type
      const detectedType = detectCardType(cleanValue);
      setCardType(detectedType);
      
      // Validate card number
      if (cleanValue.length >= 13) {
        setIsCardNumberValid(validateCardNumber(cleanValue));
      } else {
        setIsCardNumberValid(null);
      }
    }
  }

  // Handle expiry date change
  const handleExpiryDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const cleanValue = value.replace(/\D/g, '');
    
    if (cleanValue.length <= 4) {
      const formattedValue = formatExpiryDate(cleanValue);
      setExpiryDate(formattedValue);
      
      // Validate expiry date
      if (formattedValue.includes('/') && formattedValue.length === 5) {
        setIsExpiryDateValid(validateExpiryDate(formattedValue));
      } else {
        setIsExpiryDateValid(null);
      }
    }
  }

  // Handle CVV change
  const handleCvvChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    const cleanValue = value.replace(/\D/g, '');
    
    // Limit length based on card type
    const maxLength = cardType === 'amex' ? 4 : 3;
    
    if (cleanValue.length <= maxLength) {
      setCvv(cleanValue);
      
      // Validate CVV
      if (cleanValue.length === maxLength) {
        setIsCvvValid(validateCvv(cleanValue, cardType));
      } else {
        setIsCvvValid(null);
      }
    }
  }

  // Handle name change
  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setName(value);
    
    // Validate name
    if (value.length > 0) {
      setIsNameValid(validateName(value));
    } else {
      setIsNameValid(null);
    }
  }

  // Handle phone number change
  const handlePhoneNumberChange = (number: string, isValid: boolean) => {
    setPhoneNumber(number);
    setIsPhoneValid(isValid);
  }



  // This function is used by the parent component through the onSubmit prop
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    console.log('handleSubmit', "merchantId", merchantId, "checkoutId", checkoutId, "qrLinkId", qrLinkId)
    // Validate all fields before submission
    const isValid = true;
    
    if (isValid) {
      // For QR payments, include tip data if tipAmount is available
      if (qrLinkId && 'tipAmount' in paymentContext) {
        const selectedWallet = 'TELEBIRR'; // Default wallet for card payments
        const tipData = parseFloat(tipAmount) > 0 
          ? {
              amount: parseFloat(tipAmount),
              phone: tipeePhone,
              medium: tipMedium
            }
          : undefined;
        
        // Use static amount for static QR codes, otherwise use the amount from context
        const paymentAmount = isStaticAmount ? staticAmount : parseFloat(amount);
        const qrProcessPayment = paymentContext.processPayment as (medium: string, qrLinkId: string, phone: string, amount?: number, tipData?: { amount: number, phone: string, medium: string }) => Promise<void>;
      
  
        if (selectedGateway?.key=="ETHSWITCH"){
          await qrProcessPayment('ETHSWITCH', qrLinkId, phoneNumber, paymentAmount, tipData);
        }

        if (selectedGateway?.key=="CYBERSOURCE"){
           await qrProcessPayment('CYBERSOURCE', qrLinkId, phoneNumber, paymentAmount, tipData);
        }

      } else {
        // For other payment types, use the appropriate processPayment signature
        if (merchantId) {
          // MerchantPaymentContextType: processPayment(medium: string, merchantId: string, phone: string)
          const merchantProcessPayment = paymentContext.processPayment as (medium: string, merchantId: string, phone: string) => Promise<void>;
          await merchantProcessPayment('CYBERSOURCE', merchantId, phoneNumber);
        } else if (checkoutId) {
          // CheckoutPaymentContextType: processPayment(medium: string, checkoutId: string, phone: string, tipData?: { amount: number, phone: string, medium: string })
          console.log('submit payment checkoutId', checkoutId, )
          const checkoutProcessPayment = paymentContext.processPayment as (medium: string, checkoutId: string, phone: string, tipData?: { amount: number, phone: string, medium: string }) => Promise<void>;
          
          // Get tip data if available
          const tipData = checkoutId && 'tipAmount' in paymentContext && parseFloat(paymentContext.tipAmount) > 0
            ? {
                amount: parseFloat(paymentContext.tipAmount),
                phone: paymentContext.tipeePhone,
                medium: paymentContext.tipMedium
              }
            : undefined;
         
           if (selectedGateway?.key == "ETHSWITCH"){
               await checkoutProcessPayment('ETHSWITCH', checkoutId, phoneNumber);
           }

           if (selectedGateway?.key == "CYBERSOURCE"){
              await checkoutProcessPayment('CYBERSOURCE', checkoutId, phoneNumber);

           }
         
    
        } else {
          // PaymentContextType: processPayment(medium: string, details: PaymentDetails)
          const regularProcessPayment = paymentContext.processPayment as (medium: string, details: { card?: { name: string; pan: string; expiry: string; }; phone?: string; }) => Promise<void>;
          await regularProcessPayment('CYBERSOURCE', {
            card: {
              name: name || 'John Doe',
              pan: cardNumber.replace(/\s/g, ''),
              expiry: expiryDate
            },
            phone: phoneNumber
          });
        }
      }
    }
  }

   // Map gateway keys to banner image paths (update these URLs to your actual banner images)
  const bannerImages: Record<string, string> = {
    ETHSWITCH: '/card/bankslogo.png', // Your banner with EthSwitch-supported banks
    CYBERSOURCE: '/card/visamastercard.png' // Your banner with CyberSource cards
  }




  return (
    <div className="w-full max-w-[396px] pt-4">
      <form onSubmit={handleSubmit}>


   {/* SIDE BY SIDE GATEWAY SELECTOR WITH BANNERS */}
      <div className="flex flex-row md:flex-row gap-5 mb-8 pt-4 ">
        {gateways.map(gateway => (
            <div
              key={gateway.key}
              onClick={() => setSelectedGateway(gateway)}
              className={`group relative flex flex-col items-center justify-center
                ${gateways.length >1?'w-full':'w-60'} h-32 p-4 cursor-pointer rounded-xl border transition-all duration-300
                ${selectedGateway?.key === gateway.key ? 
                'border-[#2B3A67] ring-4 ring-[#eaf3fa] bg-white shadow-xl' : 
                'border-gray-300 bg-white hover:border-[#f5a414] hover:shadow-lg'}
              `}
            >
              {/* Logo */}
              <div className="relative w-full h-full mb-2">
                <Image
                  src={bannerImages[gateway.key] || gateway.icon}
                  alt={`${gateway.name} logo`}
                  fill
                  className="object-contain"
                />
              </div>
              {/* Payment Method Name */}
            <div className="text-sm font-semibold text-[#2B3A67] group-hover:text-[#f5a414]">
                {gateway.name}
              </div>
            </div>
        ))}
      </div>



      {/* Phone Number Input */}
      <div className="mt-2">
       { qrLinkId && <PhoneNumberInput 
          initialPhoneNumber={initialPhoneNumber}
          walletType={'SocialPay' as WalletType}
          onChange={handlePhoneNumberChange}
        />
        }
      </div>

        {/* Tip Section for QR payments */}
        {qrLinkId && (
          <TipSection 
            gateways={gateways}
            isTipEnabled={true}
          />
        )}
        
        {/* Tip Section for checkout payments */}
        {checkoutId && isTipEnabled && (
          <TipSection 
            gateways={gateways}
            isTipEnabled={isTipEnabled}
            contextType="checkout"
          />
        )}

          <PaymentButtons 
            onPay={handleSubmit}
            disabled={!(selectedGateway && phoneNumber.length ===9)}
            totalPrice={
              // For QR payments (existing logic)
              qrLinkId 
                ? (isStaticAmount && staticAmount 
                    ? (staticAmount + (parseFloat(tipAmount) > 0 ? parseFloat(tipAmount) : 0)).toString()
                    : (parseFloat(amount || '0') + (parseFloat(tipAmount) > 0 ? parseFloat(tipAmount) : 0)).toString()
                  )
                // For checkout payments (new logic)
                : checkoutId && 'tipAmount' in paymentContext
                  ? (parseFloat(amount || '0') + (parseFloat(paymentContext.tipAmount) > 0 ? parseFloat(paymentContext.tipAmount) : 0)).toString()
                  : amount || '0'
            }
            isLoading={isLoading}
          />
      </form>
    </div>
  )
}