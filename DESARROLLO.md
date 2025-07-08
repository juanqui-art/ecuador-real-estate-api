# 🚀 Guía de Desarrollo - GoLand + Docker

Esta guía te ayuda a configurar el entorno de desarrollo usando GoLand y Docker Desktop.

## 📋 Prerrequisitos

✅ GoLand 2025.1.3 o superior  
✅ Docker Desktop instalado y ejecutándose  
✅ Go 1.24 instalado  

## 🐳 Configuración de Docker Compose en GoLand

### 1. Iniciar Servicios desde GoLand

1. **Abrir Services Tool Window:**
   ```
   View → Tool Windows → Services
   ```

2. **Añadir Docker Compose:**
   ```
   Click en "+" → Docker → Docker Compose
   Seleccionar: docker-compose.yml
   ```

3. **Iniciar PostgreSQL:**
   ```
   Services → docker-compose → postgres → Click derecho → Start
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
   - Port: 5432
   - Database: inmobiliaria_db
   - User: admin
   - Password: password
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
   DATABASE_URL=postgresql://admin:password@localhost:5432/inmobiliaria_db
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
2. Services → Start postgres (si no está corriendo)
3. Run → "Servidor Inmobiliaria API"
4. Verificar en http://localhost:8080/api/salud
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
| Health Check | http://localhost:8080/api/salud | Estado del servicio |
| pgAdmin | http://localhost:5050 | Interfaz web PostgreSQL |
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

### Docker Logs
```
Services → postgres → Logs
```

### Git Integration
```
Git → Commit → Push (integrado)
```

## 🔧 Comandos Útiles desde GoLand Terminal

```bash
# Verificar containers
docker-compose ps

# Ver logs de PostgreSQL
docker-compose logs postgres

# Recrear base de datos
docker-compose down postgres
docker-compose up postgres

# Backup base de datos
docker-compose exec postgres pg_dump -U admin inmobiliaria_db > backup.sql
```

## ❗ Troubleshooting

### PostgreSQL no inicia
1. Verificar Docker Desktop corriendo
2. Services → postgres → Logs → Ver errores
3. Puerto 5432 no ocupado por otra aplicación

### Database connection falla
1. Verificar PostgreSQL corriendo: Services → postgres → Status
2. Test Connection en Database Tool Window
3. Verificar credenciales en docker-compose.yml

### API no conecta con BD
1. Verificar variables de entorno en Run Configuration
2. Verificar DATABASE_URL correcto
3. Ver logs de la aplicación Go

## 🎯 Próximos Pasos

1. ✅ Configurar entorno GoLand + Docker
2. ⏳ Ejecutar migraciones
3. ⏳ Probar API endpoints
4. ⏳ Crear tests unitarios
5. ⏳ Implementar funcionalidades adicionales