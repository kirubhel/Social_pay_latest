'use client'

import { Fragment, useState, useRef, useEffect } from 'react'
import { Menu, Transition } from '@headlessui/react'
import { 
  Bars3Icon, 
  BellIcon, 
  UserCircleIcon,
  MagnifyingGlassIcon,
  ChevronDownIcon,
  CogIcon,
  ArrowRightOnRectangleIcon,
  XMarkIcon,
  ClockIcon
} from '@heroicons/react/24/outline'
import { useAuthStore } from '@/stores/auth'
import { cn } from '@/lib/utils'
import Image from 'next/image'

// Mock search suggestions
const searchSuggestions = [
  { type: 'transaction', text: 'TXN-2024-001', description: 'Abebe Kebede - 2,450 ETB' },
  { type: 'customer', text: 'Sara Ahmed', description: 'Customer - 15 transactions' },
  { type: 'transaction', text: 'Failed payments', description: '3 failed transactions today' },
  { type: 'customer', text: 'Michael Tadesse', description: 'VIP Customer - 89 transactions' },
  { type: 'filter', text: 'Telebirr payments', description: 'Filter by payment method' },
]

const recentSearches = [
  'TXN-2024-001',
  'Abebe Kebede',
  'Failed payments',
  'Telebirr',
]

interface HeaderProps {
  onMobileMenuToggle?: () => void
}

export function Header({ onMobileMenuToggle }: HeaderProps) {
  const { user, logout } = useAuthStore()
  const [searchQuery, setSearchQuery] = useState('')
  const [isSearchFocused, setIsSearchFocused] = useState(false)
  const [showSuggestions, setShowSuggestions] = useState(false)
  const searchInputRef = useRef<HTMLInputElement>(null)
  const searchContainerRef = useRef<HTMLDivElement>(null)

  // Filter suggestions based on search query
  const filteredSuggestions = searchQuery.length > 0 
    ? searchSuggestions.filter(item => 
        item.text.toLowerCase().includes(searchQuery.toLowerCase()) ||
        item.description.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : []

  // Handle search input changes
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setSearchQuery(value)
    setShowSuggestions(value.length > 0 || isSearchFocused)
  }

  // Handle search focus
  const handleSearchFocus = () => {
    setIsSearchFocused(true)
    setShowSuggestions(true)
  }

  // Handle search blur
  const handleSearchBlur = () => {
    // Delay hiding suggestions to allow clicks
    setTimeout(() => {
      setIsSearchFocused(false)
      setShowSuggestions(false)
    }, 200)
  }

  // Clear search
  const clearSearch = () => {
    setSearchQuery('')
    setShowSuggestions(false)
    searchInputRef.current?.focus()
  }

  // Handle suggestion click
  const handleSuggestionClick = (suggestion: string) => {
    setSearchQuery(suggestion)
    setShowSuggestions(false)
    // Here you would typically trigger the actual search
    console.log('Searching for:', suggestion)
  }

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Escape') {
      setShowSuggestions(false)
      searchInputRef.current?.blur()
    }
    if (e.key === 'Enter') {
      setShowSuggestions(false)
      // Trigger search
      console.log('Searching for:', searchQuery)
    }
  }

  // Handle mobile menu toggle
  const handleMobileMenuToggle = () => {
    if (onMobileMenuToggle) {
      onMobileMenuToggle()
    }
  }

  // Close suggestions when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (searchContainerRef.current && !searchContainerRef.current.contains(event.target as Node)) {
        setShowSuggestions(false)
        setIsSearchFocused(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [])

  return (
    <header className="sticky top-0 z-50 flex h-16 shrink-0 items-center border-b border-gray-200/50 bg-white/95 backdrop-blur-xl px-4 shadow-sm sm:px-6 lg:px-8">
      {/* Sidebar/Logo section (left) */}
      <div className="flex items-center flex-shrink-0 lg:hidden">
      <button 
        type="button" 
          onClick={handleMobileMenuToggle}
          className="-m-2.5 p-2.5 text-gray-700 hover:text-brand-green-600 rounded-xl hover:bg-brand-green-50 transition-all duration-200"
      >
        <span className="sr-only">Open sidebar</span>
        <Bars3Icon className="h-6 w-6" aria-hidden="true" />
      </button>
        <div className="flex items-center ml-2">
          <Image src="/logo.png" alt="Social Pay" width={120} height={32} className="mr-2" />
        </div>
      </div>

      {/* Main flex: search (center) and right controls */}
      <div className="flex flex-1 items-center justify-between w-full">
        {/* Search Bar (centered) */}
        <div className="flex-1 flex justify-center">
          <div className="relative w-full max-w-lg" ref={searchContainerRef}>
  <div className="relative w-full group">
    {/* Search Icon */}
    <div className="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none">
      <MagnifyingGlassIcon className="h-5 w-5 text-gray-400 group-focus-within:text-brand-green-500 transition-colors duration-200" />
    </div>
    <input
      ref={searchInputRef}
      value={searchQuery}
      onChange={handleSearchChange}
      onFocus={handleSearchFocus}
      onBlur={handleSearchBlur}
      onKeyDown={handleKeyDown}
      type="search"
      placeholder="Search transactions, customers..."
                className="w-full pl-12 pr-12 py-2.5 rounded-xl bg-gray-50 text-gray-900 placeholder:text-gray-500 border border-gray-200 shadow-sm hover:shadow-md focus:bg-white focus:ring-2 focus:ring-brand-green-500 focus:border-transparent transition-all text-sm"
    />
    {searchQuery && (
      <button
        onClick={clearSearch}
        className="absolute inset-y-0 right-0 flex items-center pr-4 text-gray-400 hover:text-gray-600"
      >
        <XMarkIcon className="h-4 w-4" />
      </button>
    )}
    {showSuggestions && (
                <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-white rounded-xl shadow-xl border border-gray-200 max-h-80 overflow-y-auto">
        {/* Search Results */}
        {searchQuery.length > 0 && filteredSuggestions.length > 0 && (
          <div className="p-2">
            <div className="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">
              Search Results
            </div>
            {filteredSuggestions.map((suggestion, index) => (
              <button
                key={index}
                onClick={() => handleSuggestionClick(suggestion.text)}
                          className="w-full flex items-center gap-3 px-3 py-2 text-left hover:bg-brand-green-50 rounded-lg transition-colors"
              >
                <div className={cn(
                  'w-8 h-8 flex items-center justify-center rounded-lg',
                  suggestion.type === 'transaction' && 'bg-blue-100 text-blue-600',
                  suggestion.type === 'customer' && 'bg-green-100 text-green-600',
                  suggestion.type === 'filter' && 'bg-purple-100 text-purple-600'
                )}>
                  <MagnifyingGlassIcon className="h-4 w-4" />
                </div>
                <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium text-gray-900">
                    {suggestion.text}
                  </p>
                  <p className="text-xs text-gray-500 truncate">
                    {suggestion.description}
                  </p>
                </div>
              </button>
            ))}
          </div>
        )}

        {/* Recent Searches */}
        {searchQuery.length === 0 && recentSearches.length > 0 && (
          <div className="p-2">
            <div className="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider flex items-center gap-2">
              <ClockIcon className="h-3 w-3" />
              Recent Searches
            </div>
            {recentSearches.map((search, index) => (
              <button
                key={index}
                onClick={() => handleSuggestionClick(search)}
                          className="w-full flex items-center gap-3 px-3 py-2 text-left hover:bg-brand-green-50 rounded-lg transition-colors"
              >
                <ClockIcon className="h-4 w-4 text-gray-400" />
                          <span className="text-sm text-gray-700">
                  {search}
                </span>
              </button>
            ))}
          </div>
        )}

        {/* No Results */}
        {searchQuery.length > 0 && filteredSuggestions.length === 0 && (
          <div className="p-4 text-center">
            <p className="text-sm text-gray-500">No results found for "{searchQuery}"</p>
            <p className="text-xs text-gray-400 mt-1">
              Try searching for transactions, customers, or payment methods
            </p>
          </div>
        )}

        {/* Quick Actions */}
        <div className="border-t border-gray-100 p-2">
          <div className="px-3 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">
            Quick Actions
          </div>
          <button
            onClick={() => handleSuggestionClick('Advanced search')}
                      className="w-full flex items-center gap-3 px-3 py-2 text-left hover:bg-brand-green-50 rounded-lg transition-colors"
          >
            <MagnifyingGlassIcon className="h-4 w-4 text-brand-green-600" />
            <span className="text-sm font-medium text-brand-green-700">Advanced Search</span>
          </button>
        </div>
      </div>
    )}
            </div>
  </div>
</div>

        {/* Right controls */}
        <div className="flex items-center gap-x-2 sm:gap-x-4 md:gap-x-6 ml-4">
          {/* Language Selector */}
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
              <Menu.Items className="absolute right-0 z-40 mt-2 w-32 origin-top-right rounded-xl bg-white py-2 shadow-xl ring-1 ring-gray-900/5 focus:outline-none border border-gray-100">
                <Menu.Item>
                  {({ active }) => (
                    <button
                      className={cn(
                        active ? 'bg-brand-green-50 text-brand-green-700' : 'text-gray-900',
                        'flex w-full items-center gap-x-2 px-4 py-2 text-sm font-medium transition-colors duration-150'
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
                        'flex w-full items-center gap-x-2 px-4 py-2 text-sm font-medium transition-colors duration-150'
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

          {/* Notifications */}
          <button 
            type="button" 
            className="relative p-2 text-gray-400 hover:text-brand-green-600 bg-gray-50 hover:bg-brand-green-50 rounded-xl transition-all duration-200 hover:shadow-md group"
          >
            <span className="sr-only">View notifications</span>
            <BellIcon className="h-6 w-6 group-hover:animate-pulse" aria-hidden="true" />
            <span className="absolute top-1 right-1 flex h-3 w-3">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-3 w-3 bg-gradient-to-r from-red-500 to-red-600 shadow-lg"></span>
            </span>
          </button>

          {/* Profile dropdown */}
          <Menu as="div" className="relative">
            <Menu.Button className="flex items-center gap-x-3 p-2 hover:bg-brand-green-50 rounded-xl transition-all duration-200 hover:shadow-md group">
              <span className="sr-only">Open user menu</span>
              <div className="relative">
                <div className="h-8 w-8 rounded-xl bg-gradient-to-r from-brand-green-500 to-brand-green-600 flex items-center justify-center shadow-lg group-hover:shadow-xl transition-all duration-200 group-hover:scale-105">
                  <span className="text-white font-bold text-sm">
                    {user?.name?.charAt(0) || 'U'}
                  </span>
                </div>
                <div className="absolute inset-0 rounded-xl bg-gradient-to-r from-brand-green-400 to-brand-gold-400 opacity-0 group-hover:opacity-20 transition-opacity duration-200" />
              </div>
              <span className="hidden lg:flex lg:flex-col lg:items-start">
                <span className="text-sm font-bold text-gray-900 group-hover:text-brand-green-700 transition-colors duration-200">
                  {user?.name || 'User'}
                </span>
                <span className="text-xs text-gray-500 group-hover:text-brand-green-600 transition-colors duration-200">
                  Admin
                </span>
              </span>
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
              <Menu.Items className="absolute right-0 z-40 mt-2 w-56 origin-top-right rounded-xl bg-white py-2 shadow-xl ring-1 ring-gray-900/5 focus:outline-none border border-gray-100">
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
                        href="/dashboard/settings"
                        className={cn(
                          active ? 'bg-brand-green-50 text-brand-green-700' : 'text-gray-900',
                          'flex items-center gap-x-3 px-4 py-2 text-sm font-medium transition-colors duration-150'
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
                          'flex w-full items-center gap-x-3 px-4 py-2 text-sm font-medium transition-colors duration-150'
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
    </header>
  )
} 