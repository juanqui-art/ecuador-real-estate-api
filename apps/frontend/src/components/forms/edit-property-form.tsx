'use client';

import { useState } from 'react';
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
import { Alert, AlertDescription } from '@/components/ui/alert';
import { CheckCircle, AlertCircle, Save, X } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { ECUADORIAN_PROVINCES, PROPERTY_TYPES, PROPERTY_STATUS } from '@/lib/constants';
import { formatPrice } from '@/lib/utils';

// Esquema completo para referencia de tipos
const editPropertySchema = z.object({
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
  address: z.string().optional(), // No editable
  bedrooms: z.number().optional(), // No editable
  bathrooms: z.number().optional(), // No editable
  area_m2: z.number().optional(), // No editable
  parking_spaces: z.number().min(0, 'Número de parqueaderos inválido').max(20, 'Máximo 20 parqueaderos'),
  year_built: z.number().optional(), // No editable
  has_garden: z.boolean().optional(), // No editable
  has_pool: z.boolean().optional(), // No editable
  has_elevator: z.boolean().optional(), // No editable
  has_balcony: z.boolean().optional(), // No editable
  has_terrace: z.boolean().optional(), // No editable
  has_garage: z.boolean().optional(), // No editable
  is_furnished: z.boolean().optional(), // No editable
  allows_pets: z.boolean().optional(), // No editable
  contact_phone: z.string().optional(), // No editable
  contact_email: z.string().optional(), // No editable
  notes: z.string().optional(), // No editable
});

// Esquemas específicos para validación de campos editables
const editableFieldSchemas = {
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
  parking_spaces: z.number().min(0, 'Número de parqueaderos inválido').max(20, 'Máximo 20 parqueaderos'),
};

type EditPropertyFormData = z.infer<typeof editPropertySchema>;

interface Property {
  id: string;
  title: string;
  description: string;
  price: number;
  type: string;
  status: string;
  province: string;
  city: string;
  address: string | null;
  bedrooms: number;
  bathrooms: number;
  area_m2: number;
  parking_spaces: number;
  year_built?: number | null;
  has_garden: boolean;
  has_pool: boolean;
  has_elevator: boolean;
  has_balcony: boolean;
  has_terrace: boolean;
  has_garage: boolean;
  is_furnished: boolean;
  allows_pets: boolean;
  contact_phone: string | null;
  contact_email: string | null;
  notes?: string | null;
}

interface EditPropertyFormProps {
  property: Property;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function EditPropertyForm({ property, onSuccess, onCancel }: EditPropertyFormProps) {
  const queryClient = useQueryClient();
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Función robusta para extraer mensaje de error de objetos complejos de Zod
  const getErrorMessage = (error: any): string => {
    // Debug logging para identificar el problema
    console.log('Processing error:', error);
    
    try {
      // Caso simple: string
      if (typeof error === 'string') return error;
      
      // Caso nulo o undefined
      if (!error) return 'Error de validación';
      
      // Caso objeto
      if (error && typeof error === 'object') {
        // Extraer mensaje directamente si existe y es string
        if (error.message && typeof error.message === 'string') {
          return error.message;
        }
        
        // Manejar todos los códigos de error posibles
        if (error.code) {
          switch (error.code) {
            case 'too_small':
              const minText = error.inclusive ? ' o igual' : '';
              return `El valor debe ser mayor${minText} a ${error.minimum}`;
            case 'too_big':
              const maxText = error.inclusive ? ' o igual' : '';
              return `El valor debe ser menor${maxText} a ${error.maximum}`;
            case 'invalid_type':
              return 'Tipo de dato inválido';
            case 'invalid_enum_value':
              return 'Valor no válido';
            case 'invalid_string':
              return 'Formato de texto inválido';
            case 'invalid_number':
              return 'Número inválido';
            case 'invalid_date':
              return 'Fecha inválida';
            default:
              return `Error de validación: ${error.code}`;
          }
        }
        
        // Fallback para objetos sin código específico
        if (error.message) {
          return String(error.message);
        }
        
        // Último recurso: convertir objeto a string seguro
        return 'Error de validación';
      }
      
      // Fallback para cualquier otro tipo
      return String(error);
    } catch (e) {
      console.error('Error processing validation error:', e);
      return 'Error de validación';
    }
  };

  const updatePropertyMutation = useMutation({
    mutationFn: async (data: EditPropertyFormData) => {
      // Backend currently only supports these fields
      const updateData = {
        title: data.title,
        description: data.description,
        price: data.price,
        type: data.type,
        province: data.province,
        city: data.city,
        parking_spaces: data.parking_spaces,
      };
      const response = await apiClient.put(`/properties/${property.id}`, updateData);
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      onSuccess?.();
    },
    onError: (error: any) => {
      console.error('Error updating property:', error);
    },
  });

  const form = useForm({
    defaultValues: {
      title: property.title || '',
      description: property.description || '',
      price: property.price || 0,
      type: property.type as 'house' | 'apartment' | 'land' | 'commercial',
      status: property.status as 'available' | 'sold' | 'rented',
      province: property.province || '',
      city: property.city || '',
      address: property.address || '',
      bedrooms: property.bedrooms || 0,
      bathrooms: property.bathrooms || 0,
      area_m2: property.area_m2 || 0,
      parking_spaces: property.parking_spaces || 0,
      year_built: property.year_built || undefined,
      has_garden: property.has_garden || false,
      has_pool: property.has_pool || false,
      has_elevator: property.has_elevator || false,
      has_balcony: property.has_balcony || false,
      has_terrace: property.has_terrace || false,
      has_garage: property.has_garage || false,
      is_furnished: property.is_furnished || false,
      allows_pets: property.allows_pets || false,
      contact_phone: property.contact_phone || '',
      contact_email: property.contact_email || '',
      notes: property.notes || '',
    },
    onSubmit: async ({ value }) => {
      setIsSubmitting(true);
      try {
        await updatePropertyMutation.mutateAsync(value);
      } finally {
        setIsSubmitting(false);
      }
    },
    validatorAdapter: zodValidator(),
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

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">Editar Propiedad</h2>
          <p className="text-sm text-gray-600">Modifica los datos de la propiedad</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge className={getStatusColor(property.status)}>
            {getStatusLabel(property.status)}
          </Badge>
          <Badge variant="outline">
            {getTypeLabel(property.type)}
          </Badge>
        </div>
      </div>

      {/* Info Alert */}
      <Alert className="border-blue-200 bg-blue-50">
        <AlertCircle className="h-4 w-4 text-blue-600" />
        <AlertDescription className="text-blue-800">
          <strong>Información:</strong> Actualmente se pueden editar: título, descripción, precio, tipo, ubicación y parqueaderos. 
          Pronto estará disponible la edición de características adicionales y datos de contacto.
        </AlertDescription>
      </Alert>

      {/* Error Alert */}
      {updatePropertyMutation.error && (
        <Alert className="border-red-200 bg-red-50">
          <AlertCircle className="h-4 w-4 text-red-600" />
          <AlertDescription className="text-red-800">
            <strong>Error:</strong> {updatePropertyMutation.error instanceof Error 
              ? updatePropertyMutation.error.message 
              : 'No se pudo actualizar la propiedad'}
          </AlertDescription>
        </Alert>
      )}

      <form onSubmit={(e) => {
        e.preventDefault();
        form.handleSubmit();
      }}>
        <div className="space-y-6">
          {/* Información Básica */}
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-medium mb-4">Información Básica</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <form.Field
                  name="title"
                  validators={{ onChange: editableFieldSchemas.title }}
                  children={(field) => (
                    <div className="col-span-full">
                      <Label htmlFor="title">Título de la propiedad *</Label>
                      <Input
                        id="title"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        placeholder="Ej: Hermosa casa en Samborondón con piscina"
                      />
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="type"
                  validators={{ onChange: editableFieldSchemas.type }}
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="status"
                  validators={{ onChange: editableFieldSchemas.status }}
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="price"
                  validators={{ onChange: editableFieldSchemas.price }}
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="description"
                  validators={{ onChange: editableFieldSchemas.description }}
                  children={(field) => (
                    <div className="col-span-full">
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Ubicación */}
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-medium mb-4">Ubicación</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <form.Field
                  name="province"
                  validators={{ onChange: editableFieldSchemas.province }}
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="city"
                  validators={{ onChange: editableFieldSchemas.city }}
                  children={(field) => (
                    <div>
                      <Label htmlFor="city">Ciudad *</Label>
                      <Input
                        id="city"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        placeholder="Ej: Samborondón"
                      />
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="address"
                  children={(field) => (
                    <div className="col-span-full">
                      <Label htmlFor="address">Dirección completa *</Label>
                      <Input
                        id="address"
                        value={field.state.value}
                        onChange={(e) => field.handleChange(e.target.value)}
                        placeholder="Ej: Km 2.5 Vía Samborondón, Urbanización La Puntilla"
                      />
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Características */}
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-medium mb-4">Características</h3>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <form.Field
                  name="bedrooms"
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
                        disabled
                        className="bg-gray-100"
                      />
                      <p className="text-xs text-gray-500 mt-1">No editable en esta versión</p>
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="bathrooms"
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
                        disabled
                        className="bg-gray-100"
                      />
                      <p className="text-xs text-gray-500 mt-1">No editable en esta versión</p>
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="area_m2"
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
                        disabled
                        className="bg-gray-100"
                      />
                      <p className="text-xs text-gray-500 mt-1">No editable en esta versión</p>
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="parking_spaces"
                  validators={{ onChange: editableFieldSchemas.parking_spaces }}
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />

                <form.Field
                  name="year_built"
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
                      {field.state.meta.errors.map((error, index) => {
                        const errorMessage = getErrorMessage(error);
                        return (
                          <p key={index} className="text-sm text-red-500 mt-1">
                            {typeof errorMessage === 'string' ? errorMessage : 'Error de validación'}
                          </p>
                        );
                      })}
                    </div>
                  )}
                />
              </div>
            </CardContent>
          </Card>

          {/* Características Adicionales */}
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-medium mb-4">Características Adicionales</h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600 mb-3">
                  <strong>Nota:</strong> La edición de características adicionales estará disponible en la próxima versión.
                </p>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  {features.map((feature) => (
                    <form.Field
                      key={feature.name}
                      name={feature.name as keyof EditPropertyFormData}
                      children={(field) => (
                        <div className="flex items-center space-x-2">
                          <input
                            type="checkbox"
                            id={feature.name}
                            checked={field.state.value as boolean}
                            onChange={(e) => field.handleChange(e.target.checked)}
                            className="rounded border-gray-300"
                            disabled
                          />
                          <Label htmlFor={feature.name} className="text-sm font-normal text-gray-500">
                            {feature.label}
                          </Label>
                        </div>
                      )}
                    />
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Contacto */}
          <Card>
            <CardContent className="pt-6">
              <h3 className="text-lg font-medium mb-4">Contacto</h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-sm text-gray-600 mb-3">
                  <strong>Nota:</strong> La edición de datos de contacto estará disponible en la próxima versión.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <form.Field
                    name="contact_phone"
                    children={(field) => (
                      <div>
                        <Label htmlFor="contact_phone">Teléfono de contacto *</Label>
                        <Input
                          id="contact_phone"
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          placeholder="0999999999"
                          disabled
                          className="bg-gray-100"
                        />
                        {field.state.meta.errors.map((error) => (
                          <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                        ))}
                      </div>
                    )}
                  />

                  <form.Field
                    name="contact_email"
                    children={(field) => (
                      <div>
                        <Label htmlFor="contact_email">Email de contacto *</Label>
                        <Input
                          id="contact_email"
                          type="email"
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          placeholder="contacto@ejemplo.com"
                          disabled
                          className="bg-gray-100"
                        />
                        {field.state.meta.errors.map((error) => (
                          <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                        ))}
                      </div>
                    )}
                  />

                  <form.Field
                    name="notes"
                    children={(field) => (
                      <div className="col-span-full">
                        <Label htmlFor="notes">Notas adicionales</Label>
                        <Textarea
                          id="notes"
                          value={field.state.value}
                          onChange={(e) => field.handleChange(e.target.value)}
                          placeholder="Información adicional sobre la propiedad..."
                          rows={3}
                          disabled
                          className="bg-gray-100"
                        />
                        {field.state.meta.errors.map((error) => (
                          <p key={error} className="text-sm text-red-500 mt-1">{error}</p>
                        ))}
                      </div>
                    )}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Navigation Buttons */}
          <div className="flex justify-between pt-6">
            <Button
              type="button"
              variant="outline"
              onClick={onCancel}
              disabled={isSubmitting || updatePropertyMutation.isPending}
            >
              <X className="w-4 h-4 mr-2" />
              Cancelar
            </Button>

            <Button 
              type="submit" 
              disabled={isSubmitting || updatePropertyMutation.isPending}
              className="min-w-[140px]"
            >
              {isSubmitting || updatePropertyMutation.isPending ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2" />
                  Actualizando...
                </>
              ) : (
                <>
                  <Save className="w-4 h-4 mr-2" />
                  Guardar Cambios
                </>
              )}
            </Button>
          </div>
        </div>
      </form>
    </div>
  );
}