import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '@shared/types/auth';

export interface AuthState {
  user: User | null;
  access_token: string | null;
  refresh_token: string | null;
  isAuthenticated: boolean;
  
  // Actions
  login: (tokens: { access_token: string; refresh_token: string }, user: User) => void;
  logout: () => void;
  updateUser: (user: Partial<User>) => void;
  setTokens: (tokens: { access_token: string; refresh_token?: string }) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      access_token: null,
      refresh_token: null,
      isAuthenticated: false,

      login: (tokens, user) => {
        console.log('🏪 Auth Store - Login called');
        
        set({
          user,
          access_token: tokens.access_token,
          refresh_token: tokens.refresh_token,
          isAuthenticated: true,
        });
        
        // Update localStorage for api interceptor
        localStorage.setItem('access_token', tokens.access_token);
        localStorage.setItem('refresh_token', tokens.refresh_token);
        
        // Also set cookies for server-side middleware
        document.cookie = `access_token=${tokens.access_token}; path=/; max-age=900`; // 15 min
        document.cookie = `refresh_token=${tokens.refresh_token}; path=/; max-age=604800`; // 7 days
        
        console.log('🏪 Auth state updated, isAuthenticated: true');
      },

      logout: () => {
        console.log('🏪 Auth Store - Logout called');
        set({
          user: null,
          access_token: null,
          refresh_token: null,
          isAuthenticated: false,
        });
        
        // Clear localStorage
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        
        // Clear cookies
        document.cookie = 'access_token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
        document.cookie = 'refresh_token=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
      },

      updateUser: (userData) => {
        const currentUser = get().user;
        if (currentUser) {
          set({ user: { ...currentUser, ...userData } });
        }
      },

      setTokens: (tokens) => {
        set((state) => ({
          access_token: tokens.access_token,
          refresh_token: tokens.refresh_token || state.refresh_token,
        }));
        
        localStorage.setItem('access_token', tokens.access_token);
        if (tokens.refresh_token) {
          localStorage.setItem('refresh_token', tokens.refresh_token);
        }
        
        // Also update cookies
        document.cookie = `access_token=${tokens.access_token}; path=/; max-age=900`; // 15 min
        if (tokens.refresh_token) {
          document.cookie = `refresh_token=${tokens.refresh_token}; path=/; max-age=604800`; // 7 days
        }
      },
    }),
    {
      name: 'auth-storage',
      // Store only essential auth data
      partialize: (state) => ({
        user: state.user,
        access_token: state.access_token,
        refresh_token: state.refresh_token,
        isAuthenticated: state.isAuthenticated,
      }),
      // Synchronize localStorage and cookies with store on hydration
      onRehydrateStorage: () => (state) => {
        if (state) {
          console.log('🏪 Store rehydrated with state:', {
            hasUser: !!state.user,
            hasAccessToken: !!state.access_token,
            isAuthenticated: state.isAuthenticated
          });
          
          // Sync localStorage with store state
          if (state.access_token) {
            localStorage.setItem('access_token', state.access_token);
            // Also sync cookies
            document.cookie = `access_token=${state.access_token}; path=/; max-age=900`; // 15 min
          }
          if (state.refresh_token) {
            localStorage.setItem('refresh_token', state.refresh_token);
            // Also sync cookies
            document.cookie = `refresh_token=${state.refresh_token}; path=/; max-age=604800`; // 7 days
          }
        }
      },
    }
  )
);