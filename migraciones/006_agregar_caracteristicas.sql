-- Migración 006: Agregar características detalladas a propiedades
-- Fecha: 2025-01-15
-- Propósito: Agregar información detallada de construcción y estado

-- Agregar campos para características detalladas
ALTER TABLE propiedades 
ADD COLUMN ano_construccion INTEGER NULL,
ADD COLUMN pisos INTEGER NULL,
ADD COLUMN estado_propiedad VARCHAR(20) DEFAULT 'usada' NOT NULL,
ADD COLUMN amoblada BOOLEAN DEFAULT FALSE NOT NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.ano_construccion IS 'Año de construcción de la propiedad (NULL si no se conoce)';
COMMENT ON COLUMN propiedades.pisos IS 'Número de pisos de la propiedad (NULL si no aplica)';
COMMENT ON COLUMN propiedades.estado_propiedad IS 'Estado de la propiedad: nueva, usada, remodelada';
COMMENT ON COLUMN propiedades.amoblada IS 'Indica si la propiedad viene amoblada';

-- Agregar constraints para validar datos
ALTER TABLE propiedades 
ADD CONSTRAINT chk_ano_construccion_valido 
CHECK (
    ano_construccion IS NULL OR 
    (ano_construccion >= 1800 AND ano_construccion <= EXTRACT(YEAR FROM CURRENT_DATE) + 2)
);

ALTER TABLE propiedades 
ADD CONSTRAINT chk_pisos_valido 
CHECK (pisos IS NULL OR (pisos >= 1 AND pisos <= 50));

ALTER TABLE propiedades 
ADD CONSTRAINT chk_estado_propiedad_valido 
CHECK (estado_propiedad IN ('nueva', 'usada', 'remodelada'));

-- Crear índices para búsquedas y filtros
CREATE INDEX idx_propiedades_ano_construccion ON propiedades(ano_construccion) 
WHERE ano_construccion IS NOT NULL;

CREATE INDEX idx_propiedades_pisos ON propiedades(pisos) 
WHERE pisos IS NOT NULL;

CREATE INDEX idx_propiedades_estado ON propiedades(estado_propiedad);

CREATE INDEX idx_propiedades_amoblada ON propiedades(amoblada) 
WHERE amoblada = TRUE;

-- Índice compuesto para filtros de características
CREATE INDEX idx_propiedades_caracteristicas ON propiedades(estado_propiedad, amoblada, pisos) 
WHERE pisos IS NOT NULL;

-- Función para calcular la antigüedad de la propiedad
CREATE OR REPLACE FUNCTION calcular_antiguedad(ano_construccion INTEGER) 
RETURNS INTEGER AS $$
BEGIN
    IF ano_construccion IS NULL THEN
        RETURN NULL;
    END IF;
    
    RETURN EXTRACT(YEAR FROM CURRENT_DATE) - ano_construccion;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- Vista para propiedades con antigüedad calculada (útil para reportes)
CREATE OR REPLACE VIEW propiedades_con_antiguedad AS
SELECT 
    *,
    calcular_antiguedad(ano_construccion) AS antiguedad_anos,
    CASE 
        WHEN ano_construccion IS NULL THEN 'No especificado'
        WHEN calcular_antiguedad(ano_construccion) <= 2 THEN 'Nueva'
        WHEN calcular_antiguedad(ano_construccion) <= 10 THEN 'Moderna'
        WHEN calcular_antiguedad(ano_construccion) <= 25 THEN 'Establecida'
        ELSE 'Antigua'
    END AS categoria_antiguedad
FROM propiedades;

-- Comentario en la vista
COMMENT ON VIEW propiedades_con_antiguedad IS 'Vista que incluye cálculo automático de antigüedad y categorización';

-- Función para validar consistencia de datos de construcción
CREATE OR REPLACE FUNCTION validar_datos_construccion() 
RETURNS TRIGGER AS $$
BEGIN
    -- Si es nueva, debe tener año de construcción reciente
    IF NEW.estado_propiedad = 'nueva' AND NEW.ano_construccion IS NOT NULL THEN
        IF NEW.ano_construccion < EXTRACT(YEAR FROM CURRENT_DATE) - 2 THEN
            RAISE EXCEPTION 'Una propiedad "nueva" no puede tener más de 2 años de construcción';
        END IF;
    END IF;
    
    -- Si tiene más de 1 piso, debe especificar el número
    IF NEW.tipo IN ('casa', 'comercial') AND NEW.pisos IS NULL THEN
        -- Esto es solo una sugerencia, no un error
        -- Se podría implementar un warning system aquí
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar datos de construcción
CREATE TRIGGER trigger_validar_construccion
    BEFORE INSERT OR UPDATE ON propiedades
    FOR EACH ROW
    EXECUTE FUNCTION validar_datos_construccion();

-- Actualizar comentario de la tabla
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con características detalladas en migración 006';