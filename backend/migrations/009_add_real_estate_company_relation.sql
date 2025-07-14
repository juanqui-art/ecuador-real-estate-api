-- Migración 010: Agregar relación FK entre propiedades e inmobiliarias
-- Fecha: 2025-01-15
-- Propósito: Establecer relación entre propiedades y las inmobiliarias que las gestionan

-- Agregar campo inmobiliaria_id a la tabla propiedades
ALTER TABLE propiedades 
ADD COLUMN inmobiliaria_id UUID NULL;

-- Agregar comentario
COMMENT ON COLUMN propiedades.inmobiliaria_id IS 'Inmobiliaria que gestiona la propiedad (FK hacia inmobiliarias)';

-- Crear la llave foránea
ALTER TABLE propiedades 
ADD CONSTRAINT fk_propiedades_inmobiliaria 
FOREIGN KEY (inmobiliaria_id) 
REFERENCES inmobiliarias(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Crear índice para la FK (mejora performance de JOINs)
CREATE INDEX idx_propiedades_inmobiliaria_id ON propiedades(inmobiliaria_id);

-- Índice compuesto para búsquedas por inmobiliaria y estado
CREATE INDEX idx_propiedades_inmobiliaria_estado ON propiedades(inmobiliaria_id, estado) 
WHERE inmobiliaria_id IS NOT NULL;

-- Función para obtener propiedades por inmobiliaria
CREATE OR REPLACE FUNCTION obtener_propiedades_por_inmobiliaria(inmobiliaria_uuid UUID)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    provincia VARCHAR(50),
    ciudad VARCHAR(100),
    tipo VARCHAR(20),
    estado VARCHAR(20),
    area_m2 DECIMAL(10,2),
    dormitorios INTEGER,
    banos DECIMAL(3,1)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.provincia,
        p.ciudad,
        p.tipo,
        p.estado,
        p.area_m2,
        p.dormitorios,
        p.banos
    FROM propiedades p
    WHERE p.inmobiliaria_id = inmobiliaria_uuid
      AND p.estado = 'disponible'
    ORDER BY p.fecha_creacion DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de propiedades por inmobiliaria
CREATE OR REPLACE FUNCTION obtener_estadisticas_propiedades_por_inmobiliaria(inmobiliaria_uuid UUID)
RETURNS TABLE(
    total_propiedades BIGINT,
    propiedades_disponibles BIGINT,
    propiedades_vendidas BIGINT,
    propiedades_alquiladas BIGINT,
    precio_promedio DECIMAL(15,2),
    precio_minimo DECIMAL(15,2),
    precio_maximo DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_propiedades,
        COUNT(*) FILTER (WHERE estado = 'disponible') as propiedades_disponibles,
        COUNT(*) FILTER (WHERE estado = 'vendida') as propiedades_vendidas,
        COUNT(*) FILTER (WHERE estado = 'alquilada') as propiedades_alquiladas,
        ROUND(AVG(precio), 2) as precio_promedio,
        MIN(precio) as precio_minimo,
        MAX(precio) as precio_maximo
    FROM propiedades p
    WHERE p.inmobiliaria_id = inmobiliaria_uuid;
END;
$$ LANGUAGE plpgsql;

-- Vista para propiedades con información de inmobiliaria
CREATE OR REPLACE VIEW propiedades_con_inmobiliaria AS
SELECT 
    p.*,
    i.nombre as inmobiliaria_nombre,
    i.telefono as inmobiliaria_telefono,
    i.email as inmobiliaria_email,
    i.sitio_web as inmobiliaria_sitio_web,
    i.logo_url as inmobiliaria_logo
FROM propiedades p
LEFT JOIN inmobiliarias i ON p.inmobiliaria_id = i.id
WHERE i.activa = TRUE OR p.inmobiliaria_id IS NULL;

-- Comentario en la vista
COMMENT ON VIEW propiedades_con_inmobiliaria IS 'Vista que incluye información de la inmobiliaria asociada a cada propiedad';

-- Función para obtener propiedades populares por inmobiliaria
CREATE OR REPLACE FUNCTION obtener_propiedades_populares_por_inmobiliaria(
    inmobiliaria_uuid UUID,
    limite INTEGER DEFAULT 10
)
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
    WHERE p.inmobiliaria_id = inmobiliaria_uuid
      AND p.estado = 'disponible'
    ORDER BY p.visitas_contador DESC, p.destacada DESC
    LIMIT limite;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener ranking de inmobiliarias por número de propiedades
CREATE OR REPLACE FUNCTION obtener_ranking_inmobiliarias_por_propiedades()
RETURNS TABLE(
    inmobiliaria_id UUID,
    inmobiliaria_nombre VARCHAR(255),
    total_propiedades BIGINT,
    propiedades_disponibles BIGINT,
    propiedades_destacadas BIGINT,
    precio_promedio DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.id,
        i.nombre,
        COUNT(p.id) as total_propiedades,
        COUNT(p.id) FILTER (WHERE p.estado = 'disponible') as propiedades_disponibles,
        COUNT(p.id) FILTER (WHERE p.destacada = TRUE) as propiedades_destacadas,
        ROUND(AVG(p.precio), 2) as precio_promedio
    FROM inmobiliarias i
    LEFT JOIN propiedades p ON i.id = p.inmobiliaria_id
    WHERE i.activa = TRUE
    GROUP BY i.id, i.nombre
    HAVING COUNT(p.id) > 0
    ORDER BY total_propiedades DESC, propiedades_disponibles DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para validar asignación de inmobiliaria (trigger)
CREATE OR REPLACE FUNCTION validar_asignacion_inmobiliaria() 
RETURNS TRIGGER AS $$
BEGIN
    -- Verificar que la inmobiliaria existe y está activa
    IF NEW.inmobiliaria_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM inmobiliarias 
            WHERE id = NEW.inmobiliaria_id AND activa = TRUE
        ) THEN
            RAISE EXCEPTION 'La inmobiliaria especificada no existe o no está activa';
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar asignación
CREATE TRIGGER trigger_validar_inmobiliaria
    BEFORE INSERT OR UPDATE OF inmobiliaria_id ON propiedades
    FOR EACH ROW
    EXECUTE FUNCTION validar_asignacion_inmobiliaria();

-- Función para reasignar propiedades cuando se desactiva una inmobiliaria
CREATE OR REPLACE FUNCTION manejar_desactivacion_inmobiliaria() 
RETURNS TRIGGER AS $$
BEGIN
    -- Si se desactiva una inmobiliaria, quitar la asignación de sus propiedades
    IF OLD.activa = TRUE AND NEW.activa = FALSE THEN
        UPDATE propiedades 
        SET inmobiliaria_id = NULL 
        WHERE inmobiliaria_id = NEW.id;
        
        RAISE NOTICE 'Se han desasignado % propiedades de la inmobiliaria %', 
                     ROW_COUNT, NEW.nombre;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para manejar desactivación
CREATE TRIGGER trigger_manejar_desactivacion_inmobiliaria
    AFTER UPDATE OF activa ON inmobiliarias
    FOR EACH ROW
    EXECUTE FUNCTION manejar_desactivacion_inmobiliaria();

-- Comentarios en las funciones
COMMENT ON FUNCTION obtener_propiedades_por_inmobiliaria IS 'Retorna propiedades gestionadas por una inmobiliaria específica';
COMMENT ON FUNCTION obtener_estadisticas_propiedades_por_inmobiliaria IS 'Retorna estadísticas de propiedades por inmobiliaria';
COMMENT ON FUNCTION obtener_propiedades_populares_por_inmobiliaria IS 'Retorna propiedades más visitadas de una inmobiliaria';
COMMENT ON FUNCTION obtener_ranking_inmobiliarias_por_propiedades IS 'Retorna ranking de inmobiliarias por número de propiedades';
COMMENT ON FUNCTION validar_asignacion_inmobiliaria IS 'Valida que la inmobiliaria asignada existe y está activa';
COMMENT ON FUNCTION manejar_desactivacion_inmobiliaria IS 'Desasigna propiedades cuando se desactiva una inmobiliaria';

-- Actualizar comentario de la tabla propiedades
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con relación a inmobiliarias en migración 010';

-- Actualizar algunas propiedades de ejemplo con la inmobiliaria creada
-- (Solo si existen propiedades en la tabla)
UPDATE propiedades 
SET inmobiliaria_id = (
    SELECT id FROM inmobiliarias 
    WHERE nombre = 'Inmobiliaria Ejemplo S.A.' 
    LIMIT 1
)
WHERE inmobiliaria_id IS NULL 
  AND id IN (
    SELECT id FROM propiedades 
    ORDER BY fecha_creacion 
    LIMIT 3
  );