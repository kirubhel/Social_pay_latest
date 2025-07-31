import React, { forwardRef, useState, useRef } from 'react'
import { PhotoIcon, XMarkIcon, ArrowUpTrayIcon } from '@heroicons/react/24/outline'

interface FileUploadProps {
  label: string
  value: File | null
  onChange: (file: File | null) => void
  accept?: string
  maxSize?: number // in MB
  required?: boolean
  disabled?: boolean
  error?: string
  className?: string
  id?: string
}

const FileUpload = forwardRef<HTMLInputElement, FileUploadProps>(
  ({ 
    label, 
    value, 
    onChange, 
    accept = '*/*',
    maxSize = 10, // 10MB default
    required = false, 
    disabled = false, 
    error,
    className = '',
    id
  }, ref) => {
    const [isDragOver, setIsDragOver] = useState(false)
    const [dragError, setDragError] = useState('')
    const [isHovered, setIsHovered] = useState(false)
    const fileInputRef = useRef<HTMLInputElement>(null)

    const handleFileSelect = (file: File) => {
      setDragError('')
      
      // Check file size
      if (file.size > maxSize * 1024 * 1024) {
        setDragError(`File size must be less than ${maxSize}MB`)
        return
      }

      // Check file type if accept is specified
      if (accept !== '*/*') {
        const acceptedTypes = accept.split(',').map(type => type.trim())
        const fileType = file.type
        const fileExtension = '.' + file.name.split('.').pop()?.toLowerCase()
        
        const isAccepted = acceptedTypes.some(type => {
          if (type.startsWith('.')) {
            return fileExtension === type
          }
          return fileType === type || fileType.startsWith(type.replace('*', ''))
        })
        
        if (!isAccepted) {
          setDragError(`File type not supported. Accepted: ${accept}`)
          return
        }
      }

      onChange(file)
    }

    const handleDrop = (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragOver(false)
      
      const files = Array.from(e.dataTransfer.files)
      if (files.length > 0) {
        handleFileSelect(files[0])
      }
    }

    const handleDragOver = (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragOver(true)
    }

    const handleDragLeave = (e: React.DragEvent) => {
      e.preventDefault()
      setIsDragOver(false)
    }

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0] || null
      if (file) {
        handleFileSelect(file)
      }
    }

    const handleRemoveFile = () => {
      onChange(null)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }

    const formatFileSize = (bytes: number) => {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }

    const labelClasses = `
      block text-sm font-semibold mb-3 transition-all duration-300
      ${error || dragError ? 'text-red-600' : 'text-gray-700'}
    `

    return (
      <div className={`group relative ${className}`}>
        <label htmlFor={id} className={labelClasses}>
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
        
        <div className="space-y-3">
          {/* File Input */}
          <input
            ref={(node) => {
              // Handle both refs
              if (typeof ref === 'function') {
                ref(node)
              } else if (ref) {
                ref.current = node
              }
              fileInputRef.current = node
            }}
            id={id}
            type="file"
            accept={accept}
            onChange={handleInputChange}
            disabled={disabled}
            required={required}
            className="hidden"
          />
          
          {/* Upload Area */}
          <div
            className={`
              relative border-2 border-dashed rounded-xl p-8 text-center transition-all duration-300 cursor-pointer
              ${isDragOver 
                ? 'border-brand-green-500 bg-brand-green-50/80 shadow-lg shadow-brand-green-500/20' 
                : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50/80'
              }
              ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
              ${error || dragError ? 'border-red-300 bg-red-50/80' : ''}
              ${isHovered && !isDragOver ? 'shadow-md shadow-gray-200/50' : ''}
            `}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            onClick={() => !disabled && fileInputRef.current?.click()}
          >
            {/* Background gradient effect */}
            <div className={`
              absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-brand-green-500/5 to-transparent
              opacity-0 transition-opacity duration-500
              ${isDragOver ? 'opacity-100' : ''}
            `} />
            
            {value ? (
              <div className="relative z-10 space-y-4">
                <div className="flex items-center justify-center">
                  <div className="w-20 h-20 bg-gradient-to-br from-brand-green-100 to-brand-green-200 rounded-2xl flex items-center justify-center shadow-lg">
                    <PhotoIcon className="h-10 w-10 text-brand-green-600" />
                  </div>
                </div>
                <div className="space-y-2">
                  <p className="text-sm font-semibold text-gray-900">{value.name}</p>
                  <p className="text-xs text-gray-500">{formatFileSize(value.size)}</p>
                </div>
                {!disabled && (
                  <button
                    type="button"
                    onClick={(e) => {
                      e.stopPropagation()
                      handleRemoveFile()
                    }}
                    className="inline-flex items-center gap-2 text-sm text-red-600 hover:text-red-700 transition-all duration-200 hover:bg-red-50 px-3 py-1.5 rounded-lg"
                  >
                    <XMarkIcon className="h-4 w-4" />
                    Remove File
                  </button>
                )}
              </div>
            ) : (
              <div className="relative z-10 space-y-4">
                <div className="flex items-center justify-center">
                  <div className={`
                    w-20 h-20 rounded-2xl flex items-center justify-center transition-all duration-300
                    ${isDragOver 
                      ? 'bg-gradient-to-br from-brand-green-100 to-brand-green-200 shadow-lg' 
                      : 'bg-gray-100 hover:bg-gray-200'
                    }
                  `}>
                    <ArrowUpTrayIcon className={`
                      h-10 w-10 transition-all duration-300
                      ${isDragOver ? 'text-brand-green-600 scale-110' : 'text-gray-400'}
                    `} />
                  </div>
                </div>
                <div className="space-y-2">
                  <p className="text-sm font-semibold text-gray-900">
                    {disabled ? 'File upload disabled' : isDragOver ? 'Drop your file here' : 'Click to upload or drag and drop'}
                  </p>
                  <p className="text-xs text-gray-500">
                    {accept !== '*/*' ? `Accepted formats: ${accept}` : 'Any file type'} • Max {maxSize}MB
                  </p>
                </div>
              </div>
            )}
          </div>
          
          {/* Error Messages */}
          {(error || dragError) && (
            <div className="flex items-center space-x-2 text-red-600 text-sm animate-in slide-in-from-top-2 duration-300">
              <div className="w-1.5 h-1.5 bg-red-500 rounded-full animate-pulse"></div>
              <span className="font-medium">{error || dragError}</span>
            </div>
          )}
        </div>
      </div>
    )
  }
)

FileUpload.displayName = 'FileUpload'

export default FileUpload 