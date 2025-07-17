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
import { RoleBasedNavigation, RoleBasedQuickActions, RoleBasedNavStats } from '@/components/navigation/role-based-navigation';
import type { UserRole } from '@shared/types/auth';

interface DashboardLayoutProps {
  children: React.ReactNode;
}

export function DashboardLayout({ children }: DashboardLayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user } = useAuthStore();
  const logoutMutation = useLogout();

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
      <div className="flex-1 px-4 py-4 space-y-4">
        <RoleBasedNavigation />
        
        {/* Quick Actions */}
        <div className="pt-4 border-t">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
            Acciones Rápidas
          </h3>
          <RoleBasedQuickActions />
        </div>
      </div>

      {/* Stats */}
      <div className="px-4 py-2">
        <RoleBasedNavStats />
      </div>

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