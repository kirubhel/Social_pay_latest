'use client'

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { formatCurrency, formatDate } from '@/lib/utils'
import { 
  CreditCardIcon, 
  BanknotesIcon, 
  CheckCircleIcon,
  QrCodeIcon 
} from '@heroicons/react/24/outline'

// Mock data - in real app, this would come from API
const stats = {
  totalTransactions: 1248,
  totalRevenue: 125000,
  successRate: 98.5,
  activeQRCodes: 24,
}

const recentTransactions = [
  {
    id: 'tx_001',
    amount: 2500,
    currency: 'ETB',
    status: 'completed',
    medium: 'mpesa',
    reference: 'MP001234',
    createdAt: new Date().toISOString(),
  },
  {
    id: 'tx_002',
    amount: 1800,
    currency: 'ETB',
    status: 'pending',
    medium: 'telebirr',
    reference: 'TB005678',
    createdAt: new Date(Date.now() - 3600000).toISOString(),
  },
  {
    id: 'tx_003',
    amount: 5200,
    currency: 'ETB',
    status: 'completed',
    medium: 'cbe',
    reference: 'CB009876',
    createdAt: new Date(Date.now() - 7200000).toISOString(),
  },
]

export default function DashboardPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600">Welcome back to Social Pay</p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Transactions</CardTitle>
            <CreditCardIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalTransactions.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">+12% from last month</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <BanknotesIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{formatCurrency(stats.totalRevenue)}</div>
            <p className="text-xs text-muted-foreground">+8% from last month</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
            <CheckCircleIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.successRate}%</div>
            <p className="text-xs text-muted-foreground">+0.3% from last month</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active QR Codes</CardTitle>
            <QrCodeIcon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.activeQRCodes}</div>
            <p className="text-xs text-muted-foreground">+2 new this week</p>
          </CardContent>
        </Card>
      </div>

      {/* Recent Transactions */}
      <Card>
        <CardHeader>
          <CardTitle>Recent Transactions</CardTitle>
          <CardDescription>Latest payment transactions in your account</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {recentTransactions.map((transaction) => (
              <div
                key={transaction.id}
                className="flex items-center justify-between p-4 border rounded-lg"
              >
                <div className="flex items-center space-x-4">
                  <div className="flex-shrink-0">
                    <div className={`w-3 h-3 rounded-full ${
                      transaction.status === 'completed' ? 'bg-green-500' : 
                      transaction.status === 'pending' ? 'bg-yellow-500' : 'bg-red-500'
                    }`} />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-gray-900">
                      {transaction.reference}
                    </p>
                    <p className="text-sm text-gray-500">
                      {transaction.medium.toUpperCase()} • {formatDate(transaction.createdAt)}
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-sm font-medium text-gray-900">
                    {formatCurrency(transaction.amount, transaction.currency)}
                  </p>
                  <p className={`text-sm capitalize ${
                    transaction.status === 'completed' ? 'text-green-600' : 
                    transaction.status === 'pending' ? 'text-yellow-600' : 'text-red-600'
                  }`}>
                    {transaction.status}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
} 