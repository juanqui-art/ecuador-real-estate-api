-- Migration: Update properties table with user relationships
-- Description: Add user ownership and agency management to properties
-- Date: 2025-01-10

-- Add new columns for user relationships
ALTER TABLE properties ADD COLUMN IF NOT EXISTS owner_id UUID;
ALTER TABLE properties ADD COLUMN IF NOT EXISTS agent_id UUID;
ALTER TABLE properties ADD COLUMN IF NOT EXISTS agency_id UUID;
ALTER TABLE properties ADD COLUMN IF NOT EXISTS created_by UUID;
ALTER TABLE properties ADD COLUMN IF NOT EXISTS updated_by UUID;

-- Add foreign key constraints
ALTER TABLE properties ADD CONSTRAINT fk_properties_owner 
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE properties ADD CONSTRAINT fk_properties_agent 
    FOREIGN KEY (agent_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE properties ADD CONSTRAINT fk_properties_agency 
    FOREIGN KEY (agency_id) REFERENCES agencies(id) ON DELETE SET NULL;

ALTER TABLE properties ADD CONSTRAINT fk_properties_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE properties ADD CONSTRAINT fk_properties_updated_by 
    FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL;

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_owner_id ON properties(owner_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_agent_id ON properties(agent_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_agency_id ON properties(agency_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_created_by ON properties(created_by);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_updated_by ON properties(updated_by);

-- Create composite indexes for common queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_owner_status ON properties(owner_id, status) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_agency_status ON properties(agency_id, status) WHERE deleted_at IS NULL;
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_properties_agent_status ON properties(agent_id, status) WHERE deleted_at IS NULL;

-- Add constraint to ensure agent belongs to same agency as property
ALTER TABLE properties ADD CONSTRAINT chk_agent_agency_match 
    CHECK (
        (agent_id IS NULL) OR 
        (agency_id IS NULL) OR 
        (agent_id IS NOT NULL AND agency_id IS NOT NULL AND 
         EXISTS (SELECT 1 FROM users WHERE id = agent_id AND agency_id = properties.agency_id))
    );

-- Create function to check property ownership permissions
CREATE OR REPLACE FUNCTION check_property_permissions(
    user_id UUID,
    user_role VARCHAR(20),
    user_agency_id UUID,
    property_owner_id UUID,
    property_agency_id UUID,
    property_agent_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN CASE user_role
        WHEN 'admin' THEN TRUE
        WHEN 'owner' THEN user_id = property_owner_id
        WHEN 'agency' THEN user_id = property_agency_id
        WHEN 'agent' THEN (user_agency_id = property_agency_id AND user_id = property_agent_id)
        ELSE FALSE
    END;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Create function to get user properties count
CREATE OR REPLACE FUNCTION get_user_properties_count(user_uuid UUID, user_role VARCHAR(20))
RETURNS INTEGER AS $$
DECLARE
    count_result INTEGER;
BEGIN
    IF user_role = 'owner' THEN
        SELECT COUNT(*) INTO count_result
        FROM properties 
        WHERE owner_id = user_uuid AND deleted_at IS NULL;
    ELSIF user_role = 'agent' THEN
        SELECT COUNT(*) INTO count_result
        FROM properties 
        WHERE agent_id = user_uuid AND deleted_at IS NULL;
    ELSIF user_role = 'agency' THEN
        SELECT COUNT(*) INTO count_result
        FROM properties 
        WHERE agency_id = user_uuid AND deleted_at IS NULL;
    ELSE
        count_result := 0;
    END IF;
    
    RETURN count_result;
END;
$$ LANGUAGE plpgsql;

-- Create function to transfer property ownership
CREATE OR REPLACE FUNCTION transfer_property_ownership(
    property_uuid UUID,
    new_owner_id UUID,
    updated_by_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    UPDATE properties 
    SET 
        owner_id = new_owner_id,
        updated_by = updated_by_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = property_uuid AND deleted_at IS NULL;
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- Create function to assign property to agency
CREATE OR REPLACE FUNCTION assign_property_to_agency(
    property_uuid UUID,
    agency_uuid UUID,
    agent_uuid UUID DEFAULT NULL,
    updated_by_id UUID DEFAULT NULL
) RETURNS BOOLEAN AS $$
BEGIN
    -- Validate that agent belongs to agency if provided
    IF agent_uuid IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM users 
            WHERE id = agent_uuid AND agency_id = agency_uuid AND role = 'agent' AND status = 'active'
        ) THEN
            RETURN FALSE;
        END IF;
    END IF;
    
    UPDATE properties 
    SET 
        agency_id = agency_uuid,
        agent_id = agent_uuid,
        updated_by = updated_by_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = property_uuid AND deleted_at IS NULL;
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- Create view for property ownership summary
CREATE OR REPLACE VIEW property_ownership_summary AS
SELECT 
    p.id,
    p.title,
    p.status,
    p.price,
    p.province,
    p.city,
    p.type,
    p.created_at,
    -- Owner information
    u_owner.id as owner_id,
    u_owner.name as owner_name,
    u_owner.email as owner_email,
    -- Agency information
    a.id as agency_id,
    a.name as agency_name,
    a.email as agency_email,
    -- Agent information
    u_agent.id as agent_id,
    u_agent.name as agent_name,
    u_agent.email as agent_email,
    -- Creation information
    u_created.name as created_by_name,
    u_updated.name as updated_by_name
FROM properties p
LEFT JOIN users u_owner ON p.owner_id = u_owner.id
LEFT JOIN agencies a ON p.agency_id = a.id
LEFT JOIN users u_agent ON p.agent_id = u_agent.id
LEFT JOIN users u_created ON p.created_by = u_created.id
LEFT JOIN users u_updated ON p.updated_by = u_updated.id
WHERE p.deleted_at IS NULL;

-- Add comments for documentation
COMMENT ON COLUMN properties.owner_id IS 'User who owns the property';
COMMENT ON COLUMN properties.agent_id IS 'Agent assigned to manage the property';
COMMENT ON COLUMN properties.agency_id IS 'Agency managing the property';
COMMENT ON COLUMN properties.created_by IS 'User who created the property record';
COMMENT ON COLUMN properties.updated_by IS 'User who last updated the property record';
COMMENT ON FUNCTION check_property_permissions IS 'Checks if user has permission to access property based on role';
COMMENT ON FUNCTION get_user_properties_count IS 'Returns count of properties for a user based on their role';
COMMENT ON FUNCTION transfer_property_ownership IS 'Transfers property ownership to a new owner';
COMMENT ON FUNCTION assign_property_to_agency IS 'Assigns property to an agency and optionally an agent';
COMMENT ON VIEW property_ownership_summary IS 'Complete property ownership and management information';