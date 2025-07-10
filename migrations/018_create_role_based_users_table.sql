-- Migración 018: Crear tabla de usuarios con sistema de roles
-- Fecha: 2025-01-10
-- Propósito: Implementar sistema completo de roles para la plataforma inmobiliaria

-- Crear enum para roles
CREATE TYPE user_role AS ENUM ('admin', 'agency', 'agent', 'owner', 'buyer');

-- Crear tabla usuarios con sistema de roles
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) NOT NULL,
    cedula VARCHAR(10) NOT NULL UNIQUE,
    date_of_birth DATE NULL,
    role user_role NOT NULL DEFAULT 'buyer',
    active BOOLEAN DEFAULT TRUE NOT NULL,
    min_budget DECIMAL(15,2) NULL,
    max_budget DECIMAL(15,2) NULL,
    interested_provinces JSONB DEFAULT '[]' NOT NULL,
    interested_types JSONB DEFAULT '[]' NOT NULL,
    avatar_url TEXT DEFAULT '' NOT NULL,
    bio TEXT DEFAULT '' NOT NULL,
    agency_id UUID NULL,
    password_hash VARCHAR(255) NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE NOT NULL,
    email_verification_token VARCHAR(255) NULL,
    password_reset_token VARCHAR(255) NULL,
    password_reset_expires TIMESTAMP NULL,
    last_login TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Agregar comentarios para documentación
COMMENT ON TABLE users IS 'Tabla de usuarios del sistema inmobiliario con roles jerárquicos';
COMMENT ON COLUMN users.id IS 'Identificador único del usuario';
COMMENT ON COLUMN users.first_name IS 'Nombre del usuario';
COMMENT ON COLUMN users.last_name IS 'Apellido del usuario';
COMMENT ON COLUMN users.email IS 'Email único del usuario';
COMMENT ON COLUMN users.phone IS 'Número de teléfono del usuario';
COMMENT ON COLUMN users.cedula IS 'Número de cédula ecuatoriana (único)';
COMMENT ON COLUMN users.date_of_birth IS 'Fecha de nacimiento del usuario';
COMMENT ON COLUMN users.role IS 'Rol del usuario: admin(5), agency(4), agent(3), owner(2), buyer(1)';
COMMENT ON COLUMN users.active IS 'Estado del usuario en el sistema';
COMMENT ON COLUMN users.min_budget IS 'Presupuesto mínimo para compradores';
COMMENT ON COLUMN users.max_budget IS 'Presupuesto máximo para compradores';
COMMENT ON COLUMN users.interested_provinces IS 'Array JSON de provincias de interés';
COMMENT ON COLUMN users.interested_types IS 'Array JSON de tipos de propiedad de interés';
COMMENT ON COLUMN users.avatar_url IS 'URL del avatar del usuario';
COMMENT ON COLUMN users.bio IS 'Biografía o descripción del usuario';
COMMENT ON COLUMN users.agency_id IS 'Inmobiliaria asociada (requerido para agentes)';
COMMENT ON COLUMN users.password_hash IS 'Hash de la contraseña del usuario';
COMMENT ON COLUMN users.email_verified IS 'Estado de verificación del email';
COMMENT ON COLUMN users.email_verification_token IS 'Token para verificación de email';
COMMENT ON COLUMN users.password_reset_token IS 'Token para reseteo de contraseña';
COMMENT ON COLUMN users.password_reset_expires IS 'Expiración del token de reseteo';
COMMENT ON COLUMN users.last_login IS 'Fecha del último acceso';
COMMENT ON COLUMN users.created_at IS 'Fecha de registro del usuario';
COMMENT ON COLUMN users.updated_at IS 'Fecha de última actualización';

-- Crear relación con agencies (será creada después)
-- ALTER TABLE users ADD CONSTRAINT fk_users_agency FOREIGN KEY (agency_id) REFERENCES agencies(id) ON DELETE SET NULL ON UPDATE CASCADE;

-- Agregar constraints para validar datos
ALTER TABLE users ADD CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
ALTER TABLE users ADD CONSTRAINT chk_phone_format CHECK (phone ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');
ALTER TABLE users ADD CONSTRAINT chk_cedula_format CHECK (cedula ~ '^[0-9]{10}$');
ALTER TABLE users ADD CONSTRAINT chk_budget_valid CHECK (min_budget IS NULL OR max_budget IS NULL OR min_budget <= max_budget);
ALTER TABLE users ADD CONSTRAINT chk_budget_positive CHECK ((min_budget IS NULL OR min_budget > 0) AND (max_budget IS NULL OR max_budget > 0));
ALTER TABLE users ADD CONSTRAINT chk_date_of_birth_valid CHECK (date_of_birth IS NULL OR (date_of_birth >= '1900-01-01' AND date_of_birth <= CURRENT_DATE - INTERVAL '18 years'));
ALTER TABLE users ADD CONSTRAINT chk_avatar_url CHECK (avatar_url = '' OR avatar_url ~* '^https?://.*\.(jpg|jpeg|png|webp|gif)(\?.*)?$');
ALTER TABLE users ADD CONSTRAINT chk_agent_agency CHECK (role != 'agent' OR agency_id IS NOT NULL);
ALTER TABLE users ADD CONSTRAINT chk_agency_no_association CHECK (role != 'agency' OR agency_id IS NULL);

-- Crear índices para búsquedas
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_cedula ON users(cedula);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active) WHERE active = TRUE;
CREATE INDEX idx_users_agency ON users(agency_id) WHERE agency_id IS NOT NULL;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_role_active ON users(role, active);
CREATE INDEX idx_users_budget ON users(min_budget, max_budget) WHERE min_budget IS NOT NULL AND max_budget IS NOT NULL;
CREATE INDEX idx_users_email_verified ON users(email_verified);
CREATE INDEX idx_users_last_login ON users(last_login);

-- Índice para búsquedas de texto en nombre completo
CREATE INDEX idx_users_full_name ON users USING gin(to_tsvector('spanish', first_name || ' ' || last_name));

-- Índice GIN para búsquedas en provincias de interés
CREATE INDEX idx_users_interested_provinces ON users USING GIN(interested_provinces);

-- Índice GIN para búsquedas en tipos de propiedad de interés
CREATE INDEX idx_users_interested_types ON users USING GIN(interested_types);

-- Función para validar cédula ecuatoriana (algoritmo módulo 10)
CREATE OR REPLACE FUNCTION validate_ecuadorian_cedula(cedula_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    cedula_clean TEXT;
    digits INTEGER[];
    sum INTEGER := 0;
    verification_digit INTEGER;
    i INTEGER;
    multiplier INTEGER;
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
        digits[i] := CAST(SUBSTRING(cedula_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- Aplicar algoritmo módulo 10
    FOR i IN 1..9 LOOP
        IF i % 2 = 1 THEN
            multiplier := digits[i] * 2;
            IF multiplier > 9 THEN
                multiplier := multiplier - 9;
            END IF;
        ELSE
            multiplier := digits[i];
        END IF;
        sum := sum + multiplier;
    END LOOP;
    
    verification_digit := 10 - (sum % 10);
    IF verification_digit = 10 THEN
        verification_digit := 0;
    END IF;
    
    RETURN verification_digit = digits[10];
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función de validación cédula
ALTER TABLE users ADD CONSTRAINT chk_cedula_valid CHECK (validate_ecuadorian_cedula(cedula));

-- Trigger para actualizar updated_at
CREATE TRIGGER trigger_update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_fecha_modificacion();

-- Función para obtener usuarios por rol
CREATE OR REPLACE FUNCTION get_users_by_role(target_role user_role)
RETURNS TABLE(
    id UUID,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.first_name,
        u.last_name,
        u.email,
        u.phone,
        u.created_at
    FROM users u
    WHERE u.role = target_role
      AND u.active = TRUE
    ORDER BY u.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener agentes de una agency
CREATE OR REPLACE FUNCTION get_agents_by_agency(agency_uuid UUID)
RETURNS TABLE(
    id UUID,
    full_name TEXT,
    email VARCHAR(255),
    phone VARCHAR(20),
    bio TEXT,
    avatar_url TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.first_name || ' ' || u.last_name as full_name,
        u.email,
        u.phone,
        u.bio,
        u.avatar_url
    FROM users u
    WHERE u.agency_id = agency_uuid
      AND u.role = 'agent'
      AND u.active = TRUE
    ORDER BY u.first_name, u.last_name;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar compradores por presupuesto
CREATE OR REPLACE FUNCTION find_buyers_by_budget(property_price DECIMAL(15,2))
RETURNS TABLE(
    id UUID,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(20),
    min_budget DECIMAL(15,2),
    max_budget DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.first_name,
        u.last_name,
        u.email,
        u.phone,
        u.min_budget,
        u.max_budget
    FROM users u
    WHERE u.role = 'buyer'
      AND u.active = TRUE
      AND u.min_budget IS NOT NULL
      AND u.max_budget IS NOT NULL
      AND property_price >= u.min_budget
      AND property_price <= u.max_budget
    ORDER BY u.max_budget DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar usuarios por nombre
CREATE OR REPLACE FUNCTION search_users_by_name(search_query TEXT)
RETURNS TABLE(
    id UUID,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    role user_role,
    active BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id,
        u.first_name,
        u.last_name,
        u.email,
        u.role,
        u.active
    FROM users u
    WHERE to_tsvector('spanish', u.first_name || ' ' || u.last_name) @@ plainto_tsquery('spanish', search_query)
    ORDER BY ts_rank(to_tsvector('spanish', u.first_name || ' ' || u.last_name), plainto_tsquery('spanish', search_query)) DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de usuarios
CREATE OR REPLACE FUNCTION get_user_statistics()
RETURNS TABLE(
    total_users BIGINT,
    active_users BIGINT,
    admin_count BIGINT,
    agency_count BIGINT,
    agent_count BIGINT,
    owner_count BIGINT,
    buyer_count BIGINT,
    email_verified_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_users,
        COUNT(*) FILTER (WHERE active = TRUE) as active_users,
        COUNT(*) FILTER (WHERE role = 'admin') as admin_count,
        COUNT(*) FILTER (WHERE role = 'agency') as agency_count,
        COUNT(*) FILTER (WHERE role = 'agent') as agent_count,
        COUNT(*) FILTER (WHERE role = 'owner') as owner_count,
        COUNT(*) FILTER (WHERE role = 'buyer') as buyer_count,
        COUNT(*) FILTER (WHERE email_verified = TRUE) as email_verified_count
    FROM users;
END;
$$ LANGUAGE plpgsql;

-- Vista para usuarios con información completa
CREATE OR REPLACE VIEW users_with_agency AS
SELECT 
    u.*,
    a.name as agency_name,
    a.phone as agency_phone,
    a.email as agency_email
FROM users u
LEFT JOIN agencies a ON u.agency_id = a.id
WHERE u.active = TRUE;

-- Función para validar consistencia de datos de usuario
CREATE OR REPLACE FUNCTION validate_user_data() 
RETURNS TRIGGER AS $$
BEGIN
    -- Verificar que agentes tengan agency activa
    IF NEW.role = 'agent' AND NEW.agency_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM agencies 
            WHERE id = NEW.agency_id AND active = TRUE
        ) THEN
            RAISE EXCEPTION 'Los agentes deben estar asociados a una inmobiliaria activa';
        END IF;
    END IF;
    
    -- Verificar que compradores tengan presupuesto válido
    IF NEW.role = 'buyer' THEN
        IF NEW.min_budget IS NOT NULL AND NEW.max_budget IS NOT NULL THEN
            IF NEW.min_budget > NEW.max_budget THEN
                RAISE EXCEPTION 'El presupuesto mínimo no puede ser mayor al máximo';
            END IF;
        END IF;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar datos de usuario
CREATE TRIGGER trigger_validate_user_data
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION validate_user_data();

-- Comentarios en las funciones
COMMENT ON FUNCTION validate_ecuadorian_cedula IS 'Valida cédula ecuatoriana usando algoritmo módulo 10';
COMMENT ON FUNCTION get_users_by_role IS 'Retorna usuarios filtrados por rol';
COMMENT ON FUNCTION get_agents_by_agency IS 'Retorna agentes asociados a una inmobiliaria';
COMMENT ON FUNCTION find_buyers_by_budget IS 'Busca compradores que puedan pagar cierto precio';
COMMENT ON FUNCTION search_users_by_name IS 'Busca usuarios por nombre usando texto completo';
COMMENT ON FUNCTION get_user_statistics IS 'Retorna estadísticas generales de usuarios';
COMMENT ON FUNCTION validate_user_data IS 'Valida consistencia de datos de usuario';
COMMENT ON VIEW users_with_agency IS 'Vista de usuarios con información de inmobiliaria asociada';