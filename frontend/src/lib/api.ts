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
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('authToken')
      window.location.href = '/auth/login'
    }
    return Promise.reject(error)
  }
)

// Authentication API endpoints
export const authAPI = {
  // Initialize pre-session for registration (no phone needed)
  initPreSession: async () => {
    const response = await apiClient.post('/auth/init', {
      prefix: '',
      number: ''
    })
    return response.data
  },

  // Initialize pre-session for login (phone required)
  initPreSessionForLogin: async (phone: { prefix: string, number: string }) => {
    const requestData = {
      prefix: '', // Empty for login initialization
      number: phone.number // Only the phone number for lookup
    }
    console.log('Sending to /auth/init:', requestData)
    const response = await apiClient.post('/auth/init', requestData)
    return response.data
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
    const response = await apiClient.post('/auth/sign-up', userData)
    return response.data
  },

  // Verify OTP after registration
  verifyOTP: async (token: string, code: string) => {
    const response = await apiClient.post('/auth/verify-otp', {
      token,
      code
    })
    return response.data
  },

  // Sign in with phone + OTP
  signIn: async (token: string, code: string, phone: { prefix: string, number: string }) => {
    const requestData = {
      token,
      code,
      phone: {
        prefix: '', // Empty prefix to match initialization
        number: phone.number // Should match what was stored during init
      }
    }
    console.log('Sending to /auth/sign-in:', requestData)
    const response = await apiClient.post('/auth/sign-in', requestData)
    return response.data
  },

  // Check session
  checkSession: async () => {
    const response = await apiClient.get('/auth/check')
    return response.data
  },

  // Set password/2FA
  setPassword: async (password: string) => {
    const response = await apiClient.post('/auth/password', { password })
    return response.data
  },

  // Check password/2FA
  checkPassword: async (password: string) => {
    const response = await apiClient.post('/auth/password/check', { password })
    return response.data
  }
}

export default apiClient 