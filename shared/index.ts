// Export all types and schemas
export * from './types/property';

// Re-export commonly used types
export type {
  Property,
  User,
  Agency,
  UserRole,
  UserStatus,
  AgencyStatus,
  CreatePropertyInput,
  UpdatePropertyInput,
  PropertySearchFilters,
  PaginationParams,
  Pagination,
  AuthResponse,
  LoginCredentials,
  ApiError,
} from './types/property';

// Re-export schemas
export {
  PropertySchema,
  UserSchema,
  AgencySchema,
  EcuadorCoordinatesSchema,
  EcuadorRUCSchema,
  EcuadorCedulaSchema,
} from './types/property';