-- Migration 016: Add relations between English tables
-- Date: 2025-01-15
-- Purpose: Create proper foreign key relationships between English tables

-- Update properties table to have proper FK constraint to real_estate_companies
-- (Note: We need to modify the existing constraint if it exists)

-- First check if the FK constraint exists and drop it if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE constraint_name = 'properties_real_estate_company_id_fkey' 
        AND table_name = 'properties'
    ) THEN
        ALTER TABLE properties DROP CONSTRAINT properties_real_estate_company_id_fkey;
    END IF;
END $$;

-- Now add the proper foreign key constraint
ALTER TABLE properties 
ADD CONSTRAINT fk_properties_real_estate_company 
FOREIGN KEY (real_estate_company_id) 
REFERENCES real_estate_companies(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Ensure users table has proper FK constraint to real_estate_companies
-- (This should already be created from the table definition)
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.table_constraints 
        WHERE constraint_name = 'users_real_estate_company_id_fkey' 
        AND table_name = 'users'
    ) THEN
        ALTER TABLE users 
        ADD CONSTRAINT fk_users_real_estate_company 
        FOREIGN KEY (real_estate_company_id) 
        REFERENCES real_estate_companies(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE;
    END IF;
END $$;

-- Create additional useful indexes for the relationships
CREATE INDEX IF NOT EXISTS idx_properties_company_status ON properties(real_estate_company_id, status);
CREATE INDEX IF NOT EXISTS idx_users_company_type ON users(real_estate_company_id, user_type);

-- Create a view for property listings with company information
CREATE OR REPLACE VIEW property_listings AS
SELECT 
    p.id,
    p.slug,
    p.title,
    p.description,
    p.price,
    p.province,
    p.city,
    p.sector,
    p.type,
    p.status,
    p.bedrooms,
    p.bathrooms,
    p.area_m2,
    p.main_image,
    p.featured,
    p.view_count,
    p.created_at,
    p.updated_at,
    rec.id as company_id,
    rec.name as company_name,
    rec.phone as company_phone,
    rec.email as company_email,
    rec.website as company_website
FROM properties p
LEFT JOIN real_estate_companies rec ON p.real_estate_company_id = rec.id
WHERE p.status = 'available';

-- Create a view for agent profiles with company information
CREATE OR REPLACE VIEW agent_profiles AS
SELECT 
    u.id,
    u.first_name,
    u.last_name,
    (u.first_name || ' ' || u.last_name) as full_name,
    u.email,
    u.phone,
    u.bio,
    u.avatar_url,
    u.active,
    u.created_at,
    rec.id as company_id,
    rec.name as company_name,
    rec.phone as company_phone,
    rec.email as company_email,
    rec.website as company_website,
    rec.active as company_active
FROM users u
INNER JOIN real_estate_companies rec ON u.real_estate_company_id = rec.id
WHERE u.user_type = 'agent' AND u.active = TRUE;

-- Create a view for company statistics
CREATE OR REPLACE VIEW company_statistics AS
SELECT 
    rec.id,
    rec.name,
    rec.active,
    rec.created_at,
    COUNT(p.id) as total_properties,
    COUNT(p.id) FILTER (WHERE p.status = 'available') as available_properties,
    COUNT(p.id) FILTER (WHERE p.status = 'sold') as sold_properties,
    COUNT(p.id) FILTER (WHERE p.status = 'rented') as rented_properties,
    COUNT(u.id) as total_agents,
    COUNT(u.id) FILTER (WHERE u.active = TRUE) as active_agents,
    AVG(p.price) FILTER (WHERE p.status = 'available') as avg_property_price,
    MIN(p.price) FILTER (WHERE p.status = 'available') as min_property_price,
    MAX(p.price) FILTER (WHERE p.status = 'available') as max_property_price
FROM real_estate_companies rec
LEFT JOIN properties p ON rec.id = p.real_estate_company_id
LEFT JOIN users u ON rec.id = u.real_estate_company_id AND u.user_type = 'agent'
GROUP BY rec.id, rec.name, rec.active, rec.created_at;

-- Function to get properties by company
CREATE OR REPLACE FUNCTION get_properties_by_company(company_id UUID)
RETURNS TABLE(
    id VARCHAR(36),
    title VARCHAR(255),
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    status VARCHAR(50),
    bedrooms INTEGER,
    bathrooms DECIMAL(3,1),
    area_m2 DECIMAL(10,2),
    featured BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.title,
        p.price,
        p.province,
        p.city,
        p.type,
        p.status,
        p.bedrooms,
        p.bathrooms,
        p.area_m2,
        p.featured,
        p.created_at
    FROM properties p
    WHERE p.real_estate_company_id = company_id
    ORDER BY p.featured DESC, p.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to transfer properties between companies
CREATE OR REPLACE FUNCTION transfer_properties_to_company(
    from_company_id UUID, 
    to_company_id UUID
) RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER;
BEGIN
    -- Verify that both companies exist and are active
    IF NOT EXISTS (SELECT 1 FROM real_estate_companies WHERE id = from_company_id AND active = TRUE) THEN
        RAISE EXCEPTION 'Source company not found or inactive: %', from_company_id;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM real_estate_companies WHERE id = to_company_id AND active = TRUE) THEN
        RAISE EXCEPTION 'Target company not found or inactive: %', to_company_id;
    END IF;
    
    -- Update properties
    UPDATE properties 
    SET real_estate_company_id = to_company_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE real_estate_company_id = from_company_id;
    
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    
    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

-- Function to transfer agents between companies
CREATE OR REPLACE FUNCTION transfer_agents_to_company(
    from_company_id UUID, 
    to_company_id UUID
) RETURNS INTEGER AS $$
DECLARE
    updated_count INTEGER;
BEGIN
    -- Verify that both companies exist and are active
    IF NOT EXISTS (SELECT 1 FROM real_estate_companies WHERE id = from_company_id AND active = TRUE) THEN
        RAISE EXCEPTION 'Source company not found or inactive: %', from_company_id;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM real_estate_companies WHERE id = to_company_id AND active = TRUE) THEN
        RAISE EXCEPTION 'Target company not found or inactive: %', to_company_id;
    END IF;
    
    -- Update agents
    UPDATE users 
    SET real_estate_company_id = to_company_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE real_estate_company_id = from_company_id 
      AND user_type = 'agent';
    
    GET DIAGNOSTICS updated_count = ROW_COUNT;
    
    RETURN updated_count;
END;
$$ LANGUAGE plpgsql;

-- Comments on views and functions
COMMENT ON VIEW property_listings IS 'Properties with company information for public listings';
COMMENT ON VIEW agent_profiles IS 'Agent profiles with their company information';
COMMENT ON VIEW company_statistics IS 'Statistics for each real estate company';
COMMENT ON FUNCTION get_properties_by_company IS 'Get all properties managed by a specific company';
COMMENT ON FUNCTION transfer_properties_to_company IS 'Transfer properties from one company to another';
COMMENT ON FUNCTION transfer_agents_to_company IS 'Transfer agents from one company to another';

-- Update the sample property to link to the sample company
UPDATE properties 
SET real_estate_company_id = (
    SELECT id 
    FROM real_estate_companies 
    WHERE name = 'Example Real Estate S.A.' 
    LIMIT 1
)
WHERE id = 'prop-en-001';

-- Update the sample agent to link to the sample company
UPDATE users 
SET real_estate_company_id = (
    SELECT id 
    FROM real_estate_companies 
    WHERE name = 'Example Real Estate S.A.' 
    LIMIT 1
)
WHERE email = 'carlos.rodriguez@example.com' AND user_type = 'agent';