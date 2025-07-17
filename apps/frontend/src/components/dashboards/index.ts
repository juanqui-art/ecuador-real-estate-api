// Role-based dashboard components
export { AdminDashboard } from './admin-dashboard';
export { AgencyDashboard } from './agency-dashboard';
export { AgentDashboard } from './agent-dashboard';
export { OwnerDashboard } from './owner-dashboard';
export { BuyerDashboard } from './buyer-dashboard';

// Main dashboard system
export { 
  RoleBasedDashboard, 
  RoleBasedComponent, 
  useRolePermissions 
} from './role-based-dashboard';

// Dashboard utilities
export type { UserRole } from '@shared/types/auth';