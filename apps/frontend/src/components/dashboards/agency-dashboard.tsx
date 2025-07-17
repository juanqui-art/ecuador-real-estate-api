'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Building, 
  Users, 
  TrendingUp, 
  DollarSign, 
  Target,
  Calendar,
  UserPlus,
  Plus,
  Eye,
  Edit,
  MessageSquare,
  Award,
  Clock,
  Star
} from 'lucide-react';

export function AgencyDashboard() {
  // Mock data - in real app, this would come from API
  const agencyStats = {
    totalProperties: 156,
    activeListings: 134,
    soldThisMonth: 8,
    totalAgents: 12,
    monthlyRevenue: 45000,
    commissionRate: 6.5,
    averageListingTime: 45,
    clientSatisfaction: 4.8,
    monthlyTarget: 60000,
    conversionRate: 12.5,
  };

  const recentProperties = [
    { id: 1, title: 'Casa en Samborondón', price: 285000, status: 'active', views: 234, inquiries: 12 },
    { id: 2, title: 'Departamento en Urdesa', price: 180000, status: 'pending', views: 156, inquiries: 8 },
    { id: 3, title: 'Oficina en el Centro', price: 120000, status: 'sold', views: 89, inquiries: 15 },
  ];

  const topAgents = [
    { name: 'María González', sales: 6, revenue: 18000, rating: 4.9 },
    { name: 'Carlos Mendoza', sales: 5, revenue: 15000, rating: 4.7 },
    { name: 'Ana Rodríguez', sales: 4, revenue: 12000, rating: 4.8 },
  ];

  const upcomingTasks = [
    { id: 1, task: 'Reunión con cliente - Casa en Ceibos', time: '10:00 AM', priority: 'high' },
    { id: 2, task: 'Visita guiada - Departamento Centro', time: '2:00 PM', priority: 'medium' },
    { id: 3, task: 'Actualizar fotos - Casa en Urdesa', time: '4:00 PM', priority: 'low' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Dashboard Agencia</h1>
          <p className="text-gray-600">Gestión y control de tu agencia inmobiliaria</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" size="sm">
            <UserPlus className="h-4 w-4 mr-2" />
            Agregar Agente
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Nueva Propiedad
          </Button>
        </div>
      </div>

      {/* Monthly Performance */}
      <Card className="border-blue-200 bg-blue-50">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Target className="h-5 w-5 text-blue-600" />
              <div>
                <p className="font-medium text-blue-800">Meta Mensual</p>
                <p className="text-sm text-blue-700">
                  ${agencyStats.monthlyRevenue.toLocaleString()} de ${agencyStats.monthlyTarget.toLocaleString()}
                </p>
              </div>
            </div>
            <div className="text-right">
              <Badge className="bg-blue-100 text-blue-800">
                {Math.round((agencyStats.monthlyRevenue / agencyStats.monthlyTarget) * 100)}%
              </Badge>
              <Progress 
                value={(agencyStats.monthlyRevenue / agencyStats.monthlyTarget) * 100} 
                className="h-2 mt-2 w-32" 
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Propiedades Activas</CardTitle>
            <Building className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agencyStats.activeListings}</div>
            <p className="text-xs text-muted-foreground">
              {agencyStats.totalProperties} total
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Vendidas Este Mes</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agencyStats.soldThisMonth}</div>
            <p className="text-xs text-muted-foreground">
              +{agencyStats.conversionRate}% conversión
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Agentes Activos</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agencyStats.totalAgents}</div>
            <p className="text-xs text-muted-foreground">
              Rating promedio: {agencyStats.clientSatisfaction}⭐
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ingresos Mensuales</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${agencyStats.monthlyRevenue.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              Comisión: {agencyStats.commissionRate}%
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Properties */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Building className="h-5 w-5" />
              Propiedades Recientes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {recentProperties.map((property) => (
                <div key={property.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="font-medium">{property.title}</p>
                    <p className="text-sm text-gray-600">${property.price.toLocaleString()}</p>
                  </div>
                  <div className="text-right">
                    <Badge variant={
                      property.status === 'active' ? 'default' :
                      property.status === 'pending' ? 'secondary' :
                      'destructive'
                    }>
                      {property.status === 'active' ? 'Activa' :
                       property.status === 'pending' ? 'Pendiente' :
                       'Vendida'}
                    </Badge>
                    <div className="flex items-center gap-2 text-xs text-gray-500 mt-1">
                      <Eye className="h-3 w-3" />
                      <span>{property.views}</span>
                      <MessageSquare className="h-3 w-3" />
                      <span>{property.inquiries}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Building className="h-4 w-4 mr-2" />
                Ver Todas las Propiedades
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Top Agents */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="h-5 w-5" />
              Top Agentes del Mes
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topAgents.map((agent, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className={`w-8 h-8 rounded-full flex items-center justify-center text-white font-bold ${
                      index === 0 ? 'bg-yellow-500' :
                      index === 1 ? 'bg-gray-400' :
                      'bg-amber-600'
                    }`}>
                      {index + 1}
                    </div>
                    <div>
                      <p className="font-medium">{agent.name}</p>
                      <p className="text-sm text-gray-600">{agent.sales} ventas</p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="font-medium">${agent.revenue.toLocaleString()}</p>
                    <div className="flex items-center gap-1 text-sm text-gray-600">
                      <Star className="h-3 w-3 fill-yellow-400 text-yellow-400" />
                      <span>{agent.rating}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Users className="h-4 w-4 mr-2" />
                Gestionar Agentes
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Upcoming Tasks */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5" />
            Agenda de Hoy
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {upcomingTasks.map((task) => (
              <div key={task.id} className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg">
                <div className={`w-2 h-2 rounded-full ${
                  task.priority === 'high' ? 'bg-red-500' :
                  task.priority === 'medium' ? 'bg-yellow-500' :
                  'bg-green-500'
                }`} />
                <Clock className="h-4 w-4 text-gray-400" />
                <div className="flex-1">
                  <p className="text-sm font-medium">{task.task}</p>
                  <p className="text-xs text-gray-500">{task.time}</p>
                </div>
                <Button variant="ghost" size="sm">
                  <Edit className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
          
          <div className="pt-4 border-t">
            <Button variant="outline" size="sm" className="w-full">
              <Calendar className="h-4 w-4 mr-2" />
              Ver Agenda Completa
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
              <Plus className="h-5 w-5" />
              <span className="text-xs">Nueva Propiedad</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <UserPlus className="h-5 w-5" />
              <span className="text-xs">Agregar Agente</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <TrendingUp className="h-5 w-5" />
              <span className="text-xs">Reportes</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <MessageSquare className="h-5 w-5" />
              <span className="text-xs">Mensajes</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}