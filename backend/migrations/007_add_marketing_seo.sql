-- Migración 008: Agregar campos de marketing y SEO a propiedades
-- Fecha: 2025-01-15
-- Propósito: Habilitar funcionalidades de marketing, SEO y estadísticas

-- Agregar campos para marketing y SEO
ALTER TABLE propiedades 
ADD COLUMN tags JSONB DEFAULT '[]' NOT NULL,
ADD COLUMN destacada BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN visitas_contador INTEGER DEFAULT 0 NOT NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.tags IS 'Array JSON de tags para búsqueda y SEO ["lujo", "vista-al-mar", "cerca-metro"]';
COMMENT ON COLUMN propiedades.destacada IS 'Indica si la propiedad es destacada/premium en listings';
COMMENT ON COLUMN propiedades.visitas_contador IS 'Contador de vistas de la propiedad para estadísticas';

-- Agregar constraints para validar datos
ALTER TABLE propiedades 
ADD CONSTRAINT chk_visitas_positivas 
CHECK (visitas_contador >= 0);

-- Constraint para limitar número máximo de tags (evitar spam)
ALTER TABLE propiedades 
ADD CONSTRAINT chk_tags_limite 
CHECK (jsonb_array_length(tags) <= 20);

-- Función para validar que los tags sean strings válidos
CREATE OR REPLACE FUNCTION validar_tags(tags_json JSONB) 
RETURNS BOOLEAN AS $$
DECLARE
    tag_item JSONB;
    tag_text TEXT;
BEGIN
    -- Verificar que sea un array
    IF jsonb_typeof(tags_json) != 'array' THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar cada elemento del array
    FOR tag_item IN SELECT jsonb_array_elements(tags_json)
    LOOP
        -- Cada tag debe ser un string
        IF jsonb_typeof(tag_item) != 'string' THEN
            RETURN FALSE;
        END IF;
        
        -- Extraer el texto del tag
        tag_text := tag_item #>> '{}';
        
        -- El tag no puede estar vacío ni ser muy largo
        IF LENGTH(TRIM(tag_text)) = 0 OR LENGTH(tag_text) > 50 THEN
            RETURN FALSE;
        END IF;
        
        -- El tag solo puede contener letras, números, espacios y guiones
        IF tag_text !~ '^[a-záéíóúñA-ZÁÉÍÓÚÑ0-9\s\-]+$' THEN
            RETURN FALSE;
        END IF;
    END LOOP;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función de validación
ALTER TABLE propiedades 
ADD CONSTRAINT chk_tags_validos 
CHECK (validar_tags(tags));

-- Crear índices para búsquedas y filtros
CREATE INDEX idx_propiedades_destacada ON propiedades(destacada) WHERE destacada = TRUE;
CREATE INDEX idx_propiedades_visitas ON propiedades(visitas_contador DESC);

-- Índice GIN para búsquedas eficientes en tags
CREATE INDEX idx_propiedades_tags_gin ON propiedades USING GIN(tags);

-- Índice para búsquedas combinadas de propiedades destacadas con visitas
CREATE INDEX idx_propiedades_destacada_visitas ON propiedades(destacada, visitas_contador DESC) 
WHERE destacada = TRUE;

-- Función para buscar propiedades por tag
CREATE OR REPLACE FUNCTION buscar_por_tag(tag_busqueda TEXT)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    tags_encontrados JSONB,
    visitas INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.tags,
        p.visitas_contador
    FROM propiedades p
    WHERE p.tags ? tag_busqueda  -- Operador ? busca si existe el elemento en el array JSON
    ORDER BY p.destacada DESC, p.visitas_contador DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar propiedades que contengan cualquiera de varios tags
CREATE OR REPLACE FUNCTION buscar_por_tags_multiples(tags_busqueda TEXT[])
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    tags_coincidentes INTEGER,
    total_visitas INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        (
            SELECT COUNT(*)::INTEGER
            FROM jsonb_array_elements_text(p.tags) AS tag
            WHERE tag = ANY(tags_busqueda)
        ) AS tags_coincidentes,
        p.visitas_contador
    FROM propiedades p
    WHERE p.tags ?| tags_busqueda  -- Operador ?| busca si existe alguno de los elementos
    ORDER BY tags_coincidentes DESC, p.destacada DESC, p.visitas_contador DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener propiedades más populares
CREATE OR REPLACE FUNCTION obtener_propiedades_populares(limite INTEGER DEFAULT 10)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    provincia VARCHAR(50),
    ciudad VARCHAR(100),
    visitas INTEGER,
    destacada BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.provincia,
        p.ciudad,
        p.visitas_contador,
        p.destacada
    FROM propiedades p
    WHERE p.estado = 'disponible'
    ORDER BY p.visitas_contador DESC, p.destacada DESC
    LIMIT limite;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de tags más usados
CREATE OR REPLACE FUNCTION obtener_tags_populares(limite INTEGER DEFAULT 20)
RETURNS TABLE(
    tag TEXT,
    frecuencia BIGINT,
    propiedades_destacadas BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        tag_individual,
        COUNT(*) AS frecuencia,
        COUNT(*) FILTER (WHERE p.destacada = TRUE) AS propiedades_destacadas
    FROM propiedades p,
         jsonb_array_elements_text(p.tags) AS tag_individual
    GROUP BY tag_individual
    ORDER BY frecuencia DESC, propiedades_destacadas DESC
    LIMIT limite;
END;
$$ LANGUAGE plpgsql;

-- Función para incrementar contador de visitas de forma segura
CREATE OR REPLACE FUNCTION incrementar_visitas(propiedad_id UUID)
RETURNS BOOLEAN AS $$
BEGIN
    UPDATE propiedades 
    SET visitas_contador = visitas_contador + 1
    WHERE id = propiedad_id;
    
    RETURN FOUND;
END;
$$ LANGUAGE plpgsql;

-- Vista para propiedades destacadas con métricas
CREATE OR REPLACE VIEW propiedades_destacadas_metricas AS
SELECT 
    p.*,
    jsonb_array_length(p.tags) AS total_tags,
    CASE 
        WHEN p.visitas_contador > 100 THEN 'Alta popularidad'
        WHEN p.visitas_contador > 50 THEN 'Popularidad media'
        WHEN p.visitas_contador > 10 THEN 'Popularidad baja'
        ELSE 'Nueva'
    END AS categoria_popularidad
FROM propiedades p
WHERE p.destacada = TRUE
ORDER BY p.visitas_contador DESC;

-- Comentarios en funciones y vistas
COMMENT ON FUNCTION buscar_por_tag IS 'Busca propiedades que contengan un tag específico';
COMMENT ON FUNCTION buscar_por_tags_multiples IS 'Busca propiedades que contengan cualquiera de los tags especificados';
COMMENT ON FUNCTION obtener_propiedades_populares IS 'Retorna las propiedades más visitadas';
COMMENT ON FUNCTION obtener_tags_populares IS 'Retorna estadísticas de los tags más utilizados';
COMMENT ON FUNCTION incrementar_visitas IS 'Incrementa de forma segura el contador de visitas';
COMMENT ON VIEW propiedades_destacadas_metricas IS 'Vista de propiedades destacadas con métricas de popularidad';

-- Actualizar comentario de la tabla
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con marketing y SEO en migración 008';