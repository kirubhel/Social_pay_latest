export interface User {
  id: string
  email: string
  name: string
  role: 'user' | 'merchant' | 'admin'
  createdAt: string
  updatedAt: string
}

export interface Transaction {
  id: string
  merchantId: string
  amount: number
  currency: string
  status: 'pending' | 'completed' | 'failed' | 'cancelled'
  type: 'payment' | 'withdrawal' | 'refund'
  medium: 'mpesa' | 'telebirr' | 'cbe' | 'awash' | 'socialpay'
  reference: string
  description?: string
  createdAt: string
  updatedAt: string
}

export interface Merchant {
  id: string
  name: string
  email: string
  phone: string
  status: 'active' | 'inactive' | 'suspended'
  businessType: string
  address: {
    street: string
    city: string
    country: string
    postalCode: string
  }
  createdAt: string
  updatedAt: string
}

export interface PaymentMethod {
  id: string
  name: string
  type: 'mobile' | 'bank' | 'card'
  enabled: boolean
  icon: string
}

export interface APIKey {
  id: string
  key: string
  name: string
  permissions: string[]
  status: 'active' | 'revoked'
  createdAt: string
  lastUsed?: string
}

export interface QRLink {
  id: string
  merchantId: string
  amount?: number
  description: string
  type: 'fixed' | 'dynamic'
  status: 'active' | 'inactive'
  url: string
  qrCode: string
  expiresAt?: string
  createdAt: string
}

export interface DashboardStats {
  totalTransactions: number
  totalRevenue: number
  successRate: number
  activeQRCodes: number
  transactionTrends: {
    period: string
    amount: number
    count: number
  }[]
} 