import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { User } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isHydrated: boolean
  login: (user: User, token: string) => Promise<void>
  logout: () => Promise<void>
  updateUser: (user: Partial<User>) => void
  setHydrated: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isHydrated: false,
      login: async (user, token) => {
        return new Promise<void>((resolve) => {
          if (typeof window !== 'undefined') {
            localStorage.setItem('authToken', token)
          }
          set({ user, token, isAuthenticated: true })
          resolve()
        })
      },
      logout: async () => {
        return new Promise<void>((resolve) => {
          if (typeof window !== 'undefined') {
            localStorage.removeItem('authToken')
            localStorage.removeItem('refreshToken')
          }
          set({ user: null, token: null, isAuthenticated: false })
          resolve()
        })
      },
      updateUser: (userData) =>
        set((state) => ({
          user: state.user ? { ...state.user, ...userData } : null,
        })),
      setHydrated: () => set({ isHydrated: true }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        state?.setHydrated()
      },
    }
  )
) 