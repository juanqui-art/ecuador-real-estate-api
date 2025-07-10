-- Migración 020: Actualizar tabla properties para sistema de roles
-- Fecha: 2025-01-10
-- Propósito: Agregar relaciones de usuarios a las propiedades

-- Agregar columnas para sistema de roles
ALTER TABLE properties ADD COLUMN owner_id UUID NULL;
ALTER TABLE properties ADD COLUMN agent_id UUID NULL;
ALTER TABLE properties ADD COLUMN agency_id UUID NULL;
ALTER TABLE properties ADD COLUMN created_by UUID NULL;
ALTER TABLE properties ADD COLUMN updated_by UUID NULL;

-- Agregar comentarios para las nuevas columnas
COMMENT ON COLUMN properties.owner_id IS 'ID del propietario de la propiedad';
COMMENT ON COLUMN properties.agent_id IS 'ID del agente asignado a la propiedad';
COMMENT ON COLUMN properties.agency_id IS 'ID de la inmobiliaria responsable de la propiedad';
COMMENT ON COLUMN properties.created_by IS 'ID del usuario que creó la propiedad';
COMMENT ON COLUMN properties.updated_by IS 'ID del usuario que actualizó la propiedad por última vez';

-- Crear referencias de claves foráneas
ALTER TABLE properties ADD CONSTRAINT fk_properties_owner 
FOREIGN KEY (owner_id) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE properties ADD CONSTRAINT fk_properties_agent 
FOREIGN KEY (agent_id) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE properties ADD CONSTRAINT fk_properties_agency 
FOREIGN KEY (agency_id) REFERENCES agencies(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE properties ADD CONSTRAINT fk_properties_created_by 
FOREIGN KEY (created_by) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE properties ADD CONSTRAINT fk_properties_updated_by 
FOREIGN KEY (updated_by) REFERENCES users(id) 
ON DELETE SET NULL ON UPDATE CASCADE;

-- Agregar constraints de reglas de negocio
ALTER TABLE properties ADD CONSTRAINT chk_property_ownership 
CHECK (owner_id IS NOT NULL OR agency_id IS NOT NULL);

ALTER TABLE properties ADD CONSTRAINT chk_agent_has_agency 
CHECK (agent_id IS NULL OR agency_id IS NOT NULL);

-- Crear índices para las nuevas columnas
CREATE INDEX idx_properties_owner_id ON properties(owner_id) WHERE owner_id IS NOT NULL;
CREATE INDEX idx_properties_agent_id ON properties(agent_id) WHERE agent_id IS NOT NULL;
CREATE INDEX idx_properties_agency_id ON properties(agency_id) WHERE agency_id IS NOT NULL;
CREATE INDEX idx_properties_created_by ON properties(created_by) WHERE created_by IS NOT NULL;
CREATE INDEX idx_properties_updated_by ON properties(updated_by) WHERE updated_by IS NOT NULL;

-- Índices compuestos para consultas complejas
CREATE INDEX idx_properties_owner_status ON properties(owner_id, status) WHERE owner_id IS NOT NULL;
CREATE INDEX idx_properties_agency_status ON properties(agency_id, status) WHERE agency_id IS NOT NULL;
CREATE INDEX idx_properties_agent_status ON properties(agent_id, status) WHERE agent_id IS NOT NULL;

-- Función para obtener propiedades por propietario
CREATE OR REPLACE FUNCTION get_properties_by_owner(target_owner_id UUID)
RETURNS TABLE(
    id UUID,
    title VARCHAR(255),
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    status VARCHAR(50),
    bedrooms INTEGER,
    bathrooms DECIMAL(3,1),
    area_m2 DECIMAL(10,2),
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.title,
        p.price,
        p.province,
        p.city,
        p.type,
        p.status,
        p.bedrooms,
        p.bathrooms,
        p.area_m2,
        p.created_at
    FROM properties p
    WHERE p.owner_id = target_owner_id
    ORDER BY p.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener propiedades por agente
CREATE OR REPLACE FUNCTION get_properties_by_agent(target_agent_id UUID)
RETURNS TABLE(
    id UUID,
    title VARCHAR(255),
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    status VARCHAR(50),
    owner_name TEXT,
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.title,
        p.price,
        p.province,
        p.city,
        p.type,
        p.status,
        COALESCE(u.first_name || ' ' || u.last_name, 'Sin propietario') as owner_name,
        p.created_at
    FROM properties p
    LEFT JOIN users u ON p.owner_id = u.id
    WHERE p.agent_id = target_agent_id
    ORDER BY p.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para obtener propiedades por inmobiliaria
CREATE OR REPLACE FUNCTION get_properties_by_agency(target_agency_id UUID)
RETURNS TABLE(
    id UUID,
    title VARCHAR(255),
    price DECIMAL(15,2),
    province VARCHAR(100),
    city VARCHAR(100),
    type VARCHAR(50),
    status VARCHAR(50),
    agent_name TEXT,
    owner_name TEXT,
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.title,
        p.price,
        p.province,
        p.city,
        p.type,
        p.status,
        COALESCE(agent.first_name || ' ' || agent.last_name, 'Sin agente') as agent_name,
        COALESCE(owner.first_name || ' ' || owner.last_name, 'Sin propietario') as owner_name,
        p.created_at
    FROM properties p
    LEFT JOIN users agent ON p.agent_id = agent.id
    LEFT JOIN users owner ON p.owner_id = owner.id
    WHERE p.agency_id = target_agency_id
    ORDER BY p.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Función para verificar permisos de usuario sobre propiedad
CREATE OR REPLACE FUNCTION user_can_manage_property(user_id UUID, property_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    user_role user_role;
    user_agency_id UUID;
    property_owner_id UUID;
    property_agent_id UUID;
    property_agency_id UUID;
BEGIN
    -- Obtener información del usuario
    SELECT role, agency_id INTO user_role, user_agency_id
    FROM users WHERE id = user_id;
    
    -- Obtener información de la propiedad
    SELECT owner_id, agent_id, agency_id INTO property_owner_id, property_agent_id, property_agency_id
    FROM properties WHERE id = property_id;
    
    -- Verificar permisos según el rol
    CASE user_role
        WHEN 'admin' THEN
            RETURN TRUE;
        WHEN 'agency' THEN
            -- Agency puede gestionar propiedades donde es la agency
            RETURN property_agency_id = user_id;
        WHEN 'agent' THEN
            -- Agent puede gestionar propiedades asignadas a él o a su agency
            RETURN property_agent_id = user_id OR 
                   (property_agency_id = user_agency_id AND user_agency_id IS NOT NULL);
        WHEN 'owner' THEN
            -- Owner solo puede gestionar sus propias propiedades
            RETURN property_owner_id = user_id;
        WHEN 'buyer' THEN
            -- Buyers no pueden gestionar propiedades
            RETURN FALSE;
        ELSE
            RETURN FALSE;
    END CASE;
END;
$$ LANGUAGE plpgsql;

-- Función para verificar permisos de visualización
CREATE OR REPLACE FUNCTION user_can_view_property(user_id UUID, property_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    property_status VARCHAR(50);
BEGIN
    -- Obtener el estado de la propiedad
    SELECT status INTO property_status
    FROM properties WHERE id = property_id;
    
    -- Todas las propiedades disponibles son públicas
    IF property_status = 'available' THEN
        RETURN TRUE;
    END IF;
    
    -- Para propiedades no disponibles, usar permisos de gestión
    RETURN user_can_manage_property(user_id, property_id);
END;
$$ LANGUAGE plpgsql;

-- Función para obtener estadísticas de propiedades por usuario
CREATE OR REPLACE FUNCTION get_user_property_statistics(target_user_id UUID)
RETURNS TABLE(
    user_id UUID,
    user_name TEXT,
    user_role user_role,
    owned_properties BIGINT,
    managed_properties BIGINT,
    sold_properties BIGINT,
    rented_properties BIGINT,
    total_sales_value DECIMAL(15,2),
    total_rent_value DECIMAL(15,2),
    average_property_value DECIMAL(15,2)
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        u.id as user_id,
        u.first_name || ' ' || u.last_name as user_name,
        u.role as user_role,
        COUNT(p_owned.id) as owned_properties,
        COUNT(p_managed.id) as managed_properties,
        COUNT(p_sold.id) as sold_properties,
        COUNT(p_rented.id) as rented_properties,
        COALESCE(SUM(CASE WHEN p_sold.status = 'sold' THEN p_sold.price END), 0) as total_sales_value,
        COALESCE(SUM(CASE WHEN p_rented.status = 'rented' THEN p_rented.rent_price END), 0) as total_rent_value,
        COALESCE(AVG(p_owned.price), 0) as average_property_value
    FROM users u
    LEFT JOIN properties p_owned ON u.id = p_owned.owner_id
    LEFT JOIN properties p_managed ON u.id = p_managed.agent_id
    LEFT JOIN properties p_sold ON u.id = p_sold.agent_id AND p_sold.status = 'sold'
    LEFT JOIN properties p_rented ON u.id = p_rented.agent_id AND p_rented.status = 'rented'
    WHERE u.id = target_user_id
    GROUP BY u.id, u.first_name, u.last_name, u.role;
END;
$$ LANGUAGE plpgsql;

-- Vista para propiedades con información completa de usuarios
CREATE OR REPLACE VIEW properties_with_users AS
SELECT 
    p.*,
    owner.first_name || ' ' || owner.last_name as owner_name,
    owner.email as owner_email,
    owner.phone as owner_phone,
    agent.first_name || ' ' || agent.last_name as agent_name,
    agent.email as agent_email,
    agent.phone as agent_phone,
    agency.name as agency_name,
    agency.email as agency_email,
    agency.phone as agency_phone,
    creator.first_name || ' ' || creator.last_name as created_by_name,
    updater.first_name || ' ' || updater.last_name as updated_by_name
FROM properties p
LEFT JOIN users owner ON p.owner_id = owner.id
LEFT JOIN users agent ON p.agent_id = agent.id
LEFT JOIN agencies agency ON p.agency_id = agency.id
LEFT JOIN users creator ON p.created_by = creator.id
LEFT JOIN users updater ON p.updated_by = updater.id;

-- Función para asignar propiedad a agente
CREATE OR REPLACE FUNCTION assign_property_to_agent(
    property_id UUID,
    agent_id UUID,
    assigned_by UUID
) RETURNS BOOLEAN AS $$
DECLARE
    agent_agency_id UUID;
BEGIN
    -- Verificar que el agente existe y está activo
    SELECT agency_id INTO agent_agency_id
    FROM users 
    WHERE id = agent_id AND role = 'agent' AND active = TRUE;
    
    IF agent_agency_id IS NULL THEN
        RAISE EXCEPTION 'Agent not found or inactive';
    END IF;
    
    -- Actualizar la propiedad
    UPDATE properties 
    SET 
        agent_id = agent_id,
        agency_id = agent_agency_id,
        updated_by = assigned_by,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = property_id;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Función para transferir propiedad
CREATE OR REPLACE FUNCTION transfer_property_ownership(
    property_id UUID,
    new_owner_id UUID,
    transferred_by UUID
) RETURNS BOOLEAN AS $$
BEGIN
    -- Verificar que el nuevo propietario existe y está activo
    IF NOT EXISTS (
        SELECT 1 FROM users 
        WHERE id = new_owner_id AND role IN ('owner', 'buyer') AND active = TRUE
    ) THEN
        RAISE EXCEPTION 'New owner not found or invalid role';
    END IF;
    
    -- Actualizar la propiedad
    UPDATE properties 
    SET 
        owner_id = new_owner_id,
        updated_by = transferred_by,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = property_id;
    
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

-- Trigger para validar datos de propiedad con roles
CREATE OR REPLACE FUNCTION validate_property_role_data() 
RETURNS TRIGGER AS $$
BEGIN
    -- Verificar que el propietario tenga rol adecuado
    IF NEW.owner_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM users 
            WHERE id = NEW.owner_id AND role IN ('owner', 'buyer') AND active = TRUE
        ) THEN
            RAISE EXCEPTION 'Owner must have owner or buyer role and be active';
        END IF;
    END IF;
    
    -- Verificar que el agente tenga rol de agente
    IF NEW.agent_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM users 
            WHERE id = NEW.agent_id AND role = 'agent' AND active = TRUE
        ) THEN
            RAISE EXCEPTION 'Agent must have agent role and be active';
        END IF;
    END IF;
    
    -- Verificar que la agency esté activa
    IF NEW.agency_id IS NOT NULL THEN
        IF NOT EXISTS (
            SELECT 1 FROM agencies 
            WHERE id = NEW.agency_id AND active = TRUE
        ) THEN
            RAISE EXCEPTION 'Agency must be active';
        END IF;
    END IF;
    
    -- Verificar que si hay agente, también hay agency
    IF NEW.agent_id IS NOT NULL AND NEW.agency_id IS NULL THEN
        RAISE EXCEPTION 'Property with agent must have agency';
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Crear trigger para validar datos de propiedad
CREATE TRIGGER trigger_validate_property_role_data
    BEFORE INSERT OR UPDATE ON properties
    FOR EACH ROW
    EXECUTE FUNCTION validate_property_role_data();

-- Actualizar trigger existente para incluir updated_by
CREATE OR REPLACE FUNCTION update_property_updated_by() 
RETURNS TRIGGER AS $$
BEGIN
    -- Solo actualizar si no se ha establecido explícitamente
    IF NEW.updated_by IS NULL THEN
        NEW.updated_by := OLD.updated_by;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Comentarios en las funciones
COMMENT ON FUNCTION get_properties_by_owner IS 'Retorna propiedades de un propietario específico';
COMMENT ON FUNCTION get_properties_by_agent IS 'Retorna propiedades asignadas a un agente específico';
COMMENT ON FUNCTION get_properties_by_agency IS 'Retorna propiedades de una inmobiliaria específica';
COMMENT ON FUNCTION user_can_manage_property IS 'Verifica si un usuario puede gestionar una propiedad';
COMMENT ON FUNCTION user_can_view_property IS 'Verifica si un usuario puede ver una propiedad';
COMMENT ON FUNCTION get_user_property_statistics IS 'Retorna estadísticas de propiedades por usuario';
COMMENT ON FUNCTION assign_property_to_agent IS 'Asigna una propiedad a un agente';
COMMENT ON FUNCTION transfer_property_ownership IS 'Transfiere la propiedad a un nuevo propietario';
COMMENT ON FUNCTION validate_property_role_data IS 'Valida datos de propiedad relacionados con roles';
COMMENT ON VIEW properties_with_users IS 'Vista de propiedades con información completa de usuarios y agencies';

-- Migrar datos existentes (si los hay)
-- Marcar propiedades existentes como creadas por sistema
UPDATE properties 
SET created_by = NULL, updated_by = NULL 
WHERE created_by IS NULL AND updated_by IS NULL;

-- Actualizar RealEstateCompanyID existente para usar el nuevo sistema
UPDATE properties 
SET agency_id = real_estate_company_id 
WHERE real_estate_company_id IS NOT NULL AND agency_id IS NULL;

-- Comentar la columna antigua para eventual eliminación
COMMENT ON COLUMN properties.real_estate_company_id IS 'DEPRECATED: Use agency_id instead. Will be removed in future migration.';