# 🏠 Realty Core - Claude Code Commands

Comandos personalizados optimizados para el sistema inmobiliario ecuatoriano.

## 🚀 **COMANDOS DISPONIBLES**

### **📋 COMANDOS PRINCIPALES**

| Comando | Descripción | Uso Principal |
|---------|-------------|---------------|
| `realty-property` | Gestión de propiedades | CRUD, validaciones, business logic |
| `realty-api` | Endpoints REST | Crear/optimizar APIs |
| `realty-cache` | Optimización de cache | Performance, memoria, hit rates |
| `realty-test` | Testing específico | Unit, integration, benchmarks |

### **🛠️ COMANDOS ESPECIALIZADOS**

| Comando | Descripción | Uso Principal |
|---------|-------------|---------------|
| `realty-image` | Sistema de imágenes | Upload, procesamiento, storage |
| `realty-fts` | Full-text search | Búsqueda, ranking, autocomplete |
| `realty-ecuador` | Validaciones locales | Provincias, ciudades, precios |
| `realty-debug` | Debug y performance | Profiling, optimización |
| `realty-deploy` | Deployment | Docker, producción, monitoring |

## 💡 **EJEMPLOS DE USO**

### **Desarrollo de Features**
```bash
# Agregar nueva funcionalidad
claude > /project:realty-property "add parking_spaces field to Property struct"

# Crear API endpoint
claude > /project:realty-api "create property favorites endpoint"

# Optimizar cache
claude > /project:realty-cache "optimize cache for property search results"
```

### **Testing y Quality**
```bash
# Crear tests comprehensivos
claude > /project:realty-test "property service with mock repository"

# Debug performance
claude > /project:realty-debug "slow property search queries"
```

### **Funcionalidades Específicas**
```bash
# Procesamiento de imágenes
claude > /project:realty-image "add watermark to property images"

# Búsqueda avanzada
claude > /project:realty-fts "improve search ranking for Ecuador locations"

# Validaciones locales
claude > /project:realty-ecuador "validate Quito postal codes"
```

## 🎯 **WORKFLOWS COMBINADOS**

### **"Complete Feature Development"**
```bash
claude > /project:realty-property "User favorites system" && /project:realty-api "favorites endpoints" && /project:realty-test "favorites functionality"
```

### **"Performance Optimization"**
```bash
claude > /project:realty-debug "identify bottlenecks" && /project:realty-cache "optimize caching strategy" && /project:realty-fts "improve search performance"
```

### **"Production Readiness"**
```bash
claude > /project:realty-deploy "Docker production setup" && /project:realty-test "integration tests" && /project:realty-debug "performance profiling"
```

## 📊 **BENEFICIOS**

- **80% menos tokens** en operaciones comunes
- **50% menos tiempo** en development
- **Consistencia** en patterns y validaciones
- **Reutilización** de soluciones optimizadas
- **Contexto específico** para real estate

## 🔧 **INTEGRACIÓN CON MAKEFILE**

Los comandos se integran perfectamente con nuestro Makefile:

```bash
# Después de usar comandos, ejecutar:
make test-cache      # Probar cache
make test-properties # Probar propiedades
make ci             # Pipeline completo
```

## 🏠 **CONTEXTO DEL PROYECTO**

- **Versión:** v1.5.0-endpoint-expansion
- **179 tests** con 90%+ cobertura
- **9 endpoints funcionales + 48 pendientes** integración
- **Sistema completo** de imágenes con cache LRU
- **PostgreSQL FTS** en español
- **Validaciones Ecuador** integradas

---

**¡Usa estos comandos para desarrollar 10x más rápido!** 🚀