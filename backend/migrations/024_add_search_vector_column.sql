-- Migration: Add search_vector column for PostgreSQL Full-Text Search
-- Date: 2025-07-11
-- Description: Adds search_vector column to properties table for optimized FTS

-- Add search_vector column to properties table
ALTER TABLE properties 
ADD COLUMN search_vector tsvector;

-- Create index for full-text search performance
CREATE INDEX idx_properties_search_vector ON properties USING gin(search_vector);

-- Create function to update search_vector
CREATE OR REPLACE FUNCTION update_property_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('spanish', 
        COALESCE(NEW.title, '') || ' ' || 
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(NEW.province, '') || ' ' ||
        COALESCE(NEW.city, '') || ' ' ||
        COALESCE(NEW.property_type, '') || ' ' ||
        COALESCE(NEW.status, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update search_vector
DROP TRIGGER IF EXISTS trigger_update_property_search_vector ON properties;
CREATE TRIGGER trigger_update_property_search_vector
    BEFORE INSERT OR UPDATE ON properties
    FOR EACH ROW
    EXECUTE FUNCTION update_property_search_vector();

-- Populate search_vector for existing records
UPDATE properties 
SET search_vector = to_tsvector('spanish', 
    COALESCE(title, '') || ' ' || 
    COALESCE(description, '') || ' ' ||
    COALESCE(province, '') || ' ' ||
    COALESCE(city, '') || ' ' ||
    COALESCE(property_type, '') || ' ' ||
    COALESCE(status, '')
);

-- Add comment for documentation
COMMENT ON COLUMN properties.search_vector IS 'Full-text search vector for efficient search operations';
COMMENT ON INDEX idx_properties_search_vector IS 'GIN index for full-text search performance on search_vector column';
COMMENT ON FUNCTION update_property_search_vector() IS 'Automatically updates search_vector when property data changes';
COMMENT ON TRIGGER trigger_update_property_search_vector ON properties IS 'Trigger to maintain search_vector column on data changes';