import React, { forwardRef, useState } from 'react'
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline'

interface InputFieldProps {
  label: string
  type?: 'text' | 'email' | 'password' | 'tel' | 'url' | 'number'
  value: string
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  disabled?: boolean
  error?: string
  success?: boolean
  icon?: React.ReactNode
  className?: string
  autoComplete?: string
  id?: string
}

const InputField = forwardRef<HTMLInputElement, InputFieldProps>(
  ({ 
    label, 
    type = 'text', 
    value, 
    onChange, 
    placeholder, 
    required = false, 
    disabled = false, 
    error, 
    success = false,
    icon,
    className = '',
    autoComplete,
    id
  }, ref) => {
    const [showPassword, setShowPassword] = useState(false)
    const [isFocused, setIsFocused] = useState(false)

    const inputType = type === 'password' && showPassword ? 'text' : type

    const baseInputClasses = `
      w-full px-4 py-3.5 bg-white border border-gray-200 rounded-lg
      text-gray-900 placeholder-gray-400 text-sm font-medium
      transition-all duration-200 ease-out
      focus:outline-none focus:ring-0 focus:border-brand-green-500 focus:bg-white
      disabled:bg-gray-50 disabled:text-gray-400 disabled:cursor-not-allowed
      ${error ? 'border-red-300 focus:border-red-500' : ''}
      ${success ? 'border-green-300 focus:border-green-500' : ''}
      ${!error && !success && isFocused ? 'border-brand-green-500 shadow-sm' : ''}
      ${icon ? 'pl-11' : ''}
      ${type === 'password' ? 'pr-11' : ''}
    `

    const labelClasses = `
      block text-sm font-semibold text-gray-700 mb-2
      ${error ? 'text-red-600' : success ? 'text-green-600' : ''}
    `

    return (
      <div className={`space-y-2 ${className}`}>
        <label htmlFor={id} className={labelClasses}>
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
        
        <div className="relative">
          {/* Icon */}
          {icon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none">
              {icon}
            </div>
          )}
          
          {/* Input */}
          <input
            ref={ref}
            id={id}
            type={inputType}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            placeholder={placeholder}
            disabled={disabled}
            required={required}
            autoComplete={autoComplete}
            className={baseInputClasses}
          />
          
          {/* Password toggle */}
          {type === 'password' && (
            <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors duration-200"
              onClick={() => setShowPassword(!showPassword)}
              tabIndex={-1}
            >
              {showPassword ? (
                <EyeSlashIcon className="h-4 w-4" />
              ) : (
                <EyeIcon className="h-4 w-4" />
              )}
            </button>
          )}
        </div>
        
        {/* Error message */}
        {error && (
          <div className="flex items-center space-x-1 text-red-600 text-xs font-medium">
            <div className="w-1 h-1 bg-red-500 rounded-full"></div>
            <span>{error}</span>
          </div>
        )}
        
        {/* Success message */}
        {success && !error && (
          <div className="flex items-center space-x-1 text-green-600 text-xs font-medium">
            <div className="w-1 h-1 bg-green-500 rounded-full"></div>
            <span>Valid</span>
          </div>
        )}
      </div>
    )
  }
)

InputField.displayName = 'InputField'

export default InputField 