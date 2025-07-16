'use client';

import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { useAuthStore } from '@/store/auth';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Building, Users, DollarSign, TrendingUp } from 'lucide-react';

export default function SimpleDashboardPage() {
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

  const mockStats = [
    {
      title: 'Total Propiedades',
      value: '1,234',
      description: '+10% vs mes anterior',
      icon: Building,
      change: '+10%',
    },
    {
      title: 'Usuarios Activos',
      value: '567',
      description: '+5% vs mes anterior',
      icon: Users,
      change: '+5%',
    },
    {
      title: 'Ventas del Mes',
      value: '$890,123',
      description: '+15% vs mes anterior',
      icon: DollarSign,
      change: '+15%',
    },
    {
      title: 'Crecimiento',
      value: '12.5%',
      description: 'Tendencia positiva',
      icon: TrendingUp,
      change: '+2.5%',
    },
  ];

  return (
    <DashboardLayout>
      <div className="space-y-6">
        {/* Welcome Section */}
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {getRoleWelcome(user?.role || 'buyer')}
          </h1>
          <p className="text-lg text-gray-600 mt-2">
            Bienvenido, {user?.first_name} {user?.last_name}
          </p>
          <div className="mt-2">
            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              {user?.role}
            </span>
            <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              {user?.status}
            </span>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {mockStats.map((stat, index) => (
            <Card key={index}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {stat.title}
                </CardTitle>
                <stat.icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground">
                  {stat.description}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Quick Actions */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <Card>
            <CardHeader>
              <CardTitle>Acciones Rápidas</CardTitle>
              <CardDescription>
                Operaciones comunes según tu rol
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {user?.role === 'admin' && (
                  <>
                    <div className="p-2 bg-gray-50 rounded">📊 Ver todas las estadísticas</div>
                    <div className="p-2 bg-gray-50 rounded">👥 Gestionar usuarios</div>
                    <div className="p-2 bg-gray-50 rounded">🏢 Gestionar agencias</div>
                  </>
                )}
                {user?.role === 'agency' && (
                  <>
                    <div className="p-2 bg-gray-50 rounded">🏠 Gestionar propiedades</div>
                    <div className="p-2 bg-gray-50 rounded">👨‍💼 Gestionar agentes</div>
                    <div className="p-2 bg-gray-50 rounded">📈 Ver reportes</div>
                  </>
                )}
                {user?.role === 'agent' && (
                  <>
                    <div className="p-2 bg-gray-50 rounded">🏠 Mis propiedades asignadas</div>
                    <div className="p-2 bg-gray-50 rounded">📝 Crear nueva propiedad</div>
                    <div className="p-2 bg-gray-50 rounded">🤝 Mis clientes</div>
                  </>
                )}
                {user?.role === 'owner' && (
                  <>
                    <div className="p-2 bg-gray-50 rounded">🏠 Mis propiedades</div>
                    <div className="p-2 bg-gray-50 rounded">➕ Añadir nueva propiedad</div>
                    <div className="p-2 bg-gray-50 rounded">📊 Ver estadísticas</div>
                  </>
                )}
                {user?.role === 'buyer' && (
                  <>
                    <div className="p-2 bg-gray-50 rounded">🔍 Buscar propiedades</div>
                    <div className="p-2 bg-gray-50 rounded">❤️ Mis favoritos</div>
                    <div className="p-2 bg-gray-50 rounded">📞 Contactar agentes</div>
                  </>
                )}
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Estado del Sistema</CardTitle>
              <CardDescription>
                Información del backend y conexiones
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span>API Backend</span>
                  <span className="text-green-600">✅ Conectado</span>
                </div>
                <div className="flex justify-between items-center">
                  <span>Base de Datos</span>
                  <span className="text-green-600">✅ PostgreSQL</span>
                </div>
                <div className="flex justify-between items-center">
                  <span>Autenticación</span>
                  <span className="text-green-600">✅ JWT Activo</span>
                </div>
                <div className="flex justify-between items-center">
                  <span>Usuario ID</span>
                  <span className="text-sm font-mono text-gray-600">{user?.id}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span>Email</span>
                  <span className="text-sm text-gray-600">{user?.email}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </DashboardLayout>
  );
}