-- Migración 011: Crear tabla usuarios
-- Fecha: 2025-01-15
-- Propósito: Tabla para gestionar usuarios del sistema (compradores, vendedores, agentes)

-- Crear la tabla usuarios
CREATE TABLE usuarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefono VARCHAR(20) NOT NULL,
    cedula VARCHAR(10) NOT NULL UNIQUE,
    fecha_nacimiento DATE NULL,
    tipo_usuario VARCHAR(20) DEFAULT 'comprador' NOT NULL,
    activo BOOLEAN DEFAULT TRUE NOT NULL,
    presupuesto_min DECIMAL(15,2) NULL,
    presupuesto_max DECIMAL(15,2) NULL,
    provincias_interes JSONB DEFAULT '[]' NOT NULL,
    tipos_propiedad_interes JSONB DEFAULT '[]' NOT NULL,
    avatar_url TEXT DEFAULT '' NOT NULL,
    bio TEXT DEFAULT '' NOT NULL,
    inmobiliaria_id UUID NULL,
    fecha_creacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    fecha_actualizacion TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Agregar comentarios para documentación
COMMENT ON TABLE usuarios IS 'Tabla de usuarios del sistema inmobiliario';
COMMENT ON COLUMN usuarios.id IS 'Identificador único del usuario';
COMMENT ON COLUMN usuarios.nombre IS 'Nombre del usuario';
COMMENT ON COLUMN usuarios.apellido IS 'Apellido del usuario';
COMMENT ON COLUMN usuarios.email IS 'Email único del usuario';
COMMENT ON COLUMN usuarios.telefono IS 'Número de teléfono del usuario';
COMMENT ON COLUMN usuarios.cedula IS 'Número de cédula ecuatoriana (único)';
COMMENT ON COLUMN usuarios.fecha_nacimiento IS 'Fecha de nacimiento del usuario';
COMMENT ON COLUMN usuarios.tipo_usuario IS 'Tipo: comprador, vendedor, agente, admin';
COMMENT ON COLUMN usuarios.activo IS 'Estado del usuario en el sistema';
COMMENT ON COLUMN usuarios.presupuesto_min IS 'Presupuesto mínimo para compradores';
COMMENT ON COLUMN usuarios.presupuesto_max IS 'Presupuesto máximo para compradores';
COMMENT ON COLUMN usuarios.provincias_interes IS 'Array JSON de provincias de interés';
COMMENT ON COLUMN usuarios.tipos_propiedad_interes IS 'Array JSON de tipos de propiedad de interés';
COMMENT ON COLUMN usuarios.avatar_url IS 'URL del avatar del usuario';
COMMENT ON COLUMN usuarios.bio IS 'Biografía o descripción del usuario';
COMMENT ON COLUMN usuarios.inmobiliaria_id IS 'Inmobiliaria asociada (para agentes)';
COMMENT ON COLUMN usuarios.fecha_creacion IS 'Fecha de registro del usuario';
COMMENT ON COLUMN usuarios.fecha_actualizacion IS 'Fecha de última actualización';

-- Crear la relación con inmobiliarias
ALTER TABLE usuarios 
ADD CONSTRAINT fk_usuarios_inmobiliaria 
FOREIGN KEY (inmobiliaria_id) 
REFERENCES inmobiliarias(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Agregar constraints para validar datos
ALTER TABLE usuarios 
ADD CONSTRAINT chk_tipo_usuario_valido 
CHECK (tipo_usuario IN ('comprador', 'vendedor', 'agente', 'admin'));

-- Constraint para validar email
ALTER TABLE usuarios 
ADD CONSTRAINT chk_email_formato 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Constraint para validar teléfono ecuatoriano
ALTER TABLE usuarios 
ADD CONSTRAINT chk_telefono_formato 
CHECK (telefono ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');

-- Constraint para validar cédula ecuatoriana (10 dígitos)
ALTER TABLE usuarios 
ADD CONSTRAINT chk_cedula_formato 
CHECK (cedula ~ '^[0-9]{10}$');

-- Constraint para validar presupuesto
ALTER TABLE usuarios 
ADD CONSTRAINT chk_presupuesto_valido 
CHECK (
    presupuesto_min IS NULL OR 
    presupuesto_max IS NULL OR 
    presupuesto_min <= presupuesto_max
);

-- Constraint para validar que presupuestos sean positivos
ALTER TABLE usuarios 
ADD CONSTRAINT chk_presupuesto_positivo 
CHECK (
    (presupuesto_min IS NULL OR presupuesto_min > 0) AND
    (presupuesto_max IS NULL OR presupuesto_max > 0)
);

-- Constraint para validar fecha de nacimiento
ALTER TABLE usuarios 
ADD CONSTRAINT chk_fecha_nacimiento_valida 
CHECK (
    fecha_nacimiento IS NULL OR 
    (fecha_nacimiento >= '1900-01-01' AND fecha_nacimiento <= CURRENT_DATE - INTERVAL '18 years')
);

-- Constraint para validar URL del avatar
ALTER TABLE usuarios 
ADD CONSTRAINT chk_avatar_url 
CHECK (avatar_url = '' OR avatar_url ~* '^https?://.*\.(jpg|jpeg|png|webp|gif)(\\?.*)?$');

-- Constraint para validar que agentes tengan inmobiliaria
ALTER TABLE usuarios 
ADD CONSTRAINT chk_agente_inmobiliaria 
CHECK (
    tipo_usuario != 'agente' OR 
    inmobiliaria_id IS NOT NULL
);

-- Crear índices para búsquedas
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_usuarios_cedula ON usuarios(cedula);
CREATE INDEX idx_usuarios_tipo ON usuarios(tipo_usuario);
CREATE INDEX idx_usuarios_activo ON usuarios(activo) WHERE activo = TRUE;
CREATE INDEX idx_usuarios_inmobiliaria ON usuarios(inmobiliaria_id) WHERE inmobiliaria_id IS NOT NULL;
CREATE INDEX idx_usuarios_fecha_creacion ON usuarios(fecha_creacion);

-- Índice compuesto para búsquedas por tipo y estado
CREATE INDEX idx_usuarios_tipo_activo ON usuarios(tipo_usuario, activo);

-- Índice para búsquedas de presupuesto
CREATE INDEX idx_usuarios_presupuesto ON usuarios(presupuesto_min, presupuesto_max) 
WHERE presupuesto_min IS NOT NULL AND presupuesto_max IS NOT NULL;

-- Índice para búsquedas de texto en nombre completo
CREATE INDEX idx_usuarios_nombre_completo ON usuarios 
USING gin(to_tsvector('spanish', nombre || ' ' || apellido));

-- Índice GIN para búsquedas en provincias de interés
CREATE INDEX idx_usuarios_provincias_gin ON usuarios USING GIN(provincias_interes);

-- Índice GIN para búsquedas en tipos de propiedad de interés
CREATE INDEX idx_usuarios_tipos_propiedad_gin ON usuarios USING GIN(tipos_propiedad_interes);

-- Función para validar cédula ecuatoriana (algoritmo módulo 10)
CREATE OR REPLACE FUNCTION validar_cedula_ecuatoriana(cedula_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    cedula_clean TEXT;
    digitos INTEGER[];
    suma INTEGER := 0;
    digito_verificador INTEGER;
    i INTEGER;
    multiplicador INTEGER;
BEGIN
    -- Limpiar cédula (solo números)
    cedula_clean := regexp_replace(cedula_input, '[^0-9]', '', 'g');
    
    -- Verificar longitud
    IF LENGTH(cedula_clean) != 10 THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que los dos primeros dígitos sean válidos (01-24)
    IF CAST(LEFT(cedula_clean, 2) AS INTEGER) < 1 OR CAST(LEFT(cedula_clean, 2) AS INTEGER) > 24 THEN
        RETURN FALSE;
    END IF;
    
    -- Verificar que el tercer dígito sea menor a 6 (personas naturales)
    IF CAST(SUBSTRING(cedula_clean, 3, 1) AS INTEGER) >= 6 THEN
        RETURN FALSE;
    END IF;
    
    -- Convertir a array de dígitos
    FOR i IN 1..10 LOOP
        digitos[i] := CAST(SUBSTRING(cedula_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- Aplicar algoritmo módulo 10
    FOR i IN 1..9 LOOP
        IF i % 2 = 1 THEN
            multiplicador := digitos[i] * 2;
            IF multiplicador > 9 THEN
                multiplicador := multiplicador - 9;
            END IF;
        ELSE
            multiplicador := digitos[i];
        END IF;
        suma := suma + multiplicador;
    END LOOP;
    
    digito_verificador := 10 - (suma % 10);
    IF digito_verificador = 10 THEN
        digito_verificador := 0;
    END IF;
    
    RETURN digito_verificador = digitos[10];
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función de validación cédula
ALTER TABLE usuarios 
ADD CONSTRAINT chk_cedula_valida 
CHECK (validar_cedula_ecuatoriana(cedula));

-- Trigger para actualizar fecha_actualizacion
CREATE TRIGGER trigger_actualizar_fecha_usuarios
    BEFORE UPDATE ON usuarios
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_fecha_modificacion();

-- Función para obtener usuarios por tipo
CREATE OR REPLACE FUNCTION obtener_usuarios_por_tipo(tipo_buscado VARCHAR(20))
RETURNS TABLE(
    id UUID,
    nombre VARCHAR(100),
    apellido VARCHAR(100),
    email VARCHAR(255),
    telefono VARCHAR(20),
    fecha_creacion TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.nombre,
        u.apellido,
        u.email,
        u.telefono,
        u.fecha_creacion
    FROM usuarios u
    WHERE u.tipo_usuario = tipo_buscado
      AND u.activo = TRUE
    ORDER BY u.fecha_creacion DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener agentes de una inmobiliaria
CREATE OR REPLACE FUNCTION obtener_agentes_por_inmobiliaria(inmobiliaria_uuid UUID)
RETURNS TABLE(
    id UUID,
    nombre_completo TEXT,
    email VARCHAR(255),
    telefono VARCHAR(20),
    bio TEXT,
    avatar_url TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.nombre || ' ' || u.apellido as nombre_completo,
        u.email,
        u.telefono,
        u.bio,
        u.avatar_url
    FROM usuarios u
    WHERE u.inmobiliaria_id = inmobiliaria_uuid
      AND u.tipo_usuario = 'agente'
      AND u.activo = TRUE
    ORDER BY u.nombre, u.apellido;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar compradores por presupuesto
CREATE OR REPLACE FUNCTION buscar_compradores_por_presupuesto(
    precio_propiedad DECIMAL(15,2)
)
RETURNS TABLE(
    id UUID,
    nombre VARCHAR(100),
    apellido VARCHAR(100),
    email VARCHAR(255),
    telefono VARCHAR(20),
    presupuesto_min DECIMAL(15,2),
    presupuesto_max DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.nombre,
        u.apellido,
        u.email,
        u.telefono,
        u.presupuesto_min,
        u.presupuesto_max
    FROM usuarios u
    WHERE u.tipo_usuario = 'comprador'
      AND u.activo = TRUE
      AND u.presupuesto_min IS NOT NULL
      AND u.presupuesto_max IS NOT NULL
      AND precio_propiedad >= u.presupuesto_min
      AND precio_propiedad <= u.presupuesto_max
    ORDER BY u.presupuesto_max DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar usuarios por nombre
CREATE OR REPLACE FUNCTION buscar_usuarios_por_nombre(nombre_busqueda TEXT)
RETURNS TABLE(
    id UUID,
    nombre VARCHAR(100),
    apellido VARCHAR(100),
    email VARCHAR(255),
    tipo_usuario VARCHAR(20),
    activo BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.nombre,
        u.apellido,
        u.email,
        u.tipo_usuario,
        u.activo
    FROM usuarios u
    WHERE to_tsvector('spanish', u.nombre || ' ' || u.apellido) @@ plainto_tsquery('spanish', nombre_busqueda)
    ORDER BY ts_rank(to_tsvector('spanish', u.nombre || ' ' || u.apellido), plainto_tsquery('spanish', nombre_busqueda)) DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de usuarios
CREATE OR REPLACE FUNCTION obtener_estadisticas_usuarios()
RETURNS TABLE(
    total_usuarios BIGINT,
    usuarios_activos BIGINT,
    compradores BIGINT,
    vendedores BIGINT,
    agentes BIGINT,
    administradores BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_usuarios,
        COUNT(*) FILTER (WHERE activo = TRUE) as usuarios_activos,
        COUNT(*) FILTER (WHERE tipo_usuario = 'comprador') as compradores,
        COUNT(*) FILTER (WHERE tipo_usuario = 'vendedor') as vendedores,
        COUNT(*) FILTER (WHERE tipo_usuario = 'agente') as agentes,
        COUNT(*) FILTER (WHERE tipo_usuario = 'admin') as administradores
    FROM usuarios;
END;
$$ LANGUAGE plpgsql;

-- Vista para usuarios con información completa
CREATE OR REPLACE VIEW usuarios_con_inmobiliaria AS
SELECT 
    u.*,
    i.nombre as inmobiliaria_nombre,
    i.telefono as inmobiliaria_telefono,
    i.email as inmobiliaria_email
FROM usuarios u
LEFT JOIN inmobiliarias i ON u.inmobiliaria_id = i.id
WHERE u.activo = TRUE;

-- Función para validar consistencia de datos de usuario
CREATE OR REPLACE FUNCTION validar_datos_usuario() 
RETURNS TRIGGER AS $$
BEGIN
    -- Verificar que agentes tengan inmobiliaria activa
    IF NEW.tipo_usuario = 'agente' AND NEW.inmobiliaria_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM inmobiliarias 
            WHERE id = NEW.inmobiliaria_id AND activa = TRUE
        ) THEN
            RAISE EXCEPTION 'Los agentes deben estar asociados a una inmobiliaria activa';
        END IF;
    END IF;
    
    -- Verificar que compradores tengan presupuesto válido
    IF NEW.tipo_usuario = 'comprador' THEN
        IF NEW.presupuesto_min IS NOT NULL AND NEW.presupuesto_max IS NOT NULL THEN
            IF NEW.presupuesto_min > NEW.presupuesto_max THEN
                RAISE EXCEPTION 'El presupuesto mínimo no puede ser mayor al máximo';
            END IF;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar datos de usuario
CREATE TRIGGER trigger_validar_usuario
    BEFORE INSERT OR UPDATE ON usuarios
    FOR EACH ROW
    EXECUTE FUNCTION validar_datos_usuario();

-- Comentarios en las funciones
COMMENT ON FUNCTION validar_cedula_ecuatoriana IS 'Valida cédula ecuatoriana usando algoritmo módulo 10';
COMMENT ON FUNCTION obtener_usuarios_por_tipo IS 'Retorna usuarios filtrados por tipo';
COMMENT ON FUNCTION obtener_agentes_por_inmobiliaria IS 'Retorna agentes asociados a una inmobiliaria';
COMMENT ON FUNCTION buscar_compradores_por_presupuesto IS 'Busca compradores que puedan pagar cierto precio';
COMMENT ON FUNCTION buscar_usuarios_por_nombre IS 'Busca usuarios por nombre usando texto completo';
COMMENT ON FUNCTION obtener_estadisticas_usuarios IS 'Retorna estadísticas generales de usuarios';
COMMENT ON FUNCTION validar_datos_usuario IS 'Valida consistencia de datos de usuario';
COMMENT ON VIEW usuarios_con_inmobiliaria IS 'Vista de usuarios con información de inmobiliaria asociada';

-- Insertar algunos usuarios de ejemplo para testing
INSERT INTO usuarios (
    nombre, apellido, email, telefono, cedula, tipo_usuario, 
    presupuesto_min, presupuesto_max, provincias_interes, tipos_propiedad_interes
) VALUES 
(
    'Juan Carlos', 'Mendez', 'juan.mendez@email.com', '0987654321', '1234567890',
    'comprador', 80000, 150000, '["Guayas", "Pichincha"]', '["casa", "departamento"]'
),
(
    'María Elena', 'García', 'maria.garcia@email.com', '0987654322', '1234567891',
    'vendedor', NULL, NULL, '["Guayas"]', '["casa", "terreno"]'
),
(
    'Carlos Alberto', 'Rodríguez', 'carlos.rodriguez@email.com', '0987654323', '1234567892',
    'agente', NULL, NULL, '["Guayas", "Azuay"]', '["casa", "departamento", "comercial"]'
);

-- Actualizar el agente con la inmobiliaria de ejemplo
UPDATE usuarios 
SET inmobiliaria_id = (
    SELECT id FROM inmobiliarias 
    WHERE nombre = 'Inmobiliaria Ejemplo S.A.' 
    LIMIT 1
)
WHERE tipo_usuario = 'agente' 
  AND nombre = 'Carlos Alberto' 
  AND apellido = 'Rodríguez';