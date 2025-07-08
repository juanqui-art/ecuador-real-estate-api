-- Migración 012: Crear índices adicionales para optimizar performance
-- Fecha: 2025-01-15
-- Propósito: Optimizar consultas complejas y mejorar rendimiento general

-- ========================================
-- ÍNDICES PARA BÚSQUEDAS COMPLEJAS
-- ========================================

-- Índice compuesto para búsquedas geográficas con filtros
CREATE INDEX idx_propiedades_ubicacion_filtros ON propiedades(
    provincia, ciudad, tipo, estado, precio
) WHERE estado = 'disponible';

-- Índice para búsquedas por rango de área con precio
CREATE INDEX idx_propiedades_area_precio ON propiedades(
    area_m2, precio, tipo
) WHERE area_m2 > 0 AND estado = 'disponible';

-- Índice para búsquedas por número de habitaciones con precio
CREATE INDEX idx_propiedades_habitaciones_precio ON propiedades(
    dormitorios, banos, precio, tipo
) WHERE dormitorios > 0 AND estado = 'disponible';

-- Índice para búsquedas por año de construcción con características
CREATE INDEX idx_propiedades_antiguedad_caracteristicas ON propiedades(
    ano_construccion, estado_propiedad, precio
) WHERE ano_construccion IS NOT NULL AND estado = 'disponible';

-- ========================================
-- ÍNDICES PARA ESTADÍSTICAS Y REPORTES
-- ========================================

-- Índice para análisis de precios por zona
CREATE INDEX idx_propiedades_analisis_precios ON propiedades(
    provincia, ciudad, tipo, precio_m2, area_m2
) WHERE precio_m2 IS NOT NULL AND estado = 'disponible';

-- Índice para ranking de propiedades más visitadas
CREATE INDEX idx_propiedades_ranking_visitas ON propiedades(
    visitas_contador DESC, destacada, fecha_creacion DESC
) WHERE estado = 'disponible';

-- Índice para análisis temporal de creación
CREATE INDEX idx_propiedades_analisis_temporal ON propiedades(
    DATE(fecha_creacion), tipo, estado, precio
);

-- ========================================
-- ÍNDICES PARA BÚSQUEDAS AVANZADAS
-- ========================================

-- Índice para búsquedas con múltiples amenidades
CREATE INDEX idx_propiedades_amenidades_multiples ON propiedades(
    garage, piscina, seguridad, aire_acondicionado, jardin
) WHERE estado = 'disponible';

-- Índice para propiedades premium con todas las características
CREATE INDEX idx_propiedades_premium_completo ON propiedades(
    precio DESC, destacada, piscina, seguridad, garage
) WHERE precio > 100000 AND estado = 'disponible';

-- Índice para búsquedas por ubicación geográfica precisa
CREATE INDEX idx_propiedades_geolocalizacion ON propiedades(
    latitud, longitud, precision_ubicacion
) WHERE latitud IS NOT NULL AND longitud IS NOT NULL;

-- ========================================
-- ÍNDICES PARA RELACIONES Y JOINS
-- ========================================

-- Índice para joins eficientes entre propiedades e inmobiliarias
CREATE INDEX idx_propiedades_inmobiliaria_join ON propiedades(
    inmobiliaria_id, estado, fecha_creacion DESC
) WHERE inmobiliaria_id IS NOT NULL;

-- Índice para análisis de performance por inmobiliaria
CREATE INDEX idx_propiedades_inmobiliaria_metricas ON propiedades(
    inmobiliaria_id, visitas_contador, destacada, precio
) WHERE inmobiliaria_id IS NOT NULL AND estado = 'disponible';

-- ========================================
-- ÍNDICES PARA OPTIMIZACIÓN DE SORTING
-- ========================================

-- Índice para ordenamiento por precio con filtros básicos
CREATE INDEX idx_propiedades_precio_ordenado ON propiedades(
    precio ASC, area_m2 DESC, fecha_creacion DESC
) WHERE estado = 'disponible';

-- Índice para ordenamiento por relevancia (destacada + visitas)
CREATE INDEX idx_propiedades_relevancia ON propiedades(
    destacada DESC, visitas_contador DESC, precio ASC
) WHERE estado = 'disponible';

-- Índice para ordenamiento por novedad
CREATE INDEX idx_propiedades_novedad ON propiedades(
    fecha_creacion DESC, destacada DESC, estado
) WHERE estado = 'disponible';

-- ========================================
-- ÍNDICES PARA USUARIOS Y BÚSQUEDAS
-- ========================================

-- Índice para matching de usuarios con propiedades
CREATE INDEX idx_usuarios_matching_propiedades ON usuarios(
    tipo_usuario, activo, presupuesto_min, presupuesto_max
) WHERE tipo_usuario = 'comprador' AND activo = TRUE;

-- Índice para análisis de agentes por inmobiliaria
CREATE INDEX idx_usuarios_agentes_inmobiliaria ON usuarios(
    inmobiliaria_id, tipo_usuario, activo, fecha_creacion
) WHERE tipo_usuario = 'agente' AND activo = TRUE;

-- ========================================
-- ÍNDICES PARA BÚSQUEDAS DE TEXTO
-- ========================================

-- Índice de texto completo para búsquedas generales en propiedades
CREATE INDEX idx_propiedades_texto_completo ON propiedades 
USING gin(to_tsvector('spanish', 
    COALESCE(titulo, '') || ' ' || 
    COALESCE(descripcion, '') || ' ' || 
    COALESCE(provincia, '') || ' ' || 
    COALESCE(ciudad, '') || ' ' || 
    COALESCE(tipo, '')
));

-- Índice de texto completo para búsquedas en inmobiliarias
CREATE INDEX idx_inmobiliarias_texto_completo ON inmobiliarias 
USING gin(to_tsvector('spanish', 
    COALESCE(nombre, '') || ' ' || 
    COALESCE(descripcion, '') || ' ' || 
    COALESCE(direccion, '')
));

-- ========================================
-- ÍNDICES PARCIALES PARA CASOS ESPECÍFICOS
-- ========================================

-- Índice para propiedades nuevas (menos de 30 días)
CREATE INDEX idx_propiedades_nuevas ON propiedades(
    fecha_creacion DESC, precio, tipo
) WHERE fecha_creacion > CURRENT_DATE - INTERVAL '30 days' AND estado = 'disponible';

-- Índice para propiedades caras (más de 200k)
CREATE INDEX idx_propiedades_lujo ON propiedades(
    precio DESC, area_m2 DESC, provincia
) WHERE precio > 200000 AND estado = 'disponible';

-- Índice para propiedades con imágenes
CREATE INDEX idx_propiedades_con_media ON propiedades(
    imagen_principal, fecha_creacion DESC, precio
) WHERE imagen_principal != '' AND estado = 'disponible';

-- ========================================
-- ÍNDICES PARA ANÁLISIS DE PERFORMANCE
-- ========================================

-- Índice para análisis de conversión por provincia
CREATE INDEX idx_propiedades_conversion_provincia ON propiedades(
    provincia, estado, fecha_creacion, fecha_actualizacion
);

-- Índice para análisis de tiempo en mercado
CREATE INDEX idx_propiedades_tiempo_mercado ON propiedades(
    fecha_creacion, fecha_actualizacion, estado, precio
);

-- ========================================
-- FUNCIONES DE ANÁLISIS OPTIMIZADAS
-- ========================================

-- Función optimizada para búsquedas geográficas
CREATE OR REPLACE FUNCTION buscar_propiedades_por_ubicacion_optimizada(
    p_provincia VARCHAR(50),
    p_ciudad VARCHAR(100) DEFAULT NULL,
    p_tipo VARCHAR(20) DEFAULT NULL,
    p_precio_min DECIMAL(15,2) DEFAULT NULL,
    p_precio_max DECIMAL(15,2) DEFAULT NULL,
    p_limite INTEGER DEFAULT 20
)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    area_m2 DECIMAL(10,2),
    dormitorios INTEGER,
    banos DECIMAL(3,1),
    ciudad VARCHAR(100),
    visitas_contador INTEGER,
    destacada BOOLEAN
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.area_m2,
        p.dormitorios,
        p.banos,
        p.ciudad,
        p.visitas_contador,
        p.destacada
    FROM propiedades p
    WHERE p.provincia = p_provincia
      AND p.estado = 'disponible'
      AND (p_ciudad IS NULL OR p.ciudad = p_ciudad)
      AND (p_tipo IS NULL OR p.tipo = p_tipo)
      AND (p_precio_min IS NULL OR p.precio >= p_precio_min)
      AND (p_precio_max IS NULL OR p.precio <= p_precio_max)
    ORDER BY p.destacada DESC, p.visitas_contador DESC, p.precio ASC
    LIMIT p_limite;
END;
$$ LANGUAGE plpgsql;

-- Función optimizada para análisis de mercado
CREATE OR REPLACE FUNCTION analizar_mercado_por_zona_optimizada(
    p_provincia VARCHAR(50),
    p_tipo VARCHAR(20) DEFAULT NULL
)
RETURNS TABLE(
    ciudad VARCHAR(100),
    total_propiedades BIGINT,
    precio_promedio DECIMAL(15,2),
    precio_mediano DECIMAL(15,2),
    precio_min DECIMAL(15,2),
    precio_max DECIMAL(15,2),
    area_promedio DECIMAL(10,2),
    propiedades_destacadas BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.ciudad,
        COUNT(*) as total_propiedades,
        ROUND(AVG(p.precio), 2) as precio_promedio,
        ROUND(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY p.precio), 2) as precio_mediano,
        MIN(p.precio) as precio_min,
        MAX(p.precio) as precio_max,
        ROUND(AVG(p.area_m2), 2) as area_promedio,
        COUNT(*) FILTER (WHERE p.destacada = TRUE) as propiedades_destacadas
    FROM propiedades p
    WHERE p.provincia = p_provincia
      AND p.estado = 'disponible'
      AND (p_tipo IS NULL OR p.tipo = p_tipo)
    GROUP BY p.ciudad
    HAVING COUNT(*) >= 3  -- Solo ciudades con al menos 3 propiedades
    ORDER BY total_propiedades DESC, precio_promedio DESC;
END;
$$ LANGUAGE plpgsql;

-- Función optimizada para recommendations
CREATE OR REPLACE FUNCTION recomendar_propiedades_optimizada(
    p_usuario_id UUID,
    p_limite INTEGER DEFAULT 10
)
RETURNS TABLE(
    id UUID,
    titulo VARCHAR(255),
    precio DECIMAL(15,2),
    provincia VARCHAR(50),
    ciudad VARCHAR(100),
    tipo VARCHAR(20),
    score DECIMAL(5,2)
) AS $$
DECLARE
    usuario_presupuesto_min DECIMAL(15,2);
    usuario_presupuesto_max DECIMAL(15,2);
    usuario_provincias JSONB;
    usuario_tipos JSONB;
BEGIN
    -- Obtener preferencias del usuario
    SELECT u.presupuesto_min, u.presupuesto_max, u.provincias_interes, u.tipos_propiedad_interes
    INTO usuario_presupuesto_min, usuario_presupuesto_max, usuario_provincias, usuario_tipos
    FROM usuarios u
    WHERE u.id = p_usuario_id AND u.activo = TRUE;
    
    -- Si no se encontró el usuario, retornar vacío
    IF NOT FOUND THEN
        RETURN;
    END IF;
    
    RETURN QUERY
    SELECT 
        p.id,
        p.titulo,
        p.precio,
        p.provincia,
        p.ciudad,
        p.tipo,
        (
            -- Score basado en múltiples factores
            CASE WHEN p.destacada THEN 20 ELSE 0 END +
            CASE WHEN p.visitas_contador > 50 THEN 15 ELSE 0 END +
            CASE WHEN p.precio BETWEEN COALESCE(usuario_presupuesto_min, 0) 
                 AND COALESCE(usuario_presupuesto_max, 999999999) THEN 25 ELSE 0 END +
            CASE WHEN usuario_provincias ? p.provincia THEN 20 ELSE 0 END +
            CASE WHEN usuario_tipos ? p.tipo THEN 20 ELSE 0 END
        )::DECIMAL(5,2) as score
    FROM propiedades p
    WHERE p.estado = 'disponible'
    ORDER BY score DESC, p.destacada DESC, p.visitas_contador DESC
    LIMIT p_limite;
END;
$$ LANGUAGE plpgsql;

-- ========================================
-- COMENTARIOS Y DOCUMENTACIÓN
-- ========================================

-- Comentarios en las funciones optimizadas
COMMENT ON FUNCTION buscar_propiedades_por_ubicacion_optimizada IS 'Búsqueda optimizada de propiedades por ubicación con filtros';
COMMENT ON FUNCTION analizar_mercado_por_zona_optimizada IS 'Análisis optimizado del mercado inmobiliario por zona';
COMMENT ON FUNCTION recomendar_propiedades_optimizada IS 'Sistema de recomendaciones optimizado basado en preferencias del usuario';

-- Documentación de los índices
COMMENT ON INDEX idx_propiedades_ubicacion_filtros IS 'Optimiza búsquedas geográficas con filtros múltiples';
COMMENT ON INDEX idx_propiedades_area_precio IS 'Optimiza búsquedas por área y precio';
COMMENT ON INDEX idx_propiedades_habitaciones_precio IS 'Optimiza búsquedas por habitaciones y precio';
COMMENT ON INDEX idx_propiedades_geolocalizacion IS 'Optimiza búsquedas por coordenadas GPS';
COMMENT ON INDEX idx_propiedades_texto_completo IS 'Optimiza búsquedas de texto completo en propiedades';
COMMENT ON INDEX idx_inmobiliarias_texto_completo IS 'Optimiza búsquedas de texto completo en inmobiliarias';

-- ========================================
-- ANÁLISIS DE PERFORMANCE
-- ========================================

-- Crear vista para monitorear usage de índices
CREATE OR REPLACE VIEW analisis_uso_indices AS
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan as usos_indice,
    idx_tup_read as registros_leidos,
    idx_tup_fetch as registros_obtenidos
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
  AND tablename IN ('propiedades', 'inmobiliarias', 'usuarios')
ORDER BY idx_scan DESC;

COMMENT ON VIEW analisis_uso_indices IS 'Vista para monitorear el uso de índices en las tablas principales';

-- Función para obtener estadísticas de performance
CREATE OR REPLACE FUNCTION obtener_estadisticas_performance()
RETURNS TABLE(
    tabla VARCHAR(50),
    total_registros BIGINT,
    size_mb DECIMAL(10,2),
    indices_activos INTEGER,
    scans_secuenciales BIGINT,
    scans_indices BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        t.table_name::VARCHAR(50),
        t.n_tup_ins + t.n_tup_upd + t.n_tup_del as total_registros,
        ROUND((pg_total_relation_size(t.schemaname||'.'||t.tablename) / 1024.0 / 1024.0), 2) as size_mb,
        COUNT(i.indexname)::INTEGER as indices_activos,
        t.seq_scan as scans_secuenciales,
        COALESCE(SUM(i.idx_scan), 0) as scans_indices
    FROM pg_stat_user_tables t
    LEFT JOIN pg_stat_user_indexes i ON t.relname = i.relname
    WHERE t.schemaname = 'public'
      AND t.tablename IN ('propiedades', 'inmobiliarias', 'usuarios')
    GROUP BY t.table_name, t.n_tup_ins, t.n_tup_upd, t.n_tup_del, 
             t.schemaname, t.tablename, t.seq_scan
    ORDER BY total_registros DESC;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION obtener_estadisticas_performance IS 'Retorna estadísticas de performance de las tablas principales';

-- Crear un recordatorio para mantenimiento
COMMENT ON SCHEMA public IS 'Schema principal - Recordatorio: ejecutar ANALYZE después de cargar datos masivos para actualizar estadísticas de índices';