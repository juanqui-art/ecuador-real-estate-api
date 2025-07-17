// Authentication components
export { ProtectedRoute } from './protected-route';

// Role-based access control
export { 
  RoleGuard,
  AdminOnly,
  AgencyOnly,
  AgentOnly,
  OwnerOnly,
  BuyerOnly,
  CanManageUsers,
  CanManageProperties,
  CanViewAnalytics,
  CanManageAgencies,
  CanCreateProperties,
  CanDeleteProperties,
  ConditionalRoleGuard,
  RoleBasedStyles
} from './role-guard';

// Types
export type { UserRole } from '@shared/types/auth';