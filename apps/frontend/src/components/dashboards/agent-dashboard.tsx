'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Building, 
  TrendingUp, 
  DollarSign, 
  Target,
  Calendar,
  Phone,
  Mail,
  MapPin,
  Clock,
  Users,
  Star,
  CheckCircle,
  AlertCircle,
  Plus,
  Edit,
  Eye,
  MessageSquare
} from 'lucide-react';

export function AgentDashboard() {
  // Mock data - in real app, this would come from API
  const agentStats = {
    assignedProperties: 23,
    activeListings: 18,
    soldThisMonth: 3,
    totalEarnings: 12000,
    monthlyTarget: 15000,
    clientRating: 4.8,
    responsiveness: 95,
    conversionRate: 13.2,
    averageResponseTime: 15, // minutes
  };

  const myProperties = [
    { id: 1, title: 'Casa en Los Ceibos', price: 325000, status: 'active', inquiries: 5, nextVisit: '2024-01-15 10:00' },
    { id: 2, title: 'Departamento en Urdesa', price: 180000, status: 'negotiating', inquiries: 12, nextVisit: '2024-01-15 14:00' },
    { id: 3, title: 'Oficina en el Centro', price: 120000, status: 'sold', inquiries: 8, nextVisit: null },
  ];

  const todaySchedule = [
    { id: 1, time: '09:00', type: 'visit', client: 'María González', property: 'Casa en Los Ceibos', status: 'confirmed' },
    { id: 2, time: '11:30', type: 'call', client: 'Carlos Mendoza', property: 'Depto en Urdesa', status: 'pending' },
    { id: 3, time: '14:00', type: 'visit', client: 'Ana Rodríguez', property: 'Oficina Centro', status: 'confirmed' },
    { id: 4, time: '16:00', type: 'meeting', client: 'Luis Pérez', property: 'Casa en Ceibos', status: 'rescheduled' },
  ];

  const recentInquiries = [
    { id: 1, client: 'Pedro Sánchez', property: 'Casa en Los Ceibos', message: 'Quisiera agendar una visita para este fin de semana', time: '10 min' },
    { id: 2, client: 'Sofía Morales', property: 'Depto en Urdesa', message: '¿Está disponible para negociar el precio?', time: '25 min' },
    { id: 3, client: 'Roberto Castro', property: 'Oficina Centro', message: 'Necesito más información sobre los servicios incluidos', time: '1 hora' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Mi Dashboard</h1>
          <p className="text-gray-600">Gestión personal de propiedades y clientes</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" size="sm">
            <Calendar className="h-4 w-4 mr-2" />
            Mi Agenda
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Agregar Cita
          </Button>
        </div>
      </div>

      {/* Performance Alert */}
      <Card className="border-green-200 bg-green-50">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Target className="h-5 w-5 text-green-600" />
              <div>
                <p className="font-medium text-green-800">Meta Mensual</p>
                <p className="text-sm text-green-700">
                  ${agentStats.totalEarnings.toLocaleString()} de ${agentStats.monthlyTarget.toLocaleString()}
                </p>
              </div>
            </div>
            <div className="text-right">
              <Badge className="bg-green-100 text-green-800">
                {Math.round((agentStats.totalEarnings / agentStats.monthlyTarget) * 100)}%
              </Badge>
              <Progress 
                value={(agentStats.totalEarnings / agentStats.monthlyTarget) * 100} 
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
            <CardTitle className="text-sm font-medium">Propiedades Asignadas</CardTitle>
            <Building className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agentStats.assignedProperties}</div>
            <p className="text-xs text-muted-foreground">
              {agentStats.activeListings} activas
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ventas Este Mes</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agentStats.soldThisMonth}</div>
            <p className="text-xs text-muted-foreground">
              {agentStats.conversionRate}% conversión
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Calificación Cliente</CardTitle>
            <Star className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{agentStats.clientRating}</div>
            <p className="text-xs text-muted-foreground">
              {agentStats.responsiveness}% responsividad
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ganancias Mensuales</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${agentStats.totalEarnings.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              Resp. promedio: {agentStats.averageResponseTime}min
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
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
                <div key={property.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="font-medium">{property.title}</p>
                    <p className="text-sm text-gray-600">${property.price.toLocaleString()}</p>
                  </div>
                  <div className="text-right">
                    <Badge variant={
                      property.status === 'active' ? 'default' :
                      property.status === 'negotiating' ? 'secondary' :
                      'destructive'
                    }>
                      {property.status === 'active' ? 'Activa' :
                       property.status === 'negotiating' ? 'Negociando' :
                       'Vendida'}
                    </Badge>
                    <div className="flex items-center gap-2 text-xs text-gray-500 mt-1">
                      <MessageSquare className="h-3 w-3" />
                      <span>{property.inquiries} consultas</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Building className="h-4 w-4 mr-2" />
                Ver Todas Mis Propiedades
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Today's Schedule */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              Agenda de Hoy
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {todaySchedule.map((appointment) => (
                <div key={appointment.id} className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg">
                  <div className="text-center">
                    <p className="text-sm font-medium">{appointment.time}</p>
                  </div>
                  <div className={`w-2 h-2 rounded-full ${
                    appointment.status === 'confirmed' ? 'bg-green-500' :
                    appointment.status === 'pending' ? 'bg-yellow-500' :
                    'bg-red-500'
                  }`} />
                  <div className="flex-1">
                    <p className="text-sm font-medium">{appointment.client}</p>
                    <p className="text-xs text-gray-500">{appointment.property}</p>
                  </div>
                  <div className="flex items-center gap-1">
                    {appointment.type === 'visit' && <MapPin className="h-3 w-3 text-gray-400" />}
                    {appointment.type === 'call' && <Phone className="h-3 w-3 text-gray-400" />}
                    {appointment.type === 'meeting' && <Users className="h-3 w-3 text-gray-400" />}
                  </div>
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
      </div>

      {/* Recent Inquiries */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <MessageSquare className="h-5 w-5" />
            Consultas Recientes
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {recentInquiries.map((inquiry) => (
              <div key={inquiry.id} className="flex items-start gap-3 p-3 hover:bg-gray-50 rounded-lg">
                <div className="w-2 h-2 bg-blue-500 rounded-full mt-2" />
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

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Acciones Rápidas</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Plus className="h-5 w-5" />
              <span className="text-xs">Agendar Cita</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Edit className="h-5 w-5" />
              <span className="text-xs">Actualizar Propiedad</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Phone className="h-5 w-5" />
              <span className="text-xs">Llamar Cliente</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <TrendingUp className="h-5 w-5" />
              <span className="text-xs">Ver Desempeño</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}