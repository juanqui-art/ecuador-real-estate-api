'use client';

import { Button } from './button';
import { RoleGuard } from '@/components/auth/role-guard';
import { useAuthStore } from '@/store/auth';
import type { UserRole } from '@shared/types/auth';

interface RoleBasedButtonProps {
  children: React.ReactNode;
  allowedRoles?: UserRole[];
  requiredRole?: UserRole;
  requiredPermission?: string;
  fallback?: React.ReactNode;
  className?: string;
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
  size?: 'default' | 'sm' | 'lg' | 'icon';
  disabled?: boolean;
  onClick?: () => void;
  type?: 'button' | 'submit' | 'reset';
  // Role-specific styling
  roleStyles?: Partial<Record<UserRole, string>>;
}

export function RoleBasedButton({
  children,
  allowedRoles,
  requiredRole,
  requiredPermission,
  fallback,
  className = '',
  variant = 'default',
  size = 'default',
  disabled = false,
  onClick,
  type = 'button',
  roleStyles = {},
}: RoleBasedButtonProps) {
  const { user } = useAuthStore();
  
  // Apply role-specific styling
  const roleSpecificClass = user?.role && roleStyles[user.role] ? roleStyles[user.role] : '';
  const finalClassName = `${className} ${roleSpecificClass}`.trim();

  return (
    <RoleGuard
      allowedRoles={allowedRoles}
      requiredRole={requiredRole}
      requiredPermission={requiredPermission}
      fallback={fallback}
    >
      <Button
        variant={variant}
        size={size}
        disabled={disabled}
        onClick={onClick}
        type={type}
        className={finalClassName}
      >
        {children}
      </Button>
    </RoleGuard>
  );
}

// Convenience components for common role-based buttons
export function AdminButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function AgencyButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['agency']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function AgentButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['agent']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function OwnerButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['owner']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function BuyerButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['buyer']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

// Permission-based buttons
export function ManageUsersButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin', 'agency']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function ManagePropertiesButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin', 'agency', 'agent', 'owner']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function ViewAnalyticsButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin', 'agency']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function CreatePropertyButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin', 'agency', 'agent', 'owner']} {...props}>
      {children}
    </RoleBasedButton>
  );
}

export function DeletePropertyButton({ children, ...props }: Omit<RoleBasedButtonProps, 'allowedRoles'>) {
  return (
    <RoleBasedButton allowedRoles={['admin', 'agency', 'owner']} {...props}>
      {children}
    </RoleBasedButton>
  );
}