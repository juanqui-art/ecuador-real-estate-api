'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import { 
  Home, 
  DollarSign, 
  TrendingUp, 
  Eye, 
  MapPin, 
  Calendar,
  BarChart3,
  Activity
} from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { formatPrice } from '@/lib/utils';

interface PropertyStatistics {
  total_properties: number;
  average_price: number;
  by_status: {
    available?: number;
    sold?: number;
    rented?: number;
  };
  by_type: {
    house?: number;
    apartment?: number;
    land?: number;
    commercial?: number;
  };
  by_province: Record<string, number>;
}

export function PropertyStats() {
  const { data: stats, isLoading, error } = useQuery<PropertyStatistics>({
    queryKey: ['property-statistics'],
    queryFn: async () => {
      const response = await apiClient.get('/properties/statistics');
      return response.data.data; // Extraer la data anidada
    },
  });

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <Card key={i}>
            <CardHeader className="pb-3">
              <Skeleton className="h-4 w-16" />
              <Skeleton className="h-8 w-20" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-4 w-full" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (error || !stats) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-red-600">Error al cargar las estadísticas</p>
        </CardContent>
      </Card>
    );
  }

  const statCards = [
    {
      title: 'Total Propiedades',
      value: stats.total_properties,
      icon: Home,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      title: 'Disponibles',
      value: stats.by_status.available || 0,
      icon: Activity,
      color: 'text-green-600',
      bgColor: 'bg-green-50',
    },
    {
      title: 'Vendidas',
      value: stats.by_status.sold || 0,
      icon: TrendingUp,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
    },
    {
      title: 'Rentadas',
      value: stats.by_status.rented || 0,
      icon: Calendar,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
    {
      title: 'Precio Promedio',
      value: formatPrice(stats.average_price),
      icon: DollarSign,
      color: 'text-emerald-600',
      bgColor: 'bg-emerald-50',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Main Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat, index) => (
          <Card key={index} className="hover:shadow-lg transition-shadow">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-medium text-gray-600">
                  {stat.title}
                </CardTitle>
                <div className={`p-2 rounded-lg ${stat.bgColor}`}>
                  <stat.icon className={`h-4 w-4 ${stat.color}`} />
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-gray-900">
                {stat.value}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Secondary Stats */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Property Types */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Propiedades por Tipo</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">Casas</span>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{stats.by_type.house || 0}</Badge>
                  <div className="w-24 h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-blue-500 transition-all"
                      style={{ 
                        width: `${((stats.by_type.house || 0) / stats.total_properties) * 100}%` 
                      }}
                    />
                  </div>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">Departamentos</span>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{stats.by_type.apartment || 0}</Badge>
                  <div className="w-24 h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-green-500 transition-all"
                      style={{ 
                        width: `${((stats.by_type.apartment || 0) / stats.total_properties) * 100}%` 
                      }}
                    />
                  </div>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">Terrenos</span>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{stats.by_type.land || 0}</Badge>
                  <div className="w-24 h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-purple-500 transition-all"
                      style={{ 
                        width: `${((stats.by_type.land || 0) / stats.total_properties) * 100}%` 
                      }}
                    />
                  </div>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-gray-600">Comerciales</span>
                <div className="flex items-center gap-2">
                  <Badge variant="secondary">{stats.by_type.commercial || 0}</Badge>
                  <div className="w-24 h-2 bg-gray-200 rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-orange-500 transition-all"
                      style={{ 
                        width: `${((stats.by_type.commercial || 0) / stats.total_properties) * 100}%` 
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Summary Info */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Resumen General</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              <div className="text-center p-4 bg-blue-50 rounded-lg">
                <div className="text-2xl font-bold text-blue-600">
                  {stats.total_properties}
                </div>
                <div className="text-sm text-blue-600">Total Propiedades</div>
              </div>
              <div className="text-center p-4 bg-green-50 rounded-lg">
                <div className="text-2xl font-bold text-green-600">
                  {formatPrice(stats.average_price)}
                </div>
                <div className="text-sm text-green-600">Precio Promedio</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Top Provinces */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Propiedades por Provincia</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {Object.entries(stats.by_province)
              .sort(([,a], [,b]) => b - a)
              .slice(0, 6)
              .map(([province, count]) => (
                <div key={province} className="text-center p-3 bg-gray-50 rounded-lg">
                  <div className="text-lg font-bold text-gray-900">{count}</div>
                  <div className="text-sm text-gray-600">{province}</div>
                </div>
              ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}