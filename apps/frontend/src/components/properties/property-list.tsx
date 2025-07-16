'use client';

import { useQuery } from '@tanstack/react-query';
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
  Calendar
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { apiClient } from '@/lib/api-client';
import { formatPrice, formatArea, formatDate } from '@/lib/utils';

interface Property {
  id: string;
  title: string;
  description: string;
  price: number;
  type: string;
  status: string;
  province: string;
  city: string;
  address: string;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  parking_spaces: number;
  year_built?: number;
  has_garden: boolean;
  has_pool: boolean;
  has_elevator: boolean;
  has_balcony: boolean;
  has_terrace: boolean;
  has_garage: boolean;
  is_furnished: boolean;
  allows_pets: boolean;
  contact_phone: string;
  contact_email: string;
  notes?: string;
  created_at: string;
  updated_at: string;
  is_featured?: boolean;
  images?: string[];
}

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

  const { data: properties, isLoading, error } = useQuery({
    queryKey: ['properties', searchTerm, filters, page, limit],
    queryFn: async () => {
      const params = new URLSearchParams();
      
      if (searchTerm) params.append('search', searchTerm);
      if (filters.type) params.append('type', filters.type);
      if (filters.status) params.append('status', filters.status);
      if (filters.province) params.append('province', filters.province);
      if (filters.city) params.append('city', filters.city);
      if (filters.bedrooms) params.append('bedrooms', filters.bedrooms);
      if (filters.bathrooms) params.append('bathrooms', filters.bathrooms);
      if (filters.minPrice) params.append('min_price', filters.minPrice);
      if (filters.maxPrice) params.append('max_price', filters.maxPrice);
      params.append('page', page.toString());
      params.append('limit', limit.toString());

      const response = await apiClient.get(`/properties?${params.toString()}`);
      return response.data;
    },
  });

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
    <Card className="group hover:shadow-lg transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge className={getStatusColor(property.status)}>
                {getStatusLabel(property.status)}
              </Badge>
              <Badge variant="outline">{getTypeLabel(property.type)}</Badge>
              {property.is_featured && (
                <Badge variant="secondary">
                  <Star className="h-3 w-3 mr-1" />
                  Destacada
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
            <p className="text-2xl font-bold text-primary">{formatPrice(property.price)}</p>
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
          {property.has_pool && (
            <Badge variant="outline" className="text-xs">Piscina</Badge>
          )}
          {property.has_garden && (
            <Badge variant="outline" className="text-xs">Jardín</Badge>
          )}
          {property.has_elevator && (
            <Badge variant="outline" className="text-xs">Ascensor</Badge>
          )}
          {property.is_furnished && (
            <Badge variant="outline" className="text-xs">Amueblado</Badge>
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
            <Button variant="outline" size="sm">
              <Edit className="h-3 w-3 mr-1" />
              Editar
            </Button>
            <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700">
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
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <Badge className={getStatusColor(property.status)}>
                {getStatusLabel(property.status)}
              </Badge>
              <Badge variant="outline">{getTypeLabel(property.type)}</Badge>
              {property.is_featured && (
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
              <Button variant="outline" size="sm">
                <Edit className="h-3 w-3 mr-1" />
                Editar
              </Button>
              <Button variant="outline" size="sm" className="text-red-600 hover:text-red-700">
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
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-red-600 mb-4">Error al cargar las propiedades</p>
          <Button onClick={() => window.location.reload()}>
            Reintentar
          </Button>
        </CardContent>
      </Card>
    );
  }

  if (!properties?.data?.length) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-gray-500 mb-4">No se encontraron propiedades que coincidan con tu búsqueda</p>
          <Button variant="outline" onClick={() => window.location.reload()}>
            Limpiar filtros
          </Button>
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
    </div>
  );
}