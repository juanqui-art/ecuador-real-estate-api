// Export all types and schemas
export * from './types/property';
export * from './types/auth';

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

// Re-export auth types
export type {
  LoginResponse,
  RefreshTokenResponse,
  TokenValidationResponse,
  LoginRequest,
  RefreshTokenRequest,
  ChangePasswordRequest,
} from './types/auth';

// Re-export schemas
export {
  PropertySchema,
  UserSchema,
  AgencySchema,
  EcuadorCoordinatesSchema,
  EcuadorRUCSchema,
  EcuadorCedulaSchema,
} from './types/property';