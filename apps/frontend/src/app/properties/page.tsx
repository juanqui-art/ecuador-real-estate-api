'use client';

import { useState } from 'react';
import { Plus, Search, Filter, Grid, List } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogTrigger 
} from '@/components/ui/dialog';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { DashboardLayout } from '@/components/layout/dashboard-layout';
import { PropertyForm } from '@/components/forms/property-form';
import { PropertyList } from '@/components/properties/property-list';
import { PropertyFilters } from '@/components/properties/property-filters';
import { PropertyStats } from '@/components/properties/property-stats';
import { CanCreateProperties, CanViewAnalytics, RoleGuard } from '@/components/auth/role-guard';

export default function PropertiesPage() {
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [searchTerm, setSearchTerm] = useState('');
  const [showFilters, setShowFilters] = useState(false);
  const [filters, setFilters] = useState({
    type: '',
    status: '',
    minPrice: '',
    maxPrice: '',
    province: '',
    city: '',
    bedrooms: '',
    bathrooms: '',
  });

  const handlePropertyCreated = () => {
    setShowCreateDialog(false);
    // Refresh property list
  };

  const handleFiltersChange = (newFilters: typeof filters) => {
    setFilters(newFilters);
  };

  return (
    <ProtectedRoute requiredRole={['admin', 'agency', 'agent', 'seller']}>
      <DashboardLayout>
        <div className="space-y-6">
          {/* Header */}
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Propiedades</h1>
              <p className="text-gray-600">Gestiona tu portafolio inmobiliario</p>
            </div>
            
            <CanCreateProperties>
              <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
                <DialogTrigger asChild>
                  <Button>
                    <Plus className="h-4 w-4 mr-2" />
                    Nueva Propiedad
                  </Button>
                </DialogTrigger>
                <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                  <DialogHeader>
                    <DialogTitle>Crear Nueva Propiedad</DialogTitle>
                    <DialogDescription>
                      Completa la información de la propiedad para agregarla a tu portafolio.
                    </DialogDescription>
                  </DialogHeader>
                  <PropertyForm onSuccess={handlePropertyCreated} />
                </DialogContent>
              </Dialog>
            </CanCreateProperties>
          </div>

          {/* Stats Cards */}
          <CanViewAnalytics>
            <PropertyStats />
          </CanViewAnalytics>

          {/* Search and Filters */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Buscar Propiedades</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex flex-col sm:flex-row gap-4">
                <div className="flex-1 relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <Input
                    placeholder="Buscar por título, descripción, ubicación..."
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="pl-10"
                  />
                </div>
                
                <div className="flex gap-2">
                  <Button
                    variant={showFilters ? "default" : "outline"}
                    onClick={() => setShowFilters(!showFilters)}
                  >
                    <Filter className="h-4 w-4 mr-2" />
                    Filtros
                  </Button>
                  
                  <div className="flex border rounded-md">
                    <Button
                      variant={viewMode === 'grid' ? "default" : "ghost"}
                      size="sm"
                      onClick={() => setViewMode('grid')}
                      className="rounded-r-none"
                    >
                      <Grid className="h-4 w-4" />
                    </Button>
                    <Button
                      variant={viewMode === 'list' ? "default" : "ghost"}
                      size="sm"
                      onClick={() => setViewMode('list')}
                      className="rounded-l-none"
                    >
                      <List className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
              
              {showFilters && (
                <div className="mt-4 pt-4 border-t">
                  <PropertyFilters
                    filters={filters}
                    onFiltersChange={handleFiltersChange}
                  />
                </div>
              )}
            </CardContent>
          </Card>

          {/* Property List */}
          <PropertyList
            searchTerm={searchTerm}
            filters={filters}
            viewMode={viewMode}
          />
        </div>
      </DashboardLayout>
    </ProtectedRoute>
  );
}