import axios from 'axios'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8004'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Don't redirect, just reject the error
    return Promise.reject(error)
  }
)

// Authentication API endpoints
export const authAPI = {
  // Initialize pre-session for registration (no phone needed)
  initPreSession: async () => {
    try {
      const response = await apiClient.post('/auth/init', {
        prefix: '',
        number: ''
      })
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to initialize session'
        }
      }
    }
  },

  // Initialize pre-session for login (phone required)
  initPreSessionForLogin: async (phone: { prefix: string, number: string }) => {
    try {
      const requestData = {
        prefix: '', // Empty for login initialization
        number: phone.number // Only the phone number for lookup
      }
      const response = await apiClient.post('/auth/init', requestData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to initialize login session'
        }
      }
    }
  },

  // Register user
  signUp: async (userData: {
    title: string
    first_name: string
    last_name: string
    phone_prefix: string
    phone_number: string
    password: string
    password_hint?: string
    confirm_password: string
  }) => {
    try {
      const response = await apiClient.post('/auth/sign-up', userData)
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Registration failed'
        }
      }
    }
  },

  // Verify OTP after registration
  verifyOTP: async (token: string, code: string) => {
    try {
      const response = await apiClient.post('/auth/verify-otp', {
        token,
        code
      })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'OTP verification failed'
        }
      }
    }
  },

  // Sign in with phone + OTP
  signIn: async (token: string, code: string, phone: { prefix: string, number: string }) => {
    try {
      const requestData = {
        token,
        code,
        phone: {
          prefix: '', // Empty prefix to match initialization
          number: phone.number // Should match what was stored during init
        }
      }
      const response = await apiClient.post('/auth/sign-in', requestData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Sign in failed'
        }
      }
    }
  },

  // Check session
  checkSession: async () => {
    try {
      const response = await apiClient.get('/auth/check')
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Session check failed'
        }
      }
    }
  },

  // Set password/2FA
  setPassword: async (password: string) => {
    try {
      const response = await apiClient.post('/auth/password', { password })
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to set password'
        }
      }
    }
  },

  // Check password/2FA
  checkPassword: async (password: string) => {
    try {
      const response = await apiClient.post('/auth/password/check', { password })
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Password check failed'
        }
      }
    }
  },

  // Login with phone and password
  login: async ({ prefix, number, password }: { prefix: string, number: string, password: string }) => {
    try {
      const response = await apiClient.post('/auth/login', {
        prefix,
        number,
        password,
      })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Login failed'
        }
      }
    }
  },

  // Get 2FA status
  get2FAStatus: async () => {
    try {
      const response = await apiClient.get('/auth/2fa/status')
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to get 2FA status'
        }
      }
    }
  },

  // Enable 2FA
  enable2FA: async () => {
    try {
      const response = await apiClient.post('/auth/2fa/enable')
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to enable 2FA'
        }
      }
    }
  },

  // Disable 2FA
  disable2FA: async (password: string) => {
    try {
      const response = await apiClient.post('/auth/2fa/disable', { password })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to disable 2FA'
        }
      }
    }
  },

  // Verify 2FA setup
  verify2FASetup: async (code: string) => {
    try {
      const response = await apiClient.post('/auth/2fa/verify-setup', { code })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to verify 2FA setup'
        }
      }
    }
  },

  // Send 2FA code for login verification
  send2FACode: async (token: string) => {
    try {
      const response = await apiClient.post('/auth/2fa/send-login', { token })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to send 2FA code'
        }
      }
    }
  },

  // Resend 2FA code
  resend2FACode: async () => {
    try {
      const response = await apiClient.post('/auth/2fa/resend')
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to resend 2FA code'
        }
      }
    }
  },

  // Verify 2FA during login
  verify2FALogin: async (token: string, code: string) => {
    try {
      const response = await apiClient.post('/auth/2fa/verify-login', { token, code })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to verify 2FA code'
        }
      }
    }
  },

  // Update password
  updatePassword: async (currentPassword: string, newPassword: string) => {
    try {
      const response = await apiClient.post('/auth/password/update', {
        current_password: currentPassword,
        new_password: newPassword
      })
      // Handle nested response structure
      if (response.data.success && response.data.data) {
        return { success: true, data: response.data.data }
      }
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to update password'
        }
      }
    }
  },

  // Update user profile
  updateUserProfile: async (userData: {
    first_name?: string;
    last_name?: string;
    sir_name?: string;
    phone_number?: string;
  }) => {
    try {
      const response = await apiClient.post('/api/v1/user-update', userData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to update user profile'
        }
      }
    }
  },

  // Get user profile
  getUserProfile: async () => {
    try {
      const response = await apiClient.get('/api/v1/user-profile')
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to fetch user profile'
        }
      }
    }
  },

  // Get merchant details
  getMerchantDetails: async () => {
    try {
      const response = await apiClient.get('/api/v1/get/merchant/details')
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to fetch merchant details'
        }
      }
    }
  },

  // Update merchant details
  updateMerchantDetails: async (merchantData: {
    merchant_id: string;
    merchant: {
      legal_name?: string;
      trading_name?: string;
      business_reg_number?: string;
      tax_identifier?: string;
      industry_type?: string;
      business_type?: string;
      website_url?: string;
    };
    address?: {
      email?: string;
      phone_number?: string;
      personal_name?: string;
      region?: string;
      city?: string;
      sub_city?: string;
    };
  }) => {
    try {
      const response = await apiClient.post('/api/v1/merchant/update', merchantData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to update merchant details'
        }
      }
    }
  },

  // Create merchant
  createMerchant: async (merchantData: {
    legal_name?: string;
    trading_name?: string;
    business_reg_number?: string;
    tax_identifier?: string;
    industry_type?: string;
    business_type?: string;
    website_url?: string;
    address?: {
      email?: string;
      phone_number?: string;
      personal_name?: string;
      region?: string;
      city?: string;
      sub_city?: string;
    };
  }) => {
    try {
      const response = await apiClient.post('/api/merchant/create', merchantData)
      return { success: true, data: response.data }
    } catch (error: any) {
      return {
        success: false,
        error: {
          message: error.response?.data?.message || 'Failed to create merchant'
        }
      }
    }
  },
}

export default apiClient 