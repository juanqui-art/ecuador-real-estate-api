-- Migration 015: Create English users table
-- Date: 2025-01-15
-- Purpose: Complete users table in English with Ecuador validations

-- Create users table in English
CREATE TABLE IF NOT EXISTS users (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Personal information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    national_id VARCHAR(10) NOT NULL UNIQUE,  -- Ecuador cedula (10 digits)
    date_of_birth DATE,

    -- User type and status
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('buyer', 'seller', 'agent', 'admin')),
    active BOOLEAN DEFAULT TRUE NOT NULL,

    -- Search preferences (for buyers)
    min_budget DECIMAL(15,2),             -- Minimum budget
    max_budget DECIMAL(15,2),             -- Maximum budget
    preferred_provinces JSONB DEFAULT '[]', -- Provinces of interest (JSON array)
    preferred_property_types JSONB DEFAULT '[]', -- Property types of interest (JSON array)

    -- Profile information
    avatar_url TEXT DEFAULT '' NOT NULL,  -- Avatar image URL
    bio TEXT DEFAULT '' NOT NULL,         -- User biography/description

    -- Relationship with RealEstateCompany (for agents)
    real_estate_company_id UUID REFERENCES real_estate_companies(id),

    -- Notification preferences
    receive_notifications BOOLEAN DEFAULT TRUE NOT NULL,   -- Receive app notifications
    receive_newsletter BOOLEAN DEFAULT FALSE NOT NULL,     -- Receive newsletter emails

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add comments for documentation
COMMENT ON TABLE users IS 'System users (buyers, sellers, agents, admins)';
COMMENT ON COLUMN users.id IS 'Unique identifier for the user';
COMMENT ON COLUMN users.first_name IS 'User first name';
COMMENT ON COLUMN users.last_name IS 'User last name';
COMMENT ON COLUMN users.email IS 'User email address (unique)';
COMMENT ON COLUMN users.phone IS 'User phone number';
COMMENT ON COLUMN users.national_id IS 'Ecuador national ID (cedula)';
COMMENT ON COLUMN users.date_of_birth IS 'User date of birth';
COMMENT ON COLUMN users.user_type IS 'Type of user: buyer, seller, agent, admin';
COMMENT ON COLUMN users.active IS 'Indicates if the user is active';
COMMENT ON COLUMN users.min_budget IS 'Minimum budget for buyers';
COMMENT ON COLUMN users.max_budget IS 'Maximum budget for buyers';
COMMENT ON COLUMN users.preferred_provinces IS 'JSON array of preferred provinces';
COMMENT ON COLUMN users.preferred_property_types IS 'JSON array of preferred property types';
COMMENT ON COLUMN users.avatar_url IS 'URL of user avatar image';
COMMENT ON COLUMN users.bio IS 'User biography or description';
COMMENT ON COLUMN users.real_estate_company_id IS 'FK to real estate company (for agents)';
COMMENT ON COLUMN users.receive_notifications IS 'User wants to receive notifications';
COMMENT ON COLUMN users.receive_newsletter IS 'User wants to receive newsletter';

-- Add constraints for data validation
ALTER TABLE users 
ADD CONSTRAINT chk_email_format 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Constraint to validate Ecuador phone format
ALTER TABLE users 
ADD CONSTRAINT chk_phone_format 
CHECK (phone ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');

-- Constraint to validate national ID format (10 digits)
ALTER TABLE users 
ADD CONSTRAINT chk_national_id_format 
CHECK (national_id ~ '^[0-9]{10}$');

-- Constraint to validate budget logic
ALTER TABLE users 
ADD CONSTRAINT chk_budget_logic 
CHECK (min_budget IS NULL OR max_budget IS NULL OR min_budget <= max_budget);

-- Constraint to validate budget is positive
ALTER TABLE users 
ADD CONSTRAINT chk_budget_positive 
CHECK (min_budget IS NULL OR min_budget >= 0);

ALTER TABLE users 
ADD CONSTRAINT chk_max_budget_positive 
CHECK (max_budget IS NULL OR max_budget >= 0);

-- Constraint to validate avatar URL format
ALTER TABLE users 
ADD CONSTRAINT chk_avatar_url 
CHECK (avatar_url = '' OR avatar_url ~* '^https?://.*\.(jpg|jpeg|png|webp|svg)(\?.*)?$');

-- Create indexes for performance
CREATE INDEX idx_users_user_type ON users(user_type);
CREATE INDEX idx_users_active ON users(active) WHERE active = TRUE;
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_national_id ON users(national_id);
CREATE INDEX idx_users_real_estate_company_id ON users(real_estate_company_id);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Indexes for buyers' search preferences
CREATE INDEX idx_users_min_budget ON users(min_budget) WHERE user_type = 'buyer';
CREATE INDEX idx_users_max_budget ON users(max_budget) WHERE user_type = 'buyer';

-- JSONB indexes for preferences
CREATE INDEX idx_users_preferred_provinces ON users USING gin(preferred_provinces);
CREATE INDEX idx_users_preferred_property_types ON users USING gin(preferred_property_types);

-- Full-text search index for name and bio
CREATE INDEX idx_users_fulltext_search ON users 
USING gin(to_tsvector('english', first_name || ' ' || last_name || ' ' || bio));

-- Function to validate Ecuador national ID (cedula)
CREATE OR REPLACE FUNCTION validate_ecuador_national_id(national_id_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    id_clean TEXT;
    digits INTEGER[];
    sum INTEGER := 0;
    check_digit INTEGER;
    i INTEGER;
BEGIN
    -- Clean national ID (numbers only)
    id_clean := regexp_replace(national_id_input, '[^0-9]', '', 'g');
    
    -- Check length (must be 10 digits)
    IF LENGTH(id_clean) != 10 THEN
        RETURN FALSE;
    END IF;
    
    -- Convert to digit array
    FOR i IN 1..10 LOOP
        digits[i] := CAST(SUBSTRING(id_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- First two digits must be valid province (01-24)
    IF (digits[1] * 10 + digits[2]) < 1 OR (digits[1] * 10 + digits[2]) > 24 THEN
        RETURN FALSE;
    END IF;
    
    -- Third digit must be less than 6 for natural persons
    IF digits[3] >= 6 THEN
        RETURN FALSE;
    END IF;
    
    -- Apply Ecuador national ID algorithm
    FOR i IN 1..9 LOOP
        IF i % 2 = 1 THEN  -- Odd positions (1,3,5,7,9)
            sum := sum + (digits[i] * 2);
            IF (digits[i] * 2) > 9 THEN
                sum := sum - 9;
            END IF;
        ELSE  -- Even positions (2,4,6,8)
            sum := sum + digits[i];
        END IF;
    END LOOP;
    
    check_digit := (10 - (sum % 10)) % 10;
    
    RETURN check_digit = digits[10];
END;
$$ LANGUAGE plpgsql;

-- Add constraint using national ID validation function
ALTER TABLE users 
ADD CONSTRAINT chk_national_id_valid 
CHECK (validate_ecuador_national_id(national_id));

-- Trigger to automatically update updated_at
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to get users by type
CREATE OR REPLACE FUNCTION get_users_by_type(user_type_filter TEXT)
RETURNS TABLE(
    id UUID,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    user_type VARCHAR(20),
    active BOOLEAN,
    created_at TIMESTAMP WITH TIME ZONE
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.first_name,
        u.last_name,
        u.email,
        u.phone,
        u.user_type,
        u.active,
        u.created_at
    FROM users u
    WHERE u.user_type = user_type_filter
      AND u.active = TRUE
    ORDER BY u.first_name, u.last_name;
END;
$$ LANGUAGE plpgsql;

-- Function to search users by name
CREATE OR REPLACE FUNCTION search_users_by_name(search_name TEXT)
RETURNS TABLE(
    id UUID,
    full_name TEXT,
    email VARCHAR(255),
    phone VARCHAR(20),
    user_type VARCHAR(20),
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        (u.first_name || ' ' || u.last_name) as full_name,
        u.email,
        u.phone,
        u.user_type,
        ts_rank(to_tsvector('english', u.first_name || ' ' || u.last_name || ' ' || u.bio), plainto_tsquery('english', search_name)) as rank
    FROM users u
    WHERE u.active = TRUE
      AND to_tsvector('english', u.first_name || ' ' || u.last_name || ' ' || u.bio) @@ plainto_tsquery('english', search_name)
    ORDER BY rank DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get buyers that can afford a property
CREATE OR REPLACE FUNCTION get_buyers_for_property(property_price DECIMAL)
RETURNS TABLE(
    id UUID,
    full_name TEXT,
    email VARCHAR(255),
    phone VARCHAR(20),
    min_budget DECIMAL(15,2),
    max_budget DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        (u.first_name || ' ' || u.last_name) as full_name,
        u.email,
        u.phone,
        u.min_budget,
        u.max_budget
    FROM users u
    WHERE u.user_type = 'buyer'
      AND u.active = TRUE
      AND (u.min_budget IS NULL OR property_price >= u.min_budget)
      AND (u.max_budget IS NULL OR property_price <= u.max_budget)
    ORDER BY u.max_budget DESC NULLS LAST, u.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get agents by real estate company
CREATE OR REPLACE FUNCTION get_agents_by_company(company_id UUID)
RETURNS TABLE(
    id UUID,
    full_name TEXT,
    email VARCHAR(255),
    phone VARCHAR(20),
    bio TEXT,
    avatar_url TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        (u.first_name || ' ' || u.last_name) as full_name,
        u.email,
        u.phone,
        u.bio,
        u.avatar_url
    FROM users u
    WHERE u.user_type = 'agent'
      AND u.active = TRUE
      AND u.real_estate_company_id = company_id
    ORDER BY u.first_name, u.last_name;
END;
$$ LANGUAGE plpgsql;

-- Function to get user statistics
CREATE OR REPLACE FUNCTION get_user_statistics()
RETURNS TABLE(
    total_users BIGINT,
    active_users BIGINT,
    buyers BIGINT,
    sellers BIGINT,
    agents BIGINT,
    admins BIGINT,
    users_with_budget BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_users,
        COUNT(*) FILTER (WHERE active = TRUE) as active_users,
        COUNT(*) FILTER (WHERE user_type = 'buyer' AND active = TRUE) as buyers,
        COUNT(*) FILTER (WHERE user_type = 'seller' AND active = TRUE) as sellers,
        COUNT(*) FILTER (WHERE user_type = 'agent' AND active = TRUE) as agents,
        COUNT(*) FILTER (WHERE user_type = 'admin' AND active = TRUE) as admins,
        COUNT(*) FILTER (WHERE user_type = 'buyer' AND active = TRUE AND (min_budget IS NOT NULL OR max_budget IS NOT NULL)) as users_with_budget
    FROM users;
END;
$$ LANGUAGE plpgsql;

-- Comments on functions
COMMENT ON FUNCTION validate_ecuador_national_id IS 'Validates Ecuador national ID (cedula) using official algorithm';
COMMENT ON FUNCTION get_users_by_type IS 'Returns active users filtered by type';
COMMENT ON FUNCTION search_users_by_name IS 'Searches users by name using full-text search';
COMMENT ON FUNCTION get_buyers_for_property IS 'Returns buyers that can afford a specific property price';
COMMENT ON FUNCTION get_agents_by_company IS 'Returns agents associated with a specific real estate company';
COMMENT ON FUNCTION get_user_statistics IS 'Returns general user statistics by type';

-- Insert sample users for testing
INSERT INTO users (
    first_name, 
    last_name, 
    email, 
    phone, 
    national_id, 
    user_type,
    min_budget,
    max_budget,
    preferred_provinces,
    preferred_property_types
) VALUES 
(
    'Juan',
    'Pérez',
    'juan.perez@example.com',
    '0987654321',
    '1712345678',  -- Valid sample cedula
    'buyer',
    50000.00,
    200000.00,
    '["Pichincha", "Guayas"]',
    '["house", "apartment"]'
),
(
    'María',
    'González',
    'maria.gonzalez@example.com',
    '0987654322',
    '1723456789',  -- Valid sample cedula
    'seller',
    NULL,
    NULL,
    '[]',
    '[]'
),
(
    'Carlos',
    'Rodríguez',
    'carlos.rodriguez@example.com',
    '0987654323',
    '1734567890',  -- Valid sample cedula
    'agent',
    NULL,
    NULL,
    '[]',
    '[]'
),
(
    'Ana',
    'López',
    'ana.lopez@example.com',
    '0987654324',
    '1745678901',  -- Valid sample cedula
    'admin',
    NULL,
    NULL,
    '[]',
    '[]'
) ON CONFLICT (email) DO NOTHING;