'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { 
  Building, 
  Users, 
  Building2, 
  TrendingUp, 
  DollarSign, 
  Activity,
  AlertTriangle,
  Plus,
  Settings,
  BarChart3,
  UserCheck,
  FileText,
  Shield
} from 'lucide-react';

export function AdminDashboard() {
  // Mock data - in real app, this would come from API
  const systemStats = {
    totalProperties: 1247,
    totalUsers: 342,
    totalAgencies: 28,
    monthlyRevenue: 125000,
    activeListings: 856,
    pendingApprovals: 15,
    systemHealth: 98.5,
    lastBackup: '2 horas',
  };

  const recentActivity = [
    { id: 1, type: 'user', message: 'Nueva agencia registrada: InmoMax', time: '10 min' },
    { id: 2, type: 'property', message: '12 propiedades nuevas publicadas', time: '25 min' },
    { id: 3, type: 'system', message: 'Backup automático completado', time: '1 hora' },
    { id: 4, type: 'alert', message: 'Límite de almacenamiento al 85%', time: '2 horas' },
  ];

  const topAgencies = [
    { name: 'InmoMax', properties: 156, performance: 95, growth: '+12%' },
    { name: 'Casas Premium', properties: 134, performance: 92, growth: '+8%' },
    { name: 'Propiedades Elite', properties: 118, performance: 88, growth: '+15%' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Dashboard Administrativo</h1>
          <p className="text-gray-600">Gestión completa del sistema InmoEcuador</p>
        </div>
        <div className="flex gap-3">
          <Button variant="outline" size="sm">
            <Settings className="h-4 w-4 mr-2" />
            Configuración
          </Button>
          <Button size="sm">
            <Plus className="h-4 w-4 mr-2" />
            Nueva Agencia
          </Button>
        </div>
      </div>

      {/* System Health Alert */}
      <Card className="border-yellow-200 bg-yellow-50">
        <CardContent className="pt-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <AlertTriangle className="h-5 w-5 text-yellow-600" />
              <div>
                <p className="font-medium text-yellow-800">Atención requerida</p>
                <p className="text-sm text-yellow-700">
                  {systemStats.pendingApprovals} solicitudes pendientes de aprobación
                </p>
              </div>
            </div>
            <Button variant="outline" size="sm" className="border-yellow-300">
              Revisar
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Propiedades</CardTitle>
            <Building className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{systemStats.totalProperties.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              +12% vs mes anterior
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Usuarios Activos</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{systemStats.totalUsers}</div>
            <p className="text-xs text-muted-foreground">
              +8% vs mes anterior
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Agencias Registradas</CardTitle>
            <Building2 className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{systemStats.totalAgencies}</div>
            <p className="text-xs text-muted-foreground">
              +3 este mes
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Ingresos Mensuales</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${systemStats.monthlyRevenue.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground">
              +20% vs mes anterior
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* System Health */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Estado del Sistema
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Rendimiento General</span>
              <Badge className="bg-green-100 text-green-800">
                {systemStats.systemHealth}% Óptimo
              </Badge>
            </div>
            <Progress value={systemStats.systemHealth} className="h-2" />
            
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Listados Activos</span>
                <span>{systemStats.activeListings}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Último Backup</span>
                <span>{systemStats.lastBackup}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span>Usuarios Conectados</span>
                <span>127</span>
              </div>
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <Shield className="h-4 w-4 mr-2" />
                Panel de Seguridad
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Top Performing Agencies */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Top Agencias
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {topAgencies.map((agency, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <div>
                    <p className="font-medium">{agency.name}</p>
                    <p className="text-sm text-gray-600">{agency.properties} propiedades</p>
                  </div>
                  <div className="text-right">
                    <Badge variant="secondary">{agency.performance}%</Badge>
                    <p className="text-sm text-green-600 font-medium">{agency.growth}</p>
                  </div>
                </div>
              ))}
            </div>
            
            <div className="pt-4 border-t">
              <Button variant="outline" size="sm" className="w-full">
                <BarChart3 className="h-4 w-4 mr-2" />
                Ver Análisis Completo
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Activity className="h-5 w-5" />
            Actividad Reciente
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-3">
            {recentActivity.map((activity) => (
              <div key={activity.id} className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg">
                <div className={`w-2 h-2 rounded-full ${
                  activity.type === 'alert' ? 'bg-red-500' :
                  activity.type === 'user' ? 'bg-blue-500' :
                  activity.type === 'property' ? 'bg-green-500' :
                  'bg-gray-500'
                }`} />
                <div className="flex-1">
                  <p className="text-sm">{activity.message}</p>
                  <p className="text-xs text-gray-500">hace {activity.time}</p>
                </div>
              </div>
            ))}
          </div>
          
          <div className="pt-4 border-t">
            <Button variant="outline" size="sm" className="w-full">
              <FileText className="h-4 w-4 mr-2" />
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
              <UserCheck className="h-5 w-5" />
              <span className="text-xs">Aprobar Usuarios</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Building2 className="h-5 w-5" />
              <span className="text-xs">Gestionar Agencias</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <BarChart3 className="h-5 w-5" />
              <span className="text-xs">Reportes</span>
            </Button>
            <Button variant="outline" className="h-16 flex flex-col gap-2">
              <Settings className="h-5 w-5" />
              <span className="text-xs">Configuración</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}