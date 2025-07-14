'use client';

import { useState } from 'react';
import { useForm } from '@tanstack/react-form';
import { zodValidator } from '@tanstack/zod-form-adapter';
import { motion } from 'motion/react';
import { 
  Search, 
  Filter, 
  X, 
  MapPin, 
  Home, 
  Building, 
  TreePine, 
  Store,
  DollarSign,
  Calendar,
  SlidersHorizontal,
  ChevronDown
} from 'lucide-react';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { 
  Collapsible, 
  CollapsibleContent, 
  CollapsibleTrigger 
} from '@/components/ui/collapsible';

// Validation schema
const FilterSchema = z.object({
  search: z.string().optional(),
  province: z.string().optional(),
  city: z.string().optional(),
  type: z.string().optional(),
  status: z.string().optional(),
  min_price: z.number().min(0).optional(),
  max_price: z.number().min(0).optional(),
  min_bedrooms: z.number().min(0).optional(),
  max_bedrooms: z.number().min(0).optional(),
  min_bathrooms: z.number().min(0).optional(),
  max_bathrooms: z.number().min(0).optional(),
  min_area: z.number().min(0).optional(),
  max_area: z.number().min(0).optional(),
  features: z.array(z.string()).optional(),
  date_range: z.string().optional(),
});

type FilterFormData = z.infer<typeof FilterSchema>;

// Ecuador provinces
const PROVINCES = [
  'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi',
  'El Oro', 'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja',
  'Los Ríos', 'Manabí', 'Morona Santiago', 'Napo', 'Orellana', 'Pastaza',
  'Pichincha', 'Santa Elena', 'Santo Domingo', 'Sucumbíos', 'Tungurahua',
  'Zamora Chinchipe'
];

// Property types
const PROPERTY_TYPES = [
  { value: 'house', label: 'Casa', icon: Home },
  { value: 'apartment', label: 'Apartamento', icon: Building },
  { value: 'land', label: 'Terreno', icon: TreePine },
  { value: 'commercial', label: 'Comercial', icon: Store },
];

// Property status
const PROPERTY_STATUS = [
  { value: 'available', label: 'Disponible' },
  { value: 'sold', label: 'Vendido' },
  { value: 'rented', label: 'Alquilado' },
  { value: 'reserved', label: 'Reservado' },
];

// Features
const FEATURES = [
  'pool', 'garden', 'garage', 'security', 'elevator', 'terrace',
  'balcony', 'furnished', 'air_conditioning', 'heating'
];

const FEATURE_LABELS: Record<string, string> = {
  pool: 'Piscina',
  garden: 'Jardín',
  garage: 'Garaje',
  security: 'Seguridad',
  elevator: 'Ascensor',
  terrace: 'Terraza',
  balcony: 'Balcón',
  furnished: 'Amoblado',
  air_conditioning: 'Aire Acondicionado',
  heating: 'Calefacción',
};

interface PropertyFiltersProps {
  onFiltersChange: (filters: FilterFormData) => void;
  initialFilters?: Partial<FilterFormData>;
  compact?: boolean;
}

export function PropertyFilters({ 
  onFiltersChange, 
  initialFilters = {}, 
  compact = false 
}: PropertyFiltersProps) {
  const [isExpanded, setIsExpanded] = useState(!compact);
  const [activeFilters, setActiveFilters] = useState<string[]>([]);

  const form = useForm({
    defaultValues: {
      search: '',
      province: '',
      city: '',
      type: '',
      status: '',
      min_price: undefined,
      max_price: undefined,
      min_bedrooms: undefined,
      max_bedrooms: undefined,
      min_bathrooms: undefined,
      max_bathrooms: undefined,
      min_area: undefined,
      max_area: undefined,
      features: [],
      date_range: '',
      ...initialFilters,
    } as FilterFormData,
    validatorAdapter: zodValidator(),
    validators: {
      onChange: FilterSchema,
    },
    onSubmit: async ({ value }) => {
      onFiltersChange(value);
      updateActiveFilters(value);
    },
  });

  const updateActiveFilters = (filters: FilterFormData) => {
    const active: string[] = [];
    
    if (filters.search) active.push(`Búsqueda: ${filters.search}`);
    if (filters.province) active.push(`Provincia: ${filters.province}`);
    if (filters.city) active.push(`Ciudad: ${filters.city}`);
    if (filters.type) {
      const typeLabel = PROPERTY_TYPES.find(t => t.value === filters.type)?.label;
      active.push(`Tipo: ${typeLabel}`);
    }
    if (filters.status) {
      const statusLabel = PROPERTY_STATUS.find(s => s.value === filters.status)?.label;
      active.push(`Estado: ${statusLabel}`);
    }
    if (filters.min_price) active.push(`Precio mín: $${filters.min_price.toLocaleString()}`);
    if (filters.max_price) active.push(`Precio máx: $${filters.max_price.toLocaleString()}`);
    if (filters.min_bedrooms) active.push(`Dormitorios mín: ${filters.min_bedrooms}`);
    if (filters.max_bedrooms) active.push(`Dormitorios máx: ${filters.max_bedrooms}`);
    if (filters.min_area) active.push(`Área mín: ${filters.min_area}m²`);
    if (filters.max_area) active.push(`Área máx: ${filters.max_area}m²`);
    if (filters.features && filters.features.length > 0) {
      active.push(`Características: ${filters.features.length}`);
    }

    setActiveFilters(active);
  };

  const clearFilters = () => {
    form.reset();
    setActiveFilters([]);
    onFiltersChange({});
  };

  const removeFilter = (filterText: string) => {
    // Logic to remove specific filter
    const currentValues = form.getFieldValue('');
    // Implementation depends on the filter type
    form.handleSubmit();
  };

  if (compact) {
    return (
      <Card className="mb-6">
        <CardHeader className="pb-3">
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg flex items-center gap-2">
              <Filter className="h-5 w-5" />
              Filtros
            </CardTitle>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsExpanded(!isExpanded)}
            >
              <ChevronDown className={`h-4 w-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`} />
            </Button>
          </div>
        </CardHeader>
        
        <Collapsible open={isExpanded} onOpenChange={setIsExpanded}>
          <CollapsibleContent>
            <CardContent className="pt-0">
              <form.Provider>
                <form
                  onSubmit={(e) => {
                    e.preventDefault();
                    form.handleSubmit();
                  }}
                  className="space-y-4"
                >
                  {/* Quick Search */}
                  <div className="relative">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <form.Field
                      name="search"
                      children={(field) => (
                        <Input
                          placeholder="Buscar propiedades..."
                          className="pl-10"
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                        />
                      )}
                    />
                  </div>

                  {/* Basic Filters Row */}
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                    <form.Field
                      name="province"
                      children={(field) => (
                        <Select
                          value={field.state.value}
                          onValueChange={field.handleChange}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Provincia" />
                          </SelectTrigger>
                          <SelectContent>
                            {PROVINCES.map((province) => (
                              <SelectItem key={province} value={province}>
                                {province}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />

                    <form.Field
                      name="type"
                      children={(field) => (
                        <Select
                          value={field.state.value}
                          onValueChange={field.handleChange}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Tipo" />
                          </SelectTrigger>
                          <SelectContent>
                            {PROPERTY_TYPES.map((type) => (
                              <SelectItem key={type.value} value={type.value}>
                                <div className="flex items-center gap-2">
                                  <type.icon className="h-4 w-4" />
                                  {type.label}
                                </div>
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />

                    <form.Field
                      name="status"
                      children={(field) => (
                        <Select
                          value={field.state.value}
                          onValueChange={field.handleChange}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Estado" />
                          </SelectTrigger>
                          <SelectContent>
                            {PROPERTY_STATUS.map((status) => (
                              <SelectItem key={status.value} value={status.value}>
                                {status.label}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                      )}
                    />

                    <div className="flex gap-2">
                      <Button type="submit" className="flex-1">
                        <Filter className="mr-2 h-4 w-4" />
                        Filtrar
                      </Button>
                      <Button 
                        type="button" 
                        variant="outline"
                        onClick={clearFilters}
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  {/* Price Range */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <form.Field
                      name="min_price"
                      children={(field) => (
                        <div className="relative">
                          <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                          <Input
                            type="number"
                            placeholder="Precio mínimo"
                            className="pl-10"
                            value={field.state.value || ''}
                            onChange={(e) => field.handleChange(Number(e.target.value) || undefined)}
                          />
                        </div>
                      )}
                    />

                    <form.Field
                      name="max_price"
                      children={(field) => (
                        <div className="relative">
                          <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                          <Input
                            type="number"
                            placeholder="Precio máximo"
                            className="pl-10"
                            value={field.state.value || ''}
                            onChange={(e) => field.handleChange(Number(e.target.value) || undefined)}
                          />
                        </div>
                      )}
                    />
                  </div>
                </form>
              </form.Provider>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Active Filters */}
      {activeFilters.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="flex flex-wrap gap-2"
        >
          {activeFilters.map((filter, index) => (
            <Badge
              key={index}
              variant="secondary"
              className="flex items-center gap-1"
            >
              {filter}
              <Button
                variant="ghost"
                size="sm"
                className="h-auto p-0 ml-1"
                onClick={() => removeFilter(filter)}
              >
                <X className="h-3 w-3" />
              </Button>
            </Badge>
          ))}
          <Button
            variant="outline"
            size="sm"
            onClick={clearFilters}
            className="h-6"
          >
            Limpiar todo
          </Button>
        </motion.div>
      )}

      {/* Main Filter Form */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <SlidersHorizontal className="h-5 w-5" />
            Filtros Avanzados
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form.Provider>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
              }}
              className="space-y-6"
            >
              {/* Search */}
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <form.Field
                  name="search"
                  children={(field) => (
                    <Input
                      placeholder="Buscar por título, descripción, ubicación..."
                      className="pl-10"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                    />
                  )}
                />
              </div>

              <Separator />

              {/* Location */}
              <div className="space-y-4">
                <h4 className="font-medium flex items-center gap-2">
                  <MapPin className="h-4 w-4" />
                  Ubicación
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <form.Field
                    name="province"
                    children={(field) => (
                      <Select
                        value={field.state.value}
                        onValueChange={field.handleChange}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Seleccionar provincia" />
                        </SelectTrigger>
                        <SelectContent>
                          {PROVINCES.map((province) => (
                            <SelectItem key={province} value={province}>
                              {province}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                  />

                  <form.Field
                    name="city"
                    children={(field) => (
                      <Input
                        placeholder="Ciudad"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                      />
                    )}
                  />
                </div>
              </div>

              <Separator />

              {/* Property Type & Status */}
              <div className="space-y-4">
                <h4 className="font-medium">Tipo y Estado</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <form.Field
                    name="type"
                    children={(field) => (
                      <Select
                        value={field.state.value}
                        onValueChange={field.handleChange}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Tipo de propiedad" />
                        </SelectTrigger>
                        <SelectContent>
                          {PROPERTY_TYPES.map((type) => (
                            <SelectItem key={type.value} value={type.value}>
                              <div className="flex items-center gap-2">
                                <type.icon className="h-4 w-4" />
                                {type.label}
                              </div>
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                  />

                  <form.Field
                    name="status"
                    children={(field) => (
                      <Select
                        value={field.state.value}
                        onValueChange={field.handleChange}
                      >
                        <SelectTrigger>
                          <SelectValue placeholder="Estado" />
                        </SelectTrigger>
                        <SelectContent>
                          {PROPERTY_STATUS.map((status) => (
                            <SelectItem key={status.value} value={status.value}>
                              {status.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    )}
                  />
                </div>
              </div>

              <Separator />

              {/* Price Range */}
              <div className="space-y-4">
                <h4 className="font-medium flex items-center gap-2">
                  <DollarSign className="h-4 w-4" />
                  Rango de Precio
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <form.Field
                    name="min_price"
                    children={(field) => (
                      <div className="relative">
                        <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <Input
                          type="number"
                          placeholder="Precio mínimo"
                          className="pl-10"
                          value={field.state.value || ''}
                          onChange={(e) => field.handleChange(Number(e.target.value) || undefined)}
                        />
                      </div>
                    )}
                  />

                  <form.Field
                    name="max_price"
                    children={(field) => (
                      <div className="relative">
                        <DollarSign className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <Input
                          type="number"
                          placeholder="Precio máximo"
                          className="pl-10"
                          value={field.state.value || ''}
                          onChange={(e) => field.handleChange(Number(e.target.value) || undefined)}
                        />
                      </div>
                    )}
                  />
                </div>
              </div>

              {/* Actions */}
              <div className="flex gap-4 pt-4">
                <Button type="submit" className="flex-1">
                  <Filter className="mr-2 h-4 w-4" />
                  Aplicar Filtros
                </Button>
                <Button 
                  type="button" 
                  variant="outline"
                  onClick={clearFilters}
                >
                  <X className="mr-2 h-4 w-4" />
                  Limpiar
                </Button>
              </div>
            </form>
          </form.Provider>
        </CardContent>
      </Card>
    </div>
  );
}