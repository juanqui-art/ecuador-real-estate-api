-- Migration 013: Create English properties table
-- Date: 2025-01-15
-- Purpose: Complete English properties table with all features

-- Create properties table in English
CREATE TABLE IF NOT EXISTS properties (
    -- Primary identification
    id VARCHAR(36) PRIMARY KEY,
    slug VARCHAR(100) UNIQUE NOT NULL,

    -- Basic information
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(15,2) NOT NULL CHECK (price > 0),

    -- Location
    province VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    sector VARCHAR(100),
    address VARCHAR(255),

    -- Geolocation for maps
    latitude DECIMAL(10,7),  -- GPS latitude (-4.0 to 2.0 for Ecuador)
    longitude DECIMAL(11,7), -- GPS longitude (-92.0 to -75.0 for Ecuador)
    location_precision VARCHAR(20) DEFAULT 'approximate' CHECK (location_precision IN ('exact', 'approximate', 'sector')),

    -- Property characteristics
    type VARCHAR(50) NOT NULL CHECK (type IN ('house', 'apartment', 'land', 'commercial')),
    status VARCHAR(50) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'sold', 'rented', 'reserved')),
    bedrooms INTEGER DEFAULT 0 CHECK (bedrooms >= 0),
    bathrooms DECIMAL(3,1) DEFAULT 0 CHECK (bathrooms >= 0),
    area_m2 DECIMAL(10,2) DEFAULT 0 CHECK (area_m2 >= 0),

    -- Images and media
    main_image TEXT, -- URL of main image
    images JSONB DEFAULT '[]', -- Array of image URLs
    video_tour TEXT, -- Video tour URL
    tour_360 TEXT,   -- 360° virtual tour URL

    -- Additional pricing
    rent_price DECIMAL(15,2), -- Monthly rent price
    common_expenses DECIMAL(15,2), -- Common/HOA expenses
    price_per_m2 DECIMAL(15,2), -- Price per square meter

    -- Detailed characteristics
    year_built INTEGER,
    floors INTEGER,
    property_status VARCHAR(20) DEFAULT 'used' CHECK (property_status IN ('new', 'used', 'renovated')),
    furnished BOOLEAN DEFAULT FALSE,

    -- Amenities (for frontend filters)
    garage BOOLEAN DEFAULT FALSE,
    pool BOOLEAN DEFAULT FALSE,
    garden BOOLEAN DEFAULT FALSE,
    terrace BOOLEAN DEFAULT FALSE,
    balcony BOOLEAN DEFAULT FALSE,
    security BOOLEAN DEFAULT FALSE,
    elevator BOOLEAN DEFAULT FALSE,
    air_conditioning BOOLEAN DEFAULT FALSE,

    -- Marketing and SEO
    tags JSONB DEFAULT '[]', -- Search tags ["luxury", "ocean-view"]
    featured BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,

    -- Relations
    real_estate_company_id UUID REFERENCES real_estate_companies(id),

    -- Audit fields
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_properties_province ON properties(province);
CREATE INDEX IF NOT EXISTS idx_properties_city ON properties(city);
CREATE INDEX IF NOT EXISTS idx_properties_type ON properties(type);
CREATE INDEX IF NOT EXISTS idx_properties_status ON properties(status);
CREATE INDEX IF NOT EXISTS idx_properties_price ON properties(price);
CREATE INDEX IF NOT EXISTS idx_properties_area_m2 ON properties(area_m2);
CREATE INDEX IF NOT EXISTS idx_properties_bedrooms ON properties(bedrooms);
CREATE INDEX IF NOT EXISTS idx_properties_bathrooms ON properties(bathrooms);
CREATE INDEX IF NOT EXISTS idx_properties_featured ON properties(featured) WHERE featured = TRUE;
CREATE INDEX IF NOT EXISTS idx_properties_created_at ON properties(created_at);
CREATE INDEX IF NOT EXISTS idx_properties_company_id ON properties(real_estate_company_id);

-- Geospatial index for location queries
CREATE INDEX IF NOT EXISTS idx_properties_location ON properties(latitude, longitude) 
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Full-text search index for title and description
CREATE INDEX IF NOT EXISTS idx_properties_fulltext_search ON properties 
USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- JSONB indexes for tags and images
CREATE INDEX IF NOT EXISTS idx_properties_tags ON properties USING gin(tags);
CREATE INDEX IF NOT EXISTS idx_properties_images ON properties USING gin(images);

-- Add constraints for Ecuador coordinates
ALTER TABLE properties 
ADD CONSTRAINT chk_latitude_ecuador 
CHECK (latitude IS NULL OR (latitude >= -4.0 AND latitude <= 2.0));

ALTER TABLE properties 
ADD CONSTRAINT chk_longitude_ecuador 
CHECK (longitude IS NULL OR (longitude >= -92.0 AND longitude <= -75.0));

-- Add constraint for valid Ecuador provinces
ALTER TABLE properties 
ADD CONSTRAINT chk_valid_ecuador_province 
CHECK (province IN (
    'Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi',
    'El Oro', 'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja',
    'Los Ríos', 'Manabí', 'Morona Santiago', 'Napo', 'Orellana', 'Pastaza',
    'Pichincha', 'Santa Elena', 'Santo Domingo', 'Sucumbíos', 'Tungurahua', 'Zamora Chinchipe'
));

-- Comments for documentation
COMMENT ON TABLE properties IS 'Main table for real estate properties in English';
COMMENT ON COLUMN properties.id IS 'Unique UUID identifier';
COMMENT ON COLUMN properties.slug IS 'SEO-friendly URL slug';
COMMENT ON COLUMN properties.price IS 'Price in USD';
COMMENT ON COLUMN properties.area_m2 IS 'Area in square meters';
COMMENT ON COLUMN properties.bathrooms IS 'Number of bathrooms (can be decimal: 2.5)';
COMMENT ON COLUMN properties.latitude IS 'GPS latitude for Ecuador (-4.0 to 2.0)';
COMMENT ON COLUMN properties.longitude IS 'GPS longitude for Ecuador (-92.0 to -75.0)';
COMMENT ON COLUMN properties.images IS 'JSON array of image URLs';
COMMENT ON COLUMN properties.tags IS 'JSON array of search tags';

-- Trigger for auto-updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_properties_updated_at
    BEFORE UPDATE ON properties
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate price per m2
CREATE OR REPLACE FUNCTION calculate_price_per_m2()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.area_m2 > 0 THEN
        NEW.price_per_m2 = NEW.price / NEW.area_m2;
    ELSE
        NEW.price_per_m2 = NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_calculate_price_per_m2
    BEFORE INSERT OR UPDATE ON properties
    FOR EACH ROW
    EXECUTE FUNCTION calculate_price_per_m2();

-- Function to search properties with full-text search
CREATE OR REPLACE FUNCTION search_properties(search_term TEXT)
RETURNS TABLE(
    id VARCHAR(36),
    title VARCHAR(255),
    description TEXT,
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    rank REAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.title,
        p.description,
        p.price,
        p.province,
        p.city,
        p.type,
        ts_rank(to_tsvector('english', p.title || ' ' || COALESCE(p.description, '')), plainto_tsquery('english', search_term)) as rank
    FROM properties p
    WHERE p.status = 'available'
      AND to_tsvector('english', p.title || ' ' || COALESCE(p.description, '')) @@ plainto_tsquery('english', search_term)
    ORDER BY rank DESC, p.featured DESC, p.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Function to get properties within budget
CREATE OR REPLACE FUNCTION get_properties_in_budget(min_price DECIMAL, max_price DECIMAL)
RETURNS TABLE(
    id VARCHAR(36),
    title VARCHAR(255),
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    bedrooms INTEGER,
    bathrooms DECIMAL(3,1),
    area_m2 DECIMAL(10,2)
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
        p.bedrooms,
        p.bathrooms,
        p.area_m2
    FROM properties p
    WHERE p.status = 'available'
      AND (min_price IS NULL OR p.price >= min_price)
      AND (max_price IS NULL OR p.price <= max_price)
    ORDER BY p.featured DESC, p.price ASC;
END;
$$ LANGUAGE plpgsql;

-- Insert sample data for testing
INSERT INTO properties (
    id, slug, title, description, price, province, city, sector, 
    type, status, bedrooms, bathrooms, area_m2, latitude, longitude,
    main_image, featured
) VALUES 
(
    'prop-en-001',
    'beautiful-house-samborondon-prop-en-001',
    'Beautiful House in Samborondón with Pool',
    'Modern 3-story house with luxury finishes, pool and garden',
    285000.00,
    'Guayas',
    'Samborondón',
    'La Puntilla',
    'house',
    'available',
    4,
    3.5,
    320.00,
    -2.1333,
    -79.8833,
    'https://example.com/images/house1.jpg',
    TRUE
),
(
    'prop-en-002',
    'central-apartment-quito-prop-en-002',
    'Central Apartment in Quito',
    'Modern apartment in the north center of Quito',
    125000.00,
    'Pichincha',
    'Quito',
    'La Carolina',
    'apartment',
    'available',
    2,
    2.0,
    85.00,
    -0.1807,
    -78.4678,
    'https://example.com/images/apartment1.jpg',
    FALSE
),
(
    'prop-en-003',
    'land-cuenca-ideal-construction-prop-en-003',
    'Land in Cuenca - Ideal for Construction',
    'Flat land ideal for construction in a growing area',
    45000.00,
    'Azuay',
    'Cuenca',
    'El Batán',
    'land',
    'available',
    0,
    0.0,
    500.00,
    -2.9001,
    -79.0059,
    'https://example.com/images/land1.jpg',
    FALSE
) ON CONFLICT (id) DO NOTHING;

-- Comments on functions
COMMENT ON FUNCTION update_updated_at_column IS 'Automatically updates updated_at timestamp';
COMMENT ON FUNCTION calculate_price_per_m2 IS 'Automatically calculates price per square meter';
COMMENT ON FUNCTION search_properties IS 'Full-text search for properties by title and description';
COMMENT ON FUNCTION get_properties_in_budget IS 'Returns properties within a budget range';