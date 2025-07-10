-- Migración 019: Crear tabla de agencies mejorada para sistema de roles
-- Fecha: 2025-01-10
-- Propósito: Implementar tabla de agencies con funcionalidades avanzadas

-- Crear tabla agencies (inmobiliarias)
CREATE TABLE agencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    ruc VARCHAR(13) NOT NULL UNIQUE,
    address TEXT NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    website TEXT DEFAULT '' NOT NULL,
    description TEXT DEFAULT '' NOT NULL,
    logo_url TEXT DEFAULT '' NOT NULL,
    active BOOLEAN DEFAULT TRUE NOT NULL,
    license_number VARCHAR(100) NOT NULL,
    license_expiry TIMESTAMP NULL,
    commission DECIMAL(4,2) DEFAULT 3.00 NOT NULL,
    business_hours TEXT DEFAULT 'Lunes a Viernes 9:00-18:00' NOT NULL,
    social_media JSONB DEFAULT '{}' NOT NULL,
    specialties JSONB DEFAULT '[]' NOT NULL,
    service_areas JSONB DEFAULT '[]' NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Agregar comentarios para documentación
COMMENT ON TABLE agencies IS 'Tabla de inmobiliarias con funcionalidades avanzadas';
COMMENT ON COLUMN agencies.id IS 'Identificador único de la inmobiliaria';
COMMENT ON COLUMN agencies.name IS 'Nombre comercial de la inmobiliaria';
COMMENT ON COLUMN agencies.ruc IS 'RUC (Registro Único de Contribuyentes) de la empresa';
COMMENT ON COLUMN agencies.address IS 'Dirección física de la oficina principal';
COMMENT ON COLUMN agencies.phone IS 'Número de teléfono principal';
COMMENT ON COLUMN agencies.email IS 'Email de contacto principal';
COMMENT ON COLUMN agencies.website IS 'URL del sitio web de la inmobiliaria';
COMMENT ON COLUMN agencies.description IS 'Descripción de la empresa y servicios';
COMMENT ON COLUMN agencies.logo_url IS 'URL del logo de la inmobiliaria';
COMMENT ON COLUMN agencies.active IS 'Indica si la inmobiliaria está activa en el sistema';
COMMENT ON COLUMN agencies.license_number IS 'Número de licencia inmobiliaria';
COMMENT ON COLUMN agencies.license_expiry IS 'Fecha de expiración de la licencia';
COMMENT ON COLUMN agencies.commission IS 'Porcentaje de comisión por defecto (0-10%)';
COMMENT ON COLUMN agencies.business_hours IS 'Horario de atención';
COMMENT ON COLUMN agencies.social_media IS 'JSON con enlaces a redes sociales';
COMMENT ON COLUMN agencies.specialties IS 'JSON con especialidades de la inmobiliaria';
COMMENT ON COLUMN agencies.service_areas IS 'JSON con áreas de servicio (provincias)';
COMMENT ON COLUMN agencies.created_at IS 'Fecha de registro en el sistema';
COMMENT ON COLUMN agencies.updated_at IS 'Fecha de última actualización';

-- Agregar constraints para validar datos
ALTER TABLE agencies ADD CONSTRAINT chk_ruc_format CHECK (ruc ~ '^[0-9]{13}$');
ALTER TABLE agencies ADD CONSTRAINT chk_ruc_company CHECK (ruc ~ '[0-9]{10}001$');
ALTER TABLE agencies ADD CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');
ALTER TABLE agencies ADD CONSTRAINT chk_phone_format CHECK (phone ~ '^(\+593|593|0)(2|3|4|5|6|7|9)[0-9]{7,8}$');
ALTER TABLE agencies ADD CONSTRAINT chk_website_url CHECK (website = '' OR website ~* '^https?://.*');
ALTER TABLE agencies ADD CONSTRAINT chk_logo_url CHECK (logo_url = '' OR logo_url ~* '^https?://.*\.(jpg|jpeg|png|webp|svg)(\?.*)?$');
ALTER TABLE agencies ADD CONSTRAINT chk_commission_range CHECK (commission >= 0 AND commission <= 10);
ALTER TABLE agencies ADD CONSTRAINT chk_license_expiry CHECK (license_expiry IS NULL OR license_expiry > CURRENT_TIMESTAMP);

-- Crear índices para búsquedas
CREATE INDEX idx_agencies_name ON agencies(name);
CREATE INDEX idx_agencies_ruc ON agencies(ruc);
CREATE INDEX idx_agencies_active ON agencies(active) WHERE active = TRUE;
CREATE INDEX idx_agencies_license_expiry ON agencies(license_expiry);
CREATE INDEX idx_agencies_commission ON agencies(commission);
CREATE INDEX idx_agencies_created_at ON agencies(created_at);

-- Índice para búsquedas de texto en nombre y descripción
CREATE INDEX idx_agencies_text_search ON agencies USING gin(to_tsvector('spanish', name || ' ' || description));

-- Índice GIN para búsquedas en especialidades
CREATE INDEX idx_agencies_specialties ON agencies USING GIN(specialties);

-- Índice GIN para búsquedas en áreas de servicio
CREATE INDEX idx_agencies_service_areas ON agencies USING GIN(service_areas);

-- Función para validar RUC ecuatoriano (algoritmo módulo 11)
CREATE OR REPLACE FUNCTION validate_ecuadorian_ruc(ruc_input TEXT) 
RETURNS BOOLEAN AS $$
DECLARE
    ruc_clean TEXT;
    digits INTEGER[];
    sum INTEGER := 0;
    verification_digit INTEGER;
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
        digits[i] := CAST(SUBSTRING(ruc_clean, i, 1) AS INTEGER);
    END LOOP;
    
    -- Validar que el tercer dígito sea menor a 6 (personas jurídicas)
    IF digits[3] >= 6 THEN
        RETURN FALSE;
    END IF;
    
    -- Aplicar algoritmo módulo 11
    FOR i IN 1..9 LOOP
        sum := sum + (digits[i] * (10 - i));
    END LOOP;
    
    verification_digit := 11 - (sum % 11);
    
    IF verification_digit = 11 THEN
        verification_digit := 0;
    ELSIF verification_digit = 10 THEN
        verification_digit := 1;
    END IF;
    
    RETURN verification_digit = digits[10];
END;
$$ LANGUAGE plpgsql;

-- Agregar constraint usando la función de validación RUC
ALTER TABLE agencies ADD CONSTRAINT chk_ruc_valid CHECK (validate_ecuadorian_ruc(ruc));

-- Trigger para actualizar updated_at
CREATE TRIGGER trigger_update_agencies_updated_at
    BEFORE UPDATE ON agencies
    FOR EACH ROW
    EXECUTE FUNCTION actualizar_fecha_modificacion();

-- Función para obtener agencies activas
CREATE OR REPLACE FUNCTION get_active_agencies()
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    ruc VARCHAR(13),
    phone VARCHAR(20),
    email VARCHAR(255),
    website TEXT,
    commission DECIMAL(4,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id,
        a.name,
        a.ruc,
        a.phone,
        a.email,
        a.website,
        a.commission
    FROM agencies a
    WHERE a.active = TRUE
      AND (a.license_expiry IS NULL OR a.license_expiry > CURRENT_TIMESTAMP)
    ORDER BY a.name;
END;
$$ LANGUAGE plpgsql;

-- Función para buscar agencies por nombre
CREATE OR REPLACE FUNCTION search_agencies_by_name(search_query TEXT)
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    description TEXT,
    phone VARCHAR(20),
    email VARCHAR(255),
    commission DECIMAL(4,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id,
        a.name,
        a.description,
        a.phone,
        a.email,
        a.commission
    FROM agencies a
    WHERE a.active = TRUE
      AND to_tsvector('spanish', a.name || ' ' || a.description) @@ plainto_tsquery('spanish', search_query)
    ORDER BY ts_rank(to_tsvector('spanish', a.name || ' ' || a.description), plainto_tsquery('spanish', search_query)) DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener agencies por área de servicio
CREATE OR REPLACE FUNCTION get_agencies_by_service_area(target_province TEXT)
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    commission DECIMAL(4,2),
    specialties JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id,
        a.name,
        a.phone,
        a.email,
        a.commission,
        a.specialties
    FROM agencies a
    WHERE a.active = TRUE
      AND (a.license_expiry IS NULL OR a.license_expiry > CURRENT_TIMESTAMP)
      AND a.service_areas @> jsonb_build_array(target_province)
    ORDER BY a.name;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener agencies por especialidad
CREATE OR REPLACE FUNCTION get_agencies_by_specialty(target_specialty TEXT)
RETURNS TABLE(
    id UUID,
    name VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(255),
    commission DECIMAL(4,2),
    service_areas JSONB
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id,
        a.name,
        a.phone,
        a.email,
        a.commission,
        a.service_areas
    FROM agencies a
    WHERE a.active = TRUE
      AND (a.license_expiry IS NULL OR a.license_expiry > CURRENT_TIMESTAMP)
      AND a.specialties @> jsonb_build_array(target_specialty)
    ORDER BY a.name;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de agencies
CREATE OR REPLACE FUNCTION get_agency_statistics()
RETURNS TABLE(
    total_agencies BIGINT,
    active_agencies BIGINT,
    licensed_agencies BIGINT,
    expired_licenses BIGINT,
    average_commission DECIMAL(4,2),
    total_agents BIGINT,
    total_properties BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) as total_agencies,
        COUNT(*) FILTER (WHERE active = TRUE) as active_agencies,
        COUNT(*) FILTER (WHERE license_number IS NOT NULL AND license_number != '') as licensed_agencies,
        COUNT(*) FILTER (WHERE license_expiry IS NOT NULL AND license_expiry <= CURRENT_TIMESTAMP) as expired_licenses,
        ROUND(AVG(commission), 2) as average_commission,
        (SELECT COUNT(*) FROM users WHERE role = 'agent' AND agency_id IS NOT NULL) as total_agents,
        (SELECT COUNT(*) FROM properties WHERE agency_id IS NOT NULL) as total_properties
    FROM agencies;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener performance de agency
CREATE OR REPLACE FUNCTION get_agency_performance(target_agency_id UUID)
RETURNS TABLE(
    agency_id UUID,
    agency_name VARCHAR(255),
    total_properties BIGINT,
    sold_properties BIGINT,
    rented_properties BIGINT,
    total_sales_value DECIMAL(15,2),
    total_rent_value DECIMAL(15,2),
    average_property_value DECIMAL(15,2),
    total_agents BIGINT,
    active_agents BIGINT,
    conversion_rate DECIMAL(5,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        a.id as agency_id,
        a.name as agency_name,
        COUNT(p.id) as total_properties,
        COUNT(p.id) FILTER (WHERE p.status = 'sold') as sold_properties,
        COUNT(p.id) FILTER (WHERE p.status = 'rented') as rented_properties,
        COALESCE(SUM(CASE WHEN p.status = 'sold' THEN p.price END), 0) as total_sales_value,
        COALESCE(SUM(CASE WHEN p.status = 'rented' THEN p.rent_price END), 0) as total_rent_value,
        COALESCE(AVG(p.price), 0) as average_property_value,
        (SELECT COUNT(*) FROM users WHERE agency_id = a.id AND role = 'agent') as total_agents,
        (SELECT COUNT(*) FROM users WHERE agency_id = a.id AND role = 'agent' AND active = TRUE) as active_agents,
        CASE 
            WHEN COUNT(p.id) > 0 THEN 
                ROUND((COUNT(p.id) FILTER (WHERE p.status IN ('sold', 'rented')) * 100.0 / COUNT(p.id)), 2)
            ELSE 0
        END as conversion_rate
    FROM agencies a
    LEFT JOIN properties p ON a.id = p.agency_id
    WHERE a.id = target_agency_id
    GROUP BY a.id, a.name;
END;
$$ LANGUAGE plpgsql;

-- Vista para agencies con información completa
CREATE OR REPLACE VIEW agencies_with_stats AS
SELECT 
    a.*,
    (SELECT COUNT(*) FROM users WHERE agency_id = a.id AND role = 'agent' AND active = TRUE) as active_agents,
    (SELECT COUNT(*) FROM properties WHERE agency_id = a.id) as total_properties,
    (SELECT COUNT(*) FROM properties WHERE agency_id = a.id AND status = 'sold') as sold_properties,
    (SELECT COUNT(*) FROM properties WHERE agency_id = a.id AND status = 'rented') as rented_properties,
    CASE 
        WHEN a.license_expiry IS NULL THEN TRUE
        ELSE a.license_expiry > CURRENT_TIMESTAMP
    END as license_valid
FROM agencies a
WHERE a.active = TRUE;

-- Función para validar consistencia de datos de agency
CREATE OR REPLACE FUNCTION validate_agency_data() 
RETURNS TRIGGER AS $$
BEGIN
    -- Verificar que la licencia no esté expirada para agencies activas
    IF NEW.active = TRUE AND NEW.license_expiry IS NOT NULL THEN
        IF NEW.license_expiry <= CURRENT_TIMESTAMP THEN
            RAISE EXCEPTION 'No se puede activar una inmobiliaria con licencia expirada';
        END IF;
    END IF;
    
    -- Verificar que las especialidades sean válidas
    IF NEW.specialties IS NOT NULL THEN
        DECLARE
            specialty TEXT;
            valid_specialties TEXT[] := ARRAY['residencial', 'comercial', 'industrial', 'terrenos', 'lujo', 'alquiler', 'venta', 'inversion', 'desarrollo', 'consultoria'];
        BEGIN
            FOR specialty IN SELECT jsonb_array_elements_text(NEW.specialties) LOOP
                IF NOT (specialty = ANY(valid_specialties)) THEN
                    RAISE EXCEPTION 'Especialidad no válida: %', specialty;
                END IF;
            END LOOP;
        END;
    END IF;
    
    -- Verificar que las áreas de servicio sean provincias válidas
    IF NEW.service_areas IS NOT NULL THEN
        DECLARE
            service_area TEXT;
            valid_provinces TEXT[] := ARRAY['Azuay', 'Bolívar', 'Cañar', 'Carchi', 'Chimborazo', 'Cotopaxi', 'El Oro', 'Esmeraldas', 'Galápagos', 'Guayas', 'Imbabura', 'Loja', 'Los Ríos', 'Manabí', 'Morona Santiago', 'Napo', 'Orellana', 'Pastaza', 'Pichincha', 'Santa Elena', 'Santo Domingo', 'Sucumbíos', 'Tungurahua', 'Zamora Chinchipe'];
        BEGIN
            FOR service_area IN SELECT jsonb_array_elements_text(NEW.service_areas) LOOP
                IF NOT (service_area = ANY(valid_provinces)) THEN
                    RAISE EXCEPTION 'Área de servicio no válida: %', service_area;
                END IF;
            END LOOP;
        END;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar datos de agency
CREATE TRIGGER trigger_validate_agency_data
    BEFORE INSERT OR UPDATE ON agencies
    FOR EACH ROW
    EXECUTE FUNCTION validate_agency_data();

-- Ahora podemos crear la referencia desde users a agencies
ALTER TABLE users ADD CONSTRAINT fk_users_agency 
FOREIGN KEY (agency_id) REFERENCES agencies(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- Comentarios en las funciones
COMMENT ON FUNCTION validate_ecuadorian_ruc IS 'Valida RUC ecuatoriano usando algoritmo módulo 11';
COMMENT ON FUNCTION get_active_agencies IS 'Retorna lista de inmobiliarias activas con licencia válida';
COMMENT ON FUNCTION search_agencies_by_name IS 'Busca inmobiliarias por nombre usando búsqueda de texto completo';
COMMENT ON FUNCTION get_agencies_by_service_area IS 'Retorna inmobiliarias que sirven una provincia específica';
COMMENT ON FUNCTION get_agencies_by_specialty IS 'Retorna inmobiliarias con una especialidad específica';
COMMENT ON FUNCTION get_agency_statistics IS 'Retorna estadísticas generales de inmobiliarias';
COMMENT ON FUNCTION get_agency_performance IS 'Retorna métricas de performance de una inmobiliaria';
COMMENT ON FUNCTION validate_agency_data IS 'Valida consistencia de datos de inmobiliaria';
COMMENT ON VIEW agencies_with_stats IS 'Vista de inmobiliarias con estadísticas agregadas';

-- Insertar agency de ejemplo para testing
INSERT INTO agencies (
    name, 
    ruc, 
    address, 
    phone, 
    email, 
    website, 
    description,
    license_number,
    license_expiry,
    commission,
    specialties,
    service_areas
) VALUES (
    'Inmobiliaria Elite Ecuador S.A.',
    '1792146739001',  -- RUC válido de ejemplo
    'Av. 9 de Octubre 123, Edificio Torre Central, Piso 15, Guayaquil, Ecuador',
    '042234567',
    'info@eliteecuador.com',
    'https://www.eliteecuador.com',
    'Inmobiliaria líder en el Ecuador especializada en propiedades de lujo y comerciales. Más de 15 años de experiencia en el mercado inmobiliario.',
    'LIC-EC-2024-001',
    '2025-12-31 23:59:59',
    2.5,
    '["residencial", "comercial", "lujo", "venta", "alquiler"]',
    '["Guayas", "Pichincha", "Azuay", "Manabí"]'
);