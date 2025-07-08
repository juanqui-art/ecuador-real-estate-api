-- Migración 007: Agregar campos de amenidades a propiedades
-- Fecha: 2025-01-15
-- Propósito: Habilitar filtros por amenidades en el frontend

-- Agregar campos para amenidades principales
ALTER TABLE propiedades 
ADD COLUMN garage BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN piscina BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN jardin BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN terraza BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN balcon BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN seguridad BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN ascensor BOOLEAN DEFAULT FALSE NOT NULL,
ADD COLUMN aire_acondicionado BOOLEAN DEFAULT FALSE NOT NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.garage IS 'Indica si la propiedad tiene garage o estacionamiento';
COMMENT ON COLUMN propiedades.piscina IS 'Indica si la propiedad tiene piscina';
COMMENT ON COLUMN propiedades.jardin IS 'Indica si la propiedad tiene jardín';
COMMENT ON COLUMN propiedades.terraza IS 'Indica si la propiedad tiene terraza';
COMMENT ON COLUMN propiedades.balcon IS 'Indica si la propiedad tiene balcón';
COMMENT ON COLUMN propiedades.seguridad IS 'Indica si la propiedad está en conjunto con seguridad 24/7';
COMMENT ON COLUMN propiedades.ascensor IS 'Indica si el edificio tiene ascensor';
COMMENT ON COLUMN propiedades.aire_acondicionado IS 'Indica si la propiedad tiene aire acondicionado';

-- Crear índices individuales para cada amenidad (para filtros específicos)
CREATE INDEX idx_propiedades_garage ON propiedades(garage) WHERE garage = TRUE;
CREATE INDEX idx_propiedades_piscina ON propiedades(piscina) WHERE piscina = TRUE;
CREATE INDEX idx_propiedades_jardin ON propiedades(jardin) WHERE jardin = TRUE;
CREATE INDEX idx_propiedades_terraza ON propiedades(terraza) WHERE terraza = TRUE;
CREATE INDEX idx_propiedades_balcon ON propiedades(balcon) WHERE balcon = TRUE;
CREATE INDEX idx_propiedades_seguridad ON propiedades(seguridad) WHERE seguridad = TRUE;
CREATE INDEX idx_propiedades_ascensor ON propiedades(ascensor) WHERE ascensor = TRUE;
CREATE INDEX idx_propiedades_aire_acondicionado ON propiedades(aire_acondicionado) WHERE aire_acondicionado = TRUE;

-- Índice compuesto para amenidades más solicitadas (optimización para filtros combinados)
CREATE INDEX idx_propiedades_amenidades_principales ON propiedades(garage, piscina, seguridad, aire_acondicionado);

-- Índice para amenidades de espacios exteriores
CREATE INDEX idx_propiedades_espacios_exteriores ON propiedades(jardin, terraza, balcon, piscina);

-- Función para contar amenidades de una propiedad
CREATE OR REPLACE FUNCTION contar_amenidades(
    p_garage BOOLEAN DEFAULT FALSE,
    p_piscina BOOLEAN DEFAULT FALSE,
    p_jardin BOOLEAN DEFAULT FALSE,
    p_terraza BOOLEAN DEFAULT FALSE,
    p_balcon BOOLEAN DEFAULT FALSE,
    p_seguridad BOOLEAN DEFAULT FALSE,
    p_ascensor BOOLEAN DEFAULT FALSE,
    p_aire_acondicionado BOOLEAN DEFAULT FALSE
) RETURNS INTEGER AS $$
BEGIN
    RETURN (
        CASE WHEN p_garage THEN 1 ELSE 0 END +
        CASE WHEN p_piscina THEN 1 ELSE 0 END +
        CASE WHEN p_jardin THEN 1 ELSE 0 END +
        CASE WHEN p_terraza THEN 1 ELSE 0 END +
        CASE WHEN p_balcon THEN 1 ELSE 0 END +
        CASE WHEN p_seguridad THEN 1 ELSE 0 END +
        CASE WHEN p_ascensor THEN 1 ELSE 0 END +
        CASE WHEN p_aire_acondicionado THEN 1 ELSE 0 END
    );
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Función para obtener lista de amenidades activas como array
CREATE OR REPLACE FUNCTION obtener_amenidades_activas(
    p_garage BOOLEAN DEFAULT FALSE,
    p_piscina BOOLEAN DEFAULT FALSE,
    p_jardin BOOLEAN DEFAULT FALSE,
    p_terraza BOOLEAN DEFAULT FALSE,
    p_balcon BOOLEAN DEFAULT FALSE,
    p_seguridad BOOLEAN DEFAULT FALSE,
    p_ascensor BOOLEAN DEFAULT FALSE,
    p_aire_acondicionado BOOLEAN DEFAULT FALSE
) RETURNS TEXT[] AS $$
DECLARE
    amenidades TEXT[] := '{}';
BEGIN
    IF p_garage THEN amenidades := array_append(amenidades, 'Garage'); END IF;
    IF p_piscina THEN amenidades := array_append(amenidades, 'Piscina'); END IF;
    IF p_jardin THEN amenidades := array_append(amenidades, 'Jardín'); END IF;
    IF p_terraza THEN amenidades := array_append(amenidades, 'Terraza'); END IF;
    IF p_balcon THEN amenidades := array_append(amenidades, 'Balcón'); END IF;
    IF p_seguridad THEN amenidades := array_append(amenidades, 'Seguridad'); END IF;
    IF p_ascensor THEN amenidades := array_append(amenidades, 'Ascensor'); END IF;
    IF p_aire_acondicionado THEN amenidades := array_append(amenidades, 'Aire Acondicionado'); END IF;
    
    RETURN amenidades;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Vista para propiedades con resumen de amenidades (útil para listados)
CREATE OR REPLACE VIEW propiedades_con_amenidades AS
SELECT 
    *,
    contar_amenidades(garage, piscina, jardin, terraza, balcon, seguridad, ascensor, aire_acondicionado) AS total_amenidades,
    obtener_amenidades_activas(garage, piscina, jardin, terraza, balcon, seguridad, ascensor, aire_acondicionado) AS lista_amenidades
FROM propiedades;

-- Comentario en la vista
COMMENT ON VIEW propiedades_con_amenidades IS 'Vista que incluye conteo y lista de amenidades activas para cada propiedad';

-- Función para filtrar propiedades por amenidades mínimas
CREATE OR REPLACE FUNCTION buscar_por_amenidades_minimas(amenidades_minimas INTEGER DEFAULT 3)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    total_amenidades INTEGER,
    lista_amenidades TEXT[]
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        pca.total_amenidades,
        pca.lista_amenidades
    FROM propiedades_con_amenidades pca
    JOIN propiedades p ON p.id = pca.id
    WHERE pca.total_amenidades >= amenidades_minimas
    ORDER BY pca.total_amenidades DESC, p.precio ASC;
END;
$$ LANGUAGE plpgsql;

-- Función para propiedades premium (con múltiples amenidades de lujo)
CREATE OR REPLACE FUNCTION obtener_propiedades_premium()
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    provincia VARCHAR(50),
    ciudad VARCHAR(100)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.provincia,
        p.ciudad
    FROM propiedades p
    WHERE p.piscina = TRUE 
      AND p.seguridad = TRUE 
      AND (p.garage = TRUE OR p.aire_acondicionado = TRUE)
      AND p.precio > 150000  -- Umbral de precio premium
    ORDER BY p.precio DESC;
END;
$$ LANGUAGE plpgsql;

-- Comentarios en las funciones
COMMENT ON FUNCTION contar_amenidades IS 'Cuenta el número total de amenidades activas de una propiedad';
COMMENT ON FUNCTION obtener_amenidades_activas IS 'Retorna array con nombres de amenidades activas';
COMMENT ON FUNCTION buscar_por_amenidades_minimas IS 'Busca propiedades que tienen al menos X amenidades';
COMMENT ON FUNCTION obtener_propiedades_premium IS 'Retorna propiedades de lujo basadas en amenidades y precio';

-- Actualizar comentario de la tabla
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con amenidades en migración 007';