'use client';

import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { RoleBasedDashboard } from '@/components/dashboards/role-based-dashboard';
import { ProtectedRoute } from '@/components/auth/protected-route';

export default function SimpleDashboardPage() {
  return (
    <ProtectedRoute requiredRole="buyer">
      <DashboardLayout>
        <RoleBasedDashboard />
      </DashboardLayout>
    </ProtectedRoute>
  );
}