'use client';

import { motion } from 'motion/react';
import { 
  Home, 
  Building, 
  Users, 
  Building2, 
  BarChart3, 
  Settings, 
  Search,
  Heart,
  Calendar,
  MessageSquare,
  FileText,
  Shield,
  CreditCard,
  Target,
  Award,
  TrendingUp,
  UserCheck,
  Bell,
  Eye,
  Plus
} from 'lucide-react';

import { Badge } from '@/components/ui/badge';
import { useAuthStore } from '@/store/auth';
import { RoleGuard } from '@/components/auth/role-guard';
import type { UserRole } from '@shared/types/auth';

interface NavigationItem {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  roles: UserRole[];
  badge?: string;
  badgeColor?: string;
  description?: string;
}

// Role-specific navigation items
const navigationItems: NavigationItem[] = [
  // Common to all roles
  { 
    name: 'Dashboard', 
    href: '/dashboard', 
    icon: Home, 
    roles: ['admin', 'agency', 'agent', 'owner', 'buyer'],
    description: 'Panel principal'
  },
  
  // Admin-specific navigation
  {
    name: 'Administración',
    href: '/admin',
    icon: Shield,
    roles: ['admin'],
    description: 'Gestión del sistema'
  },
  {
    name: 'Usuarios',
    href: '/users',
    icon: Users,
    roles: ['admin', 'agency'],
    description: 'Gestión de usuarios'
  },
  {
    name: 'Agencias',
    href: '/agencies',
    icon: Building2,
    roles: ['admin'],
    description: 'Gestión de agencias'
  },
  {
    name: 'Reportes del Sistema',
    href: '/admin/reports',
    icon: BarChart3,
    roles: ['admin'],
    description: 'Reportes globales'
  },
  {
    name: 'Configuración Global',
    href: '/admin/settings',
    icon: Settings,
    roles: ['admin'],
    description: 'Configuración del sistema'
  },
  
  // Agency-specific navigation
  {
    name: 'Mi Agencia',
    href: '/agency',
    icon: Building2,
    roles: ['agency'],
    description: 'Gestión de agencia'
  },
  {
    name: 'Mis Agentes',
    href: '/agency/agents',
    icon: UserCheck,
    roles: ['agency'],
    description: 'Gestión de agentes'
  },
  {
    name: 'Reportes de Agencia',
    href: '/agency/reports',
    icon: BarChart3,
    roles: ['agency'],
    description: 'Análisis de desempeño'
  },
  {
    name: 'Comisiones',
    href: '/agency/commissions',
    icon: CreditCard,
    roles: ['agency'],
    description: 'Gestión de comisiones'
  },
  
  // Agent-specific navigation
  {
    name: 'Mis Clientes',
    href: '/agent/clients',
    icon: Users,
    roles: ['agent'],
    description: 'Gestión de clientes'
  },
  {
    name: 'Mi Agenda',
    href: '/agent/schedule',
    icon: Calendar,
    roles: ['agent'],
    description: 'Citas y visitas'
  },
  {
    name: 'Mis Ventas',
    href: '/agent/sales',
    icon: Award,
    roles: ['agent'],
    description: 'Historial de ventas'
  },
  {
    name: 'Mi Desempeño',
    href: '/agent/performance',
    icon: TrendingUp,
    roles: ['agent'],
    description: 'Métricas personales'
  },
  
  // Owner-specific navigation
  {
    name: 'Mis Propiedades',
    href: '/owner/properties',
    icon: Building,
    roles: ['owner'],
    description: 'Gestión de propiedades'
  },
  {
    name: 'Publicar Propiedad',
    href: '/owner/publish',
    icon: Plus,
    roles: ['owner'],
    description: 'Nueva publicación'
  },
  {
    name: 'Consultas',
    href: '/owner/inquiries',
    icon: MessageSquare,
    roles: ['owner'],
    badge: '3',
    badgeColor: 'bg-blue-500',
    description: 'Mensajes de compradores'
  },
  {
    name: 'Análisis de Mercado',
    href: '/owner/market',
    icon: Target,
    roles: ['owner'],
    description: 'Tendencias del mercado'
  },
  
  // Buyer-specific navigation
  {
    name: 'Buscar Propiedades',
    href: '/search',
    icon: Search,
    roles: ['buyer'],
    description: 'Encuentra tu hogar ideal'
  },
  {
    name: 'Mis Favoritos',
    href: '/buyer/favorites',
    icon: Heart,
    roles: ['buyer'],
    badge: '5',
    badgeColor: 'bg-red-500',
    description: 'Propiedades guardadas'
  },
  {
    name: 'Mis Visitas',
    href: '/buyer/visits',
    icon: Calendar,
    roles: ['buyer'],
    description: 'Citas programadas'
  },
  {
    name: 'Mis Alertas',
    href: '/buyer/alerts',
    icon: Bell,
    roles: ['buyer'],
    description: 'Notificaciones de precios'
  },
  {
    name: 'Historial de Búsqueda',
    href: '/buyer/history',
    icon: Eye,
    roles: ['buyer'],
    description: 'Búsquedas anteriores'
  },
  
  // Common navigation (role-based visibility)
  {
    name: 'Propiedades',
    href: '/properties',
    icon: Building,
    roles: ['admin', 'agency', 'agent', 'owner', 'buyer'],
    description: 'Ver todas las propiedades'
  },
  {
    name: 'Mensajes',
    href: '/messages',
    icon: MessageSquare,
    roles: ['admin', 'agency', 'agent', 'owner', 'buyer'],
    badge: '2',
    badgeColor: 'bg-green-500',
    description: 'Centro de mensajes'
  },
  {
    name: 'Documentos',
    href: '/documents',
    icon: FileText,
    roles: ['admin', 'agency', 'agent', 'owner'],
    description: 'Gestión de documentos'
  },
  {
    name: 'Configuración',
    href: '/settings',
    icon: Settings,
    roles: ['admin', 'agency', 'agent', 'owner', 'buyer'],
    description: 'Configuración personal'
  },
];

interface RoleBasedNavigationProps {
  className?: string;
  showDescriptions?: boolean;
  compactMode?: boolean;
}

export function RoleBasedNavigation({ 
  className = '', 
  showDescriptions = false,
  compactMode = false 
}: RoleBasedNavigationProps) {
  const { user } = useAuthStore();

  if (!user) return null;

  const filteredNavigation = navigationItems.filter(item => 
    item.roles.includes(user.role)
  );

  return (
    <nav className={`space-y-2 ${className}`}>
      {filteredNavigation.map((item) => (
        <RoleGuard key={item.name} allowedRoles={item.roles}>
          <motion.a
            href={item.href}
            className={`flex items-center px-3 py-2 text-sm font-medium rounded-md text-gray-700 hover:text-gray-900 hover:bg-gray-100 transition-colors ${
              compactMode ? 'justify-center' : ''
            }`}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            title={compactMode ? item.name : undefined}
          >
            <item.icon className={`h-5 w-5 ${compactMode ? '' : 'mr-3'}`} />
            {!compactMode && (
              <>
                <span className="flex-1">{item.name}</span>
                {item.badge && (
                  <Badge 
                    variant="secondary" 
                    className={`ml-auto text-xs ${item.badgeColor || 'bg-gray-500'} text-white`}
                  >
                    {item.badge}
                  </Badge>
                )}
              </>
            )}
          </motion.a>
          {showDescriptions && !compactMode && item.description && (
            <p className="text-xs text-gray-500 px-3 pb-2">{item.description}</p>
          )}
        </RoleGuard>
      ))}
    </nav>
  );
}

// Quick actions based on role
export function RoleBasedQuickActions() {
  const { user } = useAuthStore();

  if (!user) return null;

  return (
    <div className="space-y-2">
      <RoleGuard allowedRoles={['admin']}>
        <motion.button 
          className="w-full text-left px-3 py-2 text-sm rounded-md hover:bg-gray-100 transition-colors"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Shield className="h-4 w-4 mr-2 inline" />
          Panel de Administración
        </motion.button>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['agency']}>
        <motion.button 
          className="w-full text-left px-3 py-2 text-sm rounded-md hover:bg-gray-100 transition-colors"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <UserCheck className="h-4 w-4 mr-2 inline" />
          Gestionar Agentes
        </motion.button>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['agent']}>
        <motion.button 
          className="w-full text-left px-3 py-2 text-sm rounded-md hover:bg-gray-100 transition-colors"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Calendar className="h-4 w-4 mr-2 inline" />
          Agendar Visita
        </motion.button>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['owner']}>
        <motion.button 
          className="w-full text-left px-3 py-2 text-sm rounded-md hover:bg-gray-100 transition-colors"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Plus className="h-4 w-4 mr-2 inline" />
          Publicar Propiedad
        </motion.button>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['buyer']}>
        <motion.button 
          className="w-full text-left px-3 py-2 text-sm rounded-md hover:bg-gray-100 transition-colors"
          whileHover={{ scale: 1.02 }}
          whileTap={{ scale: 0.98 }}
        >
          <Search className="h-4 w-4 mr-2 inline" />
          Buscar Propiedades
        </motion.button>
      </RoleGuard>
    </div>
  );
}

// Role-based navigation stats
export function RoleBasedNavStats() {
  const { user } = useAuthStore();

  if (!user) return null;

  return (
    <div className="px-4 py-2 bg-gray-50 rounded-lg">
      <RoleGuard allowedRoles={['admin']}>
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span>Sistema</span>
            <span className="text-green-600">Estable</span>
          </div>
          <div className="flex justify-between text-xs">
            <span>Usuarios</span>
            <span>1,247</span>
          </div>
        </div>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['agency']}>
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span>Agentes</span>
            <span>12</span>
          </div>
          <div className="flex justify-between text-xs">
            <span>Propiedades</span>
            <span>156</span>
          </div>
        </div>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['agent']}>
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span>Asignadas</span>
            <span>23</span>
          </div>
          <div className="flex justify-between text-xs">
            <span>Vendidas</span>
            <span>5</span>
          </div>
        </div>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['owner']}>
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span>Activas</span>
            <span>3</span>
          </div>
          <div className="flex justify-between text-xs">
            <span>Consultas</span>
            <span>8</span>
          </div>
        </div>
      </RoleGuard>
      
      <RoleGuard allowedRoles={['buyer']}>
        <div className="space-y-1">
          <div className="flex justify-between text-xs">
            <span>Favoritos</span>
            <span>5</span>
          </div>
          <div className="flex justify-between text-xs">
            <span>Alertas</span>
            <span>3</span>
          </div>
        </div>
      </RoleGuard>
    </div>
  );
}