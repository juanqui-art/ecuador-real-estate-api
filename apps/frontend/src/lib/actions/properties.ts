'use server';

import { redirect } from 'next/navigation';
import { revalidatePath } from 'next/cache';
import { z } from 'zod';

/**
 * Modern Server Actions for Property CRUD Operations
 * Following Next.js 15 + React 19 Best Practices (2025)
 * 
 * Key Features:
 * - Progressive Enhancement (works without JS)
 * - Type-safe with Zod validation
 * - Proper error handling
 * - Optimistic UI support via revalidatePath
 * - Server-side validation
 * - NO AUTH MODE for development
 */

// Complete Property Schema synchronized with backend Go struct (2025)
// OPTIMIZED: Reduced from 15 to 7 required fields for better UX
const PropertySchema = z.object({
  // OBLIGATORIOS: Información básica esencial (5 campos)
  title: z.string().min(10, 'El título debe tener al menos 10 caracteres'),
  description: z.string().min(50, 'La descripción debe tener al menos 50 caracteres'),
  price: z.coerce.number().min(1000, 'El precio debe ser mayor a $1,000'),
  type: z.enum(['house', 'apartment', 'land', 'commercial'], {
    message: 'Selecciona un tipo de propiedad válido'
  }),
  status: z.enum(['available', 'sold', 'rented', 'reserved'], {
    message: 'Selecciona un estado válido'
  }),
  
  // OPCIONALES: Ubicación (puede completarse gradualmente)
  province: z.string().min(1, 'Selecciona una provincia').optional(),
  city: z.string().min(2, 'Ingresa la ciudad').optional(),
  address: z.string().min(10, 'Ingresa la dirección completa').optional(),
  sector: z.string().optional(),
  latitude: z.coerce.number().optional(),
  longitude: z.coerce.number().optional(),
  location_precision: z.string().default('approximate'),
  
  // OPCIONALES: Características (con defaults inteligentes)
  bedrooms: z.coerce.number().min(0, 'Número de dormitorios inválido').max(20, 'Máximo 20 dormitorios').default(1),
  bathrooms: z.coerce.number().min(0, 'Número de baños inválido').max(20, 'Máximo 20 baños').default(1), // Soporta 2.5
  area_m2: z.coerce.number().min(10, 'El área debe ser mayor a 10 m²').max(10000, 'Máximo 10,000 m²').optional(),
  parking_spaces: z.coerce.number().min(0, 'Número de parqueaderos inválido').max(20, 'Máximo 20 parqueaderos').default(1),
  year_built: z.coerce.number().min(1900, 'Año inválido').max(new Date().getFullYear(), 'Año no puede ser futuro').optional(),
  floors: z.coerce.number().min(1, 'Mínimo 1 piso').max(50, 'Máximo 50 pisos').optional(),
  
  // Precios adicionales
  rent_price: z.coerce.number().min(100, 'Precio de renta inválido').optional(),
  common_expenses: z.coerce.number().min(0, 'Gastos comunes inválidos').optional(),
  price_per_m2: z.coerce.number().min(10, 'Precio por m² inválido').optional(),
  
  // Multimedia
  main_image: z.string().url('URL de imagen inválida').optional(),
  images: z.array(z.string().url()).default([]),
  video_tour: z.string().url('URL de video inválida').optional(),
  tour_360: z.string().url('URL de tour 360 inválida').optional(),
  
  // Estado y clasificación
  property_status: z.enum(['new', 'used', 'renovated'], {
    message: 'Selecciona un estado de propiedad válido'
  }).default('new'),
  tags: z.array(z.string().min(2, 'Tag muy corto').max(30, 'Tag muy largo')).default([]),
  featured: z.coerce.boolean().default(false),
  view_count: z.coerce.number().default(0),
  
  // Amenidades (características adicionales) - sincronizadas con backend
  furnished: z.coerce.boolean().default(false),
  garage: z.coerce.boolean().default(false),
  pool: z.coerce.boolean().default(false),
  garden: z.coerce.boolean().default(false),
  terrace: z.coerce.boolean().default(false),
  balcony: z.coerce.boolean().default(false),
  security: z.coerce.boolean().default(false),
  elevator: z.coerce.boolean().default(false),
  air_conditioning: z.coerce.boolean().default(false),
  
  // Sistema de ownership (opcional para formularios, manejado por backend)
  real_estate_company_id: z.string().uuid().optional(),
  owner_id: z.string().uuid().optional(),
  agent_id: z.string().uuid().optional(),
  agency_id: z.string().uuid().optional(),
  
  // OBLIGATORIOS: Contacto esencial (temporal, deberá moverse a sistema de usuarios)
  contact_phone: z.string().min(10, 'Ingresa un teléfono válido'),
  contact_email: z.email('Ingresa un email válido'),
  notes: z.string().optional(),
});

// Action result type for better error handling
export type ActionResult<T = any> = {
  success: boolean;
  data?: T;
  message?: string;
  errors?: Record<string, string[]>;
};

/**
 * Modern Create Property Server Action (2025)
 * NO AUTH MODE - Direct backend communication
 */
export async function createPropertyAction(prevState: any, formData: FormData): Promise<ActionResult> {
  try {
    // NO AUTH MODE - Skip authentication check
    // const cookieStore = await cookies();
    // const accessToken = cookieStore.get('access_token')?.value;

    // Modern FormData parsing with 2025 best practices
    console.log('🔧 Server Action - Creating property with modern approach');
    
    const rawData: Record<string, any> = Object.fromEntries(formData);
    console.log('🔧 Raw FormData:', rawData);

    // Process tags field: convert comma-separated string to array
    if (rawData.tags && typeof rawData.tags === 'string') {
      rawData.tags = rawData.tags
        .split(',')
        .map((tag: string) => tag.trim())
        .filter((tag: string) => tag.length > 0);
    }

    console.log('🔧 Processed FormData:', rawData);

    // Modern validation with Zod (server-side)
    const validatedData = PropertySchema.safeParse(rawData);

    if (!validatedData.success) {
      const flattened = z.flattenError(validatedData.error);
      console.log('❌ Server Action - Validation failed:', flattened);
      return {
        success: false,
        message: 'Datos del formulario inválidos',
        errors: flattened.fieldErrors,
      };
    }

    console.log('✅ Server Action - Data validated successfully');

    // NO AUTH MODE - Direct API call without authentication
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Next.js-Server-Action/2025 (React-19)',
      },
      body: JSON.stringify(validatedData.data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      console.error('❌ Server Action - API Error:', errorData);
      return {
        success: false,
        message: errorData.message || `Error del servidor: ${response.status}`,
        errors: {},
      };
    }

    const propertyData = await response.json();
    console.log('✅ Server Action - Property created:', propertyData);

    // Modern revalidation for optimistic UI
    revalidatePath('/properties');
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad creada exitosamente',
      data: propertyData,
    };

  } catch (error) {
    // Enhanced error handling for 2025
    if (error && typeof error === 'object' && 'issues' in error) {
      // Zod validation error
      console.log('❌ Zod validation error:', error);
      return {
        success: false,
        message: 'Datos inválidos',
        errors: (error as any).flatten?.()?.fieldErrors || {},
      };
    }

    console.error('💥 Server Action - Unexpected error:', error);
    return {
      success: false,
      message: 'Error interno del servidor. Por favor intenta de nuevo.',
      errors: {},
    };
  }
}

/**
 * Modern Update Property Server Action (2025)
 * NO AUTH MODE - Direct backend communication
 */
export async function updatePropertyAction(propertyId: string, prevState: any, formData: FormData): Promise<ActionResult> {
  try {
    console.log('🔧 Server Action - Updating property:', propertyId);

    const rawData = Object.fromEntries(formData);
    const validatedData = PropertySchema.safeParse(rawData);

    if (!validatedData.success) {
      const flattened = z.flattenError(validatedData.error);
      console.log('❌ Update validation failed:', flattened);
      return {
        success: false,
        message: 'Datos del formulario inválidos',
        errors: flattened.fieldErrors,
      };
    }

    // NO AUTH MODE - Direct API call
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Next.js-Server-Action/2025 (React-19)',
      },
      body: JSON.stringify(validatedData.data),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || `Error al actualizar la propiedad: ${response.status}`,
        errors: {},
      };
    }

    const propertyData = await response.json();
    console.log('✅ Property updated successfully:', propertyData);

    // Modern revalidation
    revalidatePath('/properties');
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad actualizada exitosamente',
      data: propertyData,
    };

  } catch (error) {
    console.error('💥 Update property error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Modern Delete Property Server Action (2025)
 * NO AUTH MODE - Direct backend communication
 */
export async function deletePropertyAction(propertyId: string): Promise<ActionResult> {
  try {
    console.log('🔧 Server Action - Deleting property:', propertyId);

    // NO AUTH MODE - Direct API call
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}`, {
      method: 'DELETE',
      headers: {
        'User-Agent': 'Next.js-Server-Action/2025 (React-19)',
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || `Error al eliminar la propiedad: ${response.status}`,
      };
    }

    console.log('✅ Property deleted successfully');

    // Modern revalidation for optimistic UI
    revalidatePath('/properties');
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad eliminada exitosamente',
      data: { id: propertyId },
    };

  } catch (error) {
    console.error('💥 Delete property error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}

/**
 * Modern Upload Property Image Server Action (2025)
 * NO AUTH MODE - Direct backend communication
 */
export async function uploadPropertyImageAction(propertyId: string, formData: FormData): Promise<ActionResult> {
  try {
    console.log('🔧 Server Action - Uploading image for property:', propertyId);

    const imageFile = formData.get('image') as File;

    if (!imageFile || imageFile.size === 0) {
      return {
        success: false,
        message: 'No se seleccionó ninguna imagen',
      };
    }

    // Modern validation with enhanced limits
    if (imageFile.size > 10 * 1024 * 1024) {
      return {
        success: false,
        message: 'La imagen no puede ser mayor a 10MB',
      };
    }

    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/avif'];
    if (!allowedTypes.includes(imageFile.type)) {
      return {
        success: false,
        message: 'Formato de imagen no válido. Use JPEG, PNG, WebP o AVIF',
      };
    }

    // Prepare FormData for upload
    const uploadFormData = new FormData();
    uploadFormData.append('property_id', propertyId);
    uploadFormData.append('image', imageFile);

    // NO AUTH MODE - Direct API call
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/images`, {
      method: 'POST',
      headers: {
        'User-Agent': 'Next.js-Server-Action/2025 (React-19)',
      },
      body: uploadFormData,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || `Error al subir la imagen: ${response.status}`,
      };
    }

    const imageData = await response.json();
    console.log('✅ Image uploaded successfully:', imageData);

    // Modern revalidation
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/properties');

    return {
      success: true,
      message: 'Imagen subida exitosamente',
      data: imageData,
    };

  } catch (error) {
    console.error('💥 Upload image error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}

/**
 * Get Properties Server Action (2025)
 * For use with React 19 hooks and Suspense
 */
export async function getPropertiesAction(searchParams?: {
  search?: string;
  type?: string;
  status?: string;
  minPrice?: string;
  maxPrice?: string;
  province?: string;
}): Promise<ActionResult<any[]>> {
  try {
    console.log('🔧 Server Action - Fetching properties with filters:', searchParams);

    const params = new URLSearchParams();
    
    if (searchParams?.search) params.append('q', searchParams.search);
    if (searchParams?.type) params.append('type', searchParams.type);
    if (searchParams?.status) params.append('status', searchParams.status);
    if (searchParams?.minPrice) params.append('min_price', searchParams.minPrice);
    if (searchParams?.maxPrice) params.append('max_price', searchParams.maxPrice);
    if (searchParams?.province) params.append('province', searchParams.province);

    const url = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/filter?${params.toString()}`;
    
    const response = await fetch(url, {
      headers: {
        'User-Agent': 'Next.js-Server-Action/2025 (React-19)',
      },
      next: { revalidate: 60 }, // Cache for 60 seconds
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.status}`);
    }

    const result = await response.json();
    console.log('✅ Properties fetched successfully:', result.data?.length || 0, 'properties');
    
    return {
      success: true,
      data: result.data || result,
    };

  } catch (error) {
    console.error('💥 Get Properties Server Action Error:', error);
    return {
      success: false,
      message: 'Error cargando las propiedades.',
      data: [], // Return empty array as fallback
    };
  }
}

/**
 * Modern Progressive Enhancement Action (2025)
 * For traditional form submission (works without JS)
 * Enhanced with better error handling and UX
 */
export async function createPropertyWithRedirectAction(formData: FormData) {
  console.log('🔧 Progressive Enhancement - Processing form without JavaScript');
  
  const result = await createPropertyAction(null, formData);
  
  if (result.success) {
    console.log('✅ Progressive Enhancement - Property created, redirecting');
    // Success: redirect to properties list
    redirect('/properties?created=success&message=Propiedad+creada+exitosamente');
  } else {
    console.error('❌ Progressive Enhancement - Creation failed:', result.message);
    // Error: redirect back with error message
    const errorMsg = encodeURIComponent(result.message || 'Error al crear la propiedad');
    redirect(`/properties?created=error&message=${errorMsg}`);
  }
}

/**
 * Progressive Enhancement Update Action (2025)
 */
export async function updatePropertyWithRedirectAction(propertyId: string, formData: FormData) {
  console.log('🔧 Progressive Enhancement - Updating property:', propertyId);
  
  const result = await updatePropertyAction(propertyId, null, formData);
  
  if (result.success) {
    redirect(`/properties?updated=success&message=Propiedad+actualizada+exitosamente`);
  } else {
    const errorMsg = encodeURIComponent(result.message || 'Error al actualizar la propiedad');
    redirect(`/properties?updated=error&message=${errorMsg}`);
  }
}