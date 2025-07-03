'use client'

import { useState, useRef, Fragment, useEffect } from 'react'
import { Tab, Dialog, Transition } from '@headlessui/react'
import { cn } from '@/lib/utils'
import { authAPI } from '@/lib/api'
import {
  UserCircleIcon,
  ShieldCheckIcon,
  KeyIcon,
  LinkIcon,
  UserGroupIcon,
  DocumentTextIcon,
  CogIcon,
  GlobeAltIcon,
  EyeIcon,
  EyeSlashIcon,
  CheckCircleIcon,
  ClipboardIcon,
  CheckIcon,
  PlusIcon,
  MagnifyingGlassIcon
} from '@heroicons/react/24/outline'
import toast from 'react-hot-toast'

const tabs = [
  { name: 'General', icon: CogIcon },
  { name: 'Security', icon: ShieldCheckIcon },
  { name: 'Api', icon: KeyIcon },
  { name: 'Webhooks', icon: LinkIcon },
  { name: 'Teams', icon: UserGroupIcon },
  { name: 'Compliance', icon: DocumentTextIcon },
  { name: 'Account Settings', icon: UserCircleIcon },
  { name: 'Whitelisted IPS', icon: GlobeAltIcon },
]

export default function SettingsPage() {
  const [selectedIndex, setSelectedIndex] = useState(0)

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Settings</h1>
        <p className="mt-2 text-gray-600">
          Manage your account settings and preferences
        </p>
      </div>

      <Tab.Group selectedIndex={selectedIndex} onChange={setSelectedIndex}>
        <Tab.List className="flex space-x-1 rounded-xl bg-gray-100 p-1 mb-8 overflow-x-auto">
          {tabs.map((tab, index) => (
            <Tab
              key={tab.name}
              className={({ selected }) =>
                cn(
                  'w-full rounded-lg py-2.5 px-4 text-sm font-medium leading-5 transition-all duration-200',
                  'ring-white ring-opacity-60 ring-offset-2 ring-offset-brand-green-400 focus:outline-none focus:ring-2',
                  'flex items-center justify-center gap-2 whitespace-nowrap',
                  selected
                    ? 'bg-white text-brand-green-700 shadow-md border-b-2 border-brand-green-500'
                    : 'text-gray-600 hover:bg-white/[0.12] hover:text-brand-green-600'
                )
              }
            >
              <tab.icon className="h-4 w-4" />
              <span className="hidden sm:inline">{tab.name}</span>
            </Tab>
          ))}
        </Tab.List>

        <Tab.Panels>
          {/* General Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <GeneralSettings />
          </Tab.Panel>

          {/* Security Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <Enable2FAModalWrapper />
          </Tab.Panel>

          {/* API Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <ApiKeysPanel />
          </Tab.Panel>

          {/* Webhooks Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <WebhooksPanel />
          </Tab.Panel>

          {/* Teams Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <TeamsPanel />
          </Tab.Panel>

          {/* Compliance Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <ComingSoon tabName="Compliance" />
          </Tab.Panel>

          {/* Account Settings Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <AccountSettingsPanel />
          </Tab.Panel>

          {/* Whitelisted IPs Tab */}
          <Tab.Panel className="rounded-xl bg-white p-6 shadow-lg border border-gray-100">
            <ComingSoon tabName="Whitelisted IPs" />
          </Tab.Panel>
        </Tab.Panels>
      </Tab.Group>
    </div>
  )
}

function ComingSoon({ tabName }: { tabName: string }) {
  return (
    <div className="text-center py-12">
      <div className="mx-auto w-24 h-24 bg-gradient-to-r from-brand-green-100 to-brand-gold-100 rounded-full flex items-center justify-center mb-6">
        <CogIcon className="h-12 w-12 text-brand-green-600" />
      </div>
      <h3 className="text-xl font-semibold text-gray-900 mb-2">
        {tabName} Settings
      </h3>
      <p className="text-gray-600 mb-4">
        This section is coming soon. We're working hard to bring you the best experience.
      </p>
      <div className="inline-flex items-center text-sm text-brand-green-600">
        <div className="w-2 h-2 bg-brand-green-500 rounded-full animate-pulse mr-2"></div>
        In Development
      </div>
    </div>
  )
}

function GeneralSettings() {
  const [formData, setFormData] = useState({
    userName: 'Henok Kebede',
    emailAddress: 'example@gmail.com',
    phoneNumber: '+251 912345678',
    businessName: '',
    businessEmail: '',
    businessPhone: '',
    businessAddress: '',
    website: '',
    businessLogo: null as File | null
  })

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0] || null
    setFormData(prev => ({ ...prev, businessLogo: file }))
  }

  const handleSaveChanges = () => {
    console.log('Saving changes:', formData)
    // TODO: Implement save functionality
  }

  return (
    <div className="space-y-8">
      {/* User Information Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-6">User Information</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label htmlFor="userName" className="block text-sm font-medium text-gray-700 mb-2">
              User Name
            </label>
            <input
              type="text"
              id="userName"
              value={formData.userName}
              onChange={(e) => handleInputChange('userName', e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
          <div>
            <label htmlFor="emailAddress" className="block text-sm font-medium text-gray-700 mb-2">
              Email Address
            </label>
            <input
              type="email"
              id="emailAddress"
              value={formData.emailAddress}
              onChange={(e) => handleInputChange('emailAddress', e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
          <div className="md:col-span-2">
            <label htmlFor="phoneNumber" className="block text-sm font-medium text-gray-700 mb-2">
              Phone Number
            </label>
            <input
              type="tel"
              id="phoneNumber"
              value={formData.phoneNumber}
              onChange={(e) => handleInputChange('phoneNumber', e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
        </div>
      </div>

      {/* Business Information Section */}
      <div>
        <h2 className="text-xl font-semibold text-gray-900 mb-6">Business Information</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label htmlFor="businessName" className="block text-sm font-medium text-gray-700 mb-2">
              Business Name
            </label>
            <input
              type="text"
              id="businessName"
              value={formData.businessName}
              onChange={(e) => handleInputChange('businessName', e.target.value)}
              placeholder="Your Business Name"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 placeholder-gray-400"
            />
          </div>
          <div>
            <label htmlFor="businessEmail" className="block text-sm font-medium text-gray-700 mb-2">
              Business Email
            </label>
            <input
              type="email"
              id="businessEmail"
              value={formData.businessEmail}
              onChange={(e) => handleInputChange('businessEmail', e.target.value)}
              placeholder="Your Business Email Address"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 placeholder-gray-400"
            />
          </div>
          <div>
            <label htmlFor="businessPhone" className="block text-sm font-medium text-gray-700 mb-2">
              Business Phone Number
            </label>
            <input
              type="tel"
              id="businessPhone"
              value={formData.businessPhone}
              onChange={(e) => handleInputChange('businessPhone', e.target.value)}
              placeholder="Your Phone Number"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 placeholder-gray-400"
            />
          </div>
          <div>
            <label htmlFor="businessAddress" className="block text-sm font-medium text-gray-700 mb-2">
              Business Address
            </label>
            <input
              type="text"
              id="businessAddress"
              value={formData.businessAddress}
              onChange={(e) => handleInputChange('businessAddress', e.target.value)}
              placeholder="Your Business Address"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 placeholder-gray-400"
            />
          </div>
          <div>
            <label htmlFor="website" className="block text-sm font-medium text-gray-700 mb-2">
              Website
            </label>
            <input
              type="url"
              id="website"
              value={formData.website}
              onChange={(e) => handleInputChange('website', e.target.value)}
              placeholder="Your Website Name"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 placeholder-gray-400"
            />
          </div>
          <div>
            <label htmlFor="businessLogo" className="block text-sm font-medium text-gray-700 mb-2">
              Business Logo
            </label>
            <div className="flex items-center gap-4">
              <label htmlFor="businessLogo" className="cursor-pointer">
                <div className="px-6 py-3 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors duration-200 text-center">
                  Choose File
                </div>
                <input
                  type="file"
                  id="businessLogo"
                  accept="image/*"
                  onChange={handleFileChange}
                  className="hidden"
                />
              </label>
              <span className="text-sm text-gray-500">
                {formData.businessLogo ? formData.businessLogo.name : 'No File Chosen'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Save Button */}
      <div className="flex justify-start pt-6 border-t border-gray-200">
        <button
          onClick={handleSaveChanges}
          className="px-8 py-3 bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-brand-green-500 focus:ring-offset-2"
        >
          Save Changes
        </button>
      </div>
    </div>
  )
}

function Enable2FAModalWrapper() {
  const [modalOpen, setModalOpen] = useState(false)
  const [is2FAEnabled, setIs2FAEnabled] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [currentPassword, setCurrentPassword] = useState('')
  const [showPasswordModal, setShowPasswordModal] = useState(false)

  // Fetch 2FA status on component mount
  useEffect(() => {
    const fetch2FAStatus = async () => {
      try {
        const response = await authAPI.get2FAStatus()
        if (response.success) {
          setIs2FAEnabled(response.data.enabled || false)
        } else {
          toast.error('Failed to fetch 2FA status')
        }
      } catch (err) {
        toast.error('Failed to connect to server')
      }
    }
    fetch2FAStatus()
  }, [])

  const handleEnable2FA = async () => {
    setLoading(true)
    setError('')
    try {
      const loadingToast = toast.loading('Enabling 2FA...')
      const response = await authAPI.enable2FA()
      toast.dismiss(loadingToast)
      
      if (response.success) {
        setModalOpen(true)
        toast.success('2FA setup initiated! Please verify with the code sent to your phone.')
      } else {
        toast.error(response.error?.message || 'Failed to enable 2FA')
        setError(response.error?.message || 'Failed to enable 2FA')
      }
    } catch (err: any) {
      toast.error('An unexpected error occurred')
      setError(err.message || 'An unexpected error occurred')
    } finally {
      setLoading(false)
    }
  }

  const handleDisable2FA = async () => {
    if (!currentPassword.trim()) {
      toast.error('Please enter your current password')
      return
    }

    setLoading(true)
    setError('')
    try {
      const loadingToast = toast.loading('Disabling 2FA...')
      const response = await authAPI.disable2FA(currentPassword)
      toast.dismiss(loadingToast)
      
      if (response.success) {
        setIs2FAEnabled(false)
        setShowPasswordModal(false)
        setCurrentPassword('')
        toast.success('2FA has been disabled successfully!')
      } else {
        toast.error(response.error?.message || 'Failed to disable 2FA')
        setError(response.error?.message || 'Failed to disable 2FA')
      }
    } catch (err: any) {
      toast.error('An unexpected error occurred')
      setError(err.message || 'An unexpected error occurred')
    } finally {
      setLoading(false)
    }
  }

  const handleToggle2FA = async () => {
    if (!is2FAEnabled) {
      await handleEnable2FA()
    } else {
      setShowPasswordModal(true)
    }
  }

  return (
    <>
      <div className="max-w-2xl mx-auto">
        {/* Password Section */}
        <PasswordUpdateForm />
        <hr className="my-8 border-gray-200" />
        
        {/* 2-Step Verification Section */}
        <div>
          <h2 className="text-lg font-bold text-gray-900 mb-6">2 - Step Verification</h2>
          <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-6 border border-blue-100">
            <div className="flex items-start justify-between">
              <div className="flex items-start space-x-4">
                <div className={cn(
                  "p-3 rounded-full",
                  is2FAEnabled ? "bg-green-100" : "bg-gray-100"
                )}>
                  <ShieldCheckIcon className={cn(
                    "h-6 w-6",
                    is2FAEnabled ? "text-green-600" : "text-gray-400"
                  )} />
                </div>
                <div className="flex-1">
                  <h3 className="text-base font-semibold text-gray-900 mb-1">
                    Two-Step Verification
                  </h3>
                  <p className="text-sm text-gray-600 mb-2">
                    {is2FAEnabled
                      ? 'Your account is protected with two-step verification. You will need to enter a verification code when signing in.'
                      : 'Add an extra layer of security to your account by enabling two-step verification.'}
                  </p>
                  {is2FAEnabled && (
                    <div className="flex items-center space-x-2 text-sm">
                      <CheckCircleIcon className="h-4 w-4 text-green-500" />
                      <span className="text-green-700 font-medium">Active</span>
                    </div>
                  )}
                </div>
              </div>
              <button
                className={cn(
                  'px-6 py-2.5 font-semibold rounded-lg shadow-md transition-all duration-200 transform hover:scale-105',
                  is2FAEnabled
                    ? 'bg-red-500 hover:bg-red-600 text-white hover:shadow-lg'
                    : 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 hover:from-brand-green-600 hover:to-brand-green-700 text-white hover:shadow-lg',
                  loading && 'opacity-50 cursor-not-allowed transform-none'
                )}
                onClick={handleToggle2FA}
                disabled={loading}
              >
                {loading ? (
                  <div className="flex items-center">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                    <span>Processing...</span>
                  </div>
                ) : is2FAEnabled ? (
                  'Disable'
                ) : (
                  'Enable'
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Password Confirmation Modal for Disabling 2FA */}
      <Transition.Root show={showPasswordModal} as={Fragment}>
        <Dialog as="div" className="relative z-50" onClose={() => setShowPasswordModal(false)}>
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100"
            leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-black bg-opacity-30 transition-opacity" />
          </Transition.Child>
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100"
              leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95"
            >
              <Dialog.Panel className="mx-auto w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl">
                <div className="flex items-center justify-center mb-6">
                  <div className="p-3 bg-red-100 rounded-full">
                    <ShieldCheckIcon className="h-8 w-8 text-red-600" />
                  </div>
                </div>
                <Dialog.Title className="text-xl font-semibold text-gray-900 text-center mb-2">
                  Disable Two-Step Verification
                </Dialog.Title>
                <p className="text-sm text-gray-600 text-center mb-6">
                  Enter your current password to disable two-step verification for your account.
                </p>
                
                {error && (
                  <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
                    <p className="text-sm text-red-600 text-center">{error}</p>
                  </div>
                )}

                <div className="mb-6">
                  <label htmlFor="currentPassword" className="block text-sm font-medium text-gray-700 mb-2">
                    Current Password
                  </label>
                  <input
                    type="password"
                    id="currentPassword"
                    value={currentPassword}
                    onChange={(e) => setCurrentPassword(e.target.value)}
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500 transition-colors duration-200"
                    placeholder="Enter your current password"
                  />
                </div>

                <div className="flex space-x-3">
                  <button
                    onClick={() => {
                      setShowPasswordModal(false)
                      setCurrentPassword('')
                      setError('')
                    }}
                    className="flex-1 px-4 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors duration-200"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleDisable2FA}
                    disabled={loading || !currentPassword.trim()}
                    className={cn(
                      "flex-1 px-4 py-3 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors duration-200",
                      (loading || !currentPassword.trim()) && "opacity-50 cursor-not-allowed"
                    )}
                  >
                    {loading ? 'Disabling...' : 'Disable 2FA'}
                  </button>
                </div>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>

      {/* Verification Code Modal */}
      <VerificationCodeModal
        open={modalOpen}
        onClose={() => setModalOpen(false)}
        onVerify={async (code) => {
          try {
            const response = await authAPI.verify2FASetup(code)
            if (response.success) {
              setIs2FAEnabled(true)
              setModalOpen(false)
              toast.success('2FA has been enabled successfully! Your account is now more secure.')
              return undefined
            } else {
              return response.error?.message || 'Verification failed'
            }
          } catch (err: any) {
            return err.message || 'Verification failed'
          }
        }}
      />
    </>
  )
}

function VerificationCodeModal({
  open,
  onClose,
  onVerify
}: {
  open: boolean
  onClose: () => void
  onVerify: (code: string) => Promise<string | undefined>
}) {
  const [code, setCode] = useState(['', '', '', '', '', ''])
  const [resendLoading, setResendLoading] = useState(false)
  const [verifying, setVerifying] = useState(false)
  const [error, setError] = useState('')
  const inputRefs = useRef<(HTMLInputElement | null)[]>([])

  const handleChange = (idx: number, value: string) => {
    if (!/^[0-9]?$/.test(value)) return
    const newCode = [...code]
    newCode[idx] = value
    setCode(newCode)
    if (value && idx < 5) {
      inputRefs.current[idx + 1]?.focus()
    }
  }

  const handleKeyDown = (idx: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace' && !code[idx] && idx > 0) {
      inputRefs.current[idx - 1]?.focus()
    }
  }

  const handleResend = async () => {
    setResendLoading(true)
    try {
      const response = await authAPI.resend2FACode()
      if (response.success) {
        toast.success('New verification code sent to your phone')
      } else {
        setError(response.error?.message || 'Failed to resend code')
      }
    } catch (err: any) {
      setError(err.message || 'Failed to resend code')
    } finally {
      setResendLoading(false)
    }
  }

  const handleVerify = async () => {
    const fullCode = code.join('')
    if (fullCode.length !== 6) {
      setError('Please enter all 6 digits')
      return
    }
    setVerifying(true)
    setError('')
    try {
      const errorMessage = await onVerify(fullCode)
      if (errorMessage) {
        setError(errorMessage)
      }
    } catch (err: any) {
      setError(err.message || 'Verification failed')
    } finally {
      setVerifying(false)
    }
  }

  return (
    <Transition.Root show={open} as={Fragment}>
      <Dialog as="div" className="relative z-50" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300" enterFrom="opacity-0" enterTo="opacity-100"
          leave="ease-in duration-200" leaveFrom="opacity-100" leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-30 transition-opacity" />
        </Transition.Child>
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <Transition.Child
            as={Fragment}
            enter="ease-out duration-300" enterFrom="opacity-0 scale-95" enterTo="opacity-100 scale-100"
            leave="ease-in duration-200" leaveFrom="opacity-100 scale-100" leaveTo="opacity-0 scale-95"
          >
            <Dialog.Panel className="mx-auto w-full max-w-md rounded-2xl bg-white p-8 shadow-2xl">
              <div className="text-center mb-6">
                <div className="mx-auto w-16 h-16 bg-gradient-to-r from-green-100 to-blue-100 rounded-full flex items-center justify-center mb-4">
                  <ShieldCheckIcon className="h-8 w-8 text-green-600" />
                </div>
                <Dialog.Title className="text-xl font-semibold text-gray-900 mb-2">
                  Enter Verification Code
                </Dialog.Title>
                <p className="text-gray-600 text-sm">
                  We've sent a 6-digit verification code to your phone. Please enter it below to complete 2FA setup.
                </p>
              </div>

              <div className="flex justify-center gap-3 mb-6">
                {code.map((digit, idx) => (
                  <input
                    key={idx}
                    ref={el => { inputRefs.current[idx] = el; }}
                    type="text"
                    inputMode="numeric"
                    maxLength={1}
                    value={digit}
                    onChange={e => handleChange(idx, e.target.value)}
                    onKeyDown={e => handleKeyDown(idx, e)}
                    className="w-12 h-12 text-center text-xl font-semibold border-2 border-gray-200 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-all duration-200 hover:border-gray-300"
                  />
                ))}
              </div>

              {error && (
                <div className="mb-6 p-3 bg-red-50 border border-red-200 rounded-lg">
                  <p className="text-red-600 text-sm text-center font-medium">{error}</p>
                </div>
              )}

              <div className="flex flex-col space-y-3">
                <button
                  onClick={handleVerify}
                  disabled={verifying || code.some(d => !d)}
                  className={cn(
                    "w-full px-6 py-3 bg-gradient-to-r from-brand-green-500 to-brand-green-600 hover:from-brand-green-600 hover:to-brand-green-700 text-white font-semibold rounded-lg shadow-md transition-all duration-200 transform hover:scale-105",
                    (verifying || code.some(d => !d)) && "opacity-50 cursor-not-allowed transform-none"
                  )}
                >
                  {verifying ? (
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                      <span>Verifying...</span>
                    </div>
                  ) : (
                    'Verify Code'
                  )}
                </button>

                <div className="text-center">
                  <span className="text-sm text-gray-600">Didn't receive the code? </span>
                  <button
                    onClick={handleResend}
                    disabled={resendLoading}
                    className="text-sm font-semibold text-brand-green-600 hover:text-brand-green-500 transition-colors hover:underline disabled:opacity-50"
                  >
                    {resendLoading ? 'Resending...' : 'Resend Code'}
                  </button>
                </div>
              </div>
            </Dialog.Panel>
          </Transition.Child>
        </div>
      </Dialog>
    </Transition.Root>
  )
}

function PasswordUpdateForm() {
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showCurrent, setShowCurrent] = useState(false)
  const [showNew, setShowNew] = useState(false)
  const [showConfirm, setShowConfirm] = useState(false)
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState('')
  const [error, setError] = useState('')

  const validate = () => {
    if (!currentPassword || !newPassword || !confirmPassword) {
      setError('All fields are required.')
      return false
    }
    if (newPassword !== confirmPassword) {
      setError('New password and confirmation do not match.')
      return false
    }
    if (newPassword === currentPassword) {
      setError('New password must be different from current password.')
      return false
    }
    setError('')
    return true
  }

  const handleUpdate = async () => {
    setSuccess('')
    if (!validate()) return
    setLoading(true)
    setError('')
    // Mock API call
    await new Promise((resolve) => setTimeout(resolve, 1500))
    // Simulate success
    setLoading(false)
    setSuccess('Password updated successfully!')
    setCurrentPassword('')
    setNewPassword('')
    setConfirmPassword('')
  }

  return (
    <div className="mb-10">
      <h2 className="text-lg font-bold text-gray-900 mb-6">Password</h2>
      <div className="mb-6 space-y-6">
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">Current Password</label>
          <div className="relative">
            <input
              type={showCurrent ? 'text' : 'password'}
              value={currentPassword}
              onChange={e => setCurrentPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 pr-12"
              autoComplete="current-password"
            />
            <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700"
              onClick={() => setShowCurrent(v => !v)}
              tabIndex={-1}
            >
              {showCurrent ? <EyeSlashIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">New Password</label>
          <div className="relative">
            <input
              type={showNew ? 'text' : 'password'}
              value={newPassword}
              onChange={e => setNewPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 pr-12"
              autoComplete="new-password"
            />
            <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700"
              onClick={() => setShowNew(v => !v)}
              tabIndex={-1}
            >
              {showNew ? <EyeSlashIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">Confirm Password</label>
          <div className="relative">
            <input
              type={showConfirm ? 'text' : 'password'}
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200 pr-12"
              autoComplete="new-password"
            />
            <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-700"
              onClick={() => setShowConfirm(v => !v)}
              tabIndex={-1}
            >
              {showConfirm ? <EyeSlashIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
        {error && <div className="text-red-600 text-sm mt-2">{error}</div>}
        {success && <div className="text-green-600 text-sm mt-2">{success}</div>}
      </div>
      <button
        onClick={handleUpdate}
        disabled={loading}
        className="mt-2 px-8 py-2 bg-brand-green-600 hover:bg-brand-green-700 text-white font-semibold rounded-lg shadow transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
      >
        {loading ? 'Updating...' : 'Update Password'}
      </button>
    </div>
  )
}

function ApiKeysPanel() {
  const [copied, setCopied] = useState<{ [key: string]: boolean }>({})
  const [keys] = useState({
    public: 'CHAPUBK_TEST-auX7X3QKNb3cXPkswa7HcyI2tQFJI1',
    secret: '*************************************',
    encryption: 'Secret Key',
  })

  const handleCopy = (key: string, value: string) => {
    navigator.clipboard.writeText(value)
    setCopied(prev => ({ ...prev, [key]: true }))
    setTimeout(() => setCopied(prev => ({ ...prev, [key]: false })), 1200)
  }

  return (
    <div className="relative bg-white rounded-2xl p-8 shadow-md border border-gray-100">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h2 className="text-xl font-bold text-gray-900 mb-1">Api Keys</h2>
          <p className="text-gray-500 text-sm">Your API keys are private. Keep them safe and secure to protect your account.</p>
        </div>
        <button
          className="px-6 py-2 bg-orange-400 hover:bg-orange-500 text-white font-semibold rounded-lg shadow transition-all duration-200"
        >
          Generate key's
        </button>
      </div>
      <div className="space-y-8">
        {/* Public Key */}
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">Public Key</label>
          <div className="relative flex items-center">
            <input
              type="text"
              value={keys.public}
              readOnly
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 font-mono text-sm pr-12"
            />
            <button
              type="button"
              onClick={() => handleCopy('public', keys.public)}
              className="absolute right-3 text-gray-400 hover:text-brand-green-600"
            >
              {copied.public ? <CheckIcon className="h-5 w-5" /> : <ClipboardIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
        {/* Secret Key */}
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">Secret Key</label>
          <div className="relative flex items-center">
            <input
              type="password"
              value={keys.secret}
              readOnly
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 font-mono text-sm pr-12"
            />
            <button
              type="button"
              onClick={() => handleCopy('secret', keys.secret)}
              className="absolute right-3 text-gray-400 hover:text-brand-green-600"
            >
              {copied.secret ? <CheckIcon className="h-5 w-5" /> : <ClipboardIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
        {/* Encryption Key */}
        <div>
          <label className="block text-base font-medium text-gray-800 mb-2">Encryption Key</label>
          <div className="relative flex items-center">
            <input
              type="text"
              value={keys.encryption}
              readOnly
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 font-mono text-sm pr-12"
            />
            <button
              type="button"
              onClick={() => handleCopy('encryption', keys.encryption)}
              className="absolute right-3 text-gray-400 hover:text-brand-green-600"
            >
              {copied.encryption ? <CheckIcon className="h-5 w-5" /> : <ClipboardIcon className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

function WebhooksPanel() {
  const [webhookUrl, setWebhookUrl] = useState('')
  const [secretHash, setSecretHash] = useState('')
  const [receiveWebhook, setReceiveWebhook] = useState(false)
  const [receiveFailedWebhook, setReceiveFailedWebhook] = useState(false)
  const [success, setSuccess] = useState('')

  const handleUpdate = (e: React.FormEvent) => {
    e.preventDefault()
    setSuccess('')
    // Mock API call
    setTimeout(() => {
      setSuccess('Webhook settings updated!')
    }, 1000)
  }

  return (
    <form onSubmit={handleUpdate} className="relative bg-white rounded-2xl p-8 shadow-md border border-gray-100 max-w-2xl mx-auto">
      <h2 className="text-xl font-bold text-gray-900 mb-1">Webhooks</h2>
      <div className="text-gray-500 text-sm mb-8">Configure your webhook endpoint and secret for secure event notifications.</div>
      <div className="mb-8">
        <label className="block text-base font-medium text-gray-800 mb-2">Webhook URL</label>
        <input
          type="url"
          value={webhookUrl}
          onChange={e => setWebhookUrl(e.target.value)}
          placeholder="https://webhook.site"
          className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
        />
      </div>
      <div className="mb-8">
        <label className="block text-base font-medium text-gray-800 mb-2">Secret hash</label>
        <input
          type="text"
          value={secretHash}
          onChange={e => setSecretHash(e.target.value)}
          className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
        />
      </div>
      <div className="mb-8 space-y-2">
        <label className="flex items-center gap-2 text-gray-700 text-sm">
          <input
            type="checkbox"
            checked={receiveWebhook}
            onChange={e => setReceiveWebhook(e.target.checked)}
            className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500"
          />
          Receive Webhook
        </label>
        <label className="flex items-center gap-2 text-gray-700 text-sm">
          <input
            type="checkbox"
            checked={receiveFailedWebhook}
            onChange={e => setReceiveFailedWebhook(e.target.checked)}
            className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500"
          />
          Receive Webhook for failed Payments.
        </label>
      </div>
      {success && <div className="text-green-600 text-sm mb-4">{success}</div>}
      <button
        type="submit"
        className="mt-2 px-8 py-2 bg-brand-green-600 hover:bg-brand-green-700 text-white font-semibold rounded-lg shadow transition-all duration-200"
      >
        Update Settings
      </button>
    </form>
  )
}

function TeamsPanel() {
  const teamMembers = [
    {
      name: 'Henok Tesfaye',
      email: 'example@email.com',
      role: 'Owner',
      dateJoined: '3rd May, 2023',
      status: 'Active',
    },
    {
      name: 'Henok Tesfaye',
      email: 'example@email.com',
      role: 'Owner',
      dateJoined: '3rd May, 2023',
      status: 'Active',
    },
    {
      name: 'Henok Tesfaye',
      email: 'example@email.com',
      role: 'Owner',
      dateJoined: '3rd May, 2023',
      status: 'Active',
    },
    {
      name: 'Henok Tesfaye',
      email: 'example@email.com',
      role: 'Owner',
      dateJoined: '3rd May, 2023',
      status: 'Active',
    },
  ]
  const [modalOpen, setModalOpen] = useState(false)

  return (
    <div className="relative bg-white rounded-2xl p-8 shadow-md border border-gray-100">
      <div className="flex items-center justify-between mb-8">
        <h2 className="text-xl font-bold text-gray-900">Team Members</h2>
        <div className="flex items-center gap-2">
          <button className="p-2 rounded-lg hover:bg-gray-100 text-gray-400">
            <MagnifyingGlassIcon className="h-5 w-5" />
          </button>
          <button
            className="flex items-center gap-2 px-4 py-2 bg-brand-green-600 hover:bg-brand-green-700 text-white font-semibold rounded-lg shadow transition-all duration-200"
            onClick={() => setModalOpen(true)}
          >
            <PlusIcon className="h-5 w-5" />
            Add New Member
          </button>
        </div>
      </div>
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date Joined</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-100">
            {teamMembers.map((member, idx) => (
              <tr key={idx}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{member.name}</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-800">{member.email}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold bg-green-100 text-green-700">
                    {member.role}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{member.dateJoined}</td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className="inline-flex items-center gap-1 text-xs font-medium text-green-700">
                    <span className="w-2 h-2 rounded-full bg-green-500 inline-block"></span>
                    {member.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <AddMemberModal open={modalOpen} onClose={() => setModalOpen(false)} />
    </div>
  )
}

function AddMemberModal({ open, onClose }: { open: boolean, onClose: () => void }) {
  const [fullName, setFullName] = useState('')
  const [email, setEmail] = useState('')
  const [role, setRole] = useState('')
  const roles = ['Owner', 'Admin', 'Member']
  const isValid = fullName && email && role

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    // TODO: Add member logic
    onClose()
  }

  return (
    <Dialog open={open} onClose={onClose} className="relative z-50">
      <div className="fixed inset-0 bg-black/20" aria-hidden="true" />
      <div className="fixed inset-0 flex items-center justify-center p-4">
        <Dialog.Panel className="mx-auto w-full max-w-2xl rounded-2xl bg-white p-10 shadow-2xl">
          <Dialog.Title className="text-2xl font-semibold mb-8">Add New Member</Dialog.Title>
          <form onSubmit={handleSubmit} className="space-y-8">
            <div>
              <label className="block text-base font-medium text-gray-700 mb-2">Full Name</label>
              <input
                type="text"
                value={fullName}
                onChange={e => setFullName(e.target.value)}
                placeholder="Member's Full Name"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
              />
            </div>
            <div>
              <label className="block text-base font-medium text-gray-700 mb-2">Email Address</label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="Members Email Address"
                className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
              />
            </div>
            <div>
              <label className="block text-base font-medium text-gray-700 mb-2">Members Role</label>
              <select
                value={role}
                onChange={e => setRole(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
              >
                <option value="">Select a Role</option>
                {roles.map(r => (
                  <option key={r} value={r}>{r}</option>
                ))}
              </select>
            </div>
            <div>
              <button
                type="submit"
                disabled={!isValid}
                className="w-full px-8 py-3 bg-gray-300 text-white font-semibold rounded-lg shadow transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed"
                style={isValid ? { background: 'linear-gradient(to right, #22c55e, #16a34a)' } : {}}
              >
                Add Member
              </button>
            </div>
          </form>
        </Dialog.Panel>
      </div>
    </Dialog>
  )
}

function AccountSettingsPanel() {
  // Preference
  const [defaultCurrency, setDefaultCurrency] = useState('')
  const [callbackUrl, setCallbackUrl] = useState('')
  const [returnUrl, setReturnUrl] = useState('')

  // Transactions
  const [transactionFeePayer, setTransactionFeePayer] = useState('Charge me the transaction fees')
  const [transferFeePayer, setTransferFeePayer] = useState('Charge me the transaction fees')
  const [retryMinutes, setRetryMinutes] = useState('60')
  const [apiTransfers, setApiTransfers] = useState(true)

  // Payment Methods
  const [walletsOpen, setWalletsOpen] = useState(true)
  const [banksOpen, setBanksOpen] = useState(true)
  const [cardsEnabled, setCardsEnabled] = useState(true)
  const [wallets, setWallets] = useState<Record<string, boolean>>({
    'Tele birr': true,
    'CBE': true,
    'M-Pesa': true,
    'E Birr': false,
  })
  const [banks, setBanks] = useState<Record<string, boolean>>({
    'Abyssinia Bank': true,
    'Commercial Bank': true,
    'Awash Bank': true,
    'Dashen Bank': false,
  })

  // Notification Emails
  const [notifyImportant, setNotifyImportant] = useState(false)
  const [notifyReceipts, setNotifyReceipts] = useState(false)

  // Transaction Receipts
  const [receiptMe, setReceiptMe] = useState(false)
  const [receiptRecipients, setReceiptRecipients] = useState(false)
  const [receiptFinance, setReceiptFinance] = useState(false)

  // Transfer Receipts
  const [transferMe, setTransferMe] = useState(false)
  const [transferFinance, setTransferFinance] = useState(false)

  // Transfer Approval
  const [approveUrl, setApproveUrl] = useState(false)
  const [approveOtp, setApproveOtp] = useState(false)
  const [approveFinance, setApproveFinance] = useState(false)
  const [approvalUrl, setApprovalUrl] = useState('')
  const [approvalSecret, setApprovalSecret] = useState('')

  // Finance Email
  const [financeEmail, setFinanceEmail] = useState('')

  // Exported Files Email
  const [exportEmailType, setExportEmailType] = useState('default')
  const [customExportEmail, setCustomExportEmail] = useState('')

  // Save handlers
  const handleUpdate = (e: React.FormEvent) => {
    e.preventDefault()
    // TODO: Save logic
    alert('Settings updated!')
  }

  return (
    <form onSubmit={handleUpdate} className="space-y-10 relative">
      {/* Preference */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Preference</h3>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Default Currency</label>
            <input
              type="text"
              value={defaultCurrency}
              onChange={e => setDefaultCurrency(e.target.value)}
              placeholder="ETB, USD, etc."
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Callback URL</label>
            <input
              type="url"
              value={callbackUrl}
              onChange={e => setCallbackUrl(e.target.value)}
              placeholder="https://webhook.site"
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Return URL</label>
            <input
              type="url"
              value={returnUrl}
              onChange={e => setReturnUrl(e.target.value)}
              placeholder=""
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
        </div>
      </section>

      {/* Transactions */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Transactions</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Who should pay for transaction fees?</label>
            <select
              value={transactionFeePayer}
              onChange={e => setTransactionFeePayer(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            >
              <option>Charge me the transaction fees</option>
              <option>Charge the customer</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Who should pay for transfer fees?</label>
            <input
              type="text"
              value={transferFeePayer}
              onChange={e => setTransferFeePayer(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Allow to retry payment?</label>
            <input
              type="number"
              value={retryMinutes}
              onChange={e => setRetryMinutes(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
            />
            <span className="text-xs text-gray-400">* Value is in minutes. Default is 60 minutes.</span>
          </div>
          <div className="flex items-center gap-4 mt-8">
            <label className="block text-sm font-medium text-gray-700">Enable API Transfers</label>
            <input
              type="checkbox"
              checked={apiTransfers}
              onChange={e => setApiTransfers(e.target.checked)}
              className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500 h-5 w-5"
            />
          </div>
        </div>
      </section>

      {/* Payment Methods */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Payment Methods</h3>
        {/* Wallets */}
        <div className="mb-4 border rounded-lg">
          <button type="button" className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 rounded-t-lg" onClick={() => setWalletsOpen(w => !w)}>
            <span className="font-semibold">Wallets</span>
            <span>{walletsOpen ? '▲' : '▼'}</span>
          </button>
          {walletsOpen && (
            <div className="p-4 space-y-2">
              {Object.keys(wallets).map((w, i) => (
                <label key={w} className="flex items-center gap-2 text-gray-700 text-sm">
                  <input
                    type="checkbox"
                    checked={wallets[w]}
                    onChange={e => setWallets(prev => ({ ...prev, [w]: e.target.checked }))}
                    className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500"
                  />
                  {w}
                </label>
              ))}
              <button type="button" className="mt-2 px-6 py-2 bg-gray-200 text-gray-700 font-semibold rounded-lg shadow transition-all duration-200">Save</button>
            </div>
          )}
        </div>
        {/* Banks */}
        <div className="mb-4 border rounded-lg">
          <button type="button" className="w-full flex items-center justify-between px-4 py-3 bg-gray-50 rounded-t-lg" onClick={() => setBanksOpen(b => !b)}>
            <span className="font-semibold">Banks</span>
            <span>{banksOpen ? '▲' : '▼'}</span>
          </button>
          {banksOpen && (
            <div className="p-4 space-y-2">
              {Object.keys(banks).map((b, i) => (
                <label key={b} className="flex items-center gap-2 text-gray-700 text-sm">
                  <input
                    type="checkbox"
                    checked={banks[b]}
                    onChange={e => setBanks(prev => ({ ...prev, [b]: e.target.checked }))}
                    className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500"
                  />
                  {b}
                </label>
              ))}
              <button type="button" className="mt-2 px-6 py-2 bg-gray-200 text-gray-700 font-semibold rounded-lg shadow transition-all duration-200">Save</button>
            </div>
          )}
        </div>
        {/* Cards */}
        <div className="flex items-center gap-4 mb-2">
          <label className="block text-sm font-medium text-gray-700">Cards</label>
          <input
            type="checkbox"
            checked={cardsEnabled}
            onChange={e => setCardsEnabled(e.target.checked)}
            className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500 h-5 w-5"
          />
        </div>
      </section>

      {/* Notification Emails */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Notification Emails</h3>
        <div className="space-y-2">
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={notifyImportant} onChange={e => setNotifyImportant(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Email me for important notifications
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={notifyReceipts} onChange={e => setNotifyReceipts(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Email payment receipts for customers
          </label>
        </div>
      </section>

      {/* Transaction Receipts */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Transaction receipts</h3>
        <div className="space-y-2">
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={receiptMe} onChange={e => setReceiptMe(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to me
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={receiptRecipients} onChange={e => setReceiptRecipients(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to recipients
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={receiptFinance} onChange={e => setReceiptFinance(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to Finance Email
          </label>
        </div>
      </section>

      {/* Transfer Receipts */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Transfer receipts</h3>
        <div className="space-y-2">
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={transferMe} onChange={e => setTransferMe(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to me
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={transferFinance} onChange={e => setTransferFinance(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to Finance Email
          </label>
        </div>
      </section>

      {/* Transfer Approval */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Transfer Approval</h3>
        <div className="space-y-2 mb-4">
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={approveUrl} onChange={e => setApproveUrl(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Approve transfers using URL verification (This takes precedence over OTP verification for transfers)
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={approveOtp} onChange={e => setApproveOtp(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Approve transfers using email OTP verification (This takes precedence over URL verification for payouts)
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input type="checkbox" checked={approveFinance} onChange={e => setApproveFinance(e.target.checked)} className="rounded border-gray-300 text-brand-green-600 focus:ring-brand-green-500" />
            Send to Finance Email
          </label>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-4">
          <input
            type="text"
            value={approvalUrl}
            onChange={e => setApprovalUrl(e.target.value)}
            placeholder="Enter Approval URL"
            className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
          />
          <input
            type="text"
            value={approvalSecret}
            onChange={e => setApprovalSecret(e.target.value)}
            placeholder="Enter Approval Secret"
            className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
          />
        </div>
      </section>

      {/* Finance Email */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Finance Email</h3>
        <input
          type="email"
          value={financeEmail}
          onChange={e => setFinanceEmail(e.target.value)}
          placeholder="Enter Finance Email"
          className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
        />
      </section>

      {/* Exported Files Email */}
      <section>
        <h3 className="text-lg font-bold text-gray-900 mb-4">Email to receive exported files</h3>
        <div className="flex items-center gap-8 mb-4">
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input
              type="radio"
              checked={exportEmailType === 'default'}
              onChange={() => setExportEmailType('default')}
              className="text-brand-green-600 focus:ring-brand-green-500"
            />
            Use Default
          </label>
          <label className="flex items-center gap-2 text-gray-700 text-sm">
            <input
              type="radio"
              checked={exportEmailType === 'custom'}
              onChange={() => setExportEmailType('custom')}
              className="text-brand-green-600 focus:ring-brand-green-500"
            />
            Custom Email
          </label>
        </div>
        {exportEmailType === 'custom' && (
          <input
            type="email"
            value={customExportEmail}
            onChange={e => setCustomExportEmail(e.target.value)}
            placeholder="Enter custom email"
            className="w-full px-4 py-3 border border-gray-300 rounded-lg bg-gray-50 text-gray-700 placeholder-gray-400 focus:ring-2 focus:ring-brand-green-500 focus:border-brand-green-500 transition-colors duration-200"
          />
        )}
      </section>

      {/* Update Button */}
      <div className="flex justify-end pt-6">
        <button
          type="submit"
          className="px-8 py-3 bg-brand-green-600 hover:bg-brand-green-700 text-white font-semibold rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-brand-green-500 focus:ring-offset-2"
        >
          Update Settings
        </button>
      </div>
    </form>
  )
} 