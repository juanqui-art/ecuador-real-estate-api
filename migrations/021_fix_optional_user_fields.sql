-- Migración 021: Corregir campos opcionales de usuarios
-- Fecha: 2025-07-11
-- Propósito: Hacer avatar_url y bio opcionales para mejorar UX de registro

-- ========================================
-- PROBLEMA IDENTIFICADO
-- ========================================
-- Los campos avatar_url y bio están definidos como NOT NULL con DEFAULT ''
-- pero el service layer no los inicializa, causando errores al crear usuarios.
-- Estos campos son opcionales por naturaleza y deben permitir NULL.

-- ========================================
-- CORRECCIÓN DE CONSTRAINTS
-- ========================================

-- Hacer avatar_url opcional (permitir NULL)
ALTER TABLE users ALTER COLUMN avatar_url DROP NOT NULL;
ALTER TABLE users ALTER COLUMN avatar_url SET DEFAULT NULL;

-- Hacer bio opcional (permitir NULL)  
ALTER TABLE users ALTER COLUMN bio DROP NOT NULL;
ALTER TABLE users ALTER COLUMN bio SET DEFAULT NULL;

-- ========================================
-- ACTUALIZAR CONSTRAINTS DE VALIDACIÓN
-- ========================================

-- Actualizar constraint de avatar_url para manejar NULL
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_avatar_url;
ALTER TABLE users ADD CONSTRAINT chk_avatar_url CHECK (
    avatar_url IS NULL OR 
    avatar_url = '' OR 
    avatar_url ~* '^https?://.*\.(jpg|jpeg|png|webp|gif)(\?.*)?$'
);

-- ========================================
-- CONVERSIÓN DE DATOS EXISTENTES
-- ========================================

-- Convertir strings vacíos existentes a NULL para consistency
UPDATE users SET avatar_url = NULL WHERE avatar_url = '';
UPDATE users SET bio = NULL WHERE bio = '';

-- ========================================
-- COMENTARIOS ACTUALIZADOS
-- ========================================

COMMENT ON COLUMN users.avatar_url IS 'URL del avatar del usuario (opcional - puede ser NULL)';
COMMENT ON COLUMN users.bio IS 'Biografía o descripción del usuario (opcional - puede ser NULL)';

-- ========================================
-- FUNCIONES ACTUALIZADAS
-- ========================================

-- Actualizar función de agentes para manejar bio NULL
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
        COALESCE(u.bio, 'Sin biografía') as bio,  -- Manejar NULL
        COALESCE(u.avatar_url, '') as avatar_url   -- Manejar NULL
    FROM users u
    WHERE u.agency_id = agency_uuid
      AND u.role = 'agent'
      AND u.active = TRUE
    ORDER BY u.first_name, u.last_name;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- DOCUMENTACIÓN
-- ========================================

COMMENT ON CONSTRAINT chk_avatar_url ON users IS 'Valida formato URL de avatar o permite NULL/string vacío';

-- Registrar cambio en schema
INSERT INTO schema_migrations (version) VALUES ('021') ON CONFLICT DO NOTHING;