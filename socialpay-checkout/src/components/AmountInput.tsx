"use client"

import React, { useState } from 'react';

interface AmountInputProps {
  value: string;
  onChange: (value: string, isValid: boolean) => void;
}

const AmountInput: React.FC<AmountInputProps> = ({ value, onChange }) => {
  const [amount, setAmount] = useState(value);
  const [isFocused, setIsFocused] = useState(false);
  const [isValid, setIsValid] = useState<boolean | null>(null);

  const validateAmount = (value: string): boolean => {
    if (!value) return true;
    
    // Check format: numbers with optional 2 decimal places
    const regex = /^\d*\.?\d{0,2}$/;
    if(!regex.test(value)) return false;
    const numericValue = parseFloat(value);
    return !isNaN(numericValue) && numericValue >= 1 && numericValue <= 9999999;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    const isValidAmount = validateAmount(newValue);
    setIsValid(isValidAmount);
    if (isValidAmount) {
      setAmount(newValue);
      onChange(newValue, true);
    }
  };

  return (
    <div className="w-full">
      <label className="block text-[#505050] text-[14px] font-medium mb-2">
        Amount
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
          <input
            type="text"
            value={amount}
            onChange={handleChange}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            className="bg-transparent outline-none border-none flex-1 text-[18px] font-semibold text-[#1a3129] focus:border-none focus:outline-none active:border-none active:outline-none placeholder-[#8A8A8A]"
            placeholder="Enter amount"
            style={{ border: 'none', outline: 'none' }}
          />
          {/* Currency suffix */}
          <span className="text-[#8A8A8A] text-[15px] font-medium ml-2">ETB</span>
        </div>

        {/* Error message */}
        {isValid === false && (
          <div className="flex items-center gap-2 mt-2 text-red-500 text-sm">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
              <path fillRule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clipRule="evenodd" />
            </svg>
            <span>Please enter a valid amount</span>
          </div>
        )}
      </div>
      
      {/* Notification message */}
      <div className="flex items-center gap-2 mt-2 text-[#505050] text-[14px]">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4 text-gray-400">
          <path fillRule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clipRule="evenodd" />
        </svg>
        <span className="text-[12px]">Please enter amount to complete your transaction</span>
      </div>
    </div>
  );
};

export default AmountInput; 