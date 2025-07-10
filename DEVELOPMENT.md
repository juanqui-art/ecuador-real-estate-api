# 🚀 Guía de Desarrollo - GoLand + PostgreSQL Local

Esta guía te ayuda a configurar el entorno de desarrollo usando GoLand con PostgreSQL local.

## 📋 Prerrequisitos

✅ GoLand 2025.1.3 o superior  
✅ PostgreSQL instalado localmente  
✅ Go 1.24 instalado  

## 🐳 Configuración de PostgreSQL Local

### 1. Verificar PostgreSQL Local

1. **Verificar que PostgreSQL esté corriendo:**
   ```bash
   # Verificar servicio PostgreSQL
   brew services list | grep postgresql
   # o
   sudo service postgresql status
   ```

2. **Iniciar PostgreSQL si no está corriendo:**
   ```bash
   # macOS con Homebrew
   brew services start postgresql
   # o Linux
   sudo service postgresql start
   ```

3. **Crear base de datos del proyecto:**
   ```bash
   # Conectar a PostgreSQL
   psql postgres
   
   # Crear base de datos
   CREATE DATABASE inmobiliaria_db;
   
   # Crear usuario si es necesario
   CREATE USER admin WITH PASSWORD 'password';
   GRANT ALL PRIVILEGES ON DATABASE inmobiliaria_db TO admin;
   ```

### 2. Configurar Database Connection

1. **Abrir Database Tool Window:**
   ```
   View → Tool Windows → Database
   ```

2. **Añadir Data Source:**
   ```
   Click en "+" → Data Source → PostgreSQL
   
   Configuración:
   - Host: localhost
   - Port: 5432 (puerto estándar PostgreSQL)
   - Database: inmobiliaria_db
   - User: tu_usuario_local (ej: juanquizhpi)
   - Password: tu_password_local (o vacío si no tienes)
   ```

3. **Test Connection:** Click "Test Connection" → Debe decir "Successful"

## ⚙️ Configuración de Run/Debug

### 1. Crear Run Configuration para API

1. **Edit Configurations:**
   ```
   Run → Edit Configurations → Click "+"
   ```

2. **Seleccionar Go Build:**
   ```
   Name: Servidor Inmobiliaria API
   Run kind: Package
   Package path: realty-core/cmd/servidor
   Working directory: /Users/juanquizhpi/GolandProjects/realty-core
   ```

3. **Environment Variables:**
   ```
   DATABASE_URL=postgresql://tu_usuario:tu_password@localhost:5432/inmobiliaria_db
   PORT=8080
   LOG_LEVEL=info
   GO_ENV=development
   ```

### 2. Crear Run Configuration para Tests

```
Name: Tests Inmobiliaria
Run kind: Package
Package path: realty-core/...
Pattern: .*_test\.go
```

## 🗄️ Gestión de Base de Datos

### Ejecutar Migraciones

1. **Desde Database Tool Window:**
   ```
   Database → inmobiliaria_db → Console
   ```

2. **Ejecutar archivo SQL:**
   ```
   Drag & Drop archivo .sql desde migraciones/
   O abrir archivo y ejecutar con Ctrl+Enter
   ```

### Ver Datos

1. **Explorar tablas:**
   ```
   Database → inmobiliaria_db → schemas → public → tables
   ```

2. **Ver contenido:**
   ```
   Double-click en tabla "propiedades"
   ```

## 🔧 Flujo de Desarrollo Diario

### 1. Iniciar Desarrollo
```
1. Abrir GoLand
2. Verificar PostgreSQL local corriendo: brew services list | grep postgresql
3. Run → "Servidor Inmobiliaria API"
4. Verificar en http://localhost:8080/api/health
```

### 2. Hacer Cambios
```
1. Modificar código Go
2. GoLand auto-recompila (Hot Reload)
3. Probar cambios en browser/Postman
4. Ver logs en GoLand console
```

### 3. Debugging
```
1. Poner breakpoints en código
2. Run → Debug "Servidor Inmobiliaria API"
3. Hacer requests HTTP
4. Inspeccionar variables en GoLand
```

## 🌐 URLs de Desarrollo

| Servicio | URL | Descripción |
|----------|-----|-------------|
| API | http://localhost:8080 | API REST principal |
| Health Check | http://localhost:8080/api/health | Estado del servicio |
| PostgreSQL | localhost:5432 | Base de datos directa |

## 📊 Herramientas GoLand Útiles

### HTTP Client
```
Tools → HTTP Client → Create Request
```

### Database Console
```
Database → Console → New Console
```

### PostgreSQL Logs
```
# Ver logs PostgreSQL local
tail -f /usr/local/var/log/postgresql.log
# o
sudo journalctl -u postgresql
```

### Git Integration
```
Git → Commit → Push (integrado)
```

## 🔧 Comandos Útiles desde GoLand Terminal

```bash
# Verificar estado PostgreSQL
brew services list | grep postgresql

# Reiniciar PostgreSQL
brew services restart postgresql

# Conectar a la base de datos
psql -d inmobiliaria_db

# Backup base de datos
pg_dump inmobiliaria_db > backup.sql

# Restaurar base de datos
psql inmobiliaria_db < backup.sql
```

## ❗ Troubleshooting

### PostgreSQL no inicia
1. Verificar estado del servicio: `brew services list | grep postgresql`
2. Verificar logs: `tail -f /usr/local/var/log/postgresql.log`
3. Puerto 5432 no ocupado por otra aplicación: `lsof -i :5432`

### Database connection falla
1. Verificar PostgreSQL corriendo: `brew services list | grep postgresql`
2. Test Connection en Database Tool Window
3. Verificar credenciales locales (usuario/password)

### API no conecta con BD
1. Verificar variables de entorno en Run Configuration
2. Verificar DATABASE_URL correcto para PostgreSQL local
3. Ver logs de la aplicación Go

## 🎯 Checklist de Seguimiento Diario

### Estado Actual (2025-01-10)
- ✅ Configurar entorno GoLand + PostgreSQL local
- ✅ Sistema de propiedades completo (CRUD + FTS + paginación)
- ✅ Arquitectura completa: Domain/Service/Repository/Handlers
- ✅ Testing comprehensivo (179 tests, 90%+ cobertura)
- ✅ Sistema de imágenes implementado (13 endpoints + cache LRU)
- ✅ Sistema de usuarios y agencias (domain structures + validaciones)
- ✅ PostgreSQL FTS con ranking y autocompletado
- 🔄 Integración endpoints avanzados (imagen, usuario, agencia)

### Checklist Sesión de Trabajo

#### Al Iniciar Sesión
- [ ] Leer PROGRESS.md para contexto
- [ ] Verificar estado tests: `go test ./...`
- [ ] Verificar API corriendo: `go run cmd/server/main.go`
- [ ] Revisar git status: `git status`

#### Durante Desarrollo
- [ ] Implementar funcionalidad específica
- [ ] Crear/actualizar tests correspondientes
- [ ] Verificar cobertura: `go test -cover ./...`
- [ ] Verificar formato: `go fmt ./...`
- [ ] Verificar código: `go vet ./...`

#### Al Finalizar Feature
- [ ] Todos los tests pasan
- [ ] Cobertura >90% mantenida
- [ ] Commit con mensaje descriptivo
- [ ] Actualizar PROGRESS.md
- [ ] Actualizar CLAUDE.md si es necesario

#### Antes de Cerrar Sesión
- [ ] Commit final con estado actual
- [ ] Actualizar PROGRESS.md con próximos pasos
- [ ] Verificar que no hay cambios sin commitear
- [ ] Anotar cualquier problema o bloqueador

### Comandos Rápidos

```bash
# Status completo
git status && go test -cover ./...

# Commit rápido
git add . && git commit -m "feat: [descripción]"

# Verificar funcionalidad
curl http://localhost:8080/api/health

# Ver logs recientes
git log --oneline -10
```

### Funcionalidades Próximas

#### 1. Integración Sistema de Imágenes
- [ ] Activar ImageHandler en main.go
- [ ] Solucionar dependencias ImageService
- [ ] Probar 13 endpoints de imágenes
- [ ] Verificar cache LRU funcionando

#### 2. Integración Sistema de Usuarios
- [ ] Activar UserHandler en main.go
- [ ] Implementar JWT authentication
- [ ] Probar 10 endpoints de usuarios
- [ ] Configurar roles y permisos

#### 3. Integración Sistema de Agencias
- [ ] Activar AgencyHandler en main.go
- [ ] Probar 15 endpoints de agencias
- [ ] Configurar relaciones agencia-agente
- [ ] Sistema de comisiones

#### 4. Endpoints Avanzados de Búsqueda
- [ ] Activar endpoints FTS paginados
- [ ] Probar búsqueda avanzada con filtros
- [ ] Optimizar performance queries
- [ ] Sistema de autocompletado

### Troubleshooting Común

#### IDLE se ralentiza
1. Hacer commit frecuente del progreso
2. Consultar PROGRESS.md para contexto
3. Verificar estado con `go test ./...`
4. Continuar desde último estado conocido

#### Tests fallan
1. Verificar cambios con `git diff`
2. Ejecutar test específico: `go test ./internal/[layer] -v`
3. Revisar logs de error completos
4. Consultar tests similares en codebase

#### API no responde
1. Verificar PostgreSQL: `brew services list | grep postgresql`
2. Verificar variables entorno en Run Configuration
3. Revisar logs aplicación Go
4. Probar health check: `curl http://localhost:8080/api/health`