'use client'

import { Fragment } from 'react'
import { Menu, Transition } from '@headlessui/react'
import { 
  Bars3Icon, 
  BellIcon, 
  UserCircleIcon,
  MagnifyingGlassIcon,
  ChevronDownIcon,
  CogIcon,
  ArrowRightOnRectangleIcon
} from '@heroicons/react/24/outline'
import { useAuthStore } from '@/stores/auth'
import { cn } from '@/lib/utils'
import Image from 'next/image'

export function Header() {
  const { user, logout } = useAuthStore()

  return (
    <div className="sticky top-0 z-40 flex h-20 shrink-0 items-center gap-x-4 border-b border-gray-200/50 bg-white/80 backdrop-blur-xl px-4 shadow-lg sm:gap-x-6 sm:px-6 lg:px-8">
      {/* Mobile menu button */}
      <button 
        type="button" 
        className="-m-2.5 p-2.5 text-gray-700 hover:text-brand-green-600 lg:hidden rounded-xl hover:bg-brand-green-50 transition-all duration-200"
      >
        <span className="sr-only">Open sidebar</span>
        <Bars3Icon className="h-6 w-6" aria-hidden="true" />
      </button>

      {/* Mobile Logo */}
      <div className="flex items-center lg:hidden">
        <Image src="/logo.png" alt="Social Pay" width={32} height={32} className="mr-2" />
        <span className="text-lg font-bold bg-gradient-to-r from-brand-green-600 to-brand-gold-500 bg-clip-text text-transparent">
          SocialPay
        </span>
      </div>

      {/* Separator */}
      <div className="h-6 w-px bg-gradient-to-b from-transparent via-gray-300 to-transparent lg:hidden" aria-hidden="true" />

      <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
        {/* Enhanced Search */}
        <div className="relative flex flex-1 max-w-md">
          <div className="relative w-full group">
            <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-4">
              <MagnifyingGlassIcon className="h-5 w-5 text-gray-400 group-focus-within:text-brand-green-500 transition-colors duration-200" aria-hidden="true" />
            </div>
            <input
              className="block w-full rounded-2xl border-0 bg-gradient-to-r from-gray-50 to-gray-100 py-3 pl-12 pr-4 text-gray-900 placeholder:text-gray-500 focus:ring-2 focus:ring-brand-green-500 focus:bg-white hover:bg-white transition-all duration-200 shadow-sm hover:shadow-md focus:shadow-lg sm:text-sm sm:leading-6"
              placeholder="Search transactions, customers..."
              type="search"
            />
            <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-brand-green-500/5 to-brand-gold-500/5 opacity-0 group-focus-within:opacity-100 transition-opacity duration-200 pointer-events-none" />
          </div>
        </div>

        <div className="flex items-center gap-x-4 lg:gap-x-6">
          {/* Enhanced Language Selector */}
          <Menu as="div" className="relative">
            <Menu.Button className="flex items-center gap-x-2 px-3 py-2 text-sm font-semibold text-gray-700 hover:text-brand-green-700 bg-gray-50 hover:bg-brand-green-50 rounded-xl transition-all duration-200 hover:shadow-md group">
              <span className="text-sm">🇪🇹</span>
              <span>ENG</span>
              <ChevronDownIcon className="h-4 w-4 group-hover:text-brand-green-600 transition-colors duration-200" aria-hidden="true" />
            </Menu.Button>
            <Transition
              as={Fragment}
              enter="transition ease-out duration-200"
              enterFrom="transform opacity-0 scale-95"
              enterTo="transform opacity-100 scale-100"
              leave="transition ease-in duration-150"
              leaveFrom="transform opacity-100 scale-100"
              leaveTo="transform opacity-0 scale-95"
            >
              <Menu.Items className="absolute right-0 z-10 mt-3 w-32 origin-top-right rounded-2xl bg-white/95 backdrop-blur-xl py-3 shadow-xl ring-1 ring-gray-900/5 focus:outline-none border border-gray-100">
                <Menu.Item>
                  {({ active }) => (
                    <button
                      className={cn(
                        active ? 'bg-brand-green-50 text-brand-green-700' : 'text-gray-900',
                        'flex w-full items-center gap-x-2 px-4 py-2 text-sm font-medium transition-colors duration-150 rounded-xl mx-2'
                      )}
                    >
                      <span>🇪🇹</span>
                      <span>ENG</span>
                    </button>
                  )}
                </Menu.Item>
                <Menu.Item>
                  {({ active }) => (
                    <button
                      className={cn(
                        active ? 'bg-brand-green-50 text-brand-green-700' : 'text-gray-900',
                        'flex w-full items-center gap-x-2 px-4 py-2 text-sm font-medium transition-colors duration-150 rounded-xl mx-2'
                      )}
                    >
                      <span>🇪🇹</span>
                      <span>አማ</span>
                    </button>
                  )}
                </Menu.Item>
              </Menu.Items>
            </Transition>
          </Menu>

          {/* Enhanced Notifications */}
          <button 
            type="button" 
            className="relative -m-2.5 p-3 text-gray-400 hover:text-brand-green-600 bg-gray-50 hover:bg-brand-green-50 rounded-2xl transition-all duration-200 hover:shadow-md group"
          >
            <span className="sr-only">View notifications</span>
            <BellIcon className="h-6 w-6 group-hover:animate-pulse" aria-hidden="true" />
            {/* Enhanced notification badge */}
            <span className="absolute top-2 right-2 flex h-3 w-3">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-3 w-3 bg-gradient-to-r from-red-500 to-red-600 shadow-lg"></span>
            </span>
          </button>

          {/* Separator */}
          <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-gradient-to-b lg:from-transparent lg:via-gray-300 lg:to-transparent" aria-hidden="true" />

          {/* Enhanced Profile dropdown */}
          <Menu as="div" className="relative">
            <Menu.Button className="flex items-center gap-x-3 p-2 hover:bg-brand-green-50 rounded-2xl transition-all duration-200 hover:shadow-md group">
              <span className="sr-only">Open user menu</span>
              
              {/* Enhanced Avatar */}
              <div className="relative">
                <div className="h-10 w-10 rounded-2xl bg-gradient-to-r from-brand-green-500 to-brand-green-600 flex items-center justify-center shadow-lg group-hover:shadow-xl transition-all duration-200 group-hover:scale-105">
                  <span className="text-white font-bold text-sm">
                    {user?.name?.charAt(0) || 'U'}
                  </span>
                </div>
                <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-brand-green-400 to-brand-gold-400 opacity-0 group-hover:opacity-20 transition-opacity duration-200" />
              </div>
              
              {/* User Info */}
              <span className="hidden lg:flex lg:flex-col lg:items-start">
                <span className="text-sm font-bold text-gray-900 group-hover:text-brand-green-700 transition-colors duration-200">
                  {user?.name || 'User'}
                </span>
                <span className="text-xs text-gray-500 group-hover:text-brand-green-600 transition-colors duration-200">
                  Admin
                </span>
              </span>
              
              {/* Dropdown Arrow */}
              <ChevronDownIcon className="h-4 w-4 text-gray-400 group-hover:text-brand-green-600 transition-colors duration-200" />
            </Menu.Button>
            
            <Transition
              as={Fragment}
              enter="transition ease-out duration-200"
              enterFrom="transform opacity-0 scale-95"
              enterTo="transform opacity-100 scale-100"
              leave="transition ease-in duration-150"
              leaveFrom="transform opacity-100 scale-100"
              leaveTo="transform opacity-0 scale-95"
            >
              <Menu.Items className="absolute right-0 z-10 mt-3 w-56 origin-top-right rounded-2xl bg-white/95 backdrop-blur-xl py-3 shadow-xl ring-1 ring-gray-900/5 focus:outline-none border border-gray-100">
                {/* User Info Header */}
                <div className="px-4 py-3 border-b border-gray-100">
                  <p className="text-sm font-semibold text-gray-900">
                    {user?.name || 'User'}
                  </p>
                  <p className="text-xs text-gray-500">
                    {user?.email || 'user@socialpay.et'}
                  </p>
                </div>
                
                {/* Menu Items */}
                <div className="py-2">
                  <Menu.Item>
                    {({ active }) => (
                      <a
                        href="/settings"
                        className={cn(
                          active ? 'bg-brand-green-50 text-brand-green-700' : 'text-gray-900',
                          'flex items-center gap-x-3 px-4 py-3 text-sm font-medium transition-colors duration-150 rounded-xl mx-2'
                        )}
                      >
                        <CogIcon className="h-5 w-5" />
                        <span>Settings</span>
                      </a>
                    )}
                  </Menu.Item>
                  
                  <Menu.Item>
                    {({ active }) => (
                      <button
                        onClick={logout}
                        className={cn(
                          active ? 'bg-red-50 text-red-700' : 'text-gray-900',
                          'flex w-full items-center gap-x-3 px-4 py-3 text-sm font-medium transition-colors duration-150 rounded-xl mx-2'
                        )}
                      >
                        <ArrowRightOnRectangleIcon className="h-5 w-5" />
                        <span>Sign out</span>
                      </button>
                    )}
                  </Menu.Item>
                </div>
              </Menu.Items>
            </Transition>
          </Menu>
        </div>
      </div>
      
      {/* Gradient border */}
      <div className="absolute bottom-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-brand-green-200 to-transparent" />
    </div>
  )
} 