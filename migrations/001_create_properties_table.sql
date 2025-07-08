-- Migración: Crear tabla propiedades
-- Fecha: 2024-07-04
-- Descripción: Tabla principal para almacenar propiedades inmobiliarias

-- Crear tabla propiedades
CREATE TABLE IF NOT EXISTS propiedades (
    id VARCHAR(36) PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT,
    precio DECIMAL(15,2) NOT NULL CHECK (precio > 0),
    
    -- Ubicación
    provincia VARCHAR(100) NOT NULL,
    ciudad VARCHAR(100) NOT NULL,
    sector VARCHAR(100),
    direccion VARCHAR(255),
    
    -- Características
    tipo VARCHAR(50) NOT NULL CHECK (tipo IN ('casa', 'departamento', 'terreno', 'comercial')),
    estado VARCHAR(50) NOT NULL DEFAULT 'disponible' CHECK (estado IN ('disponible', 'vendida', 'alquilada', 'reservada')),
    dormitorios INTEGER DEFAULT 0 CHECK (dormitorios >= 0),
    banos DECIMAL(3,1) DEFAULT 0 CHECK (banos >= 0),
    area_m2 DECIMAL(10,2) DEFAULT 0 CHECK (area_m2 >= 0),
    
    -- Auditoría
    fecha_creacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    fecha_actualizacion TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Crear índices para mejorar consultas
CREATE INDEX IF NOT EXISTS idx_propiedades_provincia ON propiedades(provincia);
CREATE INDEX IF NOT EXISTS idx_propiedades_ciudad ON propiedades(ciudad);
CREATE INDEX IF NOT EXISTS idx_propiedades_tipo ON propiedades(tipo);
CREATE INDEX IF NOT EXISTS idx_propiedades_estado ON propiedades(estado);
CREATE INDEX IF NOT EXISTS idx_propiedades_precio ON propiedades(precio);
CREATE INDEX IF NOT EXISTS idx_propiedades_fecha_creacion ON propiedades(fecha_creacion);

-- Comentarios para documentar la tabla
COMMENT ON TABLE propiedades IS 'Tabla principal para propiedades inmobiliarias';
COMMENT ON COLUMN propiedades.id IS 'Identificador único UUID';
COMMENT ON COLUMN propiedades.precio IS 'Precio en USD';
COMMENT ON COLUMN propiedades.area_m2 IS 'Área en metros cuadrados';
COMMENT ON COLUMN propiedades.banos IS 'Número de baños (puede ser decimal: 2.5)';

-- Insertar datos de ejemplo para desarrollo
INSERT INTO propiedades (
    id, titulo, descripcion, precio, provincia, ciudad, sector, 
    tipo, estado, dormitorios, banos, area_m2
) VALUES 
(
    'prop-001', 
    'Hermosa casa en Samborondón', 
    'Casa moderna de 3 pisos con piscina y jardín',
    285000.00, 
    'Guayas', 
    'Samborondón', 
    'La Puntilla',
    'casa', 
    'disponible', 
    4, 
    3.5, 
    320.00
),
(
    'prop-002',
    'Departamento céntrico en Quito',
    'Departamento moderno en el centro norte de Quito',
    125000.00,
    'Pichincha',
    'Quito',
    'La Carolina',
    'departamento',
    'disponible',
    2,
    2.0,
    85.00
),
(
    'prop-003',
    'Terreno en Cuenca',
    'Terreno plano ideal para construcción',
    45000.00,
    'Azuay',
    'Cuenca',
    'El Batán',
    'terreno',
    'disponible',
    0,
    0.0,
    500.00
) ON CONFLICT (id) DO NOTHING;