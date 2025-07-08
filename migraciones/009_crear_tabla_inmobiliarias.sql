-- Migración 009: Crear tabla inmobiliarias
-- Fecha: 2025-01-15
-- Propósito: Tabla para gestionar empresas inmobiliarias

-- Crear la tabla inmobiliarias
CREATE TABLE inmobiliarias (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(255) NOT NULL,
    ruc VARCHAR(13) NOT NULL UNIQUE,
    direccion TEXT NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    sitio_web TEXT DEFAULT '' NOT NULL,
    descripcion TEXT DEFAULT '' NOT NULL,
    logo_url TEXT DEFAULT '' NOT NULL,
    activa BOOLEAN DEFAULT TRUE NOT NULL,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Agregar comentarios para documentación
COMMENT ON TABLE inmobiliarias IS 'Tabla de empresas inmobiliarias registradas en el sistema';
COMMENT ON COLUMN inmobiliarias.id IS 'Identificador único de la inmobiliaria';
COMMENT ON COLUMN inmobiliarias.nombre IS 'Nombre comercial de la inmobiliaria';
COMMENT ON COLUMN inmobiliarias.ruc IS 'RUC (Registro Único de Contribuyentes) de la empresa';
COMMENT ON COLUMN inmobiliarias.direccion IS 'Dirección física de la oficina principal';
COMMENT ON COLUMN inmobiliarias.telefono IS 'Número de teléfono principal';
COMMENT ON COLUMN inmobiliarias.email IS 'Email de contacto principal';
COMMENT ON COLUMN inmobiliarias.sitio_web IS 'URL del sitio web de la inmobiliaria';
COMMENT ON COLUMN inmobiliarias.descripcion IS 'Descripción de la empresa y servicios';
COMMENT ON COLUMN inmobiliarias.logo_url IS 'URL del logo de la inmobiliaria';
COMMENT ON COLUMN inmobiliarias.activa IS 'Indica si la inmobiliaria está activa en el sistema';
COMMENT ON COLUMN inmobiliarias.fecha_creacion IS 'Fecha de registro en el sistema';
COMMENT ON COLUMN inmobiliarias.fecha_actualizacion IS 'Fecha de última actualización';

-- Agregar constraints para validar RUC ecuatoriano
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_ruc_formato 
CHECK (ruc ~ '^[0-9]{13}$');

-- Constraint para validar que el RUC termine en 001 (empresas)
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_ruc_empresa 
CHECK (ruc ~ '[0-9]{10}001$');

-- Constraint para validar email
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_email_formato 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Constraint para validar teléfono ecuatoriano
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_telefono_formato 
CHECK (telefono ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');

-- Constraint para validar URL del sitio web
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_sitio_web_url 
CHECK (sitio_web = '' OR sitio_web ~* '^https?://.*');

-- Constraint para validar URL del logo
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_logo_url 
CHECK (logo_url = '' OR logo_url ~* '^https?://.*\.(jpg|jpeg|png|webp|svg)(\\?.*)?$');

-- Crear índices para búsquedas
CREATE INDEX idx_inmobiliarias_nombre ON inmobiliarias(nombre);
CREATE INDEX idx_inmobiliarias_activa ON inmobiliarias(activa) WHERE activa = TRUE;
CREATE INDEX idx_inmobiliarias_fecha_creacion ON inmobiliarias(fecha_creacion);

-- Índice para búsquedas de texto en nombre y descripción
CREATE INDEX idx_inmobiliarias_busqueda_texto ON inmobiliarias 
USING gin(to_tsvector('spanish', nombre || ' ' || descripcion));

-- Función para validar RUC ecuatoriano (algoritmo módulo 11)
CREATE OR REPLACE FUNCTION validar_ruc_ecuatoriano(ruc_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    ruc_clean TEXT;
    digitos INTEGER[];
    suma INTEGER := 0;
    digito_verificador INTEGER;
    i INTEGER;
BEGIN
    -- Limpiar RUC (solo números)
    ruc_clean := regexp_replace(ruc_input, '[^0-9]', '', 'g');
    
    -- Verificar longitud
    IF LENGTH(ruc_clean) != 13 THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que termine en 001 (empresas)
    IF RIGHT(ruc_clean, 3) != '001' THEN
        RETURN FALSE;
    END IF;
    
    -- Convertir a array de dígitos
    FOR i IN 1..10 LOOP
        digitos[i] := CAST(SUBSTRING(ruc_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- Validar que el tercer dígito sea menor a 6 (personas jurídicas)
    IF digitos[3] >= 6 THEN
        RETURN FALSE;
    END IF;
    
    -- Aplicar algoritmo módulo 11
    FOR i IN 1..9 LOOP
        suma := suma + (digitos[i] * (10 - i));
    END LOOP;
    
    digito_verificador := 11 - (suma % 11);
    
    IF digito_verificador = 11 THEN
        digito_verificador := 0;
    ELSIF digito_verificador = 10 THEN
        digito_verificador := 1;
    END IF;
    
    RETURN digito_verificador = digitos[10];
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función de validación RUC
ALTER TABLE inmobiliarias 
ADD CONSTRAINT chk_ruc_valido 
CHECK (validar_ruc_ecuatoriano(ruc));

-- Trigger para actualizar fecha_actualizacion automáticamente
CREATE OR REPLACE FUNCTION actualizar_fecha_modificacion() 
RETURNS TRIGGER AS $$
BEGIN
    NEW.fecha_actualizacion := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_actualizar_fecha_inmobiliarias
    BEFORE UPDATE ON inmobiliarias
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_fecha_modificacion();

-- Función para obtener inmobiliarias activas
CREATE OR REPLACE FUNCTION obtener_inmobiliarias_activas()
RETURNS TABLE(
    id UUID,
    nombre VARCHAR(255),
    ruc VARCHAR(13),
    telefono VARCHAR(20),
    email VARCHAR(255),
    sitio_web TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.id,
        i.nombre,
        i.ruc,
        i.telefono,
        i.email,
        i.sitio_web
    FROM inmobiliarias i
    WHERE i.activa = TRUE
    ORDER BY i.nombre;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar inmobiliarias por nombre
CREATE OR REPLACE FUNCTION buscar_inmobiliarias_por_nombre(nombre_busqueda TEXT)
RETURNS TABLE(
    id UUID,
    nombre VARCHAR(255),
    descripcion TEXT,
    telefono VARCHAR(20),
    email VARCHAR(255)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.id,
        i.nombre,
        i.descripcion,
        i.telefono,
        i.email
    FROM inmobiliarias i
    WHERE i.activa = TRUE
      AND to_tsvector('spanish', i.nombre || ' ' || i.descripcion) @@ plainto_tsquery('spanish', nombre_busqueda)
    ORDER BY ts_rank(to_tsvector('spanish', i.nombre || ' ' || i.descripcion), plainto_tsquery('spanish', nombre_busqueda)) DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de inmobiliarias
CREATE OR REPLACE FUNCTION obtener_estadisticas_inmobiliarias()
RETURNS TABLE(
    total_inmobiliarias BIGINT,
    inmobiliarias_activas BIGINT,
    inmobiliarias_inactivas BIGINT,
    porcentaje_activas DECIMAL(5,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_inmobiliarias,
        COUNT(*) FILTER (WHERE activa = TRUE) as inmobiliarias_activas,
        COUNT(*) FILTER (WHERE activa = FALSE) as inmobiliarias_inactivas,
        ROUND(
            (COUNT(*) FILTER (WHERE activa = TRUE) * 100.0 / NULLIF(COUNT(*), 0)), 
            2
        ) as porcentaje_activas
    FROM inmobiliarias;
END;
$$ LANGUAGE plpgsql;

-- Comentarios en las funciones
COMMENT ON FUNCTION validar_ruc_ecuatoriano IS 'Valida RUC ecuatoriano usando algoritmo módulo 11';
COMMENT ON FUNCTION actualizar_fecha_modificacion IS 'Actualiza automáticamente fecha_actualizacion en modificaciones';
COMMENT ON FUNCTION obtener_inmobiliarias_activas IS 'Retorna lista de inmobiliarias activas';
COMMENT ON FUNCTION buscar_inmobiliarias_por_nombre IS 'Busca inmobiliarias por nombre usando búsqueda de texto completo';
COMMENT ON FUNCTION obtener_estadisticas_inmobiliarias IS 'Retorna estadísticas generales de inmobiliarias';

-- Insertar inmobiliaria de ejemplo para testing
INSERT INTO inmobiliarias (
    nombre, 
    ruc, 
    direccion, 
    telefono, 
    email, 
    sitio_web, 
    descripcion
) VALUES (
    'Inmobiliaria Ejemplo S.A.',
    '1792146739001',  -- RUC válido de ejemplo
    'Av. 9 de Octubre 123, Guayaquil, Ecuador',
    '042234567',
    'info@inmobiliariaejemplo.com',
    'https://www.inmobiliariaejemplo.com',
    'Empresa líder en bienes raíces en Ecuador con más de 20 años de experiencia.'
);