-- Migration 014: Create English real estate companies table
-- Date: 2025-01-15
-- Purpose: Real estate companies table in English with Ecuador validations

-- Create real estate companies table in English
CREATE TABLE IF NOT EXISTS real_estate_companies (
    -- Primary identification
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Basic company information
    name VARCHAR(255) NOT NULL,
    ruc VARCHAR(13) NOT NULL UNIQUE,        -- Ecuador tax ID (13 digits)
    address TEXT NOT NULL,                  -- Physical address
    description TEXT DEFAULT '' NOT NULL,   -- Company description

    -- Contact information
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    website TEXT DEFAULT '' NOT NULL,       -- Company website
    logo_url TEXT DEFAULT '' NOT NULL,      -- Company logo

    -- Status
    active BOOLEAN DEFAULT TRUE NOT NULL,   -- Is company active

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add comments for documentation
COMMENT ON TABLE real_estate_companies IS 'Real estate companies registered in the system';
COMMENT ON COLUMN real_estate_companies.id IS 'Unique identifier for the real estate company';
COMMENT ON COLUMN real_estate_companies.name IS 'Commercial name of the real estate company';
COMMENT ON COLUMN real_estate_companies.ruc IS 'RUC (Registro Único de Contribuyentes) of the company';
COMMENT ON COLUMN real_estate_companies.address IS 'Physical address of the main office';
COMMENT ON COLUMN real_estate_companies.phone IS 'Main phone number';
COMMENT ON COLUMN real_estate_companies.email IS 'Main contact email';
COMMENT ON COLUMN real_estate_companies.website IS 'Company website URL';
COMMENT ON COLUMN real_estate_companies.logo_url IS 'Company logo URL';
COMMENT ON COLUMN real_estate_companies.active IS 'Indicates if the company is active in the system';
COMMENT ON COLUMN real_estate_companies.created_at IS 'Registration date in the system';
COMMENT ON COLUMN real_estate_companies.updated_at IS 'Last update date';

-- Add constraints to validate Ecuador RUC format
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_ruc_format 
CHECK (ruc ~ '^[0-9]{13}$');

-- Constraint to validate that RUC ends in 001 (companies)
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_ruc_company 
CHECK (ruc ~ '[0-9]{10}001$');

-- Constraint to validate email format
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_email_format 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Constraint to validate Ecuador phone format
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_phone_format 
CHECK (phone ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');

-- Constraint to validate website URL format
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_website_url 
CHECK (website = '' OR website ~* '^https?://.*');

-- Constraint to validate logo URL format
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_logo_url 
CHECK (logo_url = '' OR logo_url ~* '^https?://.*\.(jpg|jpeg|png|webp|svg)(\?.*)?$');

-- Create indexes for searches
CREATE INDEX idx_real_estate_companies_name ON real_estate_companies(name);
CREATE INDEX idx_real_estate_companies_active ON real_estate_companies(active) WHERE active = TRUE;
CREATE INDEX idx_real_estate_companies_created_at ON real_estate_companies(created_at);

-- Full-text search index for name and description
CREATE INDEX idx_real_estate_companies_fulltext_search ON real_estate_companies 
USING gin(to_tsvector('english', name || ' ' || description));

-- Function to validate Ecuador RUC (mod 11 algorithm)
CREATE OR REPLACE FUNCTION validate_ecuador_ruc(ruc_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    ruc_clean TEXT;
    digits INTEGER[];
    sum INTEGER := 0;
    check_digit INTEGER;
    i INTEGER;
BEGIN
    -- Clean RUC (numbers only)
    ruc_clean := regexp_replace(ruc_input, '[^0-9]', '', 'g');
    
    -- Check length
    IF LENGTH(ruc_clean) != 13 THEN
        RETURN FALSE;
    END IF;
    
    -- Check that it ends in 001 (companies)
    IF RIGHT(ruc_clean, 3) != '001' THEN
        RETURN FALSE;
    END IF;
    
    -- Convert to digit array
    FOR i IN 1..10 LOOP
        digits[i] := CAST(SUBSTRING(ruc_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- Validate that third digit is less than 6 (legal entities)
    IF digits[3] >= 6 THEN
        RETURN FALSE;
    END IF;
    
    -- Apply mod 11 algorithm
    FOR i IN 1..9 LOOP
        sum := sum + (digits[i] * (10 - i));
    END LOOP;
    
    check_digit := 11 - (sum % 11);
    
    IF check_digit = 11 THEN
        check_digit := 0;
    ELSIF check_digit = 10 THEN
        check_digit := 1;
    END IF;
    
    RETURN check_digit = digits[10];
END;
$$ LANGUAGE plpgsql;

-- Add constraint using RUC validation function
ALTER TABLE real_estate_companies 
ADD CONSTRAINT chk_ruc_valid 
CHECK (validate_ecuador_ruc(ruc));

-- Trigger to automatically update updated_at
CREATE TRIGGER trigger_real_estate_companies_updated_at
    BEFORE UPDATE ON real_estate_companies
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to get active real estate companies
CREATE OR REPLACE FUNCTION get_active_real_estate_companies()
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    ruc VARCHAR(13),
    phone VARCHAR(20),
    email VARCHAR(255),
    website TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        rec.id,
        rec.name,
        rec.ruc,
        rec.phone,
        rec.email,
        rec.website
    FROM real_estate_companies rec
    WHERE rec.active = TRUE
    ORDER BY rec.name;
END;
$$ LANGUAGE plpgsql;

-- Function to search real estate companies by name
CREATE OR REPLACE FUNCTION search_real_estate_companies_by_name(search_name TEXT)
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    description TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        rec.id,
        rec.name,
        rec.description,
        rec.phone,
        rec.email,
        ts_rank(to_tsvector('english', rec.name || ' ' || rec.description), plainto_tsquery('english', search_name)) as rank
    FROM real_estate_companies rec
    WHERE rec.active = TRUE
      AND to_tsvector('english', rec.name || ' ' || rec.description) @@ plainto_tsquery('english', search_name)
    ORDER BY rank DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get real estate company statistics
CREATE OR REPLACE FUNCTION get_real_estate_company_statistics()
RETURNS TABLE(
    total_companies BIGINT,
    active_companies BIGINT,
    inactive_companies BIGINT,
    active_percentage DECIMAL(5,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_companies,
        COUNT(*) FILTER (WHERE active = TRUE) as active_companies,
        COUNT(*) FILTER (WHERE active = FALSE) as inactive_companies,
        ROUND(
            (COUNT(*) FILTER (WHERE active = TRUE) * 100.0 / NULLIF(COUNT(*), 0)), 
            2
        ) as active_percentage
    FROM real_estate_companies;
END;
$$ LANGUAGE plpgsql;

-- Comments on functions
COMMENT ON FUNCTION validate_ecuador_ruc IS 'Validates Ecuador RUC using mod 11 algorithm';
COMMENT ON FUNCTION get_active_real_estate_companies IS 'Returns list of active real estate companies';
COMMENT ON FUNCTION search_real_estate_companies_by_name IS 'Searches real estate companies by name using full-text search';
COMMENT ON FUNCTION get_real_estate_company_statistics IS 'Returns general statistics of real estate companies';

-- Insert sample real estate company for testing
INSERT INTO real_estate_companies (
    name, 
    ruc, 
    address, 
    phone, 
    email, 
    website, 
    description
) VALUES (
    'Example Real Estate S.A.',
    '1792146739001',  -- Valid sample RUC
    'Av. 9 de Octubre 123, Guayaquil, Ecuador',
    '042234567',
    'info@examplerealestate.com',
    'https://www.examplerealestate.com',
    'Leading real estate company in Ecuador with over 20 years of experience.'
) ON CONFLICT (ruc) DO NOTHING;