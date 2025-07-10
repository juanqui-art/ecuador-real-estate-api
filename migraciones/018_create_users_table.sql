-- Migration: Create users table for role-based system
-- Description: Core users table with role-based access control
-- Date: 2025-01-10

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'agency', 'agent', 'owner', 'buyer')),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'inactive', 'suspended', 'pending')),
    agency_id UUID REFERENCES agencies(id) ON DELETE SET NULL,
    avatar TEXT,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_agency_id ON users(agency_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Create partial indexes for active users
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_active_role ON users(role) WHERE status = 'active';
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_active_agency ON users(agency_id) WHERE status = 'active' AND role = 'agent';

-- Add constraint to ensure agents have agency_id
ALTER TABLE users ADD CONSTRAINT chk_agent_agency 
    CHECK (
        (role = 'agent' AND agency_id IS NOT NULL) OR 
        (role != 'agent' AND agency_id IS NULL)
    );

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_users_updated_at();

-- Create function to get user role level
CREATE OR REPLACE FUNCTION get_user_role_level(user_role VARCHAR(20))
RETURNS INTEGER AS $$
BEGIN
    RETURN CASE user_role
        WHEN 'admin' THEN 5
        WHEN 'agency' THEN 4
        WHEN 'agent' THEN 3
        WHEN 'owner' THEN 2
        WHEN 'buyer' THEN 1
        ELSE 0
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Create function to check if user can manage property
CREATE OR REPLACE FUNCTION user_can_manage_property(
    user_id UUID,
    user_role VARCHAR(20),
    user_agency_id UUID,
    property_owner_id UUID,
    property_agency_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN CASE user_role
        WHEN 'admin' THEN TRUE
        WHEN 'owner' THEN user_id = property_owner_id
        WHEN 'agency' THEN user_id = property_agency_id
        WHEN 'agent' THEN user_agency_id = property_agency_id
        ELSE FALSE
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Add comments for documentation
COMMENT ON TABLE users IS 'Core users table with role-based access control for real estate system';
COMMENT ON COLUMN users.role IS 'User role: admin, agency, agent, owner, buyer';
COMMENT ON COLUMN users.status IS 'User status: active, inactive, suspended, pending';
COMMENT ON COLUMN users.agency_id IS 'Foreign key to agencies table (only for agents)';
COMMENT ON FUNCTION get_user_role_level IS 'Returns hierarchical level of user role (1-5)';
COMMENT ON FUNCTION user_can_manage_property IS 'Checks if user can manage a specific property based on role';