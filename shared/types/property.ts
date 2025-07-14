/**
 * TypeScript types generated from Go structs for realty-core project
 * Generated from: internal/domain/property.go, user.go, agency.go
 * 
 * Este archivo contiene todos los tipos TypeScript correspondientes a las estructuras Go
 * del backend realty-core, incluyendo validaciones Zod para runtime type checking.
 */

import { z } from 'zod';

// ==============================================================================
// ENUMS - Tipos enumerados del sistema
// ==============================================================================

/**
 * Tipos de propiedades inmobiliarias disponibles
 */
export enum PropertyType {
  HOUSE = 'house',
  APARTMENT = 'apartment',
  LAND = 'land',
  COMMERCIAL = 'commercial',
}

/**
 * Estados de disponibilidad de una propiedad
 */
export enum PropertyStatus {
  AVAILABLE = 'available',
  SOLD = 'sold',
  RENTED = 'rented',
  RESERVED = 'reserved',
}

/**
 * Condición/estado físico de la propiedad
 */
export enum PropertyCondition {
  NEW = 'new',
  USED = 'used',
  RENOVATED = 'renovated',
}

/**
 * Precisión de ubicación GPS
 */
export enum LocationPrecision {
  EXACT = 'exact',
  APPROXIMATE = 'approximate',
  SECTOR = 'sector',
}

/**
 * Roles de usuario en el sistema
 */
export enum UserRole {
  ADMIN = 'admin',
  AGENCY = 'agency',
  AGENT = 'agent',
  OWNER = 'seller', // En BD es 'seller' no 'owner'
  BUYER = 'buyer',
}

/**
 * Estados de cuenta de usuario
 */
export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
  PENDING = 'pending',
}

/**
 * Estados de agencia inmobiliaria
 */
export enum AgencyStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
  PENDING = 'pending',
}

/**
 * Provincias de Ecuador para validación
 */
export const ECUADOR_PROVINCES = [
  'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo',
  'Cotopaxi', 'El Oro', 'Esmeraldas', 'Galápagos',
  'Guayas', 'Imbabura', 'Loja', 'Los Ríos', 'Manabí',
  'Morona Santiago', 'Napo', 'Orellana', 'Pastaza',
  'Pichincha', 'Santa Elena', 'Santo Domingo',
  'Sucumbíos', 'Tungurahua', 'Zamora Chinchipe',
] as const;

export type EcuadorProvince = typeof ECUADOR_PROVINCES[number];

// ==============================================================================
// INTERFACES - Tipos principales del sistema
// ==============================================================================

/**
 * Propiedad inmobiliaria completa
 * Corresponde a la estructura Property en Go
 */
export interface Property {
  /** ID único de la propiedad */
  id: string;
  /** Slug SEO-friendly generado automáticamente */
  slug: string;
  /** Título de la propiedad */
  title: string;
  /** Descripción detallada */
  description: string;
  /** Precio de venta en USD */
  price: number;
  /** Provincia de Ecuador */
  province: string;
  /** Ciudad dentro de la provincia */
  city: string;
  /** Sector específico (opcional) */
  sector: string | null;
  /** Dirección completa (opcional) */
  address: string | null;
  /** Coordenada GPS - latitud (opcional) */
  latitude: number | null;
  /** Coordenada GPS - longitud (opcional) */
  longitude: number | null;
  /** Precisión de la ubicación GPS */
  location_precision: LocationPrecision;
  /** Tipo de propiedad */
  type: PropertyType;
  /** Estado de disponibilidad */
  status: PropertyStatus;
  /** Número de dormitorios */
  bedrooms: number;
  /** Número de baños (puede ser decimal, ej: 2.5) */
  bathrooms: number;
  /** Área en metros cuadrados */
  area_m2: number;
  /** Número de espacios de parqueadero */
  parking_spaces: number;
  /** URL de la imagen principal (opcional) */
  main_image: string | null;
  /** Array de URLs de imágenes */
  images: string[];
  /** URL del video tour (opcional) */
  video_tour: string | null;
  /** URL del tour 360° (opcional) */
  tour_360: string | null;
  /** Precio de alquiler mensual (opcional) */
  rent_price: number | null;
  /** Gastos comunes mensuales (opcional) */
  common_expenses: number | null;
  /** Precio por metro cuadrado calculado (opcional) */
  price_per_m2: number | null;
  /** Año de construcción (opcional) */
  year_built: number | null;
  /** Número de pisos de la propiedad (opcional) */
  floors: number | null;
  /** Condición/estado de la propiedad */
  property_status: PropertyCondition;
  /** Propiedad viene amoblada */
  furnished: boolean;
  /** Tiene garaje */
  garage: boolean;
  /** Tiene piscina */
  pool: boolean;
  /** Tiene jardín */
  garden: boolean;
  /** Tiene terraza */
  terrace: boolean;
  /** Tiene balcón */
  balcony: boolean;
  /** Tiene seguridad/vigilancia */
  security: boolean;
  /** Tiene ascensor */
  elevator: boolean;
  /** Tiene aire acondicionado */
  air_conditioning: boolean;
  /** Tags de búsqueda */
  tags: string[];
  /** Propiedad destacada */
  featured: boolean;
  /** Contador de visualizaciones */
  view_count: number;
  /** ID de la empresa inmobiliaria (opcional) */
  real_estate_company_id: string | null;
  /** ID del propietario (opcional) */
  owner_id: string | null;
  /** ID del agente asignado (opcional) */
  agent_id: string | null;
  /** ID de la agencia que maneja la propiedad (opcional) */
  agency_id: string | null;
  /** ID del usuario que creó la propiedad (opcional) */
  created_by: string | null;
  /** ID del usuario que actualizó la propiedad (opcional) */
  updated_by: string | null;
  /** Fecha de creación */
  created_at: string;
  /** Fecha de última actualización */
  updated_at: string;
}

/**
 * Usuario del sistema con roles y permisos
 * Corresponde a la estructura User en Go
 */
export interface User {
  /** ID único del usuario */
  id: string;
  /** Nombres */
  first_name: string;
  /** Apellidos */
  last_name: string;
  /** Email único */
  email: string;
  /** Teléfono (opcional) */
  phone: string | null;
  /** Cédula de identidad ecuatoriana (opcional) */
  cedula: string | null;
  /** Fecha de nacimiento (opcional) */
  date_of_birth: string | null;
  /** Rol del usuario en el sistema */
  role: UserRole;
  /** Estado de la cuenta */
  active: boolean;
  /** Presupuesto mínimo para búsqueda (opcional) */
  min_budget: number | null;
  /** Presupuesto máximo para búsqueda (opcional) */
  max_budget: number | null;
  /** Provincias de interés para búsqueda */
  preferred_provinces: string[];
  /** Tipos de propiedad preferidos */
  preferred_property_types: string[];
  /** URL del avatar (opcional) */
  avatar_url: string | null;
  /** Biografía del usuario (opcional) */
  bio: string | null;
  /** ID de la empresa inmobiliaria (opcional) */
  real_estate_company_id: string | null;
  /** Recibir notificaciones */
  receive_notifications: boolean;
  /** Recibir newsletter */
  receive_newsletter: boolean;
  /** ID de la agencia (solo para agentes) */
  agency_id: string | null;
  /** Email verificado */
  email_verified: boolean;
  /** Último inicio de sesión (opcional) */
  last_login: string | null;
  /** Fecha de eliminación soft (opcional) */
  deleted_at: string | null;
  /** Estado de la cuenta */
  status: UserStatus;
  /** Fecha de creación */
  created_at: string;
  /** Fecha de última actualización */
  updated_at: string;
}

/**
 * Agencia inmobiliaria
 * Corresponde a la estructura Agency en Go
 */
export interface Agency {
  /** ID único de la agencia */
  id: string;
  /** Nombre de la agencia */
  name: string;
  /** Email de contacto */
  email: string;
  /** Teléfono de contacto */
  phone: string;
  /** Dirección física */
  address: string;
  /** Ciudad */
  city: string;
  /** Provincia */
  province: string;
  /** Licencia comercial */
  license: string;
  /** RUC ecuatoriano */
  ruc: string;
  /** Sitio web (opcional) */
  website: string | null;
  /** Descripción de la agencia (opcional) */
  description: string | null;
  /** Logo de la agencia (opcional) */
  logo: string | null;
  /** URL del logo (opcional) */
  logo_url: string | null;
  /** Estado de la agencia */
  status: AgencyStatus;
  /** Agencia activa */
  active: boolean;
  /** ID del propietario/administrador */
  owner_id: string;
  /** Número de licencia */
  license_number: string;
  /** Fecha de expiración de licencia (opcional) */
  license_expiry: string | null;
  /** Porcentaje de comisión */
  commission: number;
  /** Horarios de atención */
  business_hours?: Record<string, string>;
  /** Redes sociales */
  social_media?: Record<string, string>;
  /** Especialidades */
  specialties?: string[];
  /** Áreas de servicio */
  service_areas?: string[];
  /** Fecha de creación */
  created_at: string;
  /** Fecha de última actualización */
  updated_at: string;
  /** Fecha de eliminación soft (opcional) */
  deleted_at: string | null;
}

/**
 * Parámetros de paginación para consultas
 */
export interface PaginationParams {
  /** Página actual (base 1) */
  page: number;
  /** Tamaño de página */
  page_size: number;
  /** Campo por el cual ordenar */
  sort_by: string;
  /** Ordenar descendente */
  sort_desc: boolean;
}

/**
 * Metadatos de paginación en respuestas
 */
export interface Pagination {
  /** Página actual */
  current_page: number;
  /** Tamaño de página */
  page_size: number;
  /** Total de páginas */
  total_pages: number;
  /** Total de registros */
  total_records: number;
  /** Hay página siguiente */
  has_next: boolean;
  /** Hay página anterior */
  has_prev: boolean;
}

/**
 * Respuesta paginada genérica
 */
export interface PaginatedResponse<T> {
  /** Datos de la página */
  data: T[];
  /** Metadatos de paginación */
  pagination: Pagination;
}

/**
 * Filtros de búsqueda para propiedades
 */
export interface PropertySearchFilters {
  /** Búsqueda por texto */
  query?: string;
  /** Precio mínimo */
  min_price?: number;
  /** Precio máximo */
  max_price?: number;
  /** Tipos de propiedad */
  property_types?: PropertyType[];
  /** Provincias */
  provinces?: string[];
  /** Ciudades */
  cities?: string[];
  /** Sectores */
  sectors?: string[];
  /** Dormitorios mínimos */
  min_bedrooms?: number;
  /** Dormitorios máximos */
  max_bedrooms?: number;
  /** Baños mínimos */
  min_bathrooms?: number;
  /** Baños máximos */
  max_bathrooms?: number;
  /** Área mínima */
  min_area?: number;
  /** Área máxima */
  max_area?: number;
  /** Estados de propiedad */
  status?: PropertyStatus[];
  /** Solo destacadas */
  featured?: boolean;
  /** Filtros por rol */
  owner_id?: string;
  agent_id?: string;
  agency_id?: string;
  created_by?: string;
  /** Filtros de características */
  has_pool?: boolean;
  has_garden?: boolean;
  has_terrace?: boolean;
  has_balcony?: boolean;
  has_security?: boolean;
  has_elevator?: boolean;
  has_air_condition?: boolean;
  has_parking?: boolean;
  min_parking_spaces?: number;
  furnished?: boolean;
  tags?: string[];
  /** Paginación */
  pagination?: PaginationParams;
}

/**
 * Propiedad con relaciones completas
 */
export interface PropertyWithRelations {
  /** Datos de la propiedad */
  property: Property;
  /** Propietario (opcional) */
  owner?: User;
  /** Agente asignado (opcional) */
  agent?: User;
  /** Agencia que maneja (opcional) */
  agency?: Agency;
}

/**
 * Parámetros de búsqueda de usuarios
 */
export interface UserSearchParams {
  /** Búsqueda por texto */
  query?: string;
  /** Filtrar por rol */
  role?: UserRole;
  /** Filtrar por estado */
  status?: UserStatus;
  /** Filtrar por activo */
  active?: boolean;
  /** Filtrar por provincia */
  province?: string;
  /** Filtrar por múltiples provincias */
  provinces?: string[];
  /** Filtrar por ciudad */
  city?: string;
  /** Filtrar por agencia */
  agency_id?: string;
  /** Presupuesto mínimo */
  min_budget?: number;
  /** Presupuesto máximo */
  max_budget?: number;
  /** Paginación */
  pagination?: PaginationParams;
  /** Compatibilidad con parámetros directos */
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_desc?: boolean;
}

/**
 * Estadísticas de usuarios
 */
export interface UserStats {
  /** Total de usuarios */
  total_users: number;
  /** Usuarios activos */
  active_users: number;
  /** Usuarios inactivos */
  inactive_users: number;
  /** Usuarios suspendidos */
  suspended_users: number;
  /** Usuarios pendientes */
  pending_users: number;
  /** Cantidad de administradores */
  admin_count: number;
  /** Cantidad de agencias */
  agency_count: number;
  /** Cantidad de agentes */
  agent_count: number;
  /** Cantidad de propietarios */
  owner_count: number;
  /** Cantidad de compradores */
  buyer_count: number;
  /** Emails verificados */
  email_verified: number;
  /** Usuarios con presupuesto definido */
  with_budget: number;
  /** Agentes asociados a agencias */
  associated_agents: number;
  /** Distribución por rol */
  users_by_role: Record<string, number>;
  /** Distribución por provincia */
  users_by_province: Record<string, number>;
  /** Nuevos usuarios este mes */
  new_users_this_month: number;
  /** Edad promedio */
  average_age: number;
  /** Distribución por género */
  gender_distribution: Record<string, number>;
}

// ==============================================================================
// ESQUEMAS ZOD - Validaciones runtime
// ==============================================================================

/**
 * Esquema de validación para coordenadas de Ecuador
 */
export const EcuadorCoordinatesSchema = z.object({
  latitude: z.number().min(-5.0).max(2.0).describe('Latitud dentro del territorio ecuatoriano'),
  longitude: z.number().min(-92.0).max(-75.0).describe('Longitud dentro del territorio ecuatoriano'),
});

/**
 * Esquema de validación para RUC ecuatoriano
 */
export const EcuadorRUCSchema = z.string()
  .regex(/^[0-9]{10}001$/, 'RUC debe tener 13 dígitos terminados en 001')
  .describe('RUC ecuatoriano válido');

/**
 * Esquema de validación para cédula ecuatoriana
 */
export const EcuadorCedulaSchema = z.string()
  .regex(/^[0-9]{10}$/, 'Cédula debe tener 10 dígitos')
  .describe('Cédula ecuatoriana válida');

/**
 * Esquema de validación para parámetros de paginación
 */
export const PaginationParamsSchema = z.object({
  page: z.number().min(1).default(1),
  page_size: z.number().min(1).max(100).default(20),
  sort_by: z.string().default('created_at'),
  sort_desc: z.boolean().default(true),
});

/**
 * Esquema de validación para Property
 */
export const PropertySchema = z.object({
  id: z.string().uuid(),
  slug: z.string().min(1),
  title: z.string().min(1).max(255),
  description: z.string().min(1),
  price: z.number().positive(),
  province: z.enum(ECUADOR_PROVINCES),
  city: z.string().min(1),
  sector: z.string().nullable(),
  address: z.string().nullable(),
  latitude: z.number().nullable(),
  longitude: z.number().nullable(),
  location_precision: z.nativeEnum(LocationPrecision),
  type: z.nativeEnum(PropertyType),
  status: z.nativeEnum(PropertyStatus),
  bedrooms: z.number().min(0),
  bathrooms: z.number().min(0),
  area_m2: z.number().positive(),
  parking_spaces: z.number().min(0),
  main_image: z.string().url().nullable(),
  images: z.array(z.string().url()),
  video_tour: z.string().url().nullable(),
  tour_360: z.string().url().nullable(),
  rent_price: z.number().positive().nullable(),
  common_expenses: z.number().positive().nullable(),
  price_per_m2: z.number().positive().nullable(),
  year_built: z.number().min(1900).max(new Date().getFullYear()).nullable(),
  floors: z.number().min(1).nullable(),
  property_status: z.nativeEnum(PropertyCondition),
  furnished: z.boolean(),
  garage: z.boolean(),
  pool: z.boolean(),
  garden: z.boolean(),
  terrace: z.boolean(),
  balcony: z.boolean(),
  security: z.boolean(),
  elevator: z.boolean(),
  air_conditioning: z.boolean(),
  tags: z.array(z.string()),
  featured: z.boolean(),
  view_count: z.number().min(0),
  real_estate_company_id: z.string().uuid().nullable(),
  owner_id: z.string().uuid().nullable(),
  agent_id: z.string().uuid().nullable(),
  agency_id: z.string().uuid().nullable(),
  created_by: z.string().uuid().nullable(),
  updated_by: z.string().uuid().nullable(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

/**
 * Esquema de validación para User
 */
export const UserSchema = z.object({
  id: z.string().uuid(),
  first_name: z.string().min(2).max(255),
  last_name: z.string().min(2).max(255),
  email: z.string().email(),
  phone: z.string().nullable(),
  cedula: EcuadorCedulaSchema.nullable(),
  date_of_birth: z.string().datetime().nullable(),
  role: z.nativeEnum(UserRole),
  active: z.boolean(),
  min_budget: z.number().positive().nullable(),
  max_budget: z.number().positive().nullable(),
  preferred_provinces: z.array(z.enum(ECUADOR_PROVINCES)),
  preferred_property_types: z.array(z.nativeEnum(PropertyType)),
  avatar_url: z.string().url().nullable(),
  bio: z.string().nullable(),
  real_estate_company_id: z.string().uuid().nullable(),
  receive_notifications: z.boolean(),
  receive_newsletter: z.boolean(),
  agency_id: z.string().uuid().nullable(),
  email_verified: z.boolean(),
  last_login: z.string().datetime().nullable(),
  deleted_at: z.string().datetime().nullable(),
  status: z.nativeEnum(UserStatus),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

/**
 * Esquema de validación para Agency
 */
export const AgencySchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(2).max(255),
  email: z.string().email(),
  phone: z.string().min(1),
  address: z.string().min(1),
  city: z.string().min(1),
  province: z.enum(ECUADOR_PROVINCES),
  license: z.string().min(1),
  ruc: EcuadorRUCSchema,
  website: z.string().url().nullable(),
  description: z.string().nullable(),
  logo: z.string().nullable(),
  logo_url: z.string().url().nullable(),
  status: z.nativeEnum(AgencyStatus),
  active: z.boolean(),
  owner_id: z.string().uuid(),
  license_number: z.string().min(1),
  license_expiry: z.string().datetime().nullable(),
  commission: z.number().min(0).max(100),
  business_hours: z.record(z.string()).optional(),
  social_media: z.record(z.string()).optional(),
  specialties: z.array(z.string()).optional(),
  service_areas: z.array(z.string()).optional(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
  deleted_at: z.string().datetime().nullable(),
});

/**
 * Esquema de validación para filtros de búsqueda de propiedades
 */
export const PropertySearchFiltersSchema = z.object({
  query: z.string().optional(),
  min_price: z.number().positive().optional(),
  max_price: z.number().positive().optional(),
  property_types: z.array(z.nativeEnum(PropertyType)).optional(),
  provinces: z.array(z.enum(ECUADOR_PROVINCES)).optional(),
  cities: z.array(z.string()).optional(),
  sectors: z.array(z.string()).optional(),
  min_bedrooms: z.number().min(0).optional(),
  max_bedrooms: z.number().min(0).optional(),
  min_bathrooms: z.number().min(0).optional(),
  max_bathrooms: z.number().min(0).optional(),
  min_area: z.number().positive().optional(),
  max_area: z.number().positive().optional(),
  status: z.array(z.nativeEnum(PropertyStatus)).optional(),
  featured: z.boolean().optional(),
  owner_id: z.string().uuid().optional(),
  agent_id: z.string().uuid().optional(),
  agency_id: z.string().uuid().optional(),
  created_by: z.string().uuid().optional(),
  has_pool: z.boolean().optional(),
  has_garden: z.boolean().optional(),
  has_terrace: z.boolean().optional(),
  has_balcony: z.boolean().optional(),
  has_security: z.boolean().optional(),
  has_elevator: z.boolean().optional(),
  has_air_condition: z.boolean().optional(),
  has_parking: z.boolean().optional(),
  min_parking_spaces: z.number().min(0).optional(),
  furnished: z.boolean().optional(),
  tags: z.array(z.string()).optional(),
  pagination: PaginationParamsSchema.optional(),
});

/**
 * Esquema de validación para parámetros de búsqueda de usuarios
 */
export const UserSearchParamsSchema = z.object({
  query: z.string().optional(),
  role: z.nativeEnum(UserRole).optional(),
  status: z.nativeEnum(UserStatus).optional(),
  active: z.boolean().optional(),
  province: z.enum(ECUADOR_PROVINCES).optional(),
  provinces: z.array(z.enum(ECUADOR_PROVINCES)).optional(),
  city: z.string().optional(),
  agency_id: z.string().uuid().optional(),
  min_budget: z.number().positive().optional(),
  max_budget: z.number().positive().optional(),
  pagination: PaginationParamsSchema.optional(),
  page: z.number().min(1).optional(),
  page_size: z.number().min(1).max(100).optional(),
  sort_by: z.string().optional(),
  sort_desc: z.boolean().optional(),
});

// ==============================================================================
// UTILITY TYPES - Tipos de utilidad
// ==============================================================================

/**
 * Tipo para crear una propiedad (omitiendo campos generados)
 */
export type CreatePropertyInput = Omit<
  Property,
  'id' | 'slug' | 'created_at' | 'updated_at' | 'view_count' | 'images' | 'main_image'
>;

/**
 * Tipo para actualizar una propiedad
 */
export type UpdatePropertyInput = Partial<
  Omit<Property, 'id' | 'slug' | 'created_at' | 'updated_at' | 'view_count'>
>;

/**
 * Tipo para crear un usuario
 */
export type CreateUserInput = Omit<
  User,
  'id' | 'created_at' | 'updated_at' | 'email_verified' | 'last_login' | 'deleted_at'
>;

/**
 * Tipo para actualizar un usuario
 */
export type UpdateUserInput = Partial<
  Omit<User, 'id' | 'created_at' | 'updated_at' | 'email'>
>;

/**
 * Tipo para crear una agencia
 */
export type CreateAgencyInput = Omit<
  Agency,
  'id' | 'created_at' | 'updated_at' | 'deleted_at' | 'active'
>;

/**
 * Tipo para actualizar una agencia
 */
export type UpdateAgencyInput = Partial<
  Omit<Agency, 'id' | 'created_at' | 'updated_at' | 'deleted_at' | 'ruc'>
>;

/**
 * Tipo para respuesta de autenticación
 */
export interface AuthResponse {
  /** Token de acceso JWT */
  access_token: string;
  /** Token de renovación */
  refresh_token: string;
  /** Tipo de token */
  token_type: string;
  /** Duración en segundos */
  expires_in: number;
  /** Datos del usuario autenticado */
  user: User;
}

/**
 * Tipo para login
 */
export interface LoginCredentials {
  /** Email del usuario */
  email: string;
  /** Contraseña */
  password: string;
}

/**
 * Tipo para respuesta de error de API
 */
export interface ApiError {
  /** Mensaje de error */
  message: string;
  /** Código de error */
  code?: string;
  /** Detalles adicionales */
  details?: Record<string, any>;
  /** Timestamp del error */
  timestamp: string;
}

// ==============================================================================
// CONSTANTS - Constantes útiles
// ==============================================================================

/**
 * Valores por defecto para paginación
 */
export const DEFAULT_PAGINATION: PaginationParams = {
  page: 1,
  page_size: 20,
  sort_by: 'created_at',
  sort_desc: true,
};

/**
 * Límites de validación
 */
export const VALIDATION_LIMITS = {
  MIN_TITLE_LENGTH: 1,
  MAX_TITLE_LENGTH: 255,
  MIN_DESCRIPTION_LENGTH: 1,
  MAX_DESCRIPTION_LENGTH: 5000,
  MIN_PRICE: 0.01,
  MAX_PRICE: 999999999,
  MIN_AREA: 0.01,
  MAX_AREA: 999999,
  MIN_BEDROOMS: 0,
  MAX_BEDROOMS: 50,
  MIN_BATHROOMS: 0,
  MAX_BATHROOMS: 50,
  MIN_PARKING_SPACES: 0,
  MAX_PARKING_SPACES: 50,
  MIN_YEAR_BUILT: 1900,
  MAX_YEAR_BUILT: new Date().getFullYear(),
  MIN_FLOORS: 1,
  MAX_FLOORS: 200,
  MIN_COMMISSION: 0,
  MAX_COMMISSION: 100,
} as const;

/**
 * Ciudades principales por provincia (muestra)
 */
export const MAJOR_CITIES_BY_PROVINCE: Record<string, string[]> = {
  'Guayas': ['Guayaquil', 'Samborondón', 'Durán', 'Milagro', 'Daule'],
  'Pichincha': ['Quito', 'Sangolquí', 'Cayambe', 'Tabacundo', 'Pedro Moncayo'],
  'Azuay': ['Cuenca', 'Gualaceo', 'Paute', 'Chordeleg', 'Sigsig'],
  'Manabí': ['Portoviejo', 'Manta', 'Bahía de Caráquez', 'Chone', 'Jipijapa'],
  'El Oro': ['Machala', 'Huaquillas', 'Pasaje', 'Santa Rosa', 'Arenillas'],
  'Imbabura': ['Ibarra', 'Otavalo', 'Cotacachi', 'Atuntaqui', 'Pimampiro'],
  'Tungurahua': ['Ambato', 'Baños', 'Pelileo', 'Píllaro', 'Patate'],
  'Esmeraldas': ['Esmeraldas', 'Atacames', 'Sua', 'Muisne', 'Quinindé'],
  'Loja': ['Loja', 'Catamayo', 'Cariamanga', 'Macará', 'Catacocha'],
  'Chimborazo': ['Riobamba', 'Alausí', 'Chambo', 'Colta', 'Guano'],
} as const;

/**
 * Tipos de propiedad con sus características típicas
 */
export const PROPERTY_TYPE_CHARACTERISTICS: Record<PropertyType, {
  typical_bedrooms: number[];
  typical_bathrooms: number[];
  typical_area_range: [number, number];
  common_features: string[];
}> = {
  [PropertyType.HOUSE]: {
    typical_bedrooms: [2, 3, 4, 5, 6],
    typical_bathrooms: [1, 2, 3, 4],
    typical_area_range: [80, 500],
    common_features: ['garden', 'garage', 'terrace', 'balcony'],
  },
  [PropertyType.APARTMENT]: {
    typical_bedrooms: [1, 2, 3, 4],
    typical_bathrooms: [1, 2, 3],
    typical_area_range: [45, 200],
    common_features: ['balcony', 'elevator', 'security', 'parking'],
  },
  [PropertyType.LAND]: {
    typical_bedrooms: [0],
    typical_bathrooms: [0],
    typical_area_range: [100, 10000],
    common_features: [],
  },
  [PropertyType.COMMERCIAL]: {
    typical_bedrooms: [0],
    typical_bathrooms: [1, 2, 3],
    typical_area_range: [50, 1000],
    common_features: ['parking', 'security', 'elevator'],
  },
} as const;

/**
 * Exportar tipos de utilidad para mejor developer experience
 */
export type PropertyResponse = Property;
export type UserResponse = User;
export type AgencyResponse = Agency;
export type PaginatedPropertyResponse = PaginatedResponse<Property>;
export type PaginatedUserResponse = PaginatedResponse<User>;
export type PaginatedAgencyResponse = PaginatedResponse<Agency>;