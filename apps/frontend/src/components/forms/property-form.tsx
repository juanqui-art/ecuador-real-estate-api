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
import { Separator } from '@/components/ui/separator';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { CheckCircle, ArrowRight, ArrowLeft, AlertCircle, Image as ImageIcon } from 'lucide-react';
import { apiClient } from '@/lib/api-client';
import { ECUADORIAN_PROVINCES, PROPERTY_TYPES, PROPERTY_STATUS } from '@/lib/constants';
import { formatPrice } from '@/lib/utils';
import { PropertyImageManager } from '@/components/images/property-image-manager';
import { TemporaryImageUpload } from '@/components/images/temporary-image-upload';
import { useTemporaryImages } from '@/hooks/useTemporaryImages';

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
  const [currentStep, setCurrentStep] = useState(0);
  const [createdPropertyId, setCreatedPropertyId] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Temporary images hook
  const {
    temporaryImages,
    setTemporaryImages,
    uploadImages,
    uploadProgress,
    isUploading,
    hasValidImages,
    processingCount,
    errorCount,
    clearTemporaryImages,
  } = useTemporaryImages({
    onUploadComplete: (uploadedImages) => {
      console.log('Images uploaded successfully:', uploadedImages);
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      onSuccess?.();
    },
    onUploadError: (error) => {
      console.error('Error uploading images:', error);
      // Don't fail the entire process, just show a warning
    },
  });

  const createPropertyMutation = useMutation({
    mutationFn: async (data: PropertyFormData) => {
      const response = await apiClient.post('/properties', data);
      return response.data;
    },
    onSuccess: async (createdProperty) => {
      console.log('Property created successfully:', createdProperty);
      setCreatedPropertyId(createdProperty.data.id);
      
      // If there are images, upload them
      if (temporaryImages.length > 0) {
        try {
          await uploadImages(createdProperty.data.id);
        } catch (error) {
          console.error('Error uploading images:', error);
          // Property was created successfully, just show a warning about images
          queryClient.invalidateQueries({ queryKey: ['properties'] });
          onSuccess?.();
        }
      } else {
        // No images to upload, finish successfully
        queryClient.invalidateQueries({ queryKey: ['properties'] });
        onSuccess?.();
      }
    },
  });

  const steps = [
    {
      title: 'Información Básica',
      description: 'Datos generales de la propiedad',
      icon: CheckCircle,
    },
    {
      title: 'Ubicación',
      description: 'Dirección y ubicación',
      icon: CheckCircle,
    },
    {
      title: 'Características',
      description: 'Dormitorios, baños, área, etc.',
      icon: CheckCircle,
    },
    {
      title: 'Características Adicionales',
      description: 'Jardín, piscina, garage, etc.',
      icon: CheckCircle,
    },
    {
      title: 'Imágenes',
      description: 'Fotos de la propiedad',
      icon: ImageIcon,
    },
    {
      title: 'Contacto',
      description: 'Información de contacto',
      icon: CheckCircle,
    },
  ];

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
      setIsSubmitting(true);
      try {
        await createPropertyMutation.mutateAsync(value);
      } finally {
        setIsSubmitting(false);
      }
    },
    validatorAdapter: zodValidator,
  });

  const validateCurrentStep = () => {
    const formState = form.state;
    const errors = formState.errors;

    switch (currentStep) {
      case 0: // Información Básica
        return !errors.find(error => 
          error.path === 'title' || error.path === 'description' || 
          error.path === 'price' || error.path === 'type' || error.path === 'status'
        );
      case 1: // Ubicación
        return !errors.find(error => 
          error.path === 'province' || error.path === 'city' || error.path === 'address'
        );
      case 2: // Características
        return !errors.find(error => 
          error.path === 'bedrooms' || error.path === 'bathrooms' || 
          error.path === 'area_m2' || error.path === 'parking_spaces' || error.path === 'year_built'
        );
      case 3: // Características Adicionales - No validation needed
        return true;
      case 4: // Imágenes - Optional but warn if no images
        return true;
      case 5: // Contacto
        return !errors.find(error => 
          error.path === 'contact_phone' || error.path === 'contact_email'
        );
      default:
        return true;
    }
  };

  const canProceedToNextStep = () => {
    if (currentStep === 4) {
      // Images step - can proceed but show warning if no images
      return processingCount === 0; // Wait for processing to complete
    }
    return validateCurrentStep();
  };

  const nextStep = () => {
    if (currentStep < steps.length - 1 && canProceedToNextStep()) {
      setCurrentStep(currentStep + 1);
    }
  };

  const prevStep = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const goToImagesStep = () => {
    setCurrentStep(4); // Images step
  };

  const handleFinish = () => {
    if (validateCurrentStep()) {
      form.handleSubmit();
    }
  };

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

  const renderStepContent = () => {
    switch (currentStep) {
      case 0: // Información Básica
        return (
          <Card>
            <CardContent className="pt-6">
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
        );

      case 1: // Ubicación
        return (
          <Card>
            <CardContent className="pt-6">
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
        );

      case 2: // Características
        return (
          <Card>
            <CardContent className="pt-6">
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
        );

      case 3: // Características Adicionales
        return (
          <Card>
            <CardContent className="pt-6">
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
        );

      case 4: // Imágenes
        return (
          <div className="space-y-6">
            {/* Header with tips */}
            <div className="text-center space-y-2">
              <h3 className="text-lg font-semibold text-gray-900">
                Agrega imágenes de la propiedad
              </h3>
              <p className="text-sm text-gray-600">
                Las imágenes de alta calidad aumentan hasta 40% las probabilidades de venta
              </p>
            </div>

            <TemporaryImageUpload
              images={temporaryImages}
              onImagesChange={setTemporaryImages}
              maxImages={10}
            />
            
            {/* Status Messages */}
            {temporaryImages.length === 0 && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  <strong>Recomendación:</strong> Sube al menos 3 fotos para mejorar la visibilidad de la propiedad. 
                  La primera imagen será la foto principal que se mostrará en las búsquedas.
                </AlertDescription>
              </Alert>
            )}
            
            {temporaryImages.length > 0 && temporaryImages.length < 3 && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>
                  <strong>Sugerencia:</strong> Considera agregar {3 - temporaryImages.length} imagen{3 - temporaryImages.length > 1 ? 'es' : ''} más 
                  para crear un portafolio más completo.
                </AlertDescription>
              </Alert>
            )}
            
            {temporaryImages.length >= 3 && (
              <Alert className="border-green-200 bg-green-50">
                <CheckCircle className="h-4 w-4 text-green-600" />
                <AlertDescription className="text-green-800">
                  <strong>Excelente!</strong> Tienes {temporaryImages.length} imágenes seleccionadas. 
                  Esto mejorará significativamente la presentación de tu propiedad.
                </AlertDescription>
              </Alert>
            )}
            
            {/* Image Tips */}
            <Card className="border-blue-200 bg-blue-50">
              <CardContent className="pt-4">
                <h4 className="font-medium text-blue-900 mb-2">💡 Consejos para mejores fotos:</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm text-blue-800">
                  <div className="flex items-start gap-2">
                    <span className="text-blue-600 font-bold">•</span>
                    <p>Toma fotos con buena iluminación natural</p>
                  </div>
                  <div className="flex items-start gap-2">
                    <span className="text-blue-600 font-bold">•</span>
                    <p>Incluye exteriores, interiores y áreas comunes</p>
                  </div>
                  <div className="flex items-start gap-2">
                    <span className="text-blue-600 font-bold">•</span>
                    <p>Muestra los espacios más atractivos primero</p>
                  </div>
                  <div className="flex items-start gap-2">
                    <span className="text-blue-600 font-bold">•</span>
                    <p>Evita fotos borrosas o muy oscuras</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        );

      case 5: // Contacto
        return (
          <Card>
            <CardContent className="pt-6">
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
        );

      default:
        return null;
    }
  };

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      handleFinish();
    }}>
      <div className="space-y-6">
        {/* Progress Steps */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">Crear Nueva Propiedad</h2>
            <div className="flex items-center gap-3">
              {temporaryImages.length > 0 && (
                <Badge variant="secondary" className="gap-2">
                  <ImageIcon className="w-3 h-3" />
                  {temporaryImages.length} imagen{temporaryImages.length > 1 ? 'es' : ''} seleccionada{temporaryImages.length > 1 ? 's' : ''}
                </Badge>
              )}
              <Badge variant="outline">
                Paso {currentStep + 1} de {steps.length}
              </Badge>
            </div>
          </div>
          
          <div className="flex items-center space-x-2">
            {steps.map((step, index) => (
              <div key={index} className="flex items-center">
                <div className={`flex items-center justify-center w-8 h-8 rounded-full border-2 relative ${
                  index < currentStep 
                    ? 'bg-green-500 border-green-500 text-white' 
                    : index === currentStep
                    ? 'bg-blue-500 border-blue-500 text-white'
                    : 'border-gray-300 text-gray-400'
                }`}>
                  {index < currentStep ? (
                    <CheckCircle className="w-4 h-4" />
                  ) : (
                    <span className="text-sm">{index + 1}</span>
                  )}
                  {/* Special indicator for images step */}
                  {index === 4 && temporaryImages.length > 0 && (
                    <div className="absolute -top-1 -right-1 w-4 h-4 bg-green-500 rounded-full flex items-center justify-center">
                      <span className="text-xs text-white font-bold">{temporaryImages.length}</span>
                    </div>
                  )}
                </div>
                {index < steps.length - 1 && (
                  <div className={`w-12 h-0.5 mx-2 ${
                    index < currentStep ? 'bg-green-500' : 'bg-gray-300'
                  }`} />
                )}
              </div>
            ))}
          </div>
          
          <div className="mt-4">
            <div className="flex items-center space-x-2">
              <step.icon className="w-5 h-5 text-gray-600" />
              <div>
                <h3 className="font-medium">{steps[currentStep].title}</h3>
                <p className="text-sm text-gray-600">{steps[currentStep].description}</p>
              </div>
            </div>
          </div>
        </div>

        {/* Step Content */}
        <div className="min-h-[400px] relative">
          {renderStepContent()}
          
          {/* Quick Access to Images Step */}
          {currentStep < 4 && currentStep > 0 && (
            <div className="absolute top-4 right-4">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={goToImagesStep}
                className="shadow-lg border-blue-200 hover:border-blue-300 bg-blue-50 hover:bg-blue-100"
              >
                <ImageIcon className="w-4 h-4 mr-2" />
                Agregar Imágenes
                {temporaryImages.length > 0 && (
                  <Badge variant="secondary" className="ml-2">
                    {temporaryImages.length}
                  </Badge>
                )}
              </Button>
            </div>
          )}
        </div>

        {/* Upload Progress */}
        {(isUploading || uploadProgress > 0) && (
          <Card>
            <CardContent className="pt-6">
              <div className="space-y-2">
                <div className="flex justify-between text-sm">
                  <span>Subiendo imágenes...</span>
                  <span>{Math.round(uploadProgress)}%</span>
                </div>
                <Progress value={uploadProgress} className="h-2" />
              </div>
            </CardContent>
          </Card>
        )}

        {/* Navigation Buttons */}
        <div className="flex justify-between pt-6">
          <Button
            type="button"
            variant="outline"
            onClick={prevStep}
            disabled={currentStep === 0}
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Anterior
          </Button>

          <div className="flex gap-3">
            {onCancel && (
              <Button type="button" variant="outline" onClick={onCancel}>
                Cancelar
              </Button>
            )}
            
            {currentStep < steps.length - 1 ? (
              <Button
                type="button"
                onClick={nextStep}
                disabled={!canProceedToNextStep()}
              >
                Siguiente
                <ArrowRight className="w-4 h-4 ml-2" />
              </Button>
            ) : (
              <Button 
                type="submit" 
                disabled={isSubmitting || createPropertyMutation.isPending || isUploading}
                className="min-w-[140px]"
              >
                {isSubmitting || createPropertyMutation.isPending || isUploading 
                  ? (isUploading ? 'Subiendo imágenes...' : 'Guardando...') 
                  : 'Crear Propiedad'
                }
              </Button>
            )}
          </div>
        </div>
      </div>
    </form>
  );
}