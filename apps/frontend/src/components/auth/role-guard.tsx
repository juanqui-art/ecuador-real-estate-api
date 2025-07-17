'use client';

import { useAuthStore } from '@/store/auth';
import { hasPermission, hasMinimumRole } from '@shared/types/auth';
import type { UserRole } from '@shared/types/auth';

interface RoleGuardProps {
  children: React.ReactNode;
  requiredRole?: UserRole;
  requiredPermission?: string;
  allowedRoles?: UserRole[];
  fallback?: React.ReactNode;
  showFallback?: boolean;
}

export function RoleGuard({
  children,
  requiredRole,
  requiredPermission,
  allowedRoles,
  fallback = null,
  showFallback = false
}: RoleGuardProps) {
  const { user } = useAuthStore();

  // User not authenticated
  if (!user) {
    return showFallback ? <>{fallback}</> : null;
  }

  // Check specific role requirement
  if (requiredRole && !hasMinimumRole(user.role, requiredRole)) {
    return showFallback ? <>{fallback}</> : null;
  }

  // Check allowed roles
  if (allowedRoles && !allowedRoles.includes(user.role)) {
    return showFallback ? <>{fallback}</> : null;
  }

  // Check specific permission
  if (requiredPermission && !hasPermission(user.role, requiredPermission as any)) {
    return showFallback ? <>{fallback}</> : null;
  }

  return <>{children}</>;
}

// Helper components for common permission checks
export function AdminOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function AgencyOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['agency']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function AgentOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['agent']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function OwnerOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['owner']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function BuyerOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['buyer']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

// Management permission guards
export function CanManageUsers({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin', 'agency']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function CanManageProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin', 'agency', 'agent', 'owner']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function CanViewAnalytics({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin', 'agency']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function CanManageAgencies({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

// Property-specific permissions
export function CanCreateProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin', 'agency', 'agent', 'owner']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

export function CanDeleteProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return (
    <RoleGuard allowedRoles={['admin', 'agency', 'owner']} fallback={fallback}>
      {children}
    </RoleGuard>
  );
}

// Advanced role guards with conditions
export function ConditionalRoleGuard({
  children,
  condition,
  fallback
}: {
  children: React.ReactNode;
  condition: (user: any) => boolean;
  fallback?: React.ReactNode;
}) {
  const { user } = useAuthStore();

  if (!user || !condition(user)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}

// Role-based styling helper
export function RoleBasedStyles({ 
  children, 
  roleStyles 
}: { 
  children: React.ReactNode; 
  roleStyles: Record<UserRole, string> 
}) {
  const { user } = useAuthStore();
  
  if (!user) return <>{children}</>;
  
  const className = roleStyles[user.role] || '';
  
  return (
    <div className={className}>
      {children}
    </div>
  );
}