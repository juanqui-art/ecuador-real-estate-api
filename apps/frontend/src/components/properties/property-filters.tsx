'use client';

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { 
  ECUADORIAN_PROVINCES, 
  PROPERTY_TYPES, 
  PROPERTY_STATUS, 
  PRICE_RANGES,
  BEDROOM_OPTIONS,
  BATHROOM_OPTIONS
} from '@/lib/constants';

interface PropertyFiltersProps {
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
  onFiltersChange: (filters: any) => void;
}

export function PropertyFilters({ filters, onFiltersChange }: PropertyFiltersProps) {
  const handleFilterChange = (key: string, value: string) => {
    onFiltersChange({
      ...filters,
      [key]: value,
    });
  };

  const clearFilters = () => {
    onFiltersChange({
      type: '',
      status: '',
      minPrice: '',
      maxPrice: '',
      province: '',
      city: '',
      bedrooms: '',
      bathrooms: '',
    });
  };

  const hasActiveFilters = Object.values(filters).some(value => value !== '');

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Tipo de Propiedad */}
        <div>
          <Label htmlFor="type">Tipo de Propiedad</Label>
          <Select value={filters.type} onValueChange={(value) => handleFilterChange('type', value)}>
            <SelectTrigger>
              <SelectValue placeholder="Todos los tipos" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todos los tipos</SelectItem>
              {PROPERTY_TYPES.map(type => (
                <SelectItem key={type.value} value={type.value}>
                  {type.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Estado */}
        <div>
          <Label htmlFor="status">Estado</Label>
          <Select value={filters.status} onValueChange={(value) => handleFilterChange('status', value)}>
            <SelectTrigger>
              <SelectValue placeholder="Todos los estados" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todos los estados</SelectItem>
              {PROPERTY_STATUS.map(status => (
                <SelectItem key={status.value} value={status.value}>
                  {status.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Provincia */}
        <div>
          <Label htmlFor="province">Provincia</Label>
          <Select value={filters.province} onValueChange={(value) => handleFilterChange('province', value)}>
            <SelectTrigger>
              <SelectValue placeholder="Todas las provincias" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Todas las provincias</SelectItem>
              {ECUADORIAN_PROVINCES.map(province => (
                <SelectItem key={province} value={province}>
                  {province}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Ciudad */}
        <div>
          <Label htmlFor="city">Ciudad</Label>
          <Input
            id="city"
            placeholder="Ingresa la ciudad"
            value={filters.city}
            onChange={(e) => handleFilterChange('city', e.target.value)}
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Precio Mínimo */}
        <div>
          <Label htmlFor="minPrice">Precio Mínimo (USD)</Label>
          <Input
            id="minPrice"
            type="number"
            placeholder="0"
            value={filters.minPrice}
            onChange={(e) => handleFilterChange('minPrice', e.target.value)}
          />
        </div>

        {/* Precio Máximo */}
        <div>
          <Label htmlFor="maxPrice">Precio Máximo (USD)</Label>
          <Input
            id="maxPrice"
            type="number"
            placeholder="Sin límite"
            value={filters.maxPrice}
            onChange={(e) => handleFilterChange('maxPrice', e.target.value)}
          />
        </div>

        {/* Dormitorios */}
        <div>
          <Label htmlFor="bedrooms">Dormitorios</Label>
          <Select value={filters.bedrooms} onValueChange={(value) => handleFilterChange('bedrooms', value)}>
            <SelectTrigger>
              <SelectValue placeholder="Cualquier cantidad" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Cualquier cantidad</SelectItem>
              {BEDROOM_OPTIONS.map(option => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Baños */}
        <div>
          <Label htmlFor="bathrooms">Baños</Label>
          <Select value={filters.bathrooms} onValueChange={(value) => handleFilterChange('bathrooms', value)}>
            <SelectTrigger>
              <SelectValue placeholder="Cualquier cantidad" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">Cualquier cantidad</SelectItem>
              {BATHROOM_OPTIONS.map(option => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Rangos de Precio Predefinidos */}
      <div>
        <Label>Rangos de Precio</Label>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2 mt-2">
          {PRICE_RANGES.map(range => (
            <Button
              key={range.value}
              variant="outline"
              size="sm"
              className="justify-start"
              onClick={() => {
                const [min, max] = range.value.split('-');
                handleFilterChange('minPrice', min === '0' ? '' : min);
                handleFilterChange('maxPrice', max === '+' ? '' : max || '');
              }}
            >
              {range.label}
            </Button>
          ))}
        </div>
      </div>

      {/* Acciones */}
      <div className="flex justify-between items-center pt-4 border-t">
        <div className="text-sm text-gray-500">
          {hasActiveFilters && (
            <span>Filtros activos aplicados</span>
          )}
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={clearFilters} disabled={!hasActiveFilters}>
            Limpiar Filtros
          </Button>
          <Button onClick={() => {/* Trigger search */}}>
            Aplicar Filtros
          </Button>
        </div>
      </div>
    </div>
  );
}