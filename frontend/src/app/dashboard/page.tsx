'use client';

import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { StatsCards } from '@/components/dashboard/stats-cards';
import { PropertyGrid } from '@/components/dashboard/property-grid';
import { PropertyFilters } from '@/components/dashboard/property-filters';
import { useAuthStore } from '@/store/auth';

export default function DashboardPage() {
  const { user } = useAuthStore();

  const getRoleWelcome = (role: string) => {
    switch (role) {
      case 'admin': return 'Panel de Administración';
      case 'agency': return 'Panel de Agencia';
      case 'agent': return 'Panel de Agente';
      case 'owner': return 'Panel de Propietario';
      case 'buyer': return 'Panel de Comprador';
      default: return 'Dashboard';
    }
  };

  const handleFiltersChange = (filters: any) => {
    console.log('Filters changed:', filters);
    // Here you would typically update a state or make an API call
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            {getRoleWelcome(user?.role || 'buyer')}
          </h1>
          <p className="text-gray-600">
            Bienvenido {user?.first_name} {user?.last_name} - {user?.role}
          </p>
        </div>

        {/* Stats Cards */}
        <StatsCards />

        {/* Property Filters */}
        <PropertyFilters 
          onFiltersChange={handleFiltersChange}
          compact={true}
        />

        {/* Property Grid */}
        <PropertyGrid 
          showFilters={false}
          limit={6}
        />
      </div>
    </DashboardLayout>
  );
}