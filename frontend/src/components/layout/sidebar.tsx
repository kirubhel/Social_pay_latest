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
  ChevronDownIcon,
  ChevronRightIcon,
  UserCircleIcon,
} from '@heroicons/react/24/outline'
import { useState } from 'react'

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
  { name: 'Settings', href: '/dashboard/settings', icon: CogIcon },
]

interface SidebarProps {
  onClose?: () => void
}

export function Sidebar({ onClose }: SidebarProps) {
  const pathname = usePathname()
  const [openMenu, setOpenMenu] = useState<string | null>(null)

  const isActive = (href: string) => pathname === href || pathname.startsWith(href + '/')

  const handleLinkClick = () => {
    // Close mobile menu when a link is clicked
    if (onClose) {
      onClose()
    }
  }

  // Submenus
  const inventorySubmenu = [
    { name: 'Orders', href: '/inventory/orders', icon: ClipboardDocumentListIcon },
    { name: 'Product', href: '/inventory/product', icon: QrCodeIcon },
    { name: 'Warehouse', href: '/inventory/warehouse', icon: BuildingStorefrontIcon },
    { name: 'Category', href: '/inventory/category', icon: BuildingStorefrontIcon },
  ]
  const accountsSubmenu = [
    { name: 'Profile', href: '/accounts/profile', icon: UserCircleIcon },
    { name: 'Settings', href: '/accounts/settings', icon: CogIcon },
  ]

  return (
    <div className="flex h-full w-64 flex-col bg-gradient-to-b from-white via-gray-50 to-white shadow-xl border-r border-gray-100">
      {/* Enhanced Logo Section */}
      <div className="flex h-20 items-center px-6 border-b border-gray-100 bg-gradient-to-r from-brand-green-50 to-brand-gold-50">
        <Link href="/" className="flex items-center group" onClick={handleLinkClick}>
          <div className="relative">
            <Image 
              src="/logo.png" 
              alt="Social Pay" 
              width={160} 
              height={32} 
              className="transition-transform duration-200 group-hover:scale-105" 
            />
            <div className="absolute inset-0 bg-gradient-to-r from-brand-green-500/10 to-brand-gold-500/10 rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200" />
          </div>
        </Link>
      </div>

      <nav className="flex-1 space-y-8 px-4 py-6 overflow-y-auto">
        {/* General Menu */}
        <div>
          <h3 className="px-3 mb-4 text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-2">
            <div className="w-2 h-2 bg-gradient-to-r from-brand-green-500 to-brand-gold-500 rounded-full"></div>
            General Menu
          </h3>
          <div className="space-y-2">
            {/* Dashboard */}
            <Link
              href="/dashboard"
              className={cn(
                'group relative flex items-center rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
                pathname === '/dashboard'
                  ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                  : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
              )}
              onClick={handleLinkClick}
            >
              <HomeIcon className="h-5 w-5 mr-3" />
              Dashboard
            </Link>
            {/* Inventory with submenu */}
            <div>
              <button
                type="button"
                className={cn(
                  'group relative flex items-center w-full rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
                  isActive('/inventory') || openMenu === 'inventory'
                    ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                    : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
                )}
                onClick={() => setOpenMenu(openMenu === 'inventory' ? null : 'inventory')}
              >
                <CreditCardIcon className="h-5 w-5 mr-3" />
                Inventory
                <span className="ml-auto">
                  {openMenu === 'inventory' ? <ChevronDownIcon className="h-4 w-4" /> : <ChevronRightIcon className="h-4 w-4" />}
                </span>
              </button>
              {openMenu === 'inventory' && (
                <div className="ml-8 mt-1 space-y-1">
                  {inventorySubmenu.map(item => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                        'flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-all duration-200',
                  isActive(item.href)
                          ? 'bg-brand-green-50 text-brand-green-700 font-bold'
                          : 'text-gray-700 hover:bg-brand-green-50 hover:text-brand-green-700'
                      )}
                      onClick={handleLinkClick}
                    >
                      <item.icon className="h-4 w-4 mr-2" />
                      {item.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>
            {/* Accounts with submenu */}
            <div>
              <button
                type="button"
                className={cn(
                  'group relative flex items-center w-full rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
                  isActive('/accounts') || openMenu === 'accounts'
                    ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                    : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
                )}
                onClick={() => setOpenMenu(openMenu === 'accounts' ? null : 'accounts')}
              >
                <UserGroupIcon className="h-5 w-5 mr-3" />
                Accounts
                <span className="ml-auto">
                  {openMenu === 'accounts' ? <ChevronDownIcon className="h-4 w-4" /> : <ChevronRightIcon className="h-4 w-4" />}
                </span>
              </button>
              {openMenu === 'accounts' && (
                <div className="ml-8 mt-1 space-y-1">
                  {accountsSubmenu.map(item => (
                    <Link
                      key={item.name}
                      href={item.href}
                    className={cn(
                        'flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-all duration-200',
                      isActive(item.href) 
                          ? 'bg-brand-green-50 text-brand-green-700 font-bold'
                          : 'text-gray-700 hover:bg-brand-green-50 hover:text-brand-green-700'
                    )}
                      onClick={handleLinkClick}
                    >
                      <item.icon className="h-4 w-4 mr-2" />
                      {item.name}
                    </Link>
                  ))}
                </div>
              )}
            </div>
            {/* Gateways and other menu items... */}
            <Link
              href="/gateways"
              className={cn(
                'group relative flex items-center rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
                isActive('/gateways')
                  ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                  : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
              )}
              onClick={handleLinkClick}
            >
              <GlobeAltIcon className="h-5 w-5 mr-3" />
              Gateways
              </Link>
            {/* ...rest of the menu */}
          </div>
        </div>

        {/* Management Menu */}
        <div>
          <h3 className="px-3 mb-4 text-xs font-bold text-gray-500 uppercase tracking-wider flex items-center gap-2">
            <div className="w-2 h-2 bg-gradient-to-r from-brand-gold-500 to-brand-green-500 rounded-full"></div>
            Management Menu
          </h3>
          <div className="space-y-2">
            {managementMenu.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                onClick={handleLinkClick}
                className={cn(
                  'group relative flex items-center rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
                  isActive(item.href)
                    ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                    : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
                )}
              >
                {/* Background decoration for active state */}
                {isActive(item.href) && (
                  <div className="absolute inset-0 bg-gradient-to-r from-white/10 to-transparent opacity-20" />
                )}
                
                <div className={cn(
                  'flex items-center justify-center w-8 h-8 rounded-lg mr-3 transition-all duration-200',
                  isActive(item.href) 
                    ? 'bg-white/20 backdrop-blur-sm' 
                    : 'bg-gray-100 group-hover:bg-brand-green-100 group-hover:shadow-sm'
                )}>
                  <item.icon
                    className={cn(
                      'h-5 w-5 transition-all duration-200',
                      isActive(item.href) 
                        ? 'text-white' 
                        : 'text-gray-500 group-hover:text-brand-green-600'
                    )}
                    aria-hidden="true"
                  />
                </div>
                
                <span className="relative z-10">{item.name}</span>
                
                {/* Active indicator */}
                {isActive(item.href) && (
                  <div className="absolute right-2 w-2 h-2 bg-white rounded-full animate-pulse" />
                )}
              </Link>
            ))}
          </div>
        </div>
      </nav>

      {/* Enhanced Bottom Menu */}
      <div className="border-t border-gray-100 bg-gradient-to-r from-gray-50 to-white p-4 space-y-2">
        {bottomMenu.map((item) => (
          <Link
            key={item.name}
            href={item.href}
            onClick={handleLinkClick}
            className={cn(
              'group relative flex items-center rounded-xl px-3 py-3 text-sm font-semibold transition-all duration-200 overflow-hidden',
              isActive(item.href)
                ? 'bg-gradient-to-r from-brand-green-500 to-brand-green-600 text-white shadow-lg shadow-brand-green-500/25 transform scale-105'
                : 'text-gray-700 hover:bg-gradient-to-r hover:from-brand-green-50 hover:to-brand-gold-50 hover:text-brand-green-700 hover:transform hover:scale-105 hover:shadow-md'
            )}
          >
            {/* Background decoration for active state */}
            {isActive(item.href) && (
              <div className="absolute inset-0 bg-gradient-to-r from-white/10 to-transparent opacity-20" />
            )}
            
            <div className={cn(
              'flex items-center justify-center w-8 h-8 rounded-lg mr-3 transition-all duration-200',
              isActive(item.href) 
                ? 'bg-white/20 backdrop-blur-sm' 
                : 'bg-gray-100 group-hover:bg-brand-green-100 group-hover:shadow-sm'
            )}>
              <item.icon
                className={cn(
                  'h-5 w-5 transition-all duration-200',
                  isActive(item.href) 
                    ? 'text-white' 
                    : 'text-gray-500 group-hover:text-brand-green-600'
                )}
                aria-hidden="true"
              />
            </div>
            
            <span className="relative z-10">{item.name}</span>
            
            {/* Active indicator */}
            {isActive(item.href) && (
              <div className="absolute right-2 w-2 h-2 bg-white rounded-full animate-pulse" />
            )}
          </Link>
        ))}
        
        {/* Professional footer */}
        <div className="mt-6 pt-4 border-t border-gray-200">
          <div className="text-center">
            <p className="text-xs text-gray-500 font-medium">
              Social Pay v2.0
            </p>
            <p className="text-xs text-gray-400 mt-1">
              Payment Gateway
            </p>
          </div>
        </div>
      </div>
    </div>
  )
} 