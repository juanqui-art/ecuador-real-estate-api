// Tipos de autenticación sincronizados con el backend Go
// Este archivo debe mantenerse en sync con internal/domain/user.go

export type UserRole = 'admin' | 'agency' | 'agent' | 'seller' | 'buyer';
export type UserStatus = 'active' | 'inactive' | 'suspended' | 'pending';

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: UserRole;
  agency_id?: string;
  status: UserStatus;
  phone?: string;
  profile_image?: string;
  created_at: string;
  updated_at: string;
}

// Estructura de respuesta del backend Go para login
export interface LoginResponse {
  user: User;
  tokens: {
    access_token: string;
    refresh_token: string;
  };
  expires_at: string;
  message: string;
}

// Estructura para refresh token
export interface RefreshTokenResponse {
  access_token: string;
  refresh_token?: string;
  expires_at: string;
}

// Estructura para validación de token
export interface TokenValidationResponse {
  valid: boolean;
  user?: User;
  expires_at: string;
}

// Requests
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RefreshTokenRequest {
  refresh_token: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

// Utilidades para roles
export const ROLE_HIERARCHY: UserRole[] = ['buyer', 'seller', 'agent', 'agency', 'admin'];

export const ROLE_PERMISSIONS = {
  admin: {
    canManageUsers: true,
    canManageAgencies: true,
    canManageProperties: true,
    canViewAnalytics: true,
    canViewAllProperties: true,
  },
  agency: {
    canManageUsers: true, // Only within their agency
    canManageAgencies: false,
    canManageProperties: true, // Only their agency's properties
    canViewAnalytics: true,
    canViewAllProperties: false,
  },
  agent: {
    canManageUsers: false,
    canManageAgencies: false,
    canManageProperties: true, // Only assigned properties
    canViewAnalytics: false,
    canViewAllProperties: false,
  },
  seller: {
    canManageUsers: false,
    canManageAgencies: false,
    canManageProperties: true, // Only their own properties
    canViewAnalytics: false,
    canViewAllProperties: false,
  },
  buyer: {
    canManageUsers: false,
    canManageAgencies: false,
    canManageProperties: false, // Read-only access
    canViewAnalytics: false,
    canViewAllProperties: false,
  },
} as const;

// Helper functions
export const hasMinimumRole = (userRole: UserRole, minimumRole: UserRole): boolean => {
  const userRoleIndex = ROLE_HIERARCHY.indexOf(userRole);
  const minRoleIndex = ROLE_HIERARCHY.indexOf(minimumRole);
  return userRoleIndex >= minRoleIndex;
};

export const canAccessRole = (userRole: UserRole, requiredRole: UserRole): boolean => {
  return hasMinimumRole(userRole, requiredRole);
};

export const canAccessRoles = (userRole: UserRole, requiredRoles: UserRole[]): boolean => {
  return requiredRoles.some(role => canAccessRole(userRole, role));
};

export const hasPermission = (userRole: UserRole, permission: keyof typeof ROLE_PERMISSIONS.admin): boolean => {
  return ROLE_PERMISSIONS[userRole]?.[permission] || false;
};

export const getRoleDisplayName = (role: UserRole): string => {
  const displayNames = {
    admin: 'Administrador',
    agency: 'Agencia',
    agent: 'Agente',
    seller: 'Propietario',
    buyer: 'Comprador',
  };
  return displayNames[role];
};