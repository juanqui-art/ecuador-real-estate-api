'use client';

import {ProtectedRoute} from '@/components/auth/protected-route';
import {CanCreateProperties, CanViewAnalytics} from '@/components/auth/role-guard';
import {ModernPropertyForm2025} from '@/components/forms/modern-property-form-2025';
import {DashboardLayout} from '@/components/layout/dashboard-layout';
import {PropertyFilters} from '@/components/properties/property-filters';
import {PropertyList} from '@/components/properties/property-list';
import {PropertyStats} from '@/components/properties/property-stats';
import {Badge} from '@/components/ui/badge';
import {Button} from '@/components/ui/button';
import {Card, CardContent, CardHeader, CardTitle} from '@/components/ui/card';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger
} from '@/components/ui/dialog';
import {Input} from '@/components/ui/input';
import {Filter, Grid, List, Plus, Search} from 'lucide-react';
import {useState} from 'react';

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
                                        <Plus className="h-4 w-4 mr-2"/>
                                        Nueva Propiedad
                                    </Button>
                                </DialogTrigger>
                                <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
                                    <DialogHeader>
                                        <DialogTitle className="flex items-center gap-2">
                                            🏠 Crear Nueva Propiedad
                                            <Badge variant="secondary" className="text-xs">
                                                React 19 + Server Actions
                                            </Badge>
                                        </DialogTitle>
                                        <DialogDescription>
                                            Formulario modernizado con React 19, useActionState y validación
                                            server-side.
                                            Funciona con y sin JavaScript habilitado.
                                        </DialogDescription>
                                    </DialogHeader>
                                    <ModernPropertyForm2025
                                        onSuccess={handlePropertyCreated}
                                        onCancel={() => setShowCreateDialog(false)}
                                    />
                                </DialogContent>
                            </Dialog>
                        </CanCreateProperties>
                    </div>

                    {/* Stats Cards */}
                    <CanViewAnalytics>
                        <PropertyStats/>
                    </CanViewAnalytics>

                    {/* Search and Filters */}
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-lg">Buscar Propiedades</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <div className="flex flex-col sm:flex-row gap-4">
                                <div className="flex-1 relative">
                                    <Search
                                        className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400"/>
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
                                        <Filter className="h-4 w-4 mr-2"/>
                                        Filtros
                                    </Button>

                                    <div className="flex border rounded-md">
                                        <Button
                                            variant={viewMode === 'grid' ? "default" : "ghost"}
                                            size="sm"
                                            onClick={() => setViewMode('grid')}
                                            className="rounded-r-none"
                                        >
                                            <Grid className="h-4 w-4"/>
                                        </Button>
                                        <Button
                                            variant={viewMode === 'list' ? "default" : "ghost"}
                                            size="sm"
                                            onClick={() => setViewMode('list')}
                                            className="rounded-l-none"
                                        >
                                            <List className="h-4 w-4"/>
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