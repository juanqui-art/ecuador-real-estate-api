-- Migración 005: Crear tabla de imágenes
-- Fecha: 2025-01-09
-- Propósito: Crear tabla para gestión de imágenes de propiedades con metadata completa

-- Crear tabla de imágenes
CREATE TABLE IF NOT EXISTS images (
    id VARCHAR(36) PRIMARY KEY,
    property_id VARCHAR(36) NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    original_url TEXT NOT NULL,
    alt_text TEXT DEFAULT '',
    sort_order INTEGER DEFAULT 0,
    size BIGINT DEFAULT 0,
    width INTEGER DEFAULT 0,
    height INTEGER DEFAULT 0,
    format VARCHAR(10) DEFAULT '',
    quality INTEGER DEFAULT 85,
    is_optimized BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraint para property_id
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);

-- Crear índices para optimizar consultas
CREATE INDEX IF NOT EXISTS idx_images_property_id ON images(property_id);
CREATE INDEX IF NOT EXISTS idx_images_sort_order ON images(property_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_images_format ON images(format);
CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at);
CREATE INDEX IF NOT EXISTS idx_images_size ON images(size);
CREATE INDEX IF NOT EXISTS idx_images_optimized ON images(is_optimized);

-- Agregar constraints para validación
ALTER TABLE images 
ADD CONSTRAINT chk_images_size_positive 
CHECK (size >= 0);

ALTER TABLE images 
ADD CONSTRAINT chk_images_dimensions_positive 
CHECK (width >= 0 AND height >= 0);

ALTER TABLE images 
ADD CONSTRAINT chk_images_quality_range 
CHECK (quality >= 1 AND quality <= 100);

ALTER TABLE images 
ADD CONSTRAINT chk_images_sort_order_positive 
CHECK (sort_order >= 0);

ALTER TABLE images 
ADD CONSTRAINT chk_images_format_valid 
CHECK (format IN ('jpg', 'jpeg', 'png', 'webp', 'avif', ''));

-- Agregar comentarios para documentación
COMMENT ON TABLE images IS 'Tabla de imágenes de propiedades con metadata completa';
COMMENT ON COLUMN images.id IS 'Identificador único de la imagen (UUID)';
COMMENT ON COLUMN images.property_id IS 'ID de la propiedad asociada';
COMMENT ON COLUMN images.file_name IS 'Nombre del archivo de imagen';
COMMENT ON COLUMN images.original_url IS 'URL de la imagen original optimizada';
COMMENT ON COLUMN images.alt_text IS 'Texto alternativo para accesibilidad';
COMMENT ON COLUMN images.sort_order IS 'Orden de clasificación (0 = imagen principal)';
COMMENT ON COLUMN images.size IS 'Tamaño del archivo en bytes';
COMMENT ON COLUMN images.width IS 'Ancho de la imagen en píxeles';
COMMENT ON COLUMN images.height IS 'Alto de la imagen en píxeles';
COMMENT ON COLUMN images.format IS 'Formato de la imagen (jpg, png, webp, avif)';
COMMENT ON COLUMN images.quality IS 'Calidad de compresión (1-100)';
COMMENT ON COLUMN images.is_optimized IS 'Indica si la imagen ha sido optimizada';
COMMENT ON COLUMN images.created_at IS 'Fecha de creación del registro';
COMMENT ON COLUMN images.updated_at IS 'Fecha de última actualización';

-- Función para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_images_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para actualizar updated_at automáticamente
CREATE TRIGGER trigger_update_images_updated_at
    BEFORE UPDATE ON images
    FOR EACH ROW
    EXECUTE FUNCTION update_images_updated_at();

-- Función para validar límite de imágenes por propiedad
CREATE OR REPLACE FUNCTION validate_images_per_property()
RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM images WHERE property_id = NEW.property_id) >= 50 THEN
        RAISE EXCEPTION 'Maximum number of images per property exceeded (50)';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para validar límite de imágenes por propiedad
CREATE TRIGGER trigger_validate_images_per_property
    BEFORE INSERT ON images
    FOR EACH ROW
    EXECUTE FUNCTION validate_images_per_property();

-- Vista para estadísticas de imágenes
CREATE OR REPLACE VIEW images_statistics AS
SELECT 
    COUNT(*) as total_images,
    COUNT(DISTINCT property_id) as properties_with_images,
    ROUND(AVG(size)) as average_size,
    SUM(size) as total_size,
    COUNT(CASE WHEN is_optimized = true THEN 1 END) as optimized_images,
    ROUND(COUNT(CASE WHEN is_optimized = true THEN 1 END) * 100.0 / COUNT(*), 2) as optimization_rate,
    format as image_format,
    COUNT(*) as format_count
FROM images
GROUP BY format
ORDER BY format_count DESC;

-- Comentario sobre la vista
COMMENT ON VIEW images_statistics IS 'Vista con estadísticas de imágenes agrupadas por formato';

-- Ejemplo de datos para testing (comentado)
/*
INSERT INTO images (id, property_id, file_name, original_url, alt_text, sort_order, size, width, height, format, quality, is_optimized)
VALUES (
    gen_random_uuid()::text,
    (SELECT id FROM properties LIMIT 1),
    'casa_moderna_sala.jpg',
    '/storage/images/casa_moderna_sala.jpg',
    'Sala moderna con vista al jardín',
    0,
    1048576,
    1920,
    1080,
    'jpg',
    85,
    true
);
*/