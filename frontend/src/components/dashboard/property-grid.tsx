'use client';

import { useState } from 'react';
import { motion, AnimatePresence } from 'motion/react';
import { useProperties } from '@/hooks/useProperties';
import { 
  MapPin, 
  Bed, 
  Bath, 
  Square, 
  Car, 
  Eye, 
  Heart, 
  MoreHorizontal,
  Edit,
  Trash2,
  Share2,
  Calendar,
  DollarSign,
  Building,
  Home,
  TreePine,
  Store
} from 'lucide-react';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { 
  DropdownMenu, 
  DropdownMenuContent, 
  DropdownMenuItem, 
  DropdownMenuTrigger 
} from '@/components/ui/dropdown-menu';
import { useAuthStore } from '@/store/auth';

// Mock data - will be replaced with real API data
const mockProperties = [
  {
    id: '1',
    title: 'Hermosa casa en Samborondón con piscina',
    description: 'Casa moderna de 3 pisos con acabados de lujo y piscina privada',
    price: 285000,
    province: 'Guayas',
    city: 'Samborondón',
    type: 'house',
    status: 'available',
    bedrooms: 4,
    bathrooms: 3.5,
    area_m2: 320,
    parking_spaces: 2,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300', '/api/placeholder/400/300'],
    created_at: '2024-01-15',
    features: ['pool', 'garden', 'security'],
    views: 156,
    is_featured: true,
  },
  {
    id: '2',
    title: 'Apartamento moderno en Quito Norte',
    description: 'Apartamento de 2 habitaciones en zona exclusiva',
    price: 175000,
    province: 'Pichincha',
    city: 'Quito',
    type: 'apartment',
    status: 'available',
    bedrooms: 2,
    bathrooms: 2,
    area_m2: 85,
    parking_spaces: 1,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300'],
    created_at: '2024-01-10',
    features: ['elevator', 'security'],
    views: 89,
    is_featured: false,
  },
  {
    id: '3',
    title: 'Terreno en Cuenca para construcción',
    description: 'Lote de terreno ideal para construcción residencial',
    price: 45000,
    province: 'Azuay',
    city: 'Cuenca',
    type: 'land',
    status: 'available',
    bedrooms: 0,
    bathrooms: 0,
    area_m2: 500,
    parking_spaces: 0,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300'],
    created_at: '2024-01-05',
    features: [],
    views: 34,
    is_featured: false,
  },
  {
    id: '4',
    title: 'Casa en Guayaquil cerca del Malecón',
    description: 'Casa familiar en ubicación privilegiada',
    price: 320000,
    province: 'Guayas',
    city: 'Guayaquil',
    type: 'house',
    status: 'sold',
    bedrooms: 3,
    bathrooms: 2,
    area_m2: 180,
    parking_spaces: 1,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300'],
    created_at: '2024-01-01',
    features: ['garage', 'terrace'],
    views: 203,
    is_featured: true,
  },
  {
    id: '5',
    title: 'Departamento en Ambato centro',
    description: 'Departamento céntrico con vista a la ciudad',
    price: 95000,
    province: 'Tungurahua',
    city: 'Ambato',
    type: 'apartment',
    status: 'available',
    bedrooms: 2,
    bathrooms: 1,
    area_m2: 60,
    parking_spaces: 0,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300'],
    created_at: '2023-12-28',
    features: ['balcony'],
    views: 67,
    is_featured: false,
  },
  {
    id: '6',
    title: 'Local comercial en Manta',
    description: 'Local comercial en zona de alto tráfico',
    price: 150000,
    province: 'Manabí',
    city: 'Manta',
    type: 'commercial',
    status: 'available',
    bedrooms: 0,
    bathrooms: 2,
    area_m2: 120,
    parking_spaces: 2,
    main_image: '/api/placeholder/400/300',
    images: ['/api/placeholder/400/300'],
    created_at: '2023-12-20',
    features: ['security'],
    views: 112,
    is_featured: false,
  },
];

interface PropertyCardProps {
  property: typeof mockProperties[0];
  index: number;
  canEdit: boolean;
  onEdit?: (id: string) => void;
  onDelete?: (id: string) => void;
  onToggleFavorite?: (id: string) => void;
}

function PropertyCard({ 
  property, 
  index, 
  canEdit, 
  onEdit, 
  onDelete, 
  onToggleFavorite 
}: PropertyCardProps) {
  const [isLiked, setIsLiked] = useState(false);
  const [imageError, setImageError] = useState(false);

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'house': return <Home className="h-4 w-4" />;
      case 'apartment': return <Building className="h-4 w-4" />;
      case 'land': return <TreePine className="h-4 w-4" />;
      case 'commercial': return <Store className="h-4 w-4" />;
      default: return <Building className="h-4 w-4" />;
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'house': return 'Casa';
      case 'apartment': return 'Apartamento';
      case 'land': return 'Terreno';
      case 'commercial': return 'Comercial';
      default: return 'Propiedad';
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'available':
        return <Badge className="bg-green-100 text-green-800">Disponible</Badge>;
      case 'sold':
        return <Badge className="bg-red-100 text-red-800">Vendido</Badge>;
      case 'rented':
        return <Badge className="bg-blue-100 text-blue-800">Alquilado</Badge>;
      case 'reserved':
        return <Badge className="bg-yellow-100 text-yellow-800">Reservado</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('es-EC', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
    }).format(price);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('es-EC', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4, delay: index * 0.1 }}
      whileHover={{ y: -5 }}
      className="group"
    >
      <Card className="overflow-hidden transition-all duration-300 hover:shadow-xl">
        {/* Image Section */}
        <div className="relative aspect-video overflow-hidden">
          {!imageError ? (
            <img
              src={property.main_image}
              alt={property.title}
              className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
              onError={() => setImageError(true)}
            />
          ) : (
            <div className="w-full h-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
              {getTypeIcon(property.type)}
            </div>
          )}
          
          {/* Overlays */}
          <div className="absolute top-3 left-3 flex gap-2">
            {property.is_featured && (
              <Badge className="bg-yellow-500 text-white">Destacado</Badge>
            )}
            {getStatusBadge(property.status)}
          </div>

          <div className="absolute top-3 right-3 flex gap-2">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 bg-white/80 hover:bg-white"
              onClick={() => {
                setIsLiked(!isLiked);
                onToggleFavorite?.(property.id);
              }}
            >
              <Heart
                className={`h-4 w-4 ${
                  isLiked ? 'fill-red-500 text-red-500' : 'text-gray-600'
                }`}
              />
            </Button>
            
            {canEdit && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 bg-white/80 hover:bg-white"
                  >
                    <MoreHorizontal className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuItem onClick={() => onEdit?.(property.id)}>
                    <Edit className="mr-2 h-4 w-4" />
                    Editar
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Share2 className="mr-2 h-4 w-4" />
                    Compartir
                  </DropdownMenuItem>
                  <DropdownMenuItem 
                    onClick={() => onDelete?.(property.id)}
                    className="text-red-600"
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Eliminar
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>

          {/* Price */}
          <div className="absolute bottom-3 left-3">
            <div className="bg-white/90 backdrop-blur-sm rounded-md px-3 py-1">
              <span className="text-lg font-bold text-gray-900">
                {formatPrice(property.price)}
              </span>
            </div>
          </div>
        </div>

        {/* Content Section */}
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="flex-1">
              <CardTitle className="text-lg line-clamp-1">
                {property.title}
              </CardTitle>
              <CardDescription className="flex items-center mt-1">
                <MapPin className="h-4 w-4 mr-1" />
                {property.city}, {property.province}
              </CardDescription>
            </div>
            <div className="flex items-center gap-1 text-sm text-gray-500">
              {getTypeIcon(property.type)}
              <span>{getTypeLabel(property.type)}</span>
            </div>
          </div>
        </CardHeader>

        <CardContent className="pt-0">
          <p className="text-sm text-gray-600 line-clamp-2 mb-4">
            {property.description}
          </p>

          {/* Property Details */}
          <div className="flex items-center gap-4 text-sm text-gray-500 mb-4">
            {property.bedrooms > 0 && (
              <div className="flex items-center gap-1">
                <Bed className="h-4 w-4" />
                <span>{property.bedrooms}</span>
              </div>
            )}
            {property.bathrooms > 0 && (
              <div className="flex items-center gap-1">
                <Bath className="h-4 w-4" />
                <span>{property.bathrooms}</span>
              </div>
            )}
            <div className="flex items-center gap-1">
              <Square className="h-4 w-4" />
              <span>{property.area_m2} m²</span>
            </div>
            {property.parking_spaces > 0 && (
              <div className="flex items-center gap-1">
                <Car className="h-4 w-4" />
                <span>{property.parking_spaces}</span>
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="flex items-center justify-between text-sm text-gray-500">
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-1">
                <Eye className="h-4 w-4" />
                <span>{property.views}</span>
              </div>
              <div className="flex items-center gap-1">
                <Calendar className="h-4 w-4" />
                <span>{formatDate(property.created_at)}</span>
              </div>
            </div>
            <Button variant="outline" size="sm">
              Ver Detalles
            </Button>
          </div>
        </CardContent>
      </Card>
    </motion.div>
  );
}

interface PropertyGridProps {
  showFilters?: boolean;
  limit?: number;
}

export function PropertyGrid({ showFilters = true, limit }: PropertyGridProps) {
  const { user } = useAuthStore();
  const [filters, setFilters] = useState({});
  
  // Use real API data or fallback to mock data
  const { data: apiData, isLoading, error } = useProperties(filters);
  const properties = apiData?.properties || mockProperties;

  const canEdit = user?.role === 'admin' || user?.role === 'agency' || user?.role === 'agent';

  const displayProperties = limit ? properties.slice(0, limit) : properties;

  const handleEdit = (id: string) => {
    console.log('Edit property:', id);
  };

  const handleDelete = (id: string) => {
    console.log('Delete property:', id);
  };

  const handleToggleFavorite = (id: string) => {
    console.log('Toggle favorite:', id);
  };

  // Loading state
  if (isLoading) {
    return (
      <div className="space-y-6">
        {showFilters && (
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-gray-900">Propiedades</h2>
              <p className="text-sm text-gray-600">Cargando propiedades...</p>
            </div>
          </div>
        )}
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: limit || 6 }).map((_, index) => (
            <div key={index} className="animate-pulse">
              <div className="bg-gray-200 aspect-video rounded-lg mb-4"></div>
              <div className="space-y-2">
                <div className="h-4 bg-gray-200 rounded w-3/4"></div>
                <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                <div className="h-3 bg-gray-200 rounded w-full"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="space-y-6">
        <div className="text-center py-12">
          <Building className="h-16 w-16 text-red-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            Error al cargar propiedades
          </h3>
          <p className="text-gray-600 mb-4">
            Hubo un problema al cargar las propiedades. Intenta de nuevo.
          </p>
          <Button onClick={() => window.location.reload()}>
            Reintentar
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {showFilters && (
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">Propiedades</h2>
            <p className="text-sm text-gray-600">
              {apiData?.total || displayProperties.length} propiedades encontradas
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm">
              Filtros
            </Button>
            <Button variant="outline" size="sm">
              Ordenar
            </Button>
          </div>
        </div>
      )}

      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        <AnimatePresence>
          {displayProperties.map((property, index) => (
            <PropertyCard
              key={property.id}
              property={property}
              index={index}
              canEdit={canEdit}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onToggleFavorite={handleToggleFavorite}
            />
          ))}
        </AnimatePresence>
      </motion.div>

      {displayProperties.length === 0 && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center py-12"
        >
          <Building className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            No hay propiedades disponibles
          </h3>
          <p className="text-gray-600 mb-4">
            Parece que no hay propiedades que coincidan con tus criterios.
          </p>
          <Button>
            <DollarSign className="mr-2 h-4 w-4" />
            Publicar Propiedad
          </Button>
        </motion.div>
      )}
    </div>
  );
}