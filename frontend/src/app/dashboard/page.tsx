'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  ChevronDownIcon,
  MagnifyingGlassIcon,
  FunnelIcon,
  ArrowDownTrayIcon,
  EllipsisVerticalIcon
} from '@heroicons/react/24/outline'
import { 
  ArrowTrendingUpIcon, 
  ArrowTrendingDownIcon,
  MinusIcon 
} from '@heroicons/react/20/solid'
import { cn } from '@/lib/utils'

// Mock data for the dashboard
const stats = [
  {
    title: 'Total Sales',
    value: '12',
    change: 'up from yesterday',
    trending: 'up',
    highlighted: true,
  },
  {
    title: 'Total Products',
    value: '12',
    change: 'up from yesterday',
    trending: 'up',
    highlighted: false,
  },
  {
    title: 'Total Orders',
    value: '28',
    change: '5% up from yesterday',
    trending: 'up',
    highlighted: false,
  },
  {
    title: 'Total Customers',
    value: '453',
    change: '6% up from yesterday',
    trending: 'up',
    highlighted: false,
  },
  {
    title: 'Was House',
    value: '12',
    change: 'Same amount',
    trending: 'same',
    highlighted: false,
  },
]

const orders = [
  {
    id: '#123545',
    itemName: 'Air Jordan',
    quantity: 2,
    status: 'Delivered',
    totalAmount: '8755 Br',
    billingAddress: 'Bole Ednamall E......',
    action: 'Completed',
  },
  {
    id: '#123545',
    itemName: 'All Star',
    quantity: 1,
    status: 'Pending',
    totalAmount: '6235 Br',
    billingAddress: '6 kilo University f......',
    action: 'Pending',
  },
  {
    id: '#123545',
    itemName: 'Black Puma',
    quantity: 1,
    status: 'Canceled',
    totalAmount: '3450 Br',
    billingAddress: 'Wolo Sefer infront d......',
    action: 'Canceled',
  },
]

const chartData = [
  { name: 'Product 1', value: 200 },
  { name: 'Product 1', value: 400 },
  { name: 'Product 1', value: 500 },
  { name: 'Product 1', value: 3879 },
  { name: 'Product 1', value: 600 },
  { name: 'Product 1', value: 550 },
]

export default function DashboardPage() {
  const [activeTab, setActiveTab] = useState('overview')

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-600">Today</span>
          <ChevronDownIcon className="h-4 w-4 text-gray-600" />
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('overview')}
            className={cn(
              'py-2 px-1 border-b-2 font-medium text-sm',
              activeTab === 'overview'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            )}
          >
            Overview
          </button>
          <button
            onClick={() => setActiveTab('transactions')}
            className={cn(
              'py-2 px-1 border-b-2 font-medium text-sm',
              activeTab === 'transactions'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            )}
          >
            Transactions
          </button>
        </nav>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-5">
        {stats.map((stat, index) => (
          <Card
            key={index}
            className={cn(
              'relative',
              stat.highlighted 
                ? 'bg-gradient-to-r from-orange-400 to-orange-500 text-white border-orange-500' 
                : 'border-2 border-dashed border-blue-300 bg-white'
            )}
          >
            <CardContent className="p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className={cn(
                    'text-xs font-medium mb-1',
                    stat.highlighted ? 'text-orange-100' : 'text-gray-600'
                  )}>
                    {stat.title}
                  </p>
                  <p className={cn(
                    'text-2xl font-bold',
                    stat.highlighted ? 'text-white' : 'text-gray-900'
                  )}>
                    {stat.value}
                  </p>
                  <div className="flex items-center mt-1">
                    {stat.trending === 'up' && (
                      <ArrowTrendingUpIcon className={cn(
                        'h-3 w-3 mr-1',
                        stat.highlighted ? 'text-orange-100' : 'text-green-500'
                      )} />
                    )}
                    {stat.trending === 'down' && (
                      <ArrowTrendingDownIcon className={cn(
                        'h-3 w-3 mr-1',
                        stat.highlighted ? 'text-orange-100' : 'text-red-500'
                      )} />
                    )}
                    {stat.trending === 'same' && (
                      <MinusIcon className={cn(
                        'h-3 w-3 mr-1',
                        stat.highlighted ? 'text-orange-100' : 'text-gray-500'
                      )} />
                    )}
                    <p className={cn(
                      'text-xs',
                      stat.highlighted ? 'text-orange-100' : 'text-gray-600'
                    )}>
                      {stat.change}
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Charts Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Bar Chart */}
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <CardTitle className="text-lg font-semibold">Sales by Product</CardTitle>
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-600">This Week</span>
              <ChevronDownIcon className="h-4 w-4 text-gray-600" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="h-64 flex items-end justify-center space-x-4">
              {chartData.map((item, index) => (
                <div key={index} className="flex flex-col items-center">
                  <div
                    className={cn(
                      'w-12 rounded-t-md relative',
                      index === 3 ? 'bg-green-500' : 'bg-orange-400'
                    )}
                    style={{ height: `${(item.value / 4000) * 200}px` }}
                  >
                    {index === 3 && (
                      <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs px-2 py-1 rounded">
                        {item.value}
                      </div>
                    )}
                  </div>
                  <span className="text-xs text-gray-600 mt-2">{item.name}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Pie Chart */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg font-semibold">Sales by Product</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-64 flex items-center justify-center">
              <div className="relative w-40 h-40">
                {/* Simple pie chart representation */}
                <div className="w-full h-full rounded-full border-8 border-green-500 border-r-blue-500 border-b-gray-300 border-l-yellow-500"></div>
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="w-20 h-20 bg-white rounded-full"></div>
                </div>
              </div>
            </div>
            <div className="mt-4 grid grid-cols-2 gap-2 text-sm">
              <div className="flex items-center">
                <div className="w-3 h-3 bg-green-500 rounded mr-2"></div>
                <span>Product A</span>
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 bg-blue-500 rounded mr-2"></div>
                <span>Product B</span>
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 bg-gray-300 rounded mr-2"></div>
                <span>Product C</span>
              </div>
              <div className="flex items-center">
                <div className="w-3 h-3 bg-yellow-500 rounded mr-2"></div>
                <span>Product D</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Orders Table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle className="text-lg font-semibold">Orders</CardTitle>
          <div className="flex items-center gap-2">
            <button className="flex items-center gap-2 px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50">
              <MagnifyingGlassIcon className="h-4 w-4" />
            </button>
            <button className="flex items-center gap-2 px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50">
              <FunnelIcon className="h-4 w-4" />
            </button>
            <button className="flex items-center gap-2 px-3 py-1 text-sm bg-gray-900 text-white rounded-md hover:bg-gray-800">
              <ArrowDownTrayIcon className="h-4 w-4" />
              Export
            </button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Order ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Item Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Quantity
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Total Amount
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Billing Address
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {orders.map((order, index) => (
                  <tr key={index}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-green-600 font-medium">
                      {order.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {order.itemName}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {order.quantity}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={cn(
                        'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                        order.status === 'Delivered' && 'bg-green-100 text-green-800',
                        order.status === 'Pending' && 'bg-yellow-100 text-yellow-800',
                        order.status === 'Canceled' && 'bg-red-100 text-red-800'
                      )}>
                        {order.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {order.totalAmount}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {order.billingAddress}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className="flex items-center gap-2">
                        <span className={cn(
                          'flex items-center gap-1 text-xs',
                          order.action === 'Completed' && 'text-green-600',
                          order.action === 'Pending' && 'text-yellow-600',
                          order.action === 'Canceled' && 'text-red-600'
                        )}>
                          <div className={cn(
                            'w-2 h-2 rounded-full',
                            order.action === 'Completed' && 'bg-green-500',
                            order.action === 'Pending' && 'bg-yellow-500',
                            order.action === 'Canceled' && 'bg-red-500'
                          )}></div>
                          {order.action}
                        </span>
                        <ChevronDownIcon className="h-4 w-4 text-gray-400" />
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </div>
  )
} 