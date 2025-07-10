-- Migration: Create agencies table for real estate agencies
-- Description: Agencies management table with Ecuador-specific validations
-- Date: 2025-01-10

CREATE TABLE IF NOT EXISTS agencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    city VARCHAR(100) NOT NULL,
    province VARCHAR(100) NOT NULL,
    license VARCHAR(13) UNIQUE NOT NULL, -- Ecuador RUC format
    website TEXT,
    description TEXT,
    logo TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('active', 'inactive', 'suspended', 'pending')),
    owner_id UUID NOT NULL, -- User who owns/manages the agency
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_email ON agencies(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_license ON agencies(license);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_status ON agencies(status);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_owner_id ON agencies(owner_id);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_province ON agencies(province);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_city ON agencies(city);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_created_at ON agencies(created_at);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_deleted_at ON agencies(deleted_at);

-- Create partial indexes for active agencies
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_active_province ON agencies(province) WHERE status = 'active';
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_agencies_active_city ON agencies(city) WHERE status = 'active';

-- Add constraint for Ecuador provinces
ALTER TABLE agencies ADD CONSTRAINT chk_ecuador_province 
    CHECK (province IN (
        'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi',
        'El Oro', 'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja',
        'Los Ríos', 'Manabí', 'Morona Santiago', 'Napo', 'Orellana', 'Pastaza',
        'Pichincha', 'Santa Elena', 'Santo Domingo', 'Sucumbíos', 'Tungurahua', 'Zamora Chinchipe'
    ));

-- Add constraint for Ecuador RUC format (13 digits)
ALTER TABLE agencies ADD CONSTRAINT chk_ruc_format 
    CHECK (license ~ '^[0-9]{13}$');

-- Add constraint for Ecuador phone format
ALTER TABLE agencies ADD CONSTRAINT chk_phone_format 
    CHECK (phone ~ '^(\+593|0)[0-9]{9}$');

-- Add constraint for website URL format
ALTER TABLE agencies ADD CONSTRAINT chk_website_format 
    CHECK (website IS NULL OR website ~ '^https?://[^\s/$.?#].[^\s]*$');

-- Add constraint for description length
ALTER TABLE agencies ADD CONSTRAINT chk_description_length 
    CHECK (description IS NULL OR length(description) <= 1000);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_agencies_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_agencies_updated_at
    BEFORE UPDATE ON agencies
    FOR EACH ROW EXECUTE FUNCTION update_agencies_updated_at();

-- Create function to get agency statistics
CREATE OR REPLACE FUNCTION get_agency_stats(agency_uuid UUID)
RETURNS TABLE(
    total_properties INTEGER,
    active_properties INTEGER,
    total_agents INTEGER,
    active_agents INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COALESCE(p.total, 0) as total_properties,
        COALESCE(p.active, 0) as active_properties,
        COALESCE(a.total, 0) as total_agents,
        COALESCE(a.active, 0) as active_agents
    FROM 
        (SELECT 
            COUNT(*) as total,
            COUNT(CASE WHEN status = 'available' THEN 1 END) as active
         FROM properties 
         WHERE agency_id = agency_uuid AND deleted_at IS NULL) p
    FULL OUTER JOIN
        (SELECT 
            COUNT(*) as total,
            COUNT(CASE WHEN status = 'active' THEN 1 END) as active
         FROM users 
         WHERE agency_id = agency_uuid AND deleted_at IS NULL) a
    ON TRUE;
END;
$$ LANGUAGE plpgsql;

-- Create function to check if agency can manage property
CREATE OR REPLACE FUNCTION agency_can_manage_property(
    agency_uuid UUID,
    property_agency_id UUID
) RETURNS BOOLEAN AS $$
BEGIN
    RETURN agency_uuid = property_agency_id;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Create view for agency summary
CREATE OR REPLACE VIEW agency_summary AS
SELECT 
    a.id,
    a.name,
    a.email,
    a.phone,
    a.city,
    a.province,
    a.status,
    a.created_at,
    COUNT(p.id) as total_properties,
    COUNT(CASE WHEN p.status = 'available' THEN 1 END) as available_properties,
    COUNT(u.id) as total_agents,
    COUNT(CASE WHEN u.status = 'active' THEN 1 END) as active_agents
FROM agencies a
LEFT JOIN properties p ON a.id = p.agency_id AND p.deleted_at IS NULL
LEFT JOIN users u ON a.id = u.agency_id AND u.deleted_at IS NULL
WHERE a.deleted_at IS NULL
GROUP BY a.id, a.name, a.email, a.phone, a.city, a.province, a.status, a.created_at;

-- Add comments for documentation
COMMENT ON TABLE agencies IS 'Real estate agencies table with Ecuador-specific validations';
COMMENT ON COLUMN agencies.license IS 'Ecuador RUC (Registro Único de Contribuyentes) - 13 digits';
COMMENT ON COLUMN agencies.phone IS 'Ecuador phone number format: +593xxxxxxxxx or 0xxxxxxxxx';
COMMENT ON COLUMN agencies.province IS 'Ecuador province (must be one of 24 valid provinces)';
COMMENT ON COLUMN agencies.status IS 'Agency status: active, inactive, suspended, pending';
COMMENT ON COLUMN agencies.owner_id IS 'User who owns/manages the agency';
COMMENT ON FUNCTION get_agency_stats IS 'Returns statistics for a specific agency';
COMMENT ON VIEW agency_summary IS 'Summary view of agencies with property and agent counts';