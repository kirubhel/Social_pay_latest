import React, { forwardRef, useState } from 'react'

interface TextareaFieldProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  required?: boolean
  disabled?: boolean
  error?: string
  success?: boolean
  rows?: number
  maxLength?: number
  className?: string
  id?: string
}

const TextareaField = forwardRef<HTMLTextAreaElement, TextareaFieldProps>(
  ({ 
    label, 
    value, 
    onChange, 
    placeholder, 
    required = false, 
    disabled = false, 
    error, 
    success = false,
    rows = 4,
    maxLength,
    className = '',
    id
  }, ref) => {
    const [isFocused, setIsFocused] = useState(false)
    const [isHovered, setIsHovered] = useState(false)

    const baseTextareaClasses = `
      w-full px-4 py-4 border-2 rounded-xl transition-all duration-300 ease-out
      bg-white/80 backdrop-blur-sm text-gray-900 placeholder-gray-400 resize-none
      focus:outline-none focus:ring-4 focus:ring-brand-green-500/20 focus:border-brand-green-500
      disabled:bg-gray-50/80 disabled:text-gray-500 disabled:cursor-not-allowed
      ${error ? 'border-red-300 focus:ring-red-500/20 focus:border-red-500' : ''}
      ${success ? 'border-green-300 focus:ring-green-500/20 focus:border-green-500' : ''}
      ${!error && !success ? 'border-gray-200 hover:border-gray-300' : ''}
      ${isFocused ? 'shadow-lg shadow-brand-green-500/10' : ''}
      ${isHovered && !isFocused ? 'shadow-md shadow-gray-200/50' : ''}
    `

    const labelClasses = `
      block text-sm font-semibold mb-3 transition-all duration-300
      ${error ? 'text-red-600' : success ? 'text-green-600' : 'text-gray-700'}
      ${isFocused ? 'text-brand-green-600 scale-105' : ''}
    `

    const containerClasses = `
      group relative ${className}
    `

    return (
      <div className={containerClasses}>
        <div className="flex items-center justify-between mb-3">
          <label htmlFor={id} className={labelClasses}>
            {label}
            {required && <span className="text-red-500 ml-1">*</span>}
          </label>
          {maxLength && (
            <span className={`
              text-xs font-medium transition-colors duration-200
              ${value.length > maxLength * 0.9 ? 'text-red-500' : 'text-gray-500'}
            `}>
              {value.length}/{maxLength}
            </span>
          )}
        </div>
        
        <div 
          className="relative"
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
        >
          {/* Background gradient effect */}
          <div className={`
            absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-brand-green-500/5 to-transparent
            opacity-0 transition-opacity duration-500
            ${isFocused ? 'opacity-100' : ''}
          `} />
          
          {/* Textarea field */}
          <textarea
            ref={ref}
            id={id}
            rows={rows}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            placeholder={placeholder}
            disabled={disabled}
            required={required}
            maxLength={maxLength}
            className={baseTextareaClasses}
          />
          
          {/* Focus ring effect */}
          <div className={`
            absolute inset-0 rounded-xl transition-all duration-300 pointer-events-none
            ${isFocused ? 'ring-2 ring-brand-green-500/30 ring-offset-2' : ''}
          `} />
          
          {/* Success/Error indicator */}
          {(success || error) && (
            <div className={`
              absolute right-4 top-4 z-10
              w-2 h-2 rounded-full transition-all duration-300
              ${error ? 'bg-red-500 animate-pulse' : 'bg-green-500'}
            `} />
          )}
        </div>
        
        {/* Error message */}
        {error && (
          <div className="flex items-center space-x-2 mt-2 text-red-600 text-sm animate-in slide-in-from-top-2 duration-300">
            <div className="w-1.5 h-1.5 bg-red-500 rounded-full animate-pulse"></div>
            <span className="font-medium">{error}</span>
          </div>
        )}
        
        {/* Success message */}
        {success && !error && (
          <div className="flex items-center space-x-2 mt-2 text-green-600 text-sm animate-in slide-in-from-top-2 duration-300">
            <div className="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
            <span className="font-medium">Valid</span>
          </div>
        )}
        
        {/* Character counter with progress bar */}
        {maxLength && (
          <div className="mt-2 space-y-1">
            <div className="flex justify-between items-center text-xs text-gray-500">
              <span>Character count</span>
              <span className={`
                font-medium transition-colors duration-200
                ${value.length > maxLength * 0.9 ? 'text-red-500' : ''}
              `}>
                {value.length}/{maxLength}
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-1">
              <div 
                className={`
                  h-1 rounded-full transition-all duration-300
                  ${value.length > maxLength * 0.9 ? 'bg-red-500' : 'bg-brand-green-500'}
                `}
                style={{ width: `${Math.min((value.length / maxLength) * 100, 100)}%` }}
              />
            </div>
          </div>
        )}
      </div>
    )
  }
)

TextareaField.displayName = 'TextareaField'

export default TextareaField 