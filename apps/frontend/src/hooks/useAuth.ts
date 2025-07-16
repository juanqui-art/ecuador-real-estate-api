import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { apiClient, type ApiError } from '@/lib/api-client';
import { useAuthStore } from '@/store/auth';
import type { 
  LoginRequest,
  LoginResponse, 
  RefreshTokenResponse,
  TokenValidationResponse,
  ChangePasswordRequest,
  User
} from '@shared/types/auth';

// API functions using new fetch-based client
const authApi = {
  // Login - matches Go backend response format
  login: async (credentials: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<LoginResponse>('/auth/login', credentials);
    return response.data;
  },

  // Refresh token
  refreshToken: async (refreshToken: string): Promise<RefreshTokenResponse> => {
    const response = await apiClient.post<RefreshTokenResponse>('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
  },

  // Logout
  logout: async (): Promise<void> => {
    await apiClient.post<void>('/auth/logout');
  },

  // Validate token
  validateToken: async (): Promise<TokenValidationResponse> => {
    const response = await apiClient.get<TokenValidationResponse>('/auth/validate');
    return response.data;
  },

  // Change password
  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await apiClient.post<void>('/auth/change-password', data);
  },

  // Get current user profile
  getCurrentUser: async (): Promise<User> => {
    const response = await apiClient.get<User>('/auth/me');
    return response.data;
  },
};

// React Query hooks
export const useLogin = () => {
  const router = useRouter();
  const { login } = useAuthStore();
  
  return useMutation({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      console.log('🔑 Login successful, updating auth state');
      
      // Update auth store with correct structure from Go backend
      login(
        { 
          access_token: data.tokens.access_token, 
          refresh_token: data.tokens.refresh_token 
        },
        data.user
      );
      
      // Simple redirect - let the page logic handle the actual navigation
      console.log('🔑 Auth state updated, navigation will be handled by useEffect');
    },
    onError: (error) => {
      console.error('❌ Login failed:', error);
    },
  });
};

export const useLogout = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  
  return useMutation({
    mutationFn: async () => {
      console.log('🔓 Starting logout process...');
      
      // Try to logout from backend first, but don't let it block the UI
      try {
        await authApi.logout();
        console.log('✅ Backend logout successful');
      } catch (error: any) {
        // If logout fails with 401, it means the token is already invalid
        // This is actually a successful logout from the user's perspective
        if (error?.status === 401) {
          console.log('✅ Backend logout: Token already invalid (expected)');
        } else {
          console.log('⚠️ Backend logout failed, but proceeding with client cleanup:', error.message);
        }
        // Don't throw - we want to continue with cleanup regardless
      }
    },
    onSuccess: () => {
      console.log('🧹 Performing client-side cleanup...');
      
      // Clear auth store and tokens
      logout();
      
      // Clear all queries
      queryClient.clear();
      
      console.log('✅ Logout complete - redirecting to login');
      
      // Redirect to login
      router.push('/login');
    },
    onError: (error: ApiError) => {
      // This should rarely happen now since we handle errors in mutationFn
      console.log('⚠️ Logout mutation failed, but still cleaning up:', error.message);
      
      // Always logout on error - if backend fails, we still want to clear client state
      logout();
      queryClient.clear();
      router.push('/login');
    },
  });
};

export const useValidateToken = () => {
  const { user, isAuthenticated } = useAuthStore();
  
  return useQuery({
    queryKey: ['validate-token'],
    queryFn: authApi.validateToken,
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: false,
    refetchOnWindowFocus: true,
    refetchInterval: 10 * 60 * 1000, // 10 minutes
  });
};

export const useChangePassword = () => {
  return useMutation({
    mutationFn: authApi.changePassword,
    onSuccess: () => {
      // Optionally show success message
      console.log('Password changed successfully');
    },
    onError: (error) => {
      console.error('Password change failed:', error);
    },
  });
};

export const useCurrentUser = () => {
  const { isAuthenticated } = useAuthStore();
  
  return useQuery({
    queryKey: ['current-user'],
    queryFn: authApi.getCurrentUser,
    enabled: isAuthenticated,
    staleTime: 15 * 60 * 1000, // 15 minutes
    retry: false,
  });
};

// Utility hook for automatic token refresh
export const useTokenRefresh = () => {
  const { refresh_token, setTokens, logout } = useAuthStore();
  
  return useMutation({
    mutationFn: () => {
      if (!refresh_token) {
        throw new Error('No refresh token available');
      }
      return authApi.refreshToken(refresh_token);
    },
    onSuccess: (data) => {
      setTokens({
        access_token: data.access_token,
        refresh_token: data.refresh_token,
      });
    },
    onError: (error) => {
      console.error('Token refresh failed:', error);
      logout();
    },
  });
};

// Auth guard hook
export const useAuthGuard = (requiredRole?: string[]) => {
  const { user, isAuthenticated } = useAuthStore();
  const router = useRouter();
  
  const hasRequiredRole = requiredRole 
    ? requiredRole.includes(user?.role || '')
    : true;
  
  const isAuthorized = isAuthenticated && hasRequiredRole;
  
  // Redirect if not authenticated
  if (!isAuthenticated) {
    router.push('/login');
    return { isAuthorized: false, isLoading: false };
  }
  
  // Redirect if doesn't have required role
  if (!hasRequiredRole) {
    router.push('/unauthorized');
    return { isAuthorized: false, isLoading: false };
  }
  
  return { isAuthorized, isLoading: false };
};

// Role checking utilities
export const useRoleCheck = () => {
  const { user } = useAuthStore();
  
  const hasRole = (role: string) => user?.role === role;
  const hasAnyRole = (roles: string[]) => roles.includes(user?.role || '');
  const hasMinimumRole = (minimumRole: string) => {
    const roleHierarchy = ['buyer', 'owner', 'agent', 'agency', 'admin'];
    const userRoleIndex = roleHierarchy.indexOf(user?.role || '');
    const minRoleIndex = roleHierarchy.indexOf(minimumRole);
    return userRoleIndex >= minRoleIndex;
  };
  
  return {
    hasRole,
    hasAnyRole,
    hasMinimumRole,
    isAdmin: hasRole('admin'),
    isAgency: hasRole('agency'),
    isAgent: hasRole('agent'),
    isOwner: hasRole('owner'),
    isBuyer: hasRole('buyer'),
    canCreateProperties: hasAnyRole(['admin', 'agency', 'agent', 'owner']),
    canManageUsers: hasAnyRole(['admin', 'agency']),
    canManageAgencies: hasRole('admin'),
    canViewStats: hasAnyRole(['admin', 'agency']),
  };
};