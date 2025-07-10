# 🏠 Realty Core - Sistema Inmobiliario Ecuador

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Coverage](https://img.shields.io/badge/Coverage-90%25+-brightgreen)](https://github.com/tu-usuario/realty-core)
[![Tests](https://img.shields.io/badge/Tests-200+-green)](https://github.com/tu-usuario/realty-core)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> Sistema completo de gestión inmobiliaria construido específicamente para el mercado ecuatoriano con Go 1.24, PostgreSQL FTS y arquitectura de roles avanzada.

## 🎯 **Características Destacadas**

### 🏗️ **Arquitectura Empresarial**
- **Domain-Driven Design** con capas claramente separadas
- **Sistema de roles jerárquico** (Admin → Agency → Agent → Owner → Buyer)
- **PostgreSQL FTS** optimizado para búsquedas en español
- **Cache LRU inteligente** con gestión de memoria automática

### 🇪🇨 **Optimizado para Ecuador**
- ✅ **Validación RUC** (13 dígitos) para agencias
- ✅ **Teléfonos Ecuador** (+593xxxxxxxxx formato)
- ✅ **24 Provincias** validadas automáticamente
- ✅ **Business rules** específicas del mercado local

### ⚡ **Performance de Producción**
- **26+ API Endpoints** RESTful completamente documentados
- **200+ Tests** con 90%+ de cobertura
- **Sistema de imágenes** con thumbnails y cache automático
- **Paginación inteligente** y búsqueda indexada

## 🚀 **Quick Start**

### Prerrequisitos
- Go 1.24+
- PostgreSQL 15+
- Docker & Docker Compose (opcional)

### Instalación Rápida

```bash
# Clonar repositorio
git clone https://github.com/tu-usuario/realty-core.git
cd realty-core

# Instalar dependencias
go mod download

# Iniciar base de datos con Docker
docker-compose up -d postgres

# Ejecutar migraciones automáticamente (se ejecutan al iniciar)
go run cmd/server/main.go

# El servidor estará disponible en http://localhost:8080
```

### Comandos de Desarrollo
```bash
# Tests completos
make test

# Tests específicos
make test-cache     # Tests de cache LRU
make test-images    # Tests de imágenes
make test-roles     # Tests de sistema de roles

# Desarrollo
make dev           # Servidor con hot reload
make build         # Build de producción
make ci           # Pipeline completo
```

## 📊 **Arquitectura del Sistema**

### Roles y Permisos
```
ADMIN (Nivel 5)
├── Control total del sistema
├── Gestión de todas las agencias
└── Moderación y reportes globales

AGENCY (Nivel 4)
├── Gestión de cartera de propiedades
├── Administración de agentes
└── Reportes de rendimiento

AGENT (Nivel 3)
├── Gestión de propiedades asignadas
├── Manejo de clientes y leads
└── Reportes de ventas

OWNER (Nivel 2)
├── Gestión de propiedades propias
├── Actualización de precios
└── Historial de rendimiento

BUYER (Nivel 1)
├── Búsqueda y filtrado
├── Sistema de favoritos
└── Contacto con agentes
```

### Stack Tecnológico
```
Backend:     Go 1.24 (net/http nativo)
Database:    PostgreSQL 15 + FTS
Cache:       LRU Cache personalizado
Images:      Procesamiento nativo Go
Testing:     testify + mocks
Dev Env:     Docker + GoLand
```

## 🎯 **Casos de Uso Principales**

### Para Propietarios
```go
// Crear nueva propiedad
property := domain.NewProperty(
    "Casa en Samborondón", 
    "Casa moderna con piscina",
    "Guayas", "Samborondón", "house", 
    285000.0, ownerID,
)

// Asignar a agencia
property.AssignToAgency(agencyID, agentID, userID)
```

### Para Agencias
```go
// Crear agencia
agency := domain.NewAgency(
    "Inmobiliaria Los Andes",
    "info@losandes.com",
    "0984567890",
    "Av. Amazonas 123",
    "Quito", "Pichincha",
    "1234567890123", // RUC
    ownerID,
)

// Gestionar agentes
user.SetAgency(agencyID) // Para agentes
```

### Para Compradores
```go
// Búsqueda con FTS
properties := propertyService.SearchWithFTS(
    "casa piscina Samborondón",
    filters,
    pagination,
)

// Filtros avanzados
filters := PropertyFilters{
    Province:    "Guayas",
    PriceMin:    200000,
    PriceMax:    300000,
    PropertyType: "house",
    Bedrooms:    3,
}
```

## 📈 **Estadísticas del Proyecto**

| Métrica | Valor | Detalle |
|---------|-------|---------|
| **Endpoints API** | 26+ | REST completo documentado |
| **Tests** | 200+ | 90%+ cobertura en todas las capas |
| **Migraciones DB** | 20+ | PostgreSQL con FTS y constraints |
| **Roles de Usuario** | 5 | Sistema jerárquico completo |
| **Provincias Ecuador** | 24 | Validación completa integrada |
| **Performance** | <100ms | Respuesta promedio con cache |

## 🛠️ **Herramientas de Desarrollo**

### Claude Code Commands (Optimización 10x)
```bash
# Comandos optimizados para desarrollo rápido
claude > /project:realty-property "add elevator field"
claude > /project:realty-api "create favorites endpoint"
claude > /project:realty-cache "optimize search caching"
```

### Testing Estratégico
```bash
# Tests por capa
make test-domain      # Business logic
make test-service     # Application layer  
make test-repository  # Data layer
make test-handlers    # HTTP layer
make test-integration # End-to-end
```

## 🔐 **Seguridad**

- **JWT Authentication** con roles embebidos
- **Bcrypt hashing** para passwords
- **Role-based permissions** granulares
- **SQL injection protection** con prepared statements
- **Input validation** en todas las capas
- **Rate limiting** preparado para producción

## 🌟 **Próximas Funcionalidades**

- [ ] **Multi-tenancy** completo para SaaS
- [ ] **Notificaciones** en tiempo real
- [ ] **Dashboard analytics** avanzado
- [ ] **API móvil** optimizada
- [ ] **Integración WhatsApp** Business
- [ ] **ML recommendations** de propiedades

## 📝 **Licencia**

Este proyecto está bajo la Licencia MIT. Ver [LICENSE](LICENSE) para más detalles.

## 🤝 **Contribuciones**

¡Las contribuciones son bienvenidas! Por favor:

1. Fork el repositorio
2. Crea una feature branch (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Add nueva funcionalidad'`)
4. Push a la branch (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## 💬 **Soporte**

- 📖 **Documentación:** [Wiki del proyecto](https://github.com/tu-usuario/realty-core/wiki)
- 🐛 **Issues:** [GitHub Issues](https://github.com/tu-usuario/realty-core/issues)
- 💡 **Discusiones:** [GitHub Discussions](https://github.com/tu-usuario/realty-core/discussions)

---

**Hecho con ❤️ para el mercado inmobiliario ecuatoriano**