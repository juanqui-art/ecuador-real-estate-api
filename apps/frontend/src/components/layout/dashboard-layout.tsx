'use client';

import { useState } from 'react';
import { motion } from 'motion/react';
import { 
  Home, 
  Building, 
  Users, 
  Building2, 
  BarChart3, 
  Settings, 
  Menu,
  Bell,
  LogOut,
  Search
} from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuLabel, 
  DropdownMenuSeparator, 
  DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { Input } from '@/components/ui/input';

import { useAuthStore } from '@/store/auth';
import { useLogout } from '@/hooks/useAuth';
import { PublicSearch } from '@/components/search/public-search';
import type { UserRole } from '@shared/types/auth';

interface NavigationItem {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  roles: UserRole[];
  badge?: string;
}

const navigation: NavigationItem[] = [
  { name: 'Dashboard', href: '/dashboard', icon: Home, roles: ['admin', 'agency', 'agent', 'owner', 'buyer'] },
  { name: 'Propiedades', href: '/properties', icon: Building, roles: ['admin', 'agency', 'agent', 'owner', 'buyer'] },
  { name: 'Búsqueda', href: '/search', icon: Search, roles: ['admin', 'agency', 'agent', 'owner', 'buyer'] },
  { name: 'Usuarios', href: '/users', icon: Users, roles: ['admin', 'agency'] },
  { name: 'Agencias', href: '/agencies', icon: Building2, roles: ['admin'] },
  { name: 'Analytics', href: '/analytics', icon: BarChart3, roles: ['admin', 'agency'] },
  { name: 'Configuración', href: '/settings', icon: Settings, roles: ['admin', 'agency', 'agent', 'owner', 'buyer'] },
];

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user } = useAuthStore();
  const logoutMutation = useLogout();

  const filteredNavigation = navigation.filter(item => 
    user?.role && item.roles.includes(user.role)
  );

  const getRoleColor = (role: UserRole) => {
    switch (role) {
      case 'admin': return 'bg-red-100 text-red-800';
      case 'agency': return 'bg-blue-100 text-blue-800';
      case 'agent': return 'bg-green-100 text-green-800';
      case 'owner': return 'bg-yellow-100 text-yellow-800';
      case 'buyer': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getRoleLabel = (role: UserRole) => {
    switch (role) {
      case 'admin': return 'Administrador';
      case 'agency': return 'Agencia';
      case 'agent': return 'Agente';
      case 'owner': return 'Propietario';
      case 'buyer': return 'Comprador';
      default: return 'Usuario';
    }
  };

  const SidebarContent = () => (
    <div className="flex flex-col h-full">
      {/* Logo */}
      <div className="flex items-center px-4 py-4 border-b">
        <Building className="h-8 w-8 text-blue-600" />
        <span className="ml-2 text-xl font-bold text-gray-900">InmoEcuador</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-4 py-4 space-y-2">
        {filteredNavigation.map((item) => (
          <motion.a
            key={item.name}
            href={item.href}
            className="flex items-center px-3 py-2 text-sm font-medium rounded-md text-gray-700 hover:text-gray-900 hover:bg-gray-100 transition-colors"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <item.icon className="mr-3 h-5 w-5" />
            {item.name}
            {item.badge && (
              <Badge variant="secondary" className="ml-auto">
                {item.badge}
              </Badge>
            )}
          </motion.a>
        ))}
      </nav>

      {/* User Info */}
      <div className="px-4 py-4 border-t">
        <div className="flex items-center">
          <Avatar className="h-8 w-8">
            <AvatarFallback>
              {user?.first_name?.[0]}{user?.last_name?.[0]}
            </AvatarFallback>
          </Avatar>
          <div className="ml-3 min-w-0 flex-1">
            <p className="text-sm font-medium text-gray-900 truncate">
              {user?.first_name} {user?.last_name}
            </p>
            <Badge className={`text-xs ${getRoleColor(user?.role || 'buyer')}`}>
              {getRoleLabel(user?.role || 'buyer')}
            </Badge>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Desktop Sidebar */}
      <div className="hidden lg:flex lg:flex-col lg:w-64 lg:fixed lg:inset-y-0 lg:bg-white lg:border-r lg:border-gray-200">
        <SidebarContent />
      </div>

      {/* Mobile Sidebar */}
      <Sheet open={sidebarOpen} onOpenChange={setSidebarOpen}>
        <SheetContent side="left" className="w-64 p-0">
          <SidebarContent />
        </SheetContent>
      </Sheet>

      {/* Main Content */}
      <div className="flex flex-col flex-1 lg:ml-64">
        {/* Top Bar */}
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="flex items-center justify-between px-4 py-4">
            <div className="flex items-center">
              {/* Mobile Menu Button */}
              <Sheet>
                <SheetTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="lg:hidden"
                    onClick={() => setSidebarOpen(true)}
                  >
                    <Menu className="h-5 w-5" />
                  </Button>
                </SheetTrigger>
              </Sheet>

              {/* Search */}
              <div className="ml-4 flex-1 max-w-md">
                <PublicSearch 
                  placeholder="Buscar propiedades..."
                  onResultSelect={(result) => {
                    // Navigate to property detail page
                    window.location.href = `/properties/${result.id}`;
                  }}
                />
              </div>
            </div>

            {/* Right Side */}
            <div className="flex items-center space-x-4">
              {/* Notifications */}
              <Button variant="ghost" size="icon">
                <Bell className="h-5 w-5" />
              </Button>

              {/* User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="flex items-center space-x-2">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback>
                        {user?.first_name?.[0]}{user?.last_name?.[0]}
                      </AvatarFallback>
                    </Avatar>
                    <span className="hidden md:block text-sm font-medium">
                      {user?.first_name} {user?.last_name}
                    </span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56">
                  <DropdownMenuLabel>Mi Cuenta</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>
                    <Settings className="mr-2 h-4 w-4" />
                    Configuración
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem 
                    onClick={() => logoutMutation.mutate()} 
                    className="text-red-600"
                    disabled={logoutMutation.isPending}
                  >
                    <LogOut className="mr-2 h-4 w-4" />
                    {logoutMutation.isPending ? 'Cerrando...' : 'Cerrar Sesión'}
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </header>

        {/* Main Content Area */}
        <main className="flex-1 overflow-auto">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3 }}
            className="container mx-auto px-4 py-6"
          >
            {children}
          </motion.div>
        </main>
      </div>
    </div>
  );
}