-- Migration 017: Add Full-Text Search support to properties table
-- This migration adds PostgreSQL Full-Text Search capabilities for efficient property search

-- Add tsvector column for full-text search
ALTER TABLE properties ADD COLUMN search_vector tsvector;

-- Create a function to update the search vector
CREATE OR REPLACE FUNCTION update_property_search_vector() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := 
        setweight(to_tsvector('spanish', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('spanish', coalesce(NEW.description, '')), 'B') ||
        setweight(to_tsvector('spanish', coalesce(NEW.province, '')), 'C') ||
        setweight(to_tsvector('spanish', coalesce(NEW.city, '')), 'C') ||
        setweight(to_tsvector('spanish', coalesce(NEW.sector, '')), 'D') ||
        setweight(to_tsvector('spanish', coalesce(NEW.type, '')), 'C') ||
        setweight(to_tsvector('spanish', coalesce(NEW.property_status, '')), 'D');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update search vector on insert/update
CREATE TRIGGER update_property_search_vector_trigger
    BEFORE INSERT OR UPDATE ON properties
    FOR EACH ROW EXECUTE FUNCTION update_property_search_vector();

-- Update existing records with search vectors
UPDATE properties SET search_vector = 
    setweight(to_tsvector('spanish', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('spanish', coalesce(description, '')), 'B') ||
    setweight(to_tsvector('spanish', coalesce(province, '')), 'C') ||
    setweight(to_tsvector('spanish', coalesce(city, '')), 'C') ||
    setweight(to_tsvector('spanish', coalesce(sector, '')), 'D') ||
    setweight(to_tsvector('spanish', coalesce(type, '')), 'C') ||
    setweight(to_tsvector('spanish', coalesce(property_status, '')), 'D');

-- Create GIN index for fast full-text search
CREATE INDEX idx_properties_search_vector ON properties USING gin(search_vector);

-- Create additional indexes for common search patterns
CREATE INDEX idx_properties_province_text ON properties USING gin(to_tsvector('spanish', province));
CREATE INDEX idx_properties_city_text ON properties USING gin(to_tsvector('spanish', city));
CREATE INDEX idx_properties_type_text ON properties USING gin(to_tsvector('spanish', type));

-- Create function for ranked search results
CREATE OR REPLACE FUNCTION search_properties_ranked(search_query text, result_limit int DEFAULT 50)
RETURNS TABLE (
    id uuid,
    slug text,
    title text,
    description text,
    price numeric,
    province text,
    city text,
    type text,
    rank real
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.slug,
        p.title,
        p.description,
        p.price,
        p.province,
        p.city,
        p.type,
        ts_rank_cd(p.search_vector, plainto_tsquery('spanish', search_query)) as rank
    FROM properties p
    WHERE p.search_vector @@ plainto_tsquery('spanish', search_query)
    ORDER BY 
        ts_rank_cd(p.search_vector, plainto_tsquery('spanish', search_query)) DESC,
        p.featured DESC,
        p.created_at DESC
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Create function for search suggestions/autocomplete
CREATE OR REPLACE FUNCTION get_search_suggestions(search_query text, suggestion_limit int DEFAULT 10)
RETURNS TABLE (
    suggestion text,
    category text,
    frequency int
) AS $$
BEGIN
    RETURN QUERY
    WITH suggestions AS (
        SELECT DISTINCT province as suggestion, 'province' as category, 
               COUNT(*) OVER (PARTITION BY province) as frequency
        FROM properties 
        WHERE province ILIKE '%' || search_query || '%'
        
        UNION ALL
        
        SELECT DISTINCT city as suggestion, 'city' as category,
               COUNT(*) OVER (PARTITION BY city) as frequency
        FROM properties 
        WHERE city ILIKE '%' || search_query || '%'
        
        UNION ALL
        
        SELECT DISTINCT type as suggestion, 'type' as category,
               COUNT(*) OVER (PARTITION BY type) as frequency
        FROM properties 
        WHERE type ILIKE '%' || search_query || '%'
    )
    SELECT s.suggestion, s.category, s.frequency
    FROM suggestions s
    ORDER BY s.frequency DESC, s.suggestion ASC
    LIMIT suggestion_limit;
END;
$$ LANGUAGE plpgsql;

-- Create function for advanced search with filters
CREATE OR REPLACE FUNCTION advanced_search_properties(
    search_query text DEFAULT '',
    filter_province text DEFAULT '',
    filter_city text DEFAULT '',
    filter_type text DEFAULT '',
    min_price numeric DEFAULT 0,
    max_price numeric DEFAULT 999999999,
    min_bedrooms int DEFAULT 0,
    max_bedrooms int DEFAULT 100,
    min_bathrooms numeric DEFAULT 0,
    max_bathrooms numeric DEFAULT 100,
    min_area numeric DEFAULT 0,
    max_area numeric DEFAULT 999999,
    featured_only boolean DEFAULT false,
    result_limit int DEFAULT 50
)
RETURNS TABLE (
    id uuid,
    slug text,
    title text,
    description text,
    price numeric,
    province text,
    city text,
    type text,
    bedrooms int,
    bathrooms numeric,
    area_m2 numeric,
    featured boolean,
    rank real
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id, p.slug, p.title, p.description, p.price,
        p.province, p.city, p.type, p.bedrooms, p.bathrooms, p.area_m2, p.featured,
        CASE 
            WHEN search_query = '' THEN 1.0
            ELSE ts_rank_cd(p.search_vector, plainto_tsquery('spanish', search_query))
        END as rank
    FROM properties p
    WHERE 
        (search_query = '' OR p.search_vector @@ plainto_tsquery('spanish', search_query))
        AND (filter_province = '' OR p.province = filter_province)
        AND (filter_city = '' OR p.city = filter_city)
        AND (filter_type = '' OR p.type = filter_type)
        AND p.price >= min_price
        AND p.price <= max_price
        AND p.bedrooms >= min_bedrooms
        AND p.bedrooms <= max_bedrooms
        AND p.bathrooms >= min_bathrooms
        AND p.bathrooms <= max_bathrooms
        AND p.area_m2 >= min_area
        AND p.area_m2 <= max_area
        AND (featured_only = false OR p.featured = true)
        AND p.status = 'available'
    ORDER BY 
        rank DESC,
        p.featured DESC,
        p.created_at DESC
    LIMIT result_limit;
END;
$$ LANGUAGE plpgsql;

-- Create view for search analytics
CREATE OR REPLACE VIEW search_analytics AS
SELECT 
    date_trunc('day', created_at) as search_date,
    COUNT(*) as total_properties,
    COUNT(CASE WHEN featured = true THEN 1 END) as featured_properties,
    AVG(price) as avg_price,
    province,
    city,
    type
FROM properties
WHERE status = 'available'
GROUP BY date_trunc('day', created_at), province, city, type
ORDER BY search_date DESC;

-- Grant permissions for search functions
GRANT EXECUTE ON FUNCTION search_properties_ranked(text, int) TO PUBLIC;
GRANT EXECUTE ON FUNCTION get_search_suggestions(text, int) TO PUBLIC;
GRANT EXECUTE ON FUNCTION advanced_search_properties(text, text, text, text, numeric, numeric, int, int, numeric, numeric, numeric, numeric, boolean, int) TO PUBLIC;
GRANT SELECT ON search_analytics TO PUBLIC;