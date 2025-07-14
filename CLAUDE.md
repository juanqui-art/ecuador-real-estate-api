# CLAUDE.md

Este archivo proporciona orientación a Claude Code (claude.ai/code) cuando trabaja con código en este repositorio.

## Resumen del Proyecto

Sistema de gestión de propiedades inmobiliarias en Go 1.24 para el mercado ecuatoriano. Proyecto de aprendizaje enfocado en desarrollo incremental y best practices de Go.

**Tecnologías:**
- Backend: Go 1.24 con net/http nativo
- Base de datos: PostgreSQL 
- ORM: A definir (GORM vs database/sql)
- Desarrollo local con Docker
- Testing: testify

**Objetivos:**
- CRUD completo de propiedades inmobiliarias
- Validaciones específicas para Ecuador
- Arquitectura limpia y extensible
- Aprendizaje gradual de patrones Go

## Comandos Comunes

### Desarrollo Local
```bash
# Ejecutar servidor de desarrollo
go run cmd/server/main.go

# Construir el proyecto
go build -o bin/inmobiliaria ./cmd/server

# Ejecutar tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Formatear código
go fmt ./...

# Verificar código
go vet ./...
```

### Base de Datos (GoLand + Docker Compose)
```bash
# Iniciar servicios desde GoLand
# View → Tool Windows → Services → docker-compose → postgres → Start

# Desde terminal (alternativo)
docker-compose up -d postgres

# Ver logs desde GoLand
# Services → postgres → Logs

# Conectar desde Database Tool Window
# Host: localhost, Port: 5433, DB: inmobiliaria_db, User: juanquizhpi, Pass: (vacío)

# Detener servicios
docker-compose down
```

### Dependencias
```bash
# Añadir dependencia
go get github.com/ejemplo/paquete

# Actualizar dependencias
go mod tidy

# Descargar dependencias
go mod download
```

## Arquitectura del Proyecto

**Estructura de Directorios:**
```
realty-core/
├── cmd/
│   └── server/             # Application entry point
├── internal/
│   ├── domain/             # Business logic and models
│   ├── repository/         # Data access (properties + images)
│   ├── service/            # Application services
│   ├── handlers/           # HTTP handlers + responses
│   ├── cache/              # LRU cache implementation
│   ├── storage/            # Image storage (local filesystem)
│   ├── processors/         # Image processing (resize, compress)
│   └── config/             # Configuration management
├── pkg/                   # Reusable code
├── migrations/            # Database scripts
├── tests/                 # Integration tests
└── docs/                  # Documentation
```

**Patrones Utilizados:**
- Repository Pattern para acceso a datos
- Service Layer para lógica de negocio
- Handler Pattern para HTTP
- Dependency Injection manual

**Estructura Propiedad (simplificada):**
```go
type Property struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Province    string    `json:"province"`
    City        string    `json:"city"`
    Type        string    `json:"type"` // house, apartment, land, commercial
    Status      string    `json:"status"` // available, sold, rented
    Bedrooms    int       `json:"bedrooms"`
    Bathrooms   float32   `json:"bathrooms"`
    AreaM2      float64   `json:"area_m2"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Sistema de Cache LRU:**
```go
// Cache LRU con nodos doblemente enlazados
type LRUCache struct {
    capacity    int           // Máximo número de elementos
    maxSize     int64         // Máximo tamaño en bytes
    currentSize int64         // Tamaño actual
    cache       map[string]*LRUNode
    head        *LRUNode      // Más recientemente usado
    tail        *LRUNode      // Menos recientemente usado
    mutex       sync.RWMutex  // Thread safety
    hits        int64         // Estadísticas
    misses      int64
    evictions   int64
    ttl         time.Duration // Time to live
}

// Cache específico para imágenes
type ImageCache struct {
    lru         *LRUCache
    enabled     bool
    thumbnailHits   int64
    thumbnailMisses int64
    variantHits     int64  
    variantMisses   int64
}
```

## API Endpoints

### CRUD Básico
```
GET    /api/properties         # List properties
POST   /api/properties         # Create property
GET    /api/properties/{id}    # Get property by ID
PUT    /api/properties/{id}    # Update property
DELETE /api/properties/{id}    # Delete property
GET    /api/properties/slug/{slug}  # Get property by SEO slug
```

### Búsqueda y Filtros (PostgreSQL FTS)
```
GET    /api/properties/filter  # Basic filters (query, province, price)
GET    /api/properties/search/ranked  # FTS search with ranking
GET    /api/properties/search/suggestions  # Autocomplete suggestions
POST   /api/properties/search/advanced  # Advanced multi-filter search
```

### Funcionalidades Adicionales
```
GET    /api/properties/statistics  # Property statistics
POST   /api/properties/{id}/location  # Set GPS location
POST   /api/properties/{id}/featured  # Mark as featured
GET    /api/health             # Health check
```

### Gestión de Imágenes (13 endpoints)
```
POST   /api/images                      # Upload imagen
GET    /api/images/{id}                 # Obtener metadata imagen  
GET    /api/properties/{id}/images      # Listar imágenes de propiedad
PUT    /api/images/{id}/metadata        # Actualizar metadata
DELETE /api/images/{id}                 # Eliminar imagen
POST   /api/properties/{id}/images/reorder # Reordenar imágenes
POST   /api/properties/{id}/images/main # Establecer imagen principal
GET    /api/properties/{id}/images/main # Obtener imagen principal
GET    /api/images/{id}/variant        # Obtener variante imagen
GET    /api/images/{id}/thumbnail      # Obtener thumbnail
GET    /api/images/stats               # Estadísticas de imágenes
POST   /api/images/cleanup             # Limpieza archivos temporales
GET    /api/images/cache/stats         # Estadísticas de cache
```

### Sistema de Autenticación JWT (5 endpoints) 🔐
```
POST   /api/auth/login                  # Autenticación con JWT tokens
POST   /api/auth/refresh                # Renovar access token
POST   /api/auth/logout                 # Logout seguro con token blacklisting
GET    /api/auth/validate               # Validar token actual
POST   /api/auth/change-password        # Cambiar contraseña autenticado
```

### Gestión de Usuarios (10 endpoints - PROTEGIDOS)
```
GET    /api/users                       # Búsqueda y listado (requiere auth)
POST   /api/users                       # Crear usuario (admin/agency)
GET    /api/users/{id}                  # Obtener usuario (resource access)
PUT    /api/users/{id}                  # Actualizar usuario (resource access)
DELETE /api/users/{id}                  # Eliminar usuario (resource access)
GET    /api/users/role/{role}           # Obtener usuarios por rol (requiere auth)
GET    /api/users/statistics            # Estadísticas (admin analytics)
GET    /api/users/dashboard             # Dashboard personal (autenticado)
```

### Gestión de Agencias (15 endpoints)
```
GET    /api/agencies                    # Búsqueda y listado de agencias
POST   /api/agencies                    # Crear agencia
GET    /api/agencies/{id}               # Obtener agencia por ID
PUT    /api/agencies/{id}               # Actualizar agencia
DELETE /api/agencies/{id}               # Eliminar agencia
GET    /api/agencies/active             # Obtener agencias activas
GET    /api/agencies/service-area/{area} # Agencias por área de servicio
GET    /api/agencies/specialty/{specialty} # Agencias por especialidad
GET    /api/agencies/{id}/agents        # Obtener agentes de agencia
POST   /api/agencies/{id}/license       # Gestionar licencia de agencia
GET    /api/agencies/statistics         # Estadísticas de agencias
GET    /api/agencies/{id}/performance   # Métricas de desempeño
```

### Sistema de Paginación (7 endpoints)
```
GET    /api/pagination/properties       # Propiedades paginadas
GET    /api/pagination/images           # Imágenes paginadas
GET    /api/pagination/users            # Usuarios paginados
GET    /api/pagination/agencies         # Agencias paginadas
GET    /api/pagination/search           # Búsqueda global paginada
GET    /api/pagination/stats            # Estadísticas de paginación
POST   /api/pagination/advanced         # Paginación avanzada configurable
```

## Configuración de Desarrollo

**Variables de Entorno (.env):**
```env
DATABASE_URL=postgresql://admin:password@localhost:5432/inmobiliaria_db
PORT=8080
LOG_LEVEL=info
```

**Provincias Ecuador:**
Azuay, Bolívar, Cañar, Carchi, Chimborazo, Cotopaxi, El Oro, Esmeraldas, Galápagos, Guayas, Imbabura, Loja, Los Ríos, Manabí, Morona Santiago, Napo, Orellana, Pastaza, Pichincha, Santa Elena, Santo Domingo, Sucumbíos, Tungurahua, Zamora Chinchipe

## Configuración de Desarrollo

### IDE: GoLand 2025.1.3
- **Services Tool Window:** Para gestionar Docker containers
- **Database Tool Window:** Para conexión PostgreSQL integrada
- **Run Configurations:** API configurada con variables de entorno
- **HTTP Client:** Para probar endpoints desde el IDE

### Docker Compose
- **PostgreSQL 15:** Base de datos principal
- **pgAdmin:** Interfaz web opcional (puerto 5050)
- **Volúmenes persistentes:** Datos sobreviven reinicios
- **Auto-migraciones:** Scripts SQL ejecutados automáticamente

## Estado Actual del Proyecto

**Versión:** v1.9.0-sistema-completo  
**Fecha:** 2025-01-11  
**Cobertura Tests:** 90%+ promedio (179 tests)  
**Funcionalidades:** 56+ endpoints funcionales con autenticación JWT completa  
**FASE 1 COMPLETADA:** ✅ Sistema de autenticación y autorización JWT funcional

### Funcionalidades Completadas ✅
- **Arquitectura limpia:** Domain/Service/Repository/Handlers optimizada
- **CRUD completo:** 56+ endpoints API funcionales CON AUTENTICACIÓN
- **🔐 Sistema JWT:** Access tokens (15min) + Refresh tokens (7 días)
- **🛡️ Autorización:** 5 roles con 16 permisos granulares (Admin > Agency > Agent > Owner > Buyer)
- **🔒 Middleware:** Autenticación, validación de roles, control de acceso a recursos
- **🔑 Endpoints Auth:** Login, logout, refresh, validate, change password
- **PostgreSQL FTS:** Búsqueda full-text en español con ranking
- **Sistema de Imágenes:** Upload, procesamiento, storage, cache LRU - 13 endpoints
- **Sistema de Usuarios:** Gestión completa PROTEGIDA - 10 endpoints
- **Sistema de Agencias:** Gestión completa con validación RUC - 15 endpoints
- **Sistema de Paginación:** Paginación avanzada multi-entidad - 7 endpoints
- **Sistema de Propiedades:** CRUD básico PROTEGIDO - 6 endpoints
- **Testing comprehensivo:** 179 tests con 90%+ cobertura
- **Validaciones:** Business rules específicas Ecuador
- **Código limpio:** Refactoring completo, eliminación de archivos backup
- **Compilación exitosa:** Sistema estable y funcional

### Sistemas Integrados 🏗️
1. **🔐 Autenticación (5 endpoints):** JWT, login, logout, refresh, validation
2. **Propiedades (6 endpoints):** CRUD PROTEGIDO, búsqueda pública, estadísticas
3. **Imágenes (13 endpoints):** Upload PROTEGIDO, procesamiento, cache, variantes
4. **Usuarios (10 endpoints):** Gestión PROTEGIDA con control de acceso
5. **Agencias (15 endpoints):** Gestión PROTEGIDA, performance, licencias
6. **Paginación (7 endpoints):** Paginación avanzada, búsqueda global

### FASE 1 - Sistema de Autenticación COMPLETADA 🎉
- ✅ **JWT Manager completo:** Generación, validación, refresh, blacklisting
- ✅ **Role-based Access Control:** 5 roles jerárquicos con 16 permisos
- ✅ **Middleware de seguridad:** Protección automática de endpoints
- ✅ **Resource-specific access:** Control por ownership de recursos
- ✅ **Handlers de autenticación:** Login/logout seguro con validación
- ✅ **Configuración production-ready:** Variables de entorno, secrets seguros

### PRÓXIMA FASE 2 - Dashboard Frontend 📋
- **React/Next.js 14:** Dashboard administrativo con UI/UX de élite
- **shadcn/ui + Tailwind:** Componentes modernos y responsive
- **TanStack Query:** State management y data fetching optimizado
- **Framer Motion:** Animaciones y micro-interacciones fluidas
- **Documentación:** API documentation y guías de uso
- **Deployment:** Preparación para producción
- **Dashboard:** Interfaz de administración

## Notas para el Desarrollo

- **Enfoque incremental:** Comenzar con funcionalidad básica y agregar complejidad gradualmente
- **net/http nativo:** Aprovechar las mejoras de Go 1.22+ sin librerías externas
- **Idioma español:** Todo el código y documentación en español
- **Validaciones específicas:** Enfoque en el mercado inmobiliario ecuatoriano
- **Aprendizaje Go:** Explicaciones de patrones y conceptos conforme se implementan
- **IDE-first:** Uso de GoLand para todo el flujo de desarrollo
- **Testing first:** Toda nueva funcionalidad debe incluir tests
- **Seguimiento:** PROGRESS.md actualizado con cada avance

## Ejemplo de Datos para Testing

```json
{
  "titulo": "Hermosa casa en Samborondón con piscina",
  "descripcion": "Casa moderna de 3 pisos con acabados de lujo",
  "precio": 285000,
  "provincia": "Guayas", 
  "ciudad": "Samborondón",
  "tipo": "casa",
  "dormitorios": 4,
  "banos": 3.5,
  "area_m2": 320
}
```