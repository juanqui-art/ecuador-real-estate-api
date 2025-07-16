'use server';

import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { propertySchema, propertyFilterSchema, imageUploadSchema } from '@/lib/validations/property';
import type { PropertyFormData, PropertyFilterData, ImageUploadData } from '@/lib/validations/property';

/**
 * Server action to create a new property
 */
export async function createPropertyAction(prevState: any, formData: FormData) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
        errors: {},
      };
    }

    // Extract and convert form data
    const rawData = {
      title: formData.get('title') as string,
      description: formData.get('description') as string,
      price: parseFloat(formData.get('price') as string),
      province: formData.get('province') as string,
      city: formData.get('city') as string,
      type: formData.get('type') as string,
      status: formData.get('status') as string || 'available',
      bedrooms: parseInt(formData.get('bedrooms') as string),
      bathrooms: parseFloat(formData.get('bathrooms') as string),
      area_m2: parseFloat(formData.get('area_m2') as string),
      address: formData.get('address') as string || undefined,
      latitude: formData.get('latitude') ? parseFloat(formData.get('latitude') as string) : undefined,
      longitude: formData.get('longitude') ? parseFloat(formData.get('longitude') as string) : undefined,
      featured: formData.get('featured') === 'true',
      parking_spots: formData.get('parking_spots') ? parseInt(formData.get('parking_spots') as string) : undefined,
      has_garden: formData.get('has_garden') === 'true',
      has_pool: formData.get('has_pool') === 'true',
      has_balcony: formData.get('has_balcony') === 'true',
      has_elevator: formData.get('has_elevator') === 'true',
      year_built: formData.get('year_built') ? parseInt(formData.get('year_built') as string) : undefined,
    };

    // Validate with Zod
    const validatedData = propertySchema.parse(rawData);

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify(validatedData),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al crear la propiedad',
        errors: {},
      };
    }

    const propertyData = await response.json();

    // Revalidate properties pages
    revalidatePath('/properties');
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad creada exitosamente',
      property: propertyData,
      redirect: `/properties/${propertyData.id}`,
    };

  } catch (error) {
    if (error instanceof Error && error.name === 'ZodError') {
      return {
        success: false,
        message: 'Datos inválidos',
        errors: (error as any).flatten().fieldErrors,
      };
    }

    console.error('Create property action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Server action to update an existing property
 */
export async function updatePropertyAction(propertyId: string, prevState: any, formData: FormData) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
        errors: {},
      };
    }

    // Extract and convert form data
    const rawData = {
      title: formData.get('title') as string,
      description: formData.get('description') as string,
      price: parseFloat(formData.get('price') as string),
      province: formData.get('province') as string,
      city: formData.get('city') as string,
      type: formData.get('type') as string,
      status: formData.get('status') as string,
      bedrooms: parseInt(formData.get('bedrooms') as string),
      bathrooms: parseFloat(formData.get('bathrooms') as string),
      area_m2: parseFloat(formData.get('area_m2') as string),
      address: formData.get('address') as string || undefined,
      latitude: formData.get('latitude') ? parseFloat(formData.get('latitude') as string) : undefined,
      longitude: formData.get('longitude') ? parseFloat(formData.get('longitude') as string) : undefined,
      featured: formData.get('featured') === 'true',
      parking_spots: formData.get('parking_spots') ? parseInt(formData.get('parking_spots') as string) : undefined,
      has_garden: formData.get('has_garden') === 'true',
      has_pool: formData.get('has_pool') === 'true',
      has_balcony: formData.get('has_balcony') === 'true',
      has_elevator: formData.get('has_elevator') === 'true',
      year_built: formData.get('year_built') ? parseInt(formData.get('year_built') as string) : undefined,
    };

    // Validate with Zod
    const validatedData = propertySchema.parse(rawData);

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify(validatedData),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al actualizar la propiedad',
        errors: {},
      };
    }

    const propertyData = await response.json();

    // Revalidate properties pages
    revalidatePath('/properties');
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad actualizada exitosamente',
      property: propertyData,
    };

  } catch (error) {
    if (error instanceof Error && error.name === 'ZodError') {
      return {
        success: false,
        message: 'Datos inválidos',
        errors: (error as any).flatten().fieldErrors,
      };
    }

    console.error('Update property action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Server action to delete a property
 */
export async function deletePropertyAction(propertyId: string) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
      };
    }

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al eliminar la propiedad',
      };
    }

    // Revalidate properties pages
    revalidatePath('/properties');
    revalidatePath('/dashboard');

    return {
      success: true,
      message: 'Propiedad eliminada exitosamente',
      redirect: '/properties',
    };

  } catch (error) {
    console.error('Delete property action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}

/**
 * Server action to upload property images
 */
export async function uploadPropertyImageAction(prevState: any, formData: FormData) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
        errors: {},
      };
    }

    const propertyId = formData.get('property_id') as string;
    const altText = formData.get('alt_text') as string || '';
    const imageFile = formData.get('image') as File;

    if (!imageFile || imageFile.size === 0) {
      return {
        success: false,
        message: 'No se seleccionó ninguna imagen',
        errors: {},
      };
    }

    // Validate file size (10MB limit)
    if (imageFile.size > 10 * 1024 * 1024) {
      return {
        success: false,
        message: 'La imagen no puede ser mayor a 10MB',
        errors: {},
      };
    }

    // Validate file type
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/avif'];
    if (!allowedTypes.includes(imageFile.type)) {
      return {
        success: false,
        message: 'Formato de imagen no válido. Use JPEG, PNG, WebP o AVIF',
        errors: {},
      };
    }

    // Create form data for upload
    const uploadFormData = new FormData();
    uploadFormData.append('property_id', propertyId);
    uploadFormData.append('alt_text', altText);
    uploadFormData.append('image', imageFile);

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/images`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: uploadFormData,
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al subir la imagen',
        errors: {},
      };
    }

    const imageData = await response.json();

    // Revalidate property pages
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/properties');

    return {
      success: true,
      message: 'Imagen subida exitosamente',
      image: imageData,
    };

  } catch (error) {
    console.error('Upload image action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
      errors: {},
    };
  }
}

/**
 * Server action to reorder property images
 */
export async function reorderPropertyImagesAction(propertyId: string, imageOrders: { id: string; sort_order: number }[]) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
      };
    }

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}/images/reorder`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify({ images: imageOrders }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al reordenar las imágenes',
      };
    }

    // Revalidate property pages
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/properties');

    return {
      success: true,
      message: 'Imágenes reordenadas exitosamente',
    };

  } catch (error) {
    console.error('Reorder images action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}

/**
 * Server action to set main property image
 */
export async function setMainPropertyImageAction(propertyId: string, imageId: string) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
      };
    }

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/properties/${propertyId}/images/main`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json',
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
      body: JSON.stringify({ image_id: imageId }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al establecer imagen principal',
      };
    }

    // Revalidate property pages
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/properties');

    return {
      success: true,
      message: 'Imagen principal establecida exitosamente',
    };

  } catch (error) {
    console.error('Set main image action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}

/**
 * Server action to delete property image
 */
export async function deletePropertyImageAction(imageId: string, propertyId: string) {
  try {
    const cookieStore = await cookies();
    const accessToken = cookieStore.get('access_token')?.value;

    if (!accessToken) {
      return {
        success: false,
        message: 'No autenticado',
      };
    }

    // Make API call to backend
    const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/images/${imageId}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'User-Agent': 'RealtyCore-Dashboard/1.0 (Next.js)',
      },
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      return {
        success: false,
        message: errorData.message || 'Error al eliminar la imagen',
      };
    }

    // Revalidate property pages
    revalidatePath(`/properties/${propertyId}`);
    revalidatePath('/properties');

    return {
      success: true,
      message: 'Imagen eliminada exitosamente',
    };

  } catch (error) {
    console.error('Delete image action error:', error);
    return {
      success: false,
      message: 'Error interno del servidor',
    };
  }
}