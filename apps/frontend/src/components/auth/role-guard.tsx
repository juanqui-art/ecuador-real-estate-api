'use client';

// ================ AUTHENTICATION DISABLED FOR DEVELOPMENT ================
// import { useAuthStore } from '@/store/auth';
// import { hasPermission, hasMinimumRole } from '@shared/types/auth';
// import type { UserRole } from '@shared/types/auth';

interface RoleGuardProps {
  children: React.ReactNode;
  requiredRole?: any; // UserRole;
  requiredPermission?: string;
  allowedRoles?: any[]; // UserRole[];
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
  // NO AUTH MODE: Always render children
  return <>{children}</>;

  /* ORIGINAL AUTH CODE - COMMENTED FOR DEVELOPMENT
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
  */
}

// Helper components for common permission checks (NO AUTH MODE)
export function AdminOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function AgencyOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function AgentOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function OwnerOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function BuyerOnly({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

// Management permission guards (NO AUTH MODE)
export function CanManageUsers({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function CanManageProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function CanViewAnalytics({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function CanManageAgencies({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

// Property-specific permissions (NO AUTH MODE)
export function CanCreateProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

export function CanDeleteProperties({ children, fallback }: { children: React.ReactNode; fallback?: React.ReactNode }) {
  return <>{children}</>;
}

// Advanced role guards with conditions (NO AUTH MODE)
export function ConditionalRoleGuard({
  children,
  condition,
  fallback
}: {
  children: React.ReactNode;
  condition: (user: any) => boolean;
  fallback?: React.ReactNode;
}) {
  return <>{children}</>;
}

// Role-based styling helper (NO AUTH MODE)
export function RoleBasedStyles({ 
  children, 
  roleStyles 
}: { 
  children: React.ReactNode; 
  roleStyles: Record<any, string> 
}) {
  return <>{children}</>;
}