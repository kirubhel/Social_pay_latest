'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { cn } from '@/lib/utils'
import Image from 'next/image'
import {
  HomeIcon,
  CreditCardIcon,
  QrCodeIcon,
  CogIcon,
  KeyIcon,
  ChartBarIcon,
  BuildingStorefrontIcon,
  UserGroupIcon,
  BanknotesIcon,
  ArrowUpTrayIcon,
  ChatBubbleLeftIcon,
  ClipboardDocumentListIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline'

const generalMenu = [
  { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
  { name: 'Inventory', href: '/inventory', icon: ClipboardDocumentListIcon },
  { name: 'Accounts', href: '/accounts', icon: UserGroupIcon },
  { name: 'Gateways', href: '/gateways', icon: GlobeAltIcon },
  { name: 'Transactions', href: '/transactions', icon: CreditCardIcon },
]

const managementMenu = [
  { name: 'Manage Roles', href: '/manage-roles', icon: UserGroupIcon },
  { name: 'Manage Banks', href: '/manage-banks', icon: BanknotesIcon },
  { name: 'Withdrawals', href: '/withdrawals', icon: ArrowUpTrayIcon },
  { name: 'Message', href: '/message', icon: ChatBubbleLeftIcon },
]

const bottomMenu = [
  { name: 'Feedbacks', href: '/feedbacks', icon: ChatBubbleLeftIcon },
  { name: 'Settings', href: '/settings', icon: CogIcon },
]

export function Sidebar() {
  const pathname = usePathname()

  const isActive = (href: string) => pathname === href || pathname.startsWith(href + '/')

  return (
    <div className="flex h-full w-64 flex-col bg-white shadow-sm border-r">
      {/* Logo */}
      <div className="flex h-16 items-center px-6 border-b">
        <Image src="/logo.png" alt="Social Pay" width={160} height={32} className="mr-2" />
      
      </div>

      <nav className="flex-1 space-y-6 px-4 py-4">
        {/* General Menu */}
        <div>
          <h3 className="px-2 text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
            General Menu
          </h3>
          <div className="space-y-1">
            {generalMenu.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'group flex items-center rounded-md px-2 py-2 text-sm font-medium transition-colors',
                  isActive(item.href)
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                )}
              >
                <item.icon
                  className={cn(
                    'mr-3 h-5 w-5',
                    isActive(item.href) ? 'text-blue-500' : 'text-gray-400 group-hover:text-gray-500'
                  )}
                  aria-hidden="true"
                />
                {item.name}
              </Link>
            ))}
          </div>
        </div>

        {/* Management Menu */}
        <div>
          <h3 className="px-2 text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
            Management Menu
          </h3>
          <div className="space-y-1">
            {managementMenu.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'group flex items-center rounded-md px-2 py-2 text-sm font-medium transition-colors',
                  isActive(item.href)
                    ? 'bg-blue-50 text-blue-700'
                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                )}
              >
                <item.icon
                  className={cn(
                    'mr-3 h-5 w-5',
                    isActive(item.href) ? 'text-blue-500' : 'text-gray-400 group-hover:text-gray-500'
                  )}
                  aria-hidden="true"
                />
                {item.name}
              </Link>
            ))}
          </div>
        </div>
      </nav>

      {/* Bottom Menu */}
      <div className="border-t p-4 space-y-1">
        {bottomMenu.map((item) => (
          <Link
            key={item.name}
            href={item.href}
            className={cn(
              'group flex items-center rounded-md px-2 py-2 text-sm font-medium transition-colors',
              isActive(item.href)
                ? 'bg-blue-50 text-blue-700'
                : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
            )}
          >
            <item.icon
              className={cn(
                'mr-3 h-5 w-5',
                isActive(item.href) ? 'text-blue-500' : 'text-gray-400 group-hover:text-gray-500'
              )}
              aria-hidden="true"
            />
            {item.name}
          </Link>
        ))}
      </div>
    </div>
  )
} 