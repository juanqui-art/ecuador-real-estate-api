'use client';

import { motion } from 'motion/react';
import { 
  Building, 
  TrendingUp, 
  Users, 
  DollarSign, 
  Home, 
  MapPin,
  Eye,
  Heart,
  Key,
  UserCheck
} from 'lucide-react';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useAuthStore, type UserRole } from '@/store/auth';
import { usePropertyStats } from '@/hooks/useProperties';

interface StatCardProps {
  title: string;
  value: string | number;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
  change?: {
    value: number;
    type: 'increase' | 'decrease';
  };
  trend?: 'up' | 'down' | 'stable';
  delay?: number;
}

function StatCard({ 
  title, 
  value, 
  description, 
  icon: Icon, 
  change, 
  trend = 'stable',
  delay = 0 
}: StatCardProps) {
  const getTrendColor = (trend: string) => {
    switch (trend) {
      case 'up': return 'text-green-600';
      case 'down': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up': return '↗️';
      case 'down': return '↘️';
      default: return '➡️';
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay }}
      whileHover={{ scale: 1.02 }}
      className="group"
    >
      <Card className="relative overflow-hidden transition-all duration-300 hover:shadow-lg">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="text-sm font-medium text-gray-600">
            {title}
          </CardTitle>
          <Icon className="h-4 w-4 text-gray-400 group-hover:text-blue-600 transition-colors" />
        </CardHeader>
        <CardContent>
          <div className="text-2xl font-bold text-gray-900 mb-1">
            {value}
          </div>
          <div className="flex items-center justify-between">
            <CardDescription className="text-xs">
              {description}
            </CardDescription>
            {change && (
              <Badge 
                variant="secondary" 
                className={`text-xs ${getTrendColor(change.type === 'increase' ? 'up' : 'down')}`}
              >
                {getTrendIcon(change.type === 'increase' ? 'up' : 'down')} {change.value}%
              </Badge>
            )}
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}

export function StatsCards() {
  const { user } = useAuthStore();
  const { data: stats, isLoading } = usePropertyStats();

  const getStatsForRole = (role: UserRole) => {
    switch (role) {
      case 'admin':
        return [
          {
            title: 'Total Propiedades',
            value: stats?.total_properties || '6',
            description: 'Propiedades en el sistema',
            icon: Building,
            change: { value: 12, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Disponibles',
            value: stats?.available_properties || '4',
            description: 'Propiedades disponibles',
            icon: Users,
            change: { value: 8, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Vendidas',
            value: stats?.sold_properties || '2',
            description: 'Propiedades vendidas',
            icon: Home,
            trend: 'stable' as const,
          },
          {
            title: 'Valor Total',
            value: stats?.total_value ? `$${stats.total_value.toLocaleString()}` : '$945,000',
            description: 'Valor total propiedades',
            icon: DollarSign,
            change: { value: 23, type: 'increase' as const },
            trend: 'up' as const,
          },
        ];

      case 'agency':
        return [
          {
            title: 'Mis Propiedades',
            value: '4',
            description: 'Propiedades de la agencia',
            icon: Building,
            change: { value: 15, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Agentes',
            value: '3',
            description: 'Agentes activos',
            icon: UserCheck,
            trend: 'stable' as const,
          },
          {
            title: 'Ventas Mes',
            value: '$285,000',
            description: 'Ventas del mes actual',
            icon: TrendingUp,
            change: { value: 18, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Comisiones',
            value: '$14,250',
            description: 'Comisiones generadas',
            icon: DollarSign,
            change: { value: 22, type: 'increase' as const },
            trend: 'up' as const,
          },
        ];

      case 'agent':
        return [
          {
            title: 'Asignadas',
            value: '2',
            description: 'Propiedades asignadas',
            icon: Key,
            trend: 'stable' as const,
          },
          {
            title: 'Visitadas',
            value: '12',
            description: 'Visitas este mes',
            icon: Eye,
            change: { value: 25, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Leads',
            value: '8',
            description: 'Contactos interesados',
            icon: Users,
            change: { value: 33, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Comisión',
            value: '$4,750',
            description: 'Comisión del mes',
            icon: DollarSign,
            change: { value: 12, type: 'increase' as const },
            trend: 'up' as const,
          },
        ];

      case 'owner':
        return [
          {
            title: 'Mis Propiedades',
            value: '1',
            description: 'Propiedades publicadas',
            icon: Home,
            trend: 'stable' as const,
          },
          {
            title: 'Vistas',
            value: '156',
            description: 'Vistas totales',
            icon: Eye,
            change: { value: 15, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Contactos',
            value: '7',
            description: 'Interesados este mes',
            icon: Users,
            change: { value: 40, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Valor',
            value: '$175,000',
            description: 'Valor de propiedades',
            icon: DollarSign,
            trend: 'stable' as const,
          },
        ];

      case 'buyer':
        return [
          {
            title: 'Favoritos',
            value: '3',
            description: 'Propiedades guardadas',
            icon: Heart,
            trend: 'stable' as const,
          },
          {
            title: 'Búsquedas',
            value: '12',
            description: 'Búsquedas guardadas',
            icon: MapPin,
            change: { value: 8, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Vistas',
            value: '45',
            description: 'Propiedades vistas',
            icon: Eye,
            change: { value: 20, type: 'increase' as const },
            trend: 'up' as const,
          },
          {
            title: 'Presupuesto',
            value: '$200,000',
            description: 'Presupuesto máximo',
            icon: DollarSign,
            trend: 'stable' as const,
          },
        ];

      default:
        return [];
    }
  };

  const statsData = getStatsForRole(user?.role || 'buyer');

  // Show loading skeleton
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="animate-pulse">
            <div className="bg-white rounded-lg p-6 shadow">
              <div className="flex items-center justify-between mb-2">
                <div className="h-4 bg-gray-200 rounded w-24"></div>
                <div className="h-4 w-4 bg-gray-200 rounded"></div>
              </div>
              <div className="h-8 bg-gray-200 rounded w-16 mb-2"></div>
              <div className="h-3 bg-gray-200 rounded w-32"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      {statsData.map((stat, index) => (
        <StatCard
          key={stat.title}
          title={stat.title}
          value={stat.value}
          description={stat.description}
          icon={stat.icon}
          change={stat.change}
          trend={stat.trend}
          delay={index * 0.1}
        />
      ))}
    </div>
  );
}