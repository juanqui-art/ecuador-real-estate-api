import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { useAuthStore } from '@/store/auth';

// Types
interface LoginCredentials {
  email: string;
  password: string;
}

interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    first_name: string;
    last_name: string;
    email: string;
    role: 'admin' | 'agency' | 'agent' | 'owner' | 'buyer';
    agency_id?: string;
    status: 'active' | 'inactive' | 'suspended' | 'pending';
  };
}

interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

interface TokenValidationResponse {
  valid: boolean;
  user?: AuthResponse['user'];
  expires_at: string;
}

// API functions
const authApi = {
  // Login
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', credentials);
    return response.data;
  },

  // Refresh token
  refreshToken: async (refreshToken: string): Promise<{ access_token: string; refresh_token?: string }> => {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken });
    return response.data;
  },

  // Logout
  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },

  // Validate token
  validateToken: async (): Promise<TokenValidationResponse> => {
    const response = await api.get('/auth/validate');
    return response.data;
  },

  // Change password
  changePassword: async (data: ChangePasswordRequest): Promise<void> => {
    await api.post('/auth/change-password', data);
  },

  // Get current user profile
  getCurrentUser: async (): Promise<AuthResponse['user']> => {
    const response = await api.get('/auth/me');
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
      // Update auth store
      login(
        { 
          access_token: data.access_token, 
          refresh_token: data.refresh_token 
        },
        data.user
      );
      
      // Redirect to dashboard
      router.push('/dashboard');
    },
    onError: (error) => {
      console.error('Login failed:', error);
    },
  });
};

export const useLogout = () => {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { logout } = useAuthStore();
  
  return useMutation({
    mutationFn: authApi.logout,
    onSuccess: () => {
      // Clear auth store
      logout();
      
      // Clear all queries
      queryClient.clear();
      
      // Redirect to login
      router.push('/login');
    },
    onError: (error) => {
      console.error('Logout failed:', error);
      // Still logout on error
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