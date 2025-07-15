-- Migración: Agregar campo parking_spaces a la tabla properties
-- Fecha: 2025-01-08
-- Descripción: Agregar el campo parking_spaces para rastrear espacios de estacionamiento disponibles

-- Agregar campo parking_spaces a la tabla properties
ALTER TABLE properties ADD COLUMN IF NOT EXISTS parking_spaces INTEGER DEFAULT 0 CHECK (parking_spaces >= 0);

-- Agregar comentario para documentar el campo
COMMENT ON COLUMN properties.parking_spaces IS 'Número de espacios de estacionamiento disponibles';

-- Crear índice para mejorar consultas de filtrado por parking_spaces
CREATE INDEX IF NOT EXISTS idx_properties_parking_spaces ON properties(parking_spaces);

-- Actualizar registros existentes con valor por defecto 0 si es necesario
UPDATE properties SET parking_spaces = 0 WHERE parking_spaces IS NULL;