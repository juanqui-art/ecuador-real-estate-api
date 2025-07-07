-- Migración: Agregar columna slug para SEO
-- Fecha: 2024-07-04
-- Descripción: Añadir campo slug para URLs amigables

-- Agregar columna slug
ALTER TABLE propiedades 
ADD COLUMN IF NOT EXISTS slug VARCHAR(255);

-- Crear índice único para slugs (importante para SEO)
CREATE UNIQUE INDEX IF NOT EXISTS idx_propiedades_slug ON propiedades(slug);

-- Generar slugs para propiedades existentes
-- Esta función temporal genera slugs básicos para datos existentes
DO $$
DECLARE
    prop RECORD;
    nuevo_slug VARCHAR(255);
    contador INTEGER := 1;
BEGIN
    FOR prop IN SELECT id, titulo FROM propiedades WHERE slug IS NULL OR slug = ''
    LOOP
        -- Generar slug básico: título limpio + id corto
        nuevo_slug := LOWER(
            REGEXP_REPLACE(
                REGEXP_REPLACE(prop.titulo, '[^a-zA-ZáéíóúñÁÉÍÓÚÑ0-9\s]+', '', 'g'),
                '\s+', '-', 'g'
            )
        );
        
        -- Truncar si es muy largo
        IF LENGTH(nuevo_slug) > 50 THEN
            nuevo_slug := SUBSTRING(nuevo_slug, 1, 50);
        END IF;
        
        -- Remover guiones al final
        nuevo_slug := RTRIM(nuevo_slug, '-');
        
        -- Agregar ID corto (primeros 8 caracteres)
        nuevo_slug := nuevo_slug || '-' || SUBSTRING(prop.id, 1, 8);
        
        -- Actualizar la propiedad
        UPDATE propiedades 
        SET slug = nuevo_slug 
        WHERE id = prop.id;
        
        RAISE NOTICE 'Slug generado para propiedad %: %', prop.id, nuevo_slug;
    END LOOP;
END $$;

-- Hacer obligatorio el campo slug para nuevas propiedades
ALTER TABLE propiedades 
ALTER COLUMN slug SET NOT NULL;

-- Comentario para documentar
COMMENT ON COLUMN propiedades.slug IS 'URL amigable para SEO, generada automáticamente desde el título';

-- Verificar que todas las propiedades tienen slug
DO $$
DECLARE
    count_sin_slug INTEGER;
BEGIN
    SELECT COUNT(*) INTO count_sin_slug 
    FROM propiedades 
    WHERE slug IS NULL OR slug = '';
    
    IF count_sin_slug > 0 THEN
        RAISE EXCEPTION 'Error: % propiedades sin slug después de la migración', count_sin_slug;
    ELSE
        RAISE NOTICE 'Migración completada exitosamente. Todas las propiedades tienen slug.';
    END IF;
END $$;