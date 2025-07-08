-- Migración 005: Agregar campos de precios adicionales a propiedades
-- Fecha: 2025-01-15
-- Propósito: Manejar diferentes tipos de precios (alquiler, gastos comunes, precio/m²)

-- Agregar campos para precios adicionales
ALTER TABLE propiedades 
ADD COLUMN precio_alquiler DECIMAL(15,2) NULL,
ADD COLUMN gastos_comunes DECIMAL(10,2) NULL,
ADD COLUMN precio_m2 DECIMAL(10,2) NULL;

-- Agregar comentarios para documentación
COMMENT ON COLUMN propiedades.precio_alquiler IS 'Precio mensual de alquiler de la propiedad (NULL si no está en alquiler)';
COMMENT ON COLUMN propiedades.gastos_comunes IS 'Gastos comunes mensuales de condominio/administración';
COMMENT ON COLUMN propiedades.precio_m2 IS 'Precio por metro cuadrado (calculado automáticamente)';

-- Agregar constraints para validar precios positivos
ALTER TABLE propiedades 
ADD CONSTRAINT chk_precio_alquiler_positivo 
CHECK (precio_alquiler IS NULL OR precio_alquiler > 0);

ALTER TABLE propiedades 
ADD CONSTRAINT chk_gastos_comunes_positivos 
CHECK (gastos_comunes IS NULL OR gastos_comunes >= 0);

ALTER TABLE propiedades 
ADD CONSTRAINT chk_precio_m2_positivo 
CHECK (precio_m2 IS NULL OR precio_m2 > 0);

-- Constraint para verificar que precio_m2 sea consistente con precio/area_m2
ALTER TABLE propiedades 
ADD CONSTRAINT chk_precio_m2_consistente 
CHECK (
    precio_m2 IS NULL OR 
    area_m2 = 0 OR 
    precio_m2 = ROUND(precio / NULLIF(area_m2, 0), 2)
);

-- Crear índices para búsquedas por rangos de precio
CREATE INDEX idx_propiedades_precio_alquiler ON propiedades(precio_alquiler) 
WHERE precio_alquiler IS NOT NULL;

CREATE INDEX idx_propiedades_gastos_comunes ON propiedades(gastos_comunes) 
WHERE gastos_comunes IS NOT NULL;

CREATE INDEX idx_propiedades_precio_m2 ON propiedades(precio_m2) 
WHERE precio_m2 IS NOT NULL;

-- Índice compuesto para filtros de rango de precio de venta
CREATE INDEX idx_propiedades_precio_venta_rango ON propiedades(precio, area_m2);

-- Índice compuesto para filtros de rango de precio de alquiler
CREATE INDEX idx_propiedades_precio_alquiler_rango ON propiedades(precio_alquiler, area_m2) 
WHERE precio_alquiler IS NOT NULL;

-- Función para actualizar automáticamente precio_m2 cuando cambien precio o area_m2
CREATE OR REPLACE FUNCTION actualizar_precio_m2() 
RETURNS TRIGGER AS $$
BEGIN
    -- Solo calcular si area_m2 es mayor a 0
    IF NEW.area_m2 > 0 THEN
        NEW.precio_m2 := ROUND(NEW.precio / NEW.area_m2, 2);
    ELSE
        NEW.precio_m2 := NULL;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para actualizar automáticamente precio_m2
CREATE TRIGGER trigger_actualizar_precio_m2
    BEFORE INSERT OR UPDATE OF precio, area_m2 ON propiedades
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_precio_m2();

-- Actualizar precio_m2 para registros existentes
UPDATE propiedades 
SET precio_m2 = CASE 
    WHEN area_m2 > 0 THEN ROUND(precio / area_m2, 2)
    ELSE NULL 
END
WHERE area_m2 IS NOT NULL;

-- Actualizar comentario de la tabla
COMMENT ON TABLE propiedades IS 'Tabla de propiedades inmobiliarias - Actualizada con precios adicionales en migración 005';