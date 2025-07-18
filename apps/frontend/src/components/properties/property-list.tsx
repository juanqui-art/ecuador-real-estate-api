'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { 
  MapPin, 
  Bed, 
  Bath, 
  Square, 
  Car, 
  Edit, 
  Trash2, 
  Eye,
  Star,
  Calendar,
  ImageIcon
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle, 
  DialogFooter 
} from '@/components/ui/dialog';
import { apiClient } from '@/lib/api-client';
import { EditPropertyDialog } from '@/components/dialogs/edit-property-dialog';
import { 
  Property, 
  ApiResponse
} from '@shared/types/property';
import { formatPrice, formatArea, formatDate } from '@/lib/utils';

interface PropertyListProps {
  searchTerm: string;
  filters: {
    type: string;
    status: string;
    minPrice: string;
    maxPrice: string;
    province: string;
    city: string;
    bedrooms: string;
    bathrooms: string;
  };
  viewMode: 'grid' | 'list';
}

export function PropertyList({ searchTerm, filters, viewMode }: PropertyListProps) {
  const [page, setPage] = useState(1);
  const [limit] = useState(12);
  const [selectedProperty, setSelectedProperty] = useState<Property | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  
  const queryClient = useQueryClient();

  const { data: propertiesResponse, isLoading, error } = useQuery({
    queryKey: ['properties', searchTerm, filters, page, limit],
    queryFn: async () => {
      try {
        const params = new URLSearchParams();
        
        if (searchTerm) params.append('q', searchTerm);
        if (filters.type) params.append('type', filters.type);
        if (filters.status) params.append('status', filters.status);
        if (filters.province) params.append('province', filters.province);
        if (filters.city) params.append('city', filters.city);
        if (filters.bedrooms) params.append('bedrooms', filters.bedrooms);
        if (filters.bathrooms) params.append('bathrooms', filters.bathrooms);
        if (filters.minPrice) params.append('min_price', filters.minPrice);
        if (filters.maxPrice) params.append('max_price', filters.maxPrice);
        
        // Use filter endpoint for search, regular endpoint for list
        const endpoint = searchTerm || Object.values(filters).some(Boolean) 
          ? `/api/properties/filter?${params.toString()}`
          : '/api/properties';
        
        const response = await apiClient.get<ApiResponse<Property[]>>(endpoint);
        return response.data;
      } catch (error: any) {
        console.error('Error fetching properties:', error);
        throw error;
      }
    },
    retry: (failureCount, error) => {
      // Don't retry on auth errors
      if (error?.status === 401 || error?.status === 403) {
        return false;
      }
      return failureCount < 3;
    },
  });
  
  const properties = propertiesResponse?.data || [];

  // Delete property mutation
  const deletePropertyMutation = useMutation({
    mutationFn: async (propertyId: string) => {
      try {
        const response = await apiClient.delete(`/api/properties/${propertyId}`);
        return response.data;
      } catch (error: any) {
        console.error('Error deleting property:', error);
        throw error;
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      setIsDeleteDialogOpen(false);
      setSelectedProperty(null);
    },
    onError: (error: any) => {
      console.error('Error deleting property:', error);
      // Here you could show a toast notification
      if (error?.status === 401) {
        console.error('Unauthorized - redirecting to login');
      } else if (error?.status === 403) {
        console.error('Forbidden - insufficient permissions');
      }
    },
  });

  // Handle delete property
  const handleDeleteProperty = (property: Property) => {
    setSelectedProperty(property);
    setIsDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (selectedProperty) {
      deletePropertyMutation.mutate(selectedProperty.id);
    }
  };

  // Handle edit property
  const handleEditProperty = (property: Property) => {
    setSelectedProperty(property);
    setIsEditDialogOpen(true);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available':
        return 'bg-green-100 text-green-800';
      case 'sold':
        return 'bg-red-100 text-red-800';
      case 'rented':
        return 'bg-blue-100 text-blue-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'available':
        return 'Disponible';
      case 'sold':
        return 'Vendida';
      case 'rented':
        return 'Rentada';
      default:
        return status;
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'house':
        return 'Casa';
      case 'apartment':
        return 'Departamento';
      case 'land':
        return 'Terreno';
      case 'commercial':
        return 'Comercial';
      default:
        return type;
    }
  };

  const PropertyCard = ({ property }: { property: Property }) => (
    <Card className="group hover:shadow-lg transition-shadow overflow-hidden">
      {/* Property Image */}
      <div className="relative h-48 bg-gray-100 overflow-hidden">
        {property.main_image ? (
          <img 
            src={property.main_image} 
            alt={property.title}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center bg-gray-100">
            <div className="text-center">
              <ImageIcon className="h-12 w-12 mx-auto text-gray-400 mb-2" />
              <p className="text-sm text-gray-500">Sin imagen</p>
            </div>
          </div>
        )}
        
        {/* Status Badge */}
        <div className="absolute top-3 left-3">
          <Badge className={getStatusColor(property.status)}>
            {getStatusLabel(property.status)}
          </Badge>
        </div>
        
        {/* Featured Badge */}
        {property.featured && (
          <div className="absolute top-3 right-3">
            <Badge variant="secondary">
              <Star className="h-3 w-3 mr-1" />
              Destacada
            </Badge>
          </div>
        )}
        
        {/* Price Overlay */}
        <div className="absolute bottom-3 left-3">
          <div className="bg-black bg-opacity-75 text-white px-3 py-1 rounded-md">
            <p className="text-lg font-bold">{formatPrice(property.price)}</p>
          </div>
        </div>
      </div>

      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant="outline">{getTypeLabel(property.type)}</Badge>
              {property.images && property.images.length > 0 && (
                <Badge variant="secondary" className="text-xs">
                  <ImageIcon className="h-3 w-3 mr-1" />
                  {property.images.length} foto{property.images.length > 1 ? 's' : ''}
                </Badge>
              )}
            </div>
            <h3 className="font-semibold text-lg line-clamp-2 group-hover:text-primary transition-colors">
              {property.title}
            </h3>
            <p className="text-sm text-gray-600 flex items-center mt-1">
              <MapPin className="h-3 w-3 mr-1" />
              {property.city}, {property.province}
            </p>
          </div>
          <div className="text-right">
            <p className="text-sm text-gray-500">
              {formatArea(property.area_m2)}
            </p>
          </div>
        </div>
      </CardHeader>
      
      <CardContent>
        <p className="text-sm text-gray-600 line-clamp-2 mb-4">
          {property.description}
        </p>
        
        <div className="grid grid-cols-4 gap-4 mb-4">
          <div className="flex items-center text-sm text-gray-600">
            <Bed className="h-4 w-4 mr-1" />
            {property.bedrooms}
          </div>
          <div className="flex items-center text-sm text-gray-600">
            <Bath className="h-4 w-4 mr-1" />
            {property.bathrooms}
          </div>
          <div className="flex items-center text-sm text-gray-600">
            <Square className="h-4 w-4 mr-1" />
            {property.area_m2}m²
          </div>
          <div className="flex items-center text-sm text-gray-600">
            <Car className="h-4 w-4 mr-1" />
            {property.parking_spaces}
          </div>
        </div>

        {/* Features */}
        <div className="flex flex-wrap gap-1 mb-4">
          {property.pool && (
            <Badge variant="outline" className="text-xs">Piscina</Badge>
          )}
          {property.garden && (
            <Badge variant="outline" className="text-xs">Jardín</Badge>
          )}
          {property.elevator && (
            <Badge variant="outline" className="text-xs">Ascensor</Badge>
          )}
          {property.furnished && (
            <Badge variant="outline" className="text-xs">Amueblado</Badge>
          )}
          {property.garage && (
            <Badge variant="outline" className="text-xs">Garage</Badge>
          )}
          {property.security && (
            <Badge variant="outline" className="text-xs">Seguridad</Badge>
          )}
          {property.terrace && (
            <Badge variant="outline" className="text-xs">Terraza</Badge>
          )}
          {property.balcony && (
            <Badge variant="outline" className="text-xs">Balcón</Badge>
          )}
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center text-xs text-gray-500">
            <Calendar className="h-3 w-3 mr-1" />
            {formatDate(property.created_at)}
          </div>
          <div className="flex gap-2">
            <Button variant="outline" size="sm">
              <Eye className="h-3 w-3 mr-1" />
              Ver
            </Button>
            <Button 
              variant="outline" 
              size="sm"
              onClick={() => handleEditProperty(property)}
            >
              <Edit className="h-3 w-3 mr-1" />
              Editar
            </Button>
            <Button 
              variant="outline" 
              size="sm" 
              className="text-red-600 hover:text-red-700"
              onClick={() => handleDeleteProperty(property)}
            >
              <Trash2 className="h-3 w-3" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  const PropertyRow = ({ property }: { property: Property }) => (
    <Card className="group hover:shadow-lg transition-shadow">
      <CardContent className="p-4">
        <div className="flex items-center gap-6">
          {/* Property Image */}
          <div className="flex-shrink-0">
            <div className="w-24 h-24 bg-gray-100 rounded-lg overflow-hidden">
              {property.main_image ? (
                <img 
                  src={property.main_image} 
                  alt={property.title}
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center">
                  <ImageIcon className="h-8 w-8 text-gray-400" />
                </div>
              )}
            </div>
          </div>
          
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge className={getStatusColor(property.status)}>
                {getStatusLabel(property.status)}
              </Badge>
              <Badge variant="outline">{getTypeLabel(property.type)}</Badge>
              {property.featured && (
                <Badge variant="secondary">
                  <Star className="h-3 w-3 mr-1" />
                  Destacada
                </Badge>
              )}
            </div>
            <h3 className="font-semibold text-lg mb-1 group-hover:text-primary transition-colors">
              {property.title}
            </h3>
            <p className="text-sm text-gray-600 flex items-center mb-2">
              <MapPin className="h-3 w-3 mr-1" />
              {property.address}, {property.city}, {property.province}
            </p>
            <p className="text-sm text-gray-600 line-clamp-1">
              {property.description}
            </p>
          </div>
          
          <div className="flex items-center gap-8">
            <div className="grid grid-cols-4 gap-4">
              <div className="flex items-center text-sm text-gray-600">
                <Bed className="h-4 w-4 mr-1" />
                {property.bedrooms}
              </div>
              <div className="flex items-center text-sm text-gray-600">
                <Bath className="h-4 w-4 mr-1" />
                {property.bathrooms}
              </div>
              <div className="flex items-center text-sm text-gray-600">
                <Square className="h-4 w-4 mr-1" />
                {property.area_m2}m²
              </div>
              <div className="flex items-center text-sm text-gray-600">
                <Car className="h-4 w-4 mr-1" />
                {property.parking_spaces}
              </div>
            </div>
            
            <div className="text-right">
              <p className="text-2xl font-bold text-primary">{formatPrice(property.price)}</p>
              <p className="text-sm text-gray-500">
                {formatArea(property.area_m2)}
              </p>
            </div>
            
            <div className="flex gap-2">
              <Button variant="outline" size="sm">
                <Eye className="h-3 w-3 mr-1" />
                Ver
              </Button>
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => handleEditProperty(property)}
              >
                <Edit className="h-3 w-3 mr-1" />
                Editar
              </Button>
              <Button 
                variant="outline" 
                size="sm" 
                className="text-red-600 hover:text-red-700"
                onClick={() => handleDeleteProperty(property)}
              >
                <Trash2 className="h-3 w-3" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className={viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'}>
          {Array.from({ length: 6 }).map((_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-3/4" />
                <Skeleton className="h-4 w-1/2" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-4 w-full mb-2" />
                <Skeleton className="h-4 w-2/3 mb-4" />
                <div className="flex gap-4">
                  <Skeleton className="h-8 w-16" />
                  <Skeleton className="h-8 w-16" />
                  <Skeleton className="h-8 w-16" />
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    const errorMessage = error instanceof Error ? error.message : 'Ha ocurrido un error inesperado';
    const isAuthError = error?.status === 401 || error?.status === 403;
    
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <div className="space-y-4">
            <div className="text-red-600">
              <h3 className="font-medium text-lg mb-2">
                {isAuthError ? 'Error de autenticación' : 'Error al cargar las propiedades'}
              </h3>
              <p className="text-sm text-gray-600">
                {isAuthError ? 
                  'Necesitas iniciar sesión para ver las propiedades' : 
                  errorMessage
                }
              </p>
            </div>
            <div className="space-x-3">
              {isAuthError ? (
                <Button onClick={() => window.location.href = '/login'}>
                  Iniciar sesión
                </Button>
              ) : (
                <>
                  <Button onClick={() => window.location.reload()}>
                    Reintentar
                  </Button>
                  <Button variant="outline" onClick={() => queryClient.invalidateQueries({ queryKey: ['properties'] })}>
                    Refrescar
                  </Button>
                </>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!properties?.data?.length) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <div className="space-y-4">
            <div className="text-gray-400">
              <div className="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
                <Eye className="w-8 h-8" />
              </div>
            </div>
            <div>
              <h3 className="font-medium text-lg mb-2">No se encontraron propiedades</h3>
              <p className="text-sm text-gray-600">
                {searchTerm || Object.values(filters).some(f => f) 
                  ? 'Intenta ajustar los filtros de búsqueda' 
                  : 'No hay propiedades disponibles en este momento'
                }
              </p>
            </div>
            <div className="space-x-3">
              <Button variant="outline" onClick={() => window.location.reload()}>
                Limpiar filtros
              </Button>
              <Button variant="outline" onClick={() => queryClient.invalidateQueries({ queryKey: ['properties'] })}>
                Refrescar
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className={viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'}>
        {properties.data.map((property: Property) => (
          <div key={property.id}>
            {viewMode === 'grid' ? (
              <PropertyCard property={property} />
            ) : (
              <PropertyRow property={property} />
            )}
          </div>
        ))}
      </div>

      {/* Pagination */}
      {properties.total > limit && (
        <div className="flex justify-center gap-2">
          <Button
            variant="outline"
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
          >
            Anterior
          </Button>
          <span className="flex items-center px-4 py-2 text-sm text-gray-600">
            Página {page} de {Math.ceil(properties.total / limit)}
          </span>
          <Button
            variant="outline"
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(properties.total / limit)}
          >
            Siguiente
          </Button>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Eliminar Propiedad</DialogTitle>
            <DialogDescription>
              ¿Estás seguro de que quieres eliminar la propiedad "{selectedProperty?.title}"? 
              Esta acción no se puede deshacer y también eliminará todas las imágenes asociadas.
            </DialogDescription>
          </DialogHeader>
          
          {deletePropertyMutation.error && (
            <div className="bg-red-50 border border-red-200 rounded-md p-3">
              <p className="text-sm text-red-800">
                <strong>Error:</strong> {deletePropertyMutation.error instanceof Error 
                  ? deletePropertyMutation.error.message 
                  : 'No se pudo eliminar la propiedad'}
              </p>
            </div>
          )}
          
          <DialogFooter>
            <Button 
              variant="outline" 
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={deletePropertyMutation.isPending}
            >
              Cancelar
            </Button>
            <Button 
              variant="destructive" 
              onClick={confirmDelete}
              disabled={deletePropertyMutation.isPending}
            >
              {deletePropertyMutation.isPending ? 'Eliminando...' : 'Eliminar'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Edit Property Dialog */}
      <EditPropertyDialog
        property={selectedProperty}
        open={isEditDialogOpen}
        onOpenChange={setIsEditDialogOpen}
        onSuccess={() => {
          setSelectedProperty(null);
          setIsEditDialogOpen(false);
        }}
      />
    </div>
  );
}