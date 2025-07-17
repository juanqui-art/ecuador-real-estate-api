'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Heart, 
  Search, 
  Bell, 
  MapPin,
  DollarSign,
  Calendar,
  Filter,
  Star,
  Plus,
  Eye,
  MessageSquare,
  TrendingUp,
  Home,
  Building,
  Building2,
  Car,
  Bed,
  Bath,
  Maximize,
  Clock,
  Target,
  CheckCircle,
  Bookmark
} from 'lucide-react';

export function BuyerDashboard() {
  // Mock data - in real app, this would come from API
  const buyerStats = {
    savedProperties: 12,
    recentSearches: 8,
    notifications: 3,
    budget: 250000,
    priceAlerts: 5,
    scheduledVisits: 2,
    favoriteAreas: ['Urdesa', 'Los Ceibos', 'Samborondón'],
    avgViewsPerDay: 15,
  };

  const savedProperties = [
    { 
      id: 1, 
      title: 'Casa en Los Ceibos', 
      price: 285000, 
      location: 'Los Ceibos, Guayaquil',
      type: 'house',
      bedrooms: 4,
      bathrooms: 3,
      area: 320,
      priceChange: -5000,
      daysAgo: 2,
      isNew: false
    },
    { 
      id: 2, 
      title: 'Departamento en Urdesa', 
      price: 180000, 
      location: 'Urdesa, Guayaquil',
      type: 'apartment',
      bedrooms: 2,
      bathrooms: 2,
      area: 85,
      priceChange: 0,
      daysAgo: 5,
      isNew: false
    },
    { 
      id: 3, 
      title: 'Casa en Samborondón', 
      price: 425000, 
      location: 'Samborondón, Guayas',
      type: 'house',
      bedrooms: 5,
      bathrooms: 4,
      area: 450,
      priceChange: 0,
      daysAgo: 1,
      isNew: true
    },
  ];

  const recentSearches = [
    { id: 1, query: 'Casa 3 dormitorios Urdesa', results: 23, date: '2024-01-15' },
    { id: 2, query: 'Departamento hasta $200k', results: 45, date: '2024-01-14' },
    { id: 3, query: 'Casa con piscina Los Ceibos', results: 12, date: '2024-01-13' },
  ];

  const recommendations = [
    { 
      id: 1, 
      title: 'Villa en Vía a la Costa', 
      price: 235000, 
      location: 'Vía a la Costa',
      match: 95,
      reason: 'Precio dentro de tu presupuesto y ubicación preferida',
      bedrooms: 3,
      bathrooms: 2,
      area: 280
    },
    { 
      id: 2, 
      title: 'Departamento en Entre Ríos', 
      price: 165000, 
      location: 'Entre Ríos, Guayaquil',
      match: 87,
      reason: 'Similar a propiedades que has guardado',
      bedrooms: 2,
      bathrooms: 2,
      area: 90
    },
    { 
      id: 3, 
      title: 'Casa en Ceibos Norte', 
      price: 195000, 
      location: 'Ceibos Norte, Guayaquil',
      match: 82,
      reason: 'Excelente relación precio-calidad',
      bedrooms: 3,
      bathrooms: 2,
      area: 250
    },
  ];

  const priceAlerts = [
    { property: 'Casa en Los Ceibos', oldPrice: 290000, newPrice: 285000, change: -5000 },
    { property: 'Departamento en Urdesa Centro', oldPrice: 175000, newPrice: 165000, change: -10000 },
  ];

  const scheduledVisits = [
    { id: 1, property: 'Casa en Los Ceibos', date: '2024-01-16', time: '10:00 AM', agent: 'María González' },
    { id: 2, property: 'Departamento en Urdesa', date: '2024-01-17', time: '3:00 PM', agent: 'Carlos Mendoza' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Mi Búsqueda de Hogar</h1>
          <p className="text-gray-600">Encuentra la propiedad perfecta para ti</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" size="sm">
            <Bell className="h-4 w-4 mr-2" />
            Alertas ({buyerStats.notifications})
          </Button>
          <Button size="sm">
            <Search className="h-4 w-4 mr-2" />
            Búsqueda Avanzada
          </Button>
        </div>
      </div>

      {/* Price Alerts */}
      {priceAlerts.length > 0 && (
        <Card className="border-green-200 bg-green-50">
          <CardContent className="pt-6">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <TrendingUp className="h-5 w-5 text-green-600 rotate-180" />
                <div>
                  <p className="font-medium text-green-800">¡Bajaron de precio!</p>
                  <p className="text-sm text-green-700">
                    {priceAlerts.length} propiedades guardadas tienen mejor precio
                  </p>
                </div>
              </div>
              <Button variant="outline" size="sm" className="border-green-300">
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
            <CardTitle className="text-sm font-medium">Propiedades Guardadas</CardTitle>
            <Heart className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{buyerStats.savedProperties}</div>
            <p className="text-xs text-muted-foreground">
              +3 esta semana
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Búsquedas Recientes</CardTitle>
            <Search className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{buyerStats.recentSearches}</div>
            <p className="text-xs text-muted-foreground">
              {buyerStats.avgViewsPerDay} vistas/día
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Visitas Programadas</CardTitle>
            <Calendar className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{buyerStats.scheduledVisits}</div>
            <p className="text-xs text-muted-foreground">
              Esta semana
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Presupuesto</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${buyerStats.budget.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              {buyerStats.priceAlerts} alertas activas
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Saved Properties */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Heart className="h-5 w-5" />
              Propiedades Guardadas
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {savedProperties.map((property) => (
                <div key={property.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center">
                      {property.type === 'house' ? (
                        <Home className="h-6 w-6 text-gray-400" />
                      ) : (
                        <Building className="h-6 w-6 text-gray-400" />
                      )}
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{property.title}</p>
                        {property.isNew && (
                          <Badge variant="secondary" className="text-xs">Nuevo</Badge>
                        )}
                        {property.priceChange < 0 && (
                          <Badge variant="destructive" className="text-xs">
                            -${Math.abs(property.priceChange).toLocaleString()}
                          </Badge>
                        )}
                      </div>
                      <p className="text-sm text-gray-600">${property.price.toLocaleString()}</p>
                      <div className="flex items-center gap-3 text-xs text-gray-500">
                        <div className="flex items-center gap-1">
                          <Bed className="h-3 w-3" />
                          <span>{property.bedrooms}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Bath className="h-3 w-3" />
                          <span>{property.bathrooms}</span>
                        </div>
                        <div className="flex items-center gap-1">
                          <Maximize className="h-3 w-3" />
                          <span>{property.area}m²</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm">
                      <Eye className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm">
                      <MessageSquare className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Heart className="h-4 w-4 mr-2" />
                Ver Todas las Guardadas
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Recommendations */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Target className="h-5 w-5" />
              Recomendaciones para Ti
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recommendations.map((property) => (
                <div key={property.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 bg-gray-200 rounded-lg flex items-center justify-center">
                      <Home className="h-6 w-6 text-gray-400" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="font-medium">{property.title}</p>
                        <Badge variant="secondary" className="text-xs">
                          {property.match}% match
                        </Badge>
                      </div>
                      <p className="text-sm text-gray-600">${property.price.toLocaleString()}</p>
                      <p className="text-xs text-gray-500">{property.reason}</p>
                    </div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm">
                      <Heart className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm">
                      <Eye className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Target className="h-4 w-4 mr-2" />
                Ver Más Recomendaciones
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Scheduled Visits */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Visitas Programadas
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {scheduledVisits.map((visit) => (
              <div key={visit.id} className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg">
                <div className="w-2 h-2 bg-blue-500 rounded-full" />
                <Calendar className="h-4 w-4 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm font-medium">{visit.property}</p>
                  <p className="text-xs text-gray-500">{visit.date} a las {visit.time}</p>
                  <p className="text-xs text-gray-500">Agente: {visit.agent}</p>
                </div>
                <div className="flex gap-2">
                  <Button variant="ghost" size="sm">
                    <MapPin className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm">
                    <MessageSquare className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
          
          <div className="pt-4 border-t">
            <Button variant="outline" size="sm" className="w-full">
              <Calendar className="h-4 w-4 mr-2" />
              Agendar Nueva Visita
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Recent Searches */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Search className="h-5 w-5" />
            Búsquedas Recientes
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {recentSearches.map((search) => (
              <div key={search.id} className="flex items-center justify-between p-3 hover:bg-gray-50 rounded-lg">
                <div className="flex items-center gap-3">
                  <Search className="h-4 w-4 text-gray-400" />
                  <div>
                    <p className="text-sm font-medium">{search.query}</p>
                    <p className="text-xs text-gray-500">{search.results} resultados • {search.date}</p>
                  </div>
                </div>
                <div className="flex gap-2">
                  <Button variant="ghost" size="sm">
                    <Bookmark className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm">
                    <Search className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
          
          <div className="pt-4 border-t">
            <Button variant="outline" size="sm" className="w-full">
              <Search className="h-4 w-4 mr-2" />
              Ver Historial Completo
            </Button>
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
              <Search className="h-5 w-5" />
              <span className="text-xs">Buscar Propiedades</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Filter className="h-5 w-5" />
              <span className="text-xs">Filtros Avanzados</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Bell className="h-5 w-5" />
              <span className="text-xs">Configurar Alertas</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Calendar className="h-5 w-5" />
              <span className="text-xs">Agendar Visita</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}