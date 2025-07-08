-- Migración 004: Agregar campos de imágenes y media a propiedades
-- Fecha: 2025-01-15
-- Propósito: Habilitar galería de imágenes y tours virtuales en el frontend

-- Agregar campos para imágenes y media
ALTER TABLE propiedades 
ADD COLUMN imagen_principal TEXT DEFAULT '' NOT NULL,
ADD COLUMN imagenes JSONB DEFAULT '[]' NOT NULL,
ADD COLUMN video_tour TEXT DEFAULT '' NOT NULL,
ADD COLUMN tour_360 TEXT DEFAULT '' NOT NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.imagen_principal IS 'URL de la imagen principal de la propiedad';
COMMENT ON COLUMN propiedades.imagenes IS 'Array JSON de URLs de imágenes adicionales de la galería';
COMMENT ON COLUMN propiedades.video_tour IS 'URL del video tour de la propiedad (YouTube, Vimeo, etc.)';
COMMENT ON COLUMN propiedades.tour_360 IS 'URL del tour virtual 360° de la propiedad';

-- Agregar constraints para validar URLs (formato básico)
ALTER TABLE propiedades 
ADD CONSTRAINT chk_imagen_principal_url 
CHECK (imagen_principal = '' OR imagen_principal ~* '^https?://.*\.(jpg|jpeg|png|webp)(\?.*)?$');

ALTER TABLE propiedades 
ADD CONSTRAINT chk_video_tour_url 
CHECK (video_tour = '' OR video_tour ~* '^https?://.*');

ALTER TABLE propiedades 
ADD CONSTRAINT chk_tour_360_url 
CHECK (tour_360 = '' OR tour_360 ~* '^https?://.*');

-- Crear índice para búsquedas de propiedades con imágenes
CREATE INDEX idx_propiedades_con_imagenes ON propiedades(imagen_principal) 
WHERE imagen_principal != '';

-- Crear índice para búsquedas de propiedades con video tour
CREATE INDEX idx_propiedades_con_video ON propiedades(video_tour) 
WHERE video_tour != '';

-- Crear índice para búsquedas de propiedades con tour 360
CREATE INDEX idx_propiedades_con_tour_360 ON propiedades(tour_360) 
WHERE tour_360 != '';

-- Crear índice GIN para búsquedas eficientes en el array JSON de imágenes
CREATE INDEX idx_propiedades_imagenes_gin ON propiedades USING GIN(imagenes);

-- Función helper para validar array de URLs de imágenes (opcional)
CREATE OR REPLACE FUNCTION validar_array_imagenes(imagenes_json JSONB) 
RETURNS BOOLEAN AS $$
BEGIN
    -- Verificar que sea un array
    IF jsonb_typeof(imagenes_json) != 'array' THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que no tenga más de 20 imágenes (límite razonable)
    IF jsonb_array_length(imagenes_json) > 20 THEN
        RETURN FALSE;
    END IF;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función
ALTER TABLE propiedades 
ADD CONSTRAINT chk_imagenes_validas 
CHECK (validar_array_imagenes(imagenes));

-- Actualizar comentario de la tabla
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con imágenes y media en migración 004';