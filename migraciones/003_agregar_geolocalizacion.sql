-- Migración 003: Agregar campos de geolocalización a propiedades
-- Fecha: 2025-01-15
-- Propósito: Habilitar funcionalidad de mapas en el frontend

-- Agregar campos de geolocalización
ALTER TABLE propiedades 
ADD COLUMN latitud DECIMAL(10,8) DEFAULT 0.0 NOT NULL,
ADD COLUMN longitud DECIMAL(11,8) DEFAULT 0.0 NOT NULL,
ADD COLUMN precision_ubicacion VARCHAR(20) DEFAULT 'sector' NOT NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.latitud IS 'Coordenada GPS latitud. Rango para Ecuador: -5.0 a 2.0';
COMMENT ON COLUMN propiedades.longitud IS 'Coordenada GPS longitud. Rango para Ecuador: -92.0 a -75.0';
COMMENT ON COLUMN propiedades.precision_ubicacion IS 'Precisión de ubicación: exacta, aproximada, sector';

-- Agregar constraint para validar precisión de ubicación
ALTER TABLE propiedades 
ADD CONSTRAINT chk_precision_ubicacion 
CHECK (precision_ubicacion IN ('exacta', 'aproximada', 'sector'));

-- Agregar constraint para validar coordenadas de Ecuador
ALTER TABLE propiedades 
ADD CONSTRAINT chk_coordenadas_ecuador 
CHECK (
    (latitud = 0 AND longitud = 0) OR  -- (0,0) significa no configurado
    (latitud BETWEEN -5.0 AND 2.0 AND longitud BETWEEN -92.0 AND -75.0)
);

-- Crear índice espacial para búsquedas por ubicación (será útil para filtros de mapa)
CREATE INDEX idx_propiedades_ubicacion ON propiedades(latitud, longitud) 
WHERE latitud != 0 AND longitud != 0;

-- Crear índice para precision_ubicacion (para filtros)
CREATE INDEX idx_propiedades_precision ON propiedades(precision_ubicacion);

-- Insertar comentario en tabla de migración para tracking
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con geolocalización en migración 003';