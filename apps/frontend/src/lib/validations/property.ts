import { z } from 'zod';

/**
 * Ecuador provinces for validation
 */
const ECUADOR_PROVINCES = [
  'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi', 
  'El Oro', 'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja', 
  'Los Ríos', 'Manabí', 'Morona Santiago', 'Napo', 'Orellana', 'Pastaza', 
  'Pichincha', 'Santa Elena', 'Santo Domingo', 'Sucumbíos', 'Tungurahua', 
  'Zamora Chinchipe'
] as const;

/**
 * Property types
 */
const PROPERTY_TYPES = ['house', 'apartment', 'land', 'commercial'] as const;

/**
 * Property status
 */
const PROPERTY_STATUS = ['available', 'sold', 'rented'] as const;

/**
 * Create/Update property validation schema
 */
export const propertySchema = z.object({
  title: z
    .string()
    .min(1, 'El título es requerido')
    .min(10, 'El título debe tener al menos 10 caracteres')
    .max(200, 'El título no puede tener más de 200 caracteres'),
  
  description: z
    .string()
    .min(1, 'La descripción es requerida')
    .min(20, 'La descripción debe tener al menos 20 caracteres')
    .max(2000, 'La descripción no puede tener más de 2000 caracteres'),
  
  price: z
    .number({
      required_error: 'El precio es requerido',
      invalid_type_error: 'El precio debe ser un número',
    })
    .positive('El precio debe ser mayor a 0')
    .max(10000000, 'El precio no puede ser mayor a $10,000,000'),
  
  province: z
    .enum(ECUADOR_PROVINCES, {
      required_error: 'La provincia es requerida',
      invalid_type_error: 'Provincia inválida',
    }),
  
  city: z
    .string()
    .min(1, 'La ciudad es requerida')
    .min(2, 'La ciudad debe tener al menos 2 caracteres')
    .max(100, 'La ciudad no puede tener más de 100 caracteres'),
  
  type: z
    .enum(PROPERTY_TYPES, {
      required_error: 'El tipo de propiedad es requerido',
      invalid_type_error: 'Tipo de propiedad inválido',
    }),
  
  status: z
    .enum(PROPERTY_STATUS, {
      required_error: 'El estado es requerido',
      invalid_type_error: 'Estado inválido',
    })
    .default('available'),
  
  bedrooms: z
    .number({
      required_error: 'El número de dormitorios es requerido',
      invalid_type_error: 'El número de dormitorios debe ser un número',
    })
    .int('El número de dormitorios debe ser un entero')
    .min(0, 'El número de dormitorios no puede ser negativo')
    .max(20, 'El número de dormitorios no puede ser mayor a 20'),
  
  bathrooms: z
    .number({
      required_error: 'El número de baños es requerido',
      invalid_type_error: 'El número de baños debe ser un número',
    })
    .min(0, 'El número de baños no puede ser negativo')
    .max(20, 'El número de baños no puede ser mayor a 20'),
  
  area_m2: z
    .number({
      required_error: 'El área es requerida',
      invalid_type_error: 'El área debe ser un número',
    })
    .positive('El área debe ser mayor a 0')
    .max(100000, 'El área no puede ser mayor a 100,000 m²'),
  
  // Optional fields
  address: z
    .string()
    .max(500, 'La dirección no puede tener más de 500 caracteres')
    .optional(),
  
  latitude: z
    .number()
    .min(-90, 'Latitud inválida')
    .max(90, 'Latitud inválida')
    .optional()
    .nullable(),
  
  longitude: z
    .number()
    .min(-180, 'Longitud inválida')
    .max(180, 'Longitud inválida')
    .optional()
    .nullable(),
  
  featured: z
    .boolean()
    .default(false),
  
  parking_spots: z
    .number()
    .int('Los espacios de parqueo deben ser un entero')
    .min(0, 'Los espacios de parqueo no pueden ser negativos')
    .max(50, 'Los espacios de parqueo no pueden ser mayor a 50')
    .optional()
    .nullable(),
  
  has_garden: z
    .boolean()
    .default(false),
  
  has_pool: z
    .boolean()
    .default(false),
  
  has_balcony: z
    .boolean()
    .default(false),
  
  has_elevator: z
    .boolean()
    .default(false),
  
  year_built: z
    .number()
    .int('El año de construcción debe ser un entero')
    .min(1800, 'Año de construcción muy antiguo')
    .max(new Date().getFullYear() + 5, 'Año de construcción muy futuro')
    .optional()
    .nullable(),
});

/**
 * Property search/filter validation schema
 */
export const propertyFilterSchema = z.object({
  query: z
    .string()
    .max(200, 'La búsqueda no puede tener más de 200 caracteres')
    .optional(),
  
  province: z
    .enum(ECUADOR_PROVINCES)
    .optional(),
  
  city: z
    .string()
    .max(100, 'La ciudad no puede tener más de 100 caracteres')
    .optional(),
  
  type: z
    .enum(PROPERTY_TYPES)
    .optional(),
  
  status: z
    .enum(PROPERTY_STATUS)
    .optional(),
  
  min_price: z
    .number()
    .positive('El precio mínimo debe ser mayor a 0')
    .optional(),
  
  max_price: z
    .number()
    .positive('El precio máximo debe ser mayor a 0')
    .optional(),
  
  min_bedrooms: z
    .number()
    .int('El número mínimo de dormitorios debe ser un entero')
    .min(0, 'El número mínimo de dormitorios no puede ser negativo')
    .optional(),
  
  max_bedrooms: z
    .number()
    .int('El número máximo de dormitorios debe ser un entero')
    .min(0, 'El número máximo de dormitorios no puede ser negativo')
    .optional(),
  
  min_bathrooms: z
    .number()
    .min(0, 'El número mínimo de baños no puede ser negativo')
    .optional(),
  
  max_bathrooms: z
    .number()
    .min(0, 'El número máximo de baños no puede ser negativo')
    .optional(),
  
  min_area: z
    .number()
    .positive('El área mínima debe ser mayor a 0')
    .optional(),
  
  max_area: z
    .number()
    .positive('El área máxima debe ser mayor a 0')
    .optional(),
  
  featured: z
    .boolean()
    .optional(),
  
  has_garden: z
    .boolean()
    .optional(),
  
  has_pool: z
    .boolean()
    .optional(),
  
  has_balcony: z
    .boolean()
    .optional(),
  
  has_elevator: z
    .boolean()
    .optional(),
  
  // Pagination
  page: z
    .number()
    .int('La página debe ser un entero')
    .min(1, 'La página debe ser mayor a 0')
    .default(1),
  
  limit: z
    .number()
    .int('El límite debe ser un entero')
    .min(1, 'El límite debe ser mayor a 0')
    .max(100, 'El límite no puede ser mayor a 100')
    .default(20),
  
  // Sorting
  sort_by: z
    .enum(['price', 'area_m2', 'created_at', 'bedrooms', 'bathrooms'])
    .default('created_at'),
  
  sort_order: z
    .enum(['asc', 'desc'])
    .default('desc'),
}).refine((data) => {
  // Validate price range
  if (data.min_price && data.max_price && data.min_price > data.max_price) {
    return false;
  }
  // Validate bedroom range
  if (data.min_bedrooms && data.max_bedrooms && data.min_bedrooms > data.max_bedrooms) {
    return false;
  }
  // Validate bathroom range
  if (data.min_bathrooms && data.max_bathrooms && data.min_bathrooms > data.max_bathrooms) {
    return false;
  }
  // Validate area range
  if (data.min_area && data.max_area && data.min_area > data.max_area) {
    return false;
  }
  return true;
}, {
  message: 'Los valores mínimos no pueden ser mayores que los máximos',
});

/**
 * Image upload validation schema
 */
export const imageUploadSchema = z.object({
  property_id: z
    .string()
    .uuid('ID de propiedad inválido'),
  
  alt_text: z
    .string()
    .max(200, 'El texto alternativo no puede tener más de 200 caracteres')
    .optional(),
  
  sort_order: z
    .number()
    .int('El orden debe ser un entero')
    .min(0, 'El orden no puede ser negativo')
    .optional(),
});

// Export constants for use in components
export { ECUADOR_PROVINCES, PROPERTY_TYPES, PROPERTY_STATUS };

// TypeScript types derived from schemas
export type PropertyFormData = z.infer<typeof propertySchema>;
export type PropertyFilterData = z.infer<typeof propertyFilterSchema>;
export type ImageUploadData = z.infer<typeof imageUploadSchema>;