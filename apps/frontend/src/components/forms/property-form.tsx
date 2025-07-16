'use client';

import { useForm } from '@tanstack/react-form';
import { zodValidator } from '@tanstack/zod-form-adapter';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { apiClient } from '@/lib/api-client';
import { ECUADORIAN_PROVINCES, PROPERTY_TYPES, PROPERTY_STATUS } from '@/lib/constants';
import { formatPrice } from '@/lib/utils';
import { PropertyImageManager } from '@/components/images/property-image-manager';

const propertySchema = z.object({
  title: z.string().min(10, 'El título debe tener al menos 10 caracteres'),
  description: z.string().min(50, 'La descripción debe tener al menos 50 caracteres'),
  price: z.number().min(1000, 'El precio debe ser mayor a $1,000'),
  type: z.enum(['house', 'apartment', 'land', 'commercial'], {
    errorMap: () => ({ message: 'Selecciona un tipo de propiedad' })
  }),
  status: z.enum(['available', 'sold', 'rented'], {
    errorMap: () => ({ message: 'Selecciona un estado' })
  }),
  province: z.string().min(1, 'Selecciona una provincia'),
  city: z.string().min(2, 'Ingresa la ciudad'),
  address: z.string().min(10, 'Ingresa la dirección completa'),
  bedrooms: z.number().min(0, 'Número de dormitorios inválido').max(20, 'Máximo 20 dormitorios'),
  bathrooms: z.number().min(0, 'Número de baños inválido').max(20, 'Máximo 20 baños'),
  area_m2: z.number().min(10, 'El área debe ser mayor a 10 m²').max(10000, 'Máximo 10,000 m²'),
  parking_spaces: z.number().min(0, 'Número de parqueaderos inválido').max(20, 'Máximo 20 parqueaderos'),
  year_built: z.number().min(1900, 'Año inválido').max(new Date().getFullYear(), 'Año no puede ser futuro').optional(),
  has_garden: z.boolean().default(false),
  has_pool: z.boolean().default(false),
  has_elevator: z.boolean().default(false),
  has_balcony: z.boolean().default(false),
  has_terrace: z.boolean().default(false),
  has_garage: z.boolean().default(false),
  is_furnished: z.boolean().default(false),
  allows_pets: z.boolean().default(false),
  contact_phone: z.string().min(10, 'Ingresa un teléfono válido'),
  contact_email: z.string().email('Ingresa un email válido'),
  notes: z.string().optional(),
});

type PropertyFormData = z.infer<typeof propertySchema>;

interface PropertyFormProps {
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: Partial<PropertyFormData>;
  propertyId?: string; // Para edición
}

export function PropertyForm({ onSuccess, onCancel, initialData, propertyId }: PropertyFormProps) {
  const queryClient = useQueryClient();

  const createPropertyMutation = useMutation({
    mutationFn: async (data: PropertyFormData) => {
      const response = await apiClient.post('/properties', data);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      onSuccess?.();
    },
  });

  const form = useForm({
    defaultValues: {
      title: initialData?.title || '',
      description: initialData?.description || '',
      price: initialData?.price || 0,
      type: initialData?.type || 'house',
      status: initialData?.status || 'available',
      province: initialData?.province || '',
      city: initialData?.city || '',
      address: initialData?.address || '',
      bedrooms: initialData?.bedrooms || 1,
      bathrooms: initialData?.bathrooms || 1,
      area_m2: initialData?.area_m2 || 0,
      parking_spaces: initialData?.parking_spaces || 0,
      year_built: initialData?.year_built || undefined,
      has_garden: initialData?.has_garden || false,
      has_pool: initialData?.has_pool || false,
      has_elevator: initialData?.has_elevator || false,
      has_balcony: initialData?.has_balcony || false,
      has_terrace: initialData?.has_terrace || false,
      has_garage: initialData?.has_garage || false,
      is_furnished: initialData?.is_furnished || false,
      allows_pets: initialData?.allows_pets || false,
      contact_phone: initialData?.contact_phone || '',
      contact_email: initialData?.contact_email || '',
      notes: initialData?.notes || '',
    },
    onSubmit: async ({ value }) => {
      await createPropertyMutation.mutateAsync(value);
    },
    validatorAdapter: zodValidator,
  });

  const features = [
    { name: 'has_garden', label: 'Jardín' },
    { name: 'has_pool', label: 'Piscina' },
    { name: 'has_elevator', label: 'Ascensor' },
    { name: 'has_balcony', label: 'Balcón' },
    { name: 'has_terrace', label: 'Terraza' },
    { name: 'has_garage', label: 'Garaje' },
    { name: 'is_furnished', label: 'Amueblado' },
    { name: 'allows_pets', label: 'Permite mascotas' },
  ];

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      form.handleSubmit();
    }}>
      <div className="space-y-6">
        {/* Información Básica */}
        <Card>
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4">Información Básica</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <form.Field
                name="title"
                validators={{ onChange: propertySchema.shape.title }}
                children={(field) => (
                  <div className="col-span-full">
                    <Label htmlFor="title">Título de la propiedad *</Label>
                    <Input
                      id="title"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="Ej: Hermosa casa en Samborondón con piscina"
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="type"
                validators={{ onChange: propertySchema.shape.type }}
                children={(field) => (
                  <div>
                    <Label htmlFor="type">Tipo de propiedad *</Label>
                    <Select value={field.state.value} onValueChange={field.handleChange}>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecciona el tipo" />
                      </SelectTrigger>
                      <SelectContent>
                        {PROPERTY_TYPES.map(type => (
                          <SelectItem key={type.value} value={type.value}>
                            {type.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="status"
                validators={{ onChange: propertySchema.shape.status }}
                children={(field) => (
                  <div>
                    <Label htmlFor="status">Estado *</Label>
                    <Select value={field.state.value} onValueChange={field.handleChange}>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecciona el estado" />
                      </SelectTrigger>
                      <SelectContent>
                        {PROPERTY_STATUS.map(status => (
                          <SelectItem key={status.value} value={status.value}>
                            {status.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="price"
                validators={{ onChange: propertySchema.shape.price }}
                children={(field) => (
                  <div>
                    <Label htmlFor="price">Precio (USD) *</Label>
                    <Input
                      id="price"
                      type="number"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(Number(e.target.value))}
                      placeholder="285000"
                    />
                    {field.state.value > 0 && (
                      <p className="text-sm text-gray-500 mt-1">
                        {formatPrice(field.state.value)}
                      </p>
                    )}
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />
            </div>

            <form.Field
              name="description"
              validators={{ onChange: propertySchema.shape.description }}
              children={(field) => (
                <div className="mt-4">
                  <Label htmlFor="description">Descripción *</Label>
                  <Textarea
                    id="description"
                    value={field.state.value}
                    onChange={(e) => field.handleChange(e.target.value)}
                    placeholder="Describe las características principales de la propiedad..."
                    rows={4}
                  />
                  <p className="text-sm text-gray-500 mt-1">
                    {field.state.value.length}/500 caracteres
                  </p>
                  {field.state.meta.errors.map((error) => (
                    <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                  ))}
                </div>
              )}
            />
          </CardContent>
        </Card>

        {/* Ubicación */}
        <Card>
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4">Ubicación</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <form.Field
                name="province"
                validators={{ onChange: propertySchema.shape.province }}
                children={(field) => (
                  <div>
                    <Label htmlFor="province">Provincia *</Label>
                    <Select value={field.state.value} onValueChange={field.handleChange}>
                      <SelectTrigger>
                        <SelectValue placeholder="Selecciona la provincia" />
                      </SelectTrigger>
                      <SelectContent>
                        {ECUADORIAN_PROVINCES.map(province => (
                          <SelectItem key={province} value={province}>
                            {province}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="city"
                validators={{ onChange: propertySchema.shape.city }}
                children={(field) => (
                  <div>
                    <Label htmlFor="city">Ciudad *</Label>
                    <Input
                      id="city"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="Ej: Samborondón"
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="address"
                validators={{ onChange: propertySchema.shape.address }}
                children={(field) => (
                  <div className="col-span-full">
                    <Label htmlFor="address">Dirección completa *</Label>
                    <Input
                      id="address"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="Ej: Km 2.5 Vía Samborondón, Urbanización La Puntilla"
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />
            </div>
          </CardContent>
        </Card>

        {/* Características */}
        <Card>
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4">Características</h3>
            
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <form.Field
                name="bedrooms"
                validators={{ onChange: propertySchema.shape.bedrooms }}
                children={(field) => (
                  <div>
                    <Label htmlFor="bedrooms">Dormitorios *</Label>
                    <Input
                      id="bedrooms"
                      type="number"
                      min="0"
                      max="20"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(Number(e.target.value))}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="bathrooms"
                validators={{ onChange: propertySchema.shape.bathrooms }}
                children={(field) => (
                  <div>
                    <Label htmlFor="bathrooms">Baños *</Label>
                    <Input
                      id="bathrooms"
                      type="number"
                      min="0"
                      max="20"
                      step="0.5"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(Number(e.target.value))}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="area_m2"
                validators={{ onChange: propertySchema.shape.area_m2 }}
                children={(field) => (
                  <div>
                    <Label htmlFor="area_m2">Área (m²) *</Label>
                    <Input
                      id="area_m2"
                      type="number"
                      min="10"
                      max="10000"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(Number(e.target.value))}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="parking_spaces"
                validators={{ onChange: propertySchema.shape.parking_spaces }}
                children={(field) => (
                  <div>
                    <Label htmlFor="parking_spaces">Parqueaderos</Label>
                    <Input
                      id="parking_spaces"
                      type="number"
                      min="0"
                      max="20"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(Number(e.target.value))}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="year_built"
                validators={{ onChange: propertySchema.shape.year_built }}
                children={(field) => (
                  <div>
                    <Label htmlFor="year_built">Año de construcción</Label>
                    <Input
                      id="year_built"
                      type="number"
                      min="1900"
                      max={new Date().getFullYear()}
                      value={field.state.value || ''}
                      onChange={(e) => field.handleChange(e.target.value ? Number(e.target.value) : undefined)}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />
            </div>
          </CardContent>
        </Card>

        {/* Características Adicionales */}
        <Card>
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4">Características Adicionales</h3>
            
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {features.map((feature) => (
                <form.Field
                  key={feature.name}
                  name={feature.name as keyof PropertyFormData}
                  children={(field) => (
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        id={feature.name}
                        checked={field.state.value as boolean}
                        onChange={(e) => field.handleChange(e.target.checked)}
                        className="rounded border-gray-300"
                      />
                      <Label htmlFor={feature.name} className="text-sm font-normal">
                        {feature.label}
                      </Label>
                    </div>
                  )}
                />
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Información de Contacto */}
        <Card>
          <CardContent className="pt-6">
            <h3 className="text-lg font-semibold mb-4">Información de Contacto</h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <form.Field
                name="contact_phone"
                validators={{ onChange: propertySchema.shape.contact_phone }}
                children={(field) => (
                  <div>
                    <Label htmlFor="contact_phone">Teléfono de contacto *</Label>
                    <Input
                      id="contact_phone"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="0999999999"
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="contact_email"
                validators={{ onChange: propertySchema.shape.contact_email }}
                children={(field) => (
                  <div>
                    <Label htmlFor="contact_email">Email de contacto *</Label>
                    <Input
                      id="contact_email"
                      type="email"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="contacto@ejemplo.com"
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />

              <form.Field
                name="notes"
                validators={{ onChange: propertySchema.shape.notes }}
                children={(field) => (
                  <div className="col-span-full">
                    <Label htmlFor="notes">Notas adicionales</Label>
                    <Textarea
                      id="notes"
                      value={field.state.value}
                      onChange={(e) => field.handleChange(e.target.value)}
                      placeholder="Información adicional sobre la propiedad..."
                      rows={3}
                    />
                    {field.state.meta.errors.map((error) => (
                      <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                    ))}
                  </div>
                )}
              />
            </div>
          </CardContent>
        </Card>

        {/* Imágenes - Solo mostrar si ya existe la propiedad */}
        {propertyId && (
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-semibold mb-4">Imágenes de la Propiedad</h3>
              <PropertyImageManager propertyId={propertyId} />
            </CardContent>
          </Card>
        )}

        {/* Botones */}
        <div className="flex justify-end gap-3 pt-6">
          {onCancel && (
            <Button type="button" variant="outline" onClick={onCancel}>
              Cancelar
            </Button>
          )}
          <Button 
            type="submit" 
            disabled={createPropertyMutation.isPending}
            className="min-w-[120px]"
          >
            {createPropertyMutation.isPending ? 'Guardando...' : 'Crear Propiedad'}
          </Button>
        </div>
      </div>
    </form>
  );
}