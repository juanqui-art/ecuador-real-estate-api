'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Building, 
  TrendingUp, 
  DollarSign, 
  Eye,
  MessageSquare,
  Calendar,
  Users,
  Star,
  Plus,
  Edit,
  Camera,
  AlertCircle,
  CheckCircle,
  Clock,
  MapPin,
  Phone,
  Mail,
  Target,
  BarChart3
} from 'lucide-react';

export function OwnerDashboard() {
  // Mock data - in real app, this would come from API
  const ownerStats = {
    totalProperties: 5,
    activeListings: 3,
    totalViews: 1245,
    totalInquiries: 28,
    averagePrice: 195000,
    bestPerforming: 'Casa en Los Ceibos',
    avgDaysOnMarket: 35,
    priceChanges: 2,
  };

  const myProperties = [
    { 
      id: 1, 
      title: 'Casa en Los Ceibos', 
      price: 325000, 
      status: 'active', 
      views: 456, 
      inquiries: 12,
      daysOnMarket: 25,
      lastPriceUpdate: '2024-01-10',
      needsAction: false
    },
    { 
      id: 2, 
      title: 'Departamento en Urdesa', 
      price: 180000, 
      status: 'active', 
      views: 289, 
      inquiries: 8,
      daysOnMarket: 42,
      lastPriceUpdate: '2024-01-05',
      needsAction: true
    },
    { 
      id: 3, 
      title: 'Oficina en el Centro', 
      price: 120000, 
      status: 'negotiating', 
      views: 167, 
      inquiries: 15,
      daysOnMarket: 18,
      lastPriceUpdate: '2024-01-12',
      needsAction: false
    },
    { 
      id: 4, 
      title: 'Terreno en Vía a la Costa', 
      price: 85000, 
      status: 'paused', 
      views: 234, 
      inquiries: 5,
      daysOnMarket: 60,
      lastPriceUpdate: '2023-12-15',
      needsAction: true
    },
    { 
      id: 5, 
      title: 'Casa en Samborondón', 
      price: 425000, 
      status: 'sold', 
      views: 523, 
      inquiries: 22,
      daysOnMarket: 33,
      lastPriceUpdate: '2024-01-08',
      needsAction: false
    },
  ];

  const recentInquiries = [
    { 
      id: 1, 
      property: 'Casa en Los Ceibos', 
      client: 'María González', 
      message: 'Estoy interesada en agendar una visita para este fin de semana',
      time: '2 horas',
      type: 'visit_request',
      priority: 'high'
    },
    { 
      id: 2, 
      property: 'Departamento en Urdesa', 
      client: 'Carlos Mendoza', 
      message: '¿Está disponible para negociar el precio?',
      time: '5 horas',
      type: 'price_negotiation',
      priority: 'medium'
    },
    { 
      id: 3, 
      property: 'Oficina en el Centro', 
      client: 'Ana Rodríguez', 
      message: 'Necesito más información sobre los gastos comunes',
      time: '1 día',
      type: 'information',
      priority: 'low'
    },
  ];

  const marketInsights = [
    { metric: 'Precio promedio zona', value: '$285,000', change: '+5.2%', trend: 'up' },
    { metric: 'Días promedio venta', value: '42 días', change: '-8%', trend: 'down' },
    { metric: 'Propiedades similares', value: '23 activas', change: '+12%', trend: 'up' },
    { metric: 'Interés del mercado', value: 'Alto', change: '+15%', trend: 'up' },
  ];

  const suggestions = [
    { 
      property: 'Departamento en Urdesa',
      type: 'price_adjustment',
      message: 'Considera reducir el precio en 5% para aumentar el interés',
      priority: 'high'
    },
    { 
      property: 'Terreno en Vía a la Costa',
      type: 'photos_update',
      message: 'Las fotos tienen más de 6 meses, considera actualizarlas',
      priority: 'medium'
    },
    { 
      property: 'Casa en Los Ceibos',
      type: 'description_update',
      message: 'Añade información sobre las mejoras recientes',
      priority: 'low'
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Mis Propiedades</h1>
          <p className="text-gray-600">Gestión y seguimiento de tus propiedades en venta</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" size="sm">
            <BarChart3 className="h-4 w-4 mr-2" />
            Análisis de Mercado
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Publicar Propiedad
          </Button>
        </div>
      </div>

      {/* Action Required Alert */}
      {myProperties.some(p => p.needsAction) && (
        <Card className="border-orange-200 bg-orange-50">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <AlertCircle className="h-5 w-5 text-orange-600" />
                <div>
                  <p className="font-medium text-orange-800">Atención requerida</p>
                  <p className="text-sm text-orange-700">
                    {myProperties.filter(p => p.needsAction).length} propiedades necesitan tu atención
                  </p>
                </div>
              </div>
              <Button variant="outline" size="sm" className="border-orange-300">
                Ver Detalles
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Propiedades Activas</CardTitle>
            <Building className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{ownerStats.activeListings}</div>
            <p className="text-xs text-muted-foreground">
              {ownerStats.totalProperties} total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Visualizaciones</CardTitle>
            <Eye className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{ownerStats.totalViews.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              +15% vs mes anterior
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Consultas</CardTitle>
            <MessageSquare className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{ownerStats.totalInquiries}</div>
            <p className="text-xs text-muted-foreground">
              +8% vs mes anterior
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Precio Promedio</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${ownerStats.averagePrice.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              {ownerStats.avgDaysOnMarket} días promedio
            </p>
          </CardContent>
        </Card>
      </div>

      {/* My Properties */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Building className="h-5 w-5" />
            Mis Propiedades
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {myProperties.map((property) => (
              <div key={property.id} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center">
                    <Building className="h-6 w-6 text-gray-400" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <p className="font-medium">{property.title}</p>
                      {property.needsAction && (
                        <AlertCircle className="h-4 w-4 text-orange-500" />
                      )}
                    </div>
                    <p className="text-sm text-gray-600">${property.price.toLocaleString()}</p>
                    <p className="text-xs text-gray-500">{property.daysOnMarket} días en el mercado</p>
                  </div>
                </div>
                
                <div className="flex items-center gap-6">
                  <div className="text-center">
                    <div className="flex items-center gap-2 text-sm">
                      <Eye className="h-4 w-4 text-gray-400" />
                      <span>{property.views}</span>
                    </div>
                    <p className="text-xs text-gray-500">vistas</p>
                  </div>
                  
                  <div className="text-center">
                    <div className="flex items-center gap-2 text-sm">
                      <MessageSquare className="h-4 w-4 text-gray-400" />
                      <span>{property.inquiries}</span>
                    </div>
                    <p className="text-xs text-gray-500">consultas</p>
                  </div>
                  
                  <div className="text-center">
                    <Badge variant={
                      property.status === 'active' ? 'default' :
                      property.status === 'negotiating' ? 'secondary' :
                      property.status === 'sold' ? 'destructive' :
                      'outline'
                    }>
                      {property.status === 'active' ? 'Activa' :
                       property.status === 'negotiating' ? 'Negociando' :
                       property.status === 'sold' ? 'Vendida' :
                       'Pausada'}
                    </Badge>
                  </div>
                  
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm">
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm">
                      <Eye className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Inquiries */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <MessageSquare className="h-5 w-5" />
              Consultas Recientes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentInquiries.map((inquiry) => (
                <div key={inquiry.id} className="flex items-start gap-3 p-3 hover:bg-gray-50 rounded-lg">
                  <div className={`w-2 h-2 rounded-full mt-2 ${
                    inquiry.priority === 'high' ? 'bg-red-500' :
                    inquiry.priority === 'medium' ? 'bg-yellow-500' :
                    'bg-green-500'
                  }`} />
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <p className="text-sm font-medium">{inquiry.client}</p>
                      <Badge variant="outline" className="text-xs">
                        {inquiry.property}
                      </Badge>
                    </div>
                    <p className="text-sm text-gray-600 mb-1">{inquiry.message}</p>
                    <p className="text-xs text-gray-500">hace {inquiry.time}</p>
                  </div>
                  <div className="flex gap-1">
                    <Button variant="ghost" size="sm">
                      <Phone className="h-3 w-3" />
                    </Button>
                    <Button variant="ghost" size="sm">
                      <Mail className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <MessageSquare className="h-4 w-4 mr-2" />
                Ver Todas las Consultas
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Market Insights */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Análisis de Mercado
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {marketInsights.map((insight, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="text-sm font-medium">{insight.metric}</p>
                    <p className="text-lg font-bold">{insight.value}</p>
                  </div>
                  <div className="text-right">
                    <div className={`flex items-center gap-1 text-sm ${
                      insight.trend === 'up' ? 'text-green-600' : 'text-red-600'
                    }`}>
                      <TrendingUp className={`h-3 w-3 ${
                        insight.trend === 'down' ? 'rotate-180' : ''
                      }`} />
                      <span>{insight.change}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <BarChart3 className="h-4 w-4 mr-2" />
                Reporte Detallado
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Suggestions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="h-5 w-5" />
            Sugerencias para Mejorar
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {suggestions.map((suggestion, index) => (
              <div key={index} className="flex items-start gap-3 p-3 hover:bg-gray-50 rounded-lg">
                <div className={`w-2 h-2 rounded-full mt-2 ${
                  suggestion.priority === 'high' ? 'bg-red-500' :
                  suggestion.priority === 'medium' ? 'bg-yellow-500' :
                  'bg-green-500'
                }`} />
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <Badge variant="outline" className="text-xs">
                      {suggestion.property}
                    </Badge>
                    <Badge variant="secondary" className="text-xs">
                      {suggestion.type === 'price_adjustment' ? 'Precio' :
                       suggestion.type === 'photos_update' ? 'Fotos' :
                       'Descripción'}
                    </Badge>
                  </div>
                  <p className="text-sm text-gray-600">{suggestion.message}</p>
                </div>
                <Button variant="ghost" size="sm">
                  <Edit className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Acciones Rápidas</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Plus className="h-5 w-5" />
              <span className="text-xs">Publicar Propiedad</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Camera className="h-5 w-5" />
              <span className="text-xs">Actualizar Fotos</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Edit className="h-5 w-5" />
              <span className="text-xs">Editar Precio</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <BarChart3 className="h-5 w-5" />
              <span className="text-xs">Ver Estadísticas</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}