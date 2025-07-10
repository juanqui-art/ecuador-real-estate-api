# 🏠 Ecuador Real Estate API

Una API REST moderna y robusta para gestión de propiedades inmobiliarias en Ecuador, construida con Go y PostgreSQL.

## 🚀 Características Principales

- **CRUD Completo** - Gestión completa de propiedades inmobiliarias
- **Búsqueda Avanzada** - PostgreSQL Full-Text Search en español
- **Sistema de Imágenes** - Upload, procesamiento y optimización automática
- **Cache Inteligente** - Sistema LRU para máximo rendimiento
- **API RESTful** - 26+ endpoints bien documentados
- **Paginación** - Manejo eficiente de grandes datasets
- **Localización Ecuador** - Provincias, ciudades y validaciones locales

## 🛠️ Tecnologías

- **Backend**: Go 1.24+ con `net/http` nativo
- **Base de Datos**: PostgreSQL 15 con Full-Text Search
- **Cache**: LRU Cache personalizado con TTL
- **Procesamiento**: Imágenes con resize y compresión
- **Testing**: 157 tests con 90%+ cobertura
- **Desarrollo**: Docker Compose para entorno local

## 📋 Prerrequisitos

- Go 1.24 o superior
- PostgreSQL 15+
- Docker y Docker Compose (opcional)
- Make (para comandos optimizados)

## 🔧 Instalación

### 1. Clonar el repositorio
```bash
git clone https://github.com/tuusuario/ecuador-real-estate-api.git
cd ecuador-real-estate-api
```

### 2. Configurar variables de entorno
```bash
cp .env.example .env
# Editar .env con tus configuraciones
```

### 3. Iniciar base de datos
```bash
# Con Docker Compose
make db-up

# O manualmente con PostgreSQL local
createdb inmobiliaria_db
```

### 4. Ejecutar migraciones
```bash
# Las migraciones se ejecutan automáticamente al iniciar
make run
```

### 5. Instalar dependencias y ejecutar
```bash
make dev
```

## 🎯 Uso Rápido

### Iniciar el servidor
```bash
make run
# Servidor disponible en http://localhost:8080
```

### Comandos útiles
```bash
make help          # Ver todos los comandos disponibles
make test-cache    # Tests del sistema de cache
make test-images   # Tests del sistema de imágenes
make ci           # Pipeline completo de CI
make build-prod   # Build para producción
```

## 📚 Documentación de la API

### Endpoints Principales

#### 🏠 Propiedades
```bash
GET    /api/properties              # Listar propiedades
POST   /api/properties              # Crear propiedad
GET    /api/properties/{id}         # Obtener por ID
PUT    /api/properties/{id}         # Actualizar
DELETE /api/properties/{id}         # Eliminar
```

#### 🔍 Búsqueda
```bash
GET    /api/properties/filter            # Filtros básicos
GET    /api/properties/search/ranked     # Búsqueda FTS
GET    /api/properties/search/suggestions # Autocompletado
POST   /api/properties/search/advanced   # Búsqueda avanzada
```

#### 🖼️ Imágenes
```bash
POST   /api/images                  # Upload imagen
GET    /api/images/{id}/thumbnail   # Obtener thumbnail
GET    /api/images/{id}/variant     # Obtener variante
GET    /api/images/cache/stats      # Estadísticas cache
```

### Ejemplos de Uso

#### Crear una propiedad
```bash
curl -X POST http://localhost:8080/api/properties \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Casa en Samborondón",
    "descripcion": "Hermosa casa con piscina",
    "precio": 285000,
    "provincia": "Guayas",
    "ciudad": "Samborondón",
    "tipo": "casa",
    "dormitorios": 4,
    "banos": 3.5,
    "area_m2": 320
  }'
```

#### Buscar propiedades
```bash
curl "http://localhost:8080/api/properties/search/ranked?q=casa+piscina&limit=10"
```

#### Subir imagen
```bash
curl -X POST http://localhost:8080/api/images \
  -F "image=@casa.jpg" \
  -F "property_id=123" \
  -F "alt_text=Fachada principal"
```

## 🏗️ Arquitectura

```
internal/
├── domain/          # Modelos y lógica de negocio
├── repository/      # Acceso a datos (PostgreSQL)
├── service/         # Lógica de aplicación
├── handlers/        # Controladores HTTP
├── cache/           # Sistema LRU cache
├── storage/         # Almacenamiento de archivos
├── processors/      # Procesamiento de imágenes
└── config/          # Configuración
```

## 🧪 Testing

```bash
make test           # Todos los tests
make test-cache     # Tests de cache
make test-images    # Tests de imágenes
make test-coverage  # Reporte de cobertura
make test-bench     # Benchmarks
```

## 📊 Métricas del Proyecto

- **157 tests** con 90%+ cobertura
- **32 archivos Go** organizados en capas
- **26+ endpoints API** funcionales
- **16,000+ líneas** de código
- **Arquitectura limpia** con separación de responsabilidades

## 🇪🇨 Localización Ecuador

### Provincias Soportadas
Azuay, Bolívar, Cañar, Carchi, Chimborazo, Cotopaxi, El Oro, Esmeraldas, Galápagos, Guayas, Imbabura, Loja, Los Ríos, Manabí, Morona Santiago, Napo, Orellana, Pastaza, Pichincha, Santa Elena, Santo Domingo, Sucumbíos, Tungurahua, Zamora Chinchipe

### Tipos de Propiedades
- **Casa** - Viviendas unifamiliares
- **Apartamento** - Departamentos y condominios
- **Terreno** - Lotes y terrenos
- **Comercial** - Locales comerciales y oficinas

## 🔧 Desarrollo

### Estructura del Proyecto
```bash
make info     # Información del proyecto
make status   # Estado actual del desarrollo
make dev      # Workflow completo de desarrollo
```

### Contribuir
1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## 📈 Roadmap

- [x] **v0.1** - CRUD básico de propiedades
- [x] **v0.2** - PostgreSQL FTS y búsqueda avanzada
- [x] **v0.3** - Sistema de testing comprehensivo
- [x] **v0.4** - Sistema de imágenes y cache LRU
- [ ] **v0.5** - Sistema de usuarios y autenticación
- [ ] **v0.6** - Dashboard y analytics
- [ ] **v0.7** - Multi-tenancy y SaaS

## 📝 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para más detalles.

## 🤝 Soporte

- **Issues**: [GitHub Issues](https://github.com/tuusuario/ecuador-real-estate-api/issues)
- **Documentación**: [Wiki del proyecto](https://github.com/tuusuario/ecuador-real-estate-api/wiki)
- **Email**: tu-email@ejemplo.com

---

**Desarrollado con ❤️ en Ecuador usando Go**