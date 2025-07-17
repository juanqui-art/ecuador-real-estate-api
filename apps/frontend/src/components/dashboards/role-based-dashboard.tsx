'use client';

import { useAuthStore } from '@/store/auth';
import { AdminDashboard } from './admin-dashboard';
import { AgencyDashboard } from './agency-dashboard';
import { AgentDashboard } from './agent-dashboard';
import { OwnerDashboard } from './owner-dashboard';
import { BuyerDashboard } from './buyer-dashboard';
import { Card, CardContent } from '@/components/ui/card';
import { AlertCircle, Shield } from 'lucide-react';
import type { UserRole } from '@shared/types/auth';

// Default dashboard for unauthorized users
function UnauthorizedDashboard() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <CardContent className="pt-6">
          <div className="flex flex-col items-center space-y-4">
            <Shield className="h-12 w-12 text-gray-400" />
            <div className="text-center">
              <h2 className="text-lg font-semibold text-gray-900">Acceso Requerido</h2>
              <p className="text-sm text-gray-600">
                Necesitas iniciar sesión para acceder al dashboard
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// Error dashboard for invalid roles
function ErrorDashboard({ role }: { role: UserRole }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md border-red-200">
        <CardContent className="pt-6">
          <div className="flex flex-col items-center space-y-4">
            <AlertCircle className="h-12 w-12 text-red-500" />
            <div className="text-center">
              <h2 className="text-lg font-semibold text-red-900">Error de Configuración</h2>
              <p className="text-sm text-red-600">
                El rol "{role}" no tiene un dashboard configurado
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export function RoleBasedDashboard() {
  const { user, isAuthenticated } = useAuthStore();

  // User not authenticated
  if (!isAuthenticated || !user) {
    return <UnauthorizedDashboard />;
  }

  // Render appropriate dashboard based on user role
  switch (user.role) {
    case 'admin':
      return <AdminDashboard />;
    
    case 'agency':
      return <AgencyDashboard />;
    
    case 'agent':
      return <AgentDashboard />;
    
    case 'owner':
      return <OwnerDashboard />;
    
    case 'buyer':
      return <BuyerDashboard />;
    
    default:
      return <ErrorDashboard role={user.role} />;
  }
}

// Helper component for role-based access control
export function RoleBasedComponent({ 
  allowedRoles, 
  userRole, 
  children, 
  fallback = null 
}: {
  allowedRoles: UserRole[];
  userRole: UserRole;
  children: React.ReactNode;
  fallback?: React.ReactNode;
}) {
  if (!allowedRoles.includes(userRole)) {
    return <>{fallback}</>;
  }
  
  return <>{children}</>;
}

// Hook for checking role permissions
export function useRolePermissions() {
  const { user } = useAuthStore();
  
  const hasRole = (role: UserRole): boolean => {
    return user?.role === role;
  };
  
  const hasAnyRole = (roles: UserRole[]): boolean => {
    return user?.role ? roles.includes(user.role) : false;
  };
  
  const isAdmin = (): boolean => hasRole('admin');
  const isAgency = (): boolean => hasRole('agency');
  const isAgent = (): boolean => hasRole('agent');
  const isOwner = (): boolean => hasRole('owner');
  const isBuyer = (): boolean => hasRole('buyer');
  
  const canManageUsers = (): boolean => {
    return hasAnyRole(['admin', 'agency']);
  };
  
  const canManageProperties = (): boolean => {
    return hasAnyRole(['admin', 'agency', 'agent', 'owner']);
  };
  
  const canViewAnalytics = (): boolean => {
    return hasAnyRole(['admin', 'agency']);
  };
  
  const canManageAgencies = (): boolean => {
    return hasRole('admin');
  };
  
  return {
    userRole: user?.role,
    hasRole,
    hasAnyRole,
    isAdmin,
    isAgency,
    isAgent,
    isOwner,
    isBuyer,
    canManageUsers,
    canManageProperties,
    canViewAnalytics,
    canManageAgencies,
  };
}