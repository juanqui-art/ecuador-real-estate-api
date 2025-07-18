'use client';

import { useAuthStore } from '@/store/auth';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Loading } from '@/components/ui/loading';
import type { UserRole } from '@shared/types/auth';
import { canAccessRoles } from '@shared/types/auth';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: UserRole | UserRole[];
  fallback?: React.ReactNode;
}

export function ProtectedRoute({ 
  children, 
  requiredRole, 
  fallback = <Loading /> 
}: ProtectedRouteProps) {
  const { isAuthenticated, user } = useAuthStore();
  const router = useRouter();
  const [hasCheckedAuth, setHasCheckedAuth] = useState(false);

  useEffect(() => {
    // Pequeño delay para permitir que el store se hidrate
    const timer = setTimeout(() => {
      setHasCheckedAuth(true);
      
      if (!isAuthenticated) {
        const currentPath = window.location.pathname;
        const redirectParam = encodeURIComponent(currentPath);
        router.push(`/login?redirect=${redirectParam}`);
        return;
      }

      if (requiredRole) {
        const rolesArray = Array.isArray(requiredRole) ? requiredRole : [requiredRole];
        if (!canAccessRoles(user?.role || 'buyer', rolesArray)) {
          router.push('/unauthorized');
          return;
        }
      }
    }, 100);

    return () => clearTimeout(timer);
  }, [isAuthenticated, user, requiredRole, router]);

  // Mostrar loading mientras se verifica auth
  if (!hasCheckedAuth) {
    return fallback;
  }

  // Verificar autenticación
  if (!isAuthenticated) {
    return fallback;
  }

  // Verificar rol si es requerido
  if (requiredRole) {
    const rolesArray = Array.isArray(requiredRole) ? requiredRole : [requiredRole];
    if (!canAccessRoles(user?.role || 'buyer', rolesArray)) {
      return fallback;
    }
  }

  return <>{children}</>;
}