'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  ChevronDownIcon,
  MagnifyingGlassIcon,
  FunnelIcon,
  ArrowDownTrayIcon,
  EllipsisVerticalIcon,
  CalendarIcon,
  ClockIcon,
  CurrencyDollarIcon,
  ShoppingCartIcon,
  UsersIcon
} from '@heroicons/react/24/outline'
import { 
  ArrowTrendingUpIcon, 
  ArrowTrendingDownIcon,
  MinusIcon
} from '@heroicons/react/20/solid'
import { cn } from '@/lib/utils'

const stats = [
  {
    title: 'Total Revenue',
    value: '1,245,670',
    currency: 'ETB',
    change: '+12.5%',
    changeText: 'from last month',
    trending: 'up',
    highlighted: true,
    icon: CurrencyDollarIcon,
  },
  {
    title: 'Total Transactions',
    value: '2,847',
    change: '+8.2%',
    changeText: 'from last month',
    trending: 'up',
    highlighted: false,
    icon: ShoppingCartIcon,
  },
  {
    title: 'Active Customers',
    value: '1,432',
    change: '+23.1%',
    changeText: 'from last month',
    trending: 'up',
    highlighted: false,
    icon: UsersIcon,
  },
  {
    title: 'Success Rate',
    value: '98.7',
    suffix: '%',
    change: '+0.3%',
    changeText: 'from last month',
    trending: 'up',
    highlighted: false,
    icon: ArrowTrendingUpIcon,
  },
  {
    title: 'Pending Orders',
    value: '142',
    change: 'Same as yesterday',
    changeText: '',
    trending: 'same',
    highlighted: false,
    icon: ClockIcon,
  },
]

const recentTransactions = [
  {
    id: 'TXN-2024-001',
    customer: 'Abebe Kebede',
    amount: '2,450 ETB',
    status: 'Completed',
    method: 'Telebirr',
    time: '2 minutes ago',
    avatar: 'AK',
  },
  {
    id: 'TXN-2024-002',
    customer: 'Sara Ahmed',
    amount: '890 ETB',
    status: 'Processing',
    method: 'CBE Birr',
    time: '5 minutes ago',
    avatar: 'SA',
  },
  {
    id: 'TXN-2024-003',
    customer: 'Michael Tadesse',
    amount: '5,200 ETB',
    status: 'Failed',
    method: 'Bank Transfer',
    time: '12 minutes ago',
    avatar: 'MT',
  },
]

const chartData = [
  { day: 'Mon', revenue: 12400, transactions: 45 },
  { day: 'Tue', revenue: 15200, transactions: 62 },
  { day: 'Wed', revenue: 18900, transactions: 78 },
  { day: 'Thu', revenue: 22100, transactions: 89 },
  { day: 'Fri', revenue: 28500, transactions: 94 },
  { day: 'Sat', revenue: 19200, transactions: 67 },
  { day: 'Sun', revenue: 16800, transactions: 58 },
]

const paymentMethods = [
  { name: 'Telebirr', percentage: 42, color: 'bg-brand-green-500' },
  { name: 'CBE Birr', percentage: 28, color: 'bg-blue-500' },
  { name: 'M-Birr', percentage: 18, color: 'bg-brand-gold-500' },
  { name: 'Bank Transfer', percentage: 12, color: 'bg-purple-500' },
]

export default function DashboardPage() {
  const [activeTab, setActiveTab] = useState('overview')
  const [selectedPeriod, setSelectedPeriod] = useState('today')

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-white to-gray-50">
      <div className="max-w-7xl mx-auto space-y-8 p-6">
        {/* Enhanced Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 className="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
              Dashboard Overview
            </h1>
            <p className="text-gray-600 mt-1">
              Monitor your payment gateway performance in real-time
            </p>
          </div>
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-4 py-2 bg-white border border-gray-200 rounded-xl shadow-sm hover:shadow-md transition-shadow">
              <CalendarIcon className="h-4 w-4 text-gray-500" />
              <select 
                value={selectedPeriod}
                onChange={(e) => setSelectedPeriod(e.target.value)}
                className="text-sm font-medium text-gray-700 bg-transparent border-none outline-none cursor-pointer"
              >
                <option value="today">Today</option>
                <option value="week">This Week</option>
                <option value="month">This Month</option>
                <option value="year">This Year</option>
              </select>
              <ChevronDownIcon className="h-4 w-4 text-gray-500" />
            </div>
          </div>
        </div>

        {/* Enhanced Tabs */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-1">
          <nav className="flex space-x-1">
            <button
              onClick={() => setActiveTab('overview')}
              className={cn(
                'flex-1 py-3 px-6 text-sm font-semibold rounded-xl transition-all duration-200',
                activeTab === 'overview'
                  ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
              )}
            >
              Overview
            </button>
            <button
              onClick={() => setActiveTab('transactions')}
              className={cn(
                'flex-1 py-3 px-6 text-sm font-semibold rounded-xl transition-all duration-200',
                activeTab === 'transactions'
                  ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
              )}
            >
              Transactions
            </button>
          </nav>
        </div>

        {/* Enhanced Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
          {stats.map((stat, index) => (
            <Card
              key={index}
              className={cn(
                'group relative overflow-hidden border-0 shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1',
                stat.highlighted 
                  ? 'bg-gradient-to-br from-brand-green-500 via-brand-green-600 to-brand-green-700 text-white' 
                  : 'bg-white border border-gray-100 hover:border-gray-200'
              )}
            >
              <div className={cn(
                'absolute top-0 right-0 w-20 h-20 rounded-full -mr-10 -mt-10 opacity-10',
                stat.highlighted ? 'bg-white' : 'bg-brand-green-500'
              )} />
              
              <CardContent className="p-6">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-4">
                      <div className={cn(
                        'p-2 rounded-xl',
                        stat.highlighted 
                          ? 'bg-white/20 backdrop-blur-sm' 
                          : 'bg-brand-green-50'
                      )}>
                        <stat.icon className={cn(
                          'h-5 w-5',
                          stat.highlighted ? 'text-white' : 'text-brand-green-600'
                        )} />
                      </div>
                      <p className={cn(
                        'text-sm font-medium',
                        stat.highlighted ? 'text-white/90' : 'text-gray-600'
                      )}>
                        {stat.title}
                      </p>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex items-baseline gap-1">
                        {stat.currency && (
                          <span className={cn(
                            'text-sm font-medium',
                            stat.highlighted ? 'text-white/80' : 'text-gray-500'
                          )}>
                            {stat.currency}
                          </span>
                        )}
                        <p className={cn(
                          'text-2xl font-bold',
                          stat.highlighted ? 'text-white' : 'text-gray-900'
                        )}>
                          {stat.value}
                        </p>
                        {stat.suffix && (
                          <span className={cn(
                            'text-lg font-semibold',
                            stat.highlighted ? 'text-white/80' : 'text-gray-600'
                          )}>
                            {stat.suffix}
                          </span>
                        )}
                      </div>
                      
                      <div className="flex items-center gap-1">
                        {stat.trending === 'up' && (
                          <ArrowTrendingUpIcon className={cn(
                            'h-4 w-4',
                            stat.highlighted ? 'text-green-200' : 'text-green-500'
                          )} />
                        )}
                        {stat.trending === 'down' && (
                          <ArrowTrendingDownIcon className={cn(
                            'h-4 w-4',
                            stat.highlighted ? 'text-red-200' : 'text-red-500'
                          )} />
                        )}
                        {stat.trending === 'same' && (
                          <MinusIcon className={cn(
                            'h-4 w-4',
                            stat.highlighted ? 'text-white/60' : 'text-gray-500'
                          )} />
                        )}
                        <span className={cn(
                          'text-sm font-medium',
                          stat.trending === 'up' && (stat.highlighted ? 'text-green-200' : 'text-green-600'),
                          stat.trending === 'down' && (stat.highlighted ? 'text-red-200' : 'text-red-600'),
                          stat.trending === 'same' && (stat.highlighted ? 'text-white/60' : 'text-gray-500')
                        )}>
                          {stat.change}
                        </span>
                      </div>
                      
                      {stat.changeText && (
                        <p className={cn(
                          'text-xs',
                          stat.highlighted ? 'text-white/70' : 'text-gray-500'
                        )}>
                          {stat.changeText}
                        </p>
                      )}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Charts and Analytics Section */}
        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
          {/* Revenue Chart */}
          <Card className="xl:col-span-2 border-0 shadow-lg">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <div>
                <CardTitle className="text-xl font-bold text-gray-900">Revenue Overview</CardTitle>
                <CardDescription className="text-gray-600">Daily revenue and transaction trends</CardDescription>
              </div>
              <div className="flex items-center gap-2">
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-3 h-3 bg-brand-green-500 rounded-full"></div>
                  <span className="text-gray-600">Revenue</span>
                </div>
                <div className="flex items-center gap-2 text-sm">
                  <div className="w-3 h-3 bg-brand-gold-500 rounded-full"></div>
                  <span className="text-gray-600">Transactions</span>
                </div>
              </div>
            </CardHeader>
            <CardContent className="pt-4">
              <div className="h-80 flex items-end justify-between px-4 gap-3">
                {chartData.map((item, index) => (
                  <div key={index} className="flex flex-col items-center group cursor-pointer">
                    <div className="relative flex flex-col items-end gap-1 mb-4">
                      <div
                        className="w-8 bg-gradient-to-t from-brand-green-500 to-brand-green-400 rounded-t-md relative group-hover:from-brand-green-600 group-hover:to-brand-green-500 transition-colors duration-200"
                        style={{ height: `${(item.revenue / 30000) * 250}px` }}
                      >
                        <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-900 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                          {item.revenue.toLocaleString()} ETB
                        </div>
                      </div>
                      <div
                        className="w-8 bg-gradient-to-t from-brand-gold-500 to-brand-gold-400 rounded-t-md group-hover:from-brand-gold-600 group-hover:to-brand-gold-500 transition-colors duration-200"
                        style={{ height: `${(item.transactions / 100) * 150}px` }}
                      >
                        <div className="absolute -top-16 left-1/2 transform -translate-x-1/2 bg-gray-900 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                          {item.transactions} transactions
                        </div>
                      </div>
                    </div>
                    <span className="text-sm font-medium text-gray-600 group-hover:text-gray-900 transition-colors">
                      {item.day}
                    </span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Payment Methods Distribution */}
          <Card className="border-0 shadow-lg">
            <CardHeader>
              <CardTitle className="text-xl font-bold text-gray-900">Payment Methods</CardTitle>
              <CardDescription className="text-gray-600">Distribution by usage</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {paymentMethods.map((method, index) => (
                  <div key={index} className="space-y-2">
                    <div className="flex justify-between items-center">
                      <span className="font-medium text-gray-900">{method.name}</span>
                      <span className="text-sm font-semibold text-gray-600">{method.percentage}%</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
                      <div
                        className={cn(method.color, 'h-full rounded-full transition-all duration-1000 ease-out')}
                        style={{ width: `${method.percentage}%` }}
                      />
                    </div>
                  </div>
                ))}
              </div>

              <div className="mt-6 pt-6 border-t border-gray-100 space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Total Volume</span>
                  <span className="font-semibold text-gray-900">125,890 ETB</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Active Methods</span>
                  <span className="font-semibold text-gray-900">4 of 6</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Success Rate</span>
                  <span className="font-semibold text-green-600">98.7%</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Recent Transactions */}
        <Card className="border-0 shadow-lg">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-xl font-bold text-gray-900">Recent Transactions</CardTitle>
              <CardDescription className="text-gray-600">Latest payment activities</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <button className="flex items-center gap-2 px-4 py-2 text-sm border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors">
                <MagnifyingGlassIcon className="h-4 w-4" />
                Search
              </button>
              <button className="flex items-center gap-2 px-4 py-2 text-sm border border-gray-200 rounded-xl hover:bg-gray-50 transition-colors">
                <FunnelIcon className="h-4 w-4" />
                Filter
              </button>
              <button className="flex items-center gap-2 px-4 py-2 text-sm bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white rounded-xl hover:from-brand-green-600 hover:to-brand-green-700 transition-all shadow-lg shadow-brand-green-500/25">
                <ArrowDownTrayIcon className="h-4 w-4" />
                Export
              </button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentTransactions.map((transaction, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between p-4 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors group cursor-pointer"
                >
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 bg-gradient-to-r from-brand-green-500 to-brand-green-600 rounded-xl flex items-center justify-center text-white font-bold text-sm shadow-lg">
                      {transaction.avatar}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900 group-hover:text-brand-green-600 transition-colors">
                        {transaction.customer}
                      </p>
                      <p className="text-sm text-gray-600">{transaction.id}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-center gap-6">
                    <div className="text-right">
                      <p className="font-semibold text-gray-900">{transaction.amount}</p>
                      <p className="text-sm text-gray-600">{transaction.method}</p>
                    </div>
                    
                    <div className="text-right">
                      <span className={cn(
                        'inline-flex px-3 py-1 text-xs font-semibold rounded-full',
                        transaction.status === 'Completed' && 'bg-green-100 text-green-700',
                        transaction.status === 'Processing' && 'bg-yellow-100 text-yellow-700',
                        transaction.status === 'Failed' && 'bg-red-100 text-red-700'
                      )}>
                        {transaction.status}
                      </span>
                      <p className="text-sm text-gray-500 mt-1">{transaction.time}</p>
                    </div>
                    
                    <EllipsisVerticalIcon className="h-5 w-5 text-gray-400 group-hover:text-gray-600 transition-colors" />
                  </div>
                </div>
              ))}
            </div>
            
            <div className="mt-6 text-center">
              <button className="text-brand-green-600 hover:text-brand-green-700 font-semibold text-sm transition-colors">
                View All Transactions →
              </button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
} 