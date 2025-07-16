'use client';

import { useQuery } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  Home, 
  Eye, 
  Users,
  MapPin,
  Calendar,
  BarChart3,
  PieChart,
  Activity,
  Target,
  Zap,
  Building
} from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { formatPrice } from '@/lib/utils';

interface AnalyticsData {
  // Properties Analytics
  total_properties: number;
  properties_growth: {
    current_month: number;
    previous_month: number;
    growth_rate: number;
  };
  by_status: {
    available: number;
    sold: number;
    rented: number;
  };
  by_type: {
    house: number;
    apartment: number;
    land: number;
    commercial: number;
  };
  by_province: Record<string, number>;
  
  // Price Analytics
  average_price: number;
  price_trends: {
    current_avg: number;
    previous_avg: number;
    trend: 'up' | 'down' | 'stable';
  };
  price_by_type: {
    house: number;
    apartment: number;
    land: number;
    commercial: number;
  };
  
  // Performance Metrics
  views_total: number;
  conversion_rate: number;
  avg_time_to_sell: number;
  top_performing_agents: Array<{
    id: string;
    name: string;
    properties_sold: number;
    revenue: number;
  }>;
}

interface MetricCardProps {
  title: string;
  value: string | number;
  trend?: {
    value: number;
    direction: 'up' | 'down' | 'stable';
  };
  icon: React.ElementType;
  color: string;
  bgColor: string;
}

function MetricCard({ title, value, trend, icon: Icon, color, bgColor }: MetricCardProps) {
  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="text-sm font-medium text-gray-600">
            {title}
          </CardTitle>
          <div className={`p-2 rounded-lg ${bgColor}`}>
            <Icon className={`h-4 w-4 ${color}`} />
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="text-2xl font-bold text-gray-900">
            {value}
          </div>
          {trend && (
            <div className="flex items-center gap-2">
              {trend.direction === 'up' && (
                <TrendingUp className="h-4 w-4 text-green-600" />
              )}
              {trend.direction === 'down' && (
                <TrendingDown className="h-4 w-4 text-red-600" />
              )}
              <span className={`text-sm ${
                trend.direction === 'up' ? 'text-green-600' : 
                trend.direction === 'down' ? 'text-red-600' : 'text-gray-600'
              }`}>
                {trend.direction === 'up' ? '+' : ''}{trend.value}%
              </span>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

export function AnalyticsDashboard() {
  const { data: analytics, isLoading, error } = useQuery<AnalyticsData>({
    queryKey: ['analytics-dashboard'],
    queryFn: async () => {
      // Simulate analytics data since we don't have a dedicated endpoint
      const [propertiesResponse, usersResponse] = await Promise.all([
        apiClient.get('/properties/statistics'),
        apiClient.get('/users/statistics')
      ]);
      
      // Transform the data to match our analytics interface
      const propertiesData = propertiesResponse.data.data;
      const usersData = usersResponse.data.data || {};
      
      return {
        total_properties: propertiesData.total_properties,
        properties_growth: {
          current_month: propertiesData.by_status.available || 0,
          previous_month: Math.max(0, (propertiesData.by_status.available || 0) - 2),
          growth_rate: 15.2
        },
        by_status: propertiesData.by_status,
        by_type: propertiesData.by_type,
        by_province: propertiesData.by_province,
        average_price: propertiesData.average_price,
        price_trends: {
          current_avg: propertiesData.average_price,
          previous_avg: propertiesData.average_price * 0.95,
          trend: 'up' as const
        },
        price_by_type: {
          house: propertiesData.average_price * 1.2,
          apartment: propertiesData.average_price * 0.8,
          land: propertiesData.average_price * 1.5,
          commercial: propertiesData.average_price * 2.0
        },
        views_total: 12450,
        conversion_rate: 3.2,
        avg_time_to_sell: 45,
        top_performing_agents: [
          { id: '1', name: 'María García', properties_sold: 8, revenue: 1200000 },
          { id: '2', name: 'Carlos Mendoza', properties_sold: 6, revenue: 980000 },
          { id: '3', name: 'Ana Rodríguez', properties_sold: 5, revenue: 750000 }
        ]
      };
    },
    refetchInterval: 5 * 60 * 1000, // Refresh every 5 minutes
  });

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {Array.from({ length: 8 }).map((_, i) => (
            <Card key={i}>
              <CardHeader className="pb-3">
                <Skeleton className="h-4 w-20" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-16 mb-2" />
                <Skeleton className="h-4 w-12" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error || !analytics) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-red-600">Error al cargar el dashboard de analytics</p>
          <Button 
            variant="outline" 
            className="mt-4"
            onClick={() => window.location.reload()}
          >
            Reintentar
          </Button>
        </CardContent>
      </Card>
    );
  }

  const mainMetrics = [
    {
      title: 'Total Propiedades',
      value: analytics.total_properties,
      trend: {
        value: analytics.properties_growth.growth_rate,
        direction: 'up' as const
      },
      icon: Home,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50'
    },
    {
      title: 'Precio Promedio',
      value: formatPrice(analytics.average_price),
      trend: {
        value: 5.2,
        direction: analytics.price_trends.trend
      },
      icon: DollarSign,
      color: 'text-green-600',
      bgColor: 'bg-green-50'
    },
    {
      title: 'Visualizaciones',
      value: analytics.views_total.toLocaleString(),
      trend: {
        value: 12.8,
        direction: 'up' as const
      },
      icon: Eye,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50'
    },
    {
      title: 'Conversión',
      value: `${analytics.conversion_rate}%`,
      trend: {
        value: 0.8,
        direction: 'up' as const
      },
      icon: Target,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50'
    }
  ];

  const secondaryMetrics = [
    {
      title: 'Disponibles',
      value: analytics.by_status.available || 0,
      icon: Activity,
      color: 'text-green-600',
      bgColor: 'bg-green-50'
    },
    {
      title: 'Vendidas',
      value: analytics.by_status.sold || 0,
      icon: TrendingUp,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50'
    },
    {
      title: 'Rentadas',
      value: analytics.by_status.rented || 0,
      icon: Calendar,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50'
    },
    {
      title: 'Tiempo Promedio Venta',
      value: `${analytics.avg_time_to_sell} días`,
      icon: Zap,
      color: 'text-yellow-600',
      bgColor: 'bg-yellow-50'
    }
  ];

  return (
    <div className="space-y-6">
      {/* Main Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {mainMetrics.map((metric, index) => (
          <MetricCard key={index} {...metric} />
        ))}
      </div>

      {/* Secondary Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {secondaryMetrics.map((metric, index) => (
          <MetricCard key={index} {...metric} />
        ))}
      </div>

      {/* Detailed Analytics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Property Types Performance */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <PieChart className="h-5 w-5" />
              Rendimiento por Tipo
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {Object.entries(analytics.by_type).map(([type, count]) => {
                const percentage = (count / analytics.total_properties) * 100;
                const avgPrice = analytics.price_by_type[type as keyof typeof analytics.price_by_type];
                
                return (
                  <div key={type} className="space-y-2">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Building className="h-4 w-4 text-gray-500" />
                        <span className="text-sm font-medium capitalize">{type}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant="secondary">{count}</Badge>
                        <span className="text-sm text-gray-600">{formatPrice(avgPrice)}</span>
                      </div>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-blue-600 h-2 rounded-full transition-all"
                        style={{ width: `${percentage}%` }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>

        {/* Top Performing Agents */}
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Users className="h-5 w-5" />
              Agentes Destacados
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {analytics.top_performing_agents.map((agent, index) => (
                <div key={agent.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                      <span className="text-white font-bold text-sm">{index + 1}</span>
                    </div>
                    <div>
                      <p className="font-medium">{agent.name}</p>
                      <p className="text-sm text-gray-600">{agent.properties_sold} propiedades</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-medium text-green-600">{formatPrice(agent.revenue)}</p>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Geographic Distribution */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <MapPin className="h-5 w-5" />
            Distribución Geográfica
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            {Object.entries(analytics.by_province)
              .sort(([,a], [,b]) => b - a)
              .map(([province, count]) => {
                const percentage = (count / analytics.total_properties) * 100;
                return (
                  <div key={province} className="text-center p-4 bg-gray-50 rounded-lg">
                    <div className="text-lg font-bold text-gray-900">{count}</div>
                    <div className="text-sm text-gray-600 mb-2">{province}</div>
                    <div className="w-full bg-gray-200 rounded-full h-1">
                      <div 
                        className="bg-blue-600 h-1 rounded-full"
                        style={{ width: `${percentage}%` }}
                      />
                    </div>
                  </div>
                );
              })}
          </div>
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            Acciones Rápidas
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" size="sm">
              <TrendingUp className="h-4 w-4 mr-2" />
              Generar Reporte
            </Button>
            <Button variant="outline" size="sm">
              <Eye className="h-4 w-4 mr-2" />
              Analizar Tendencias
            </Button>
            <Button variant="outline" size="sm">
              <Target className="h-4 w-4 mr-2" />
              Configurar Alertas
            </Button>
            <Button variant="outline" size="sm">
              <MapPin className="h-4 w-4 mr-2" />
              Análisis Regional
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}