# Docker Deployment Guide - Realty Core

Esta guía describe cómo desplegar Realty Core en producción usando Docker y Docker Compose.

## 🚀 Inicio Rápido

### Prerequisitos

- Docker 20.10+ 
- Docker Compose 2.0+
- Al menos 2GB RAM disponible
- 10GB espacio en disco

### Despliegue Básico

```bash
# 1. Clonar el repositorio
git clone <repository-url>
cd realty-core

# 2. Configurar variables de entorno
cp .env.production.example .env.production
# Editar .env.production con tus valores

# 3. Construir y desplegar
./scripts/deploy.sh deploy
```

## 📁 Estructura de Archivos

```
realty-core/
├── Dockerfile                    # Imagen base (scratch)
├── Dockerfile.distroless         # Imagen distroless (más segura)
├── .dockerignore                 # Archivos excluidos del build
├── docker-compose.production.yml # Configuración de producción
├── .env.production.example       # Variables de entorno ejemplo
├── nginx/
│   └── nginx.conf               # Configuración Nginx
└── scripts/
    ├── docker-build.sh          # Script de construcción
    └── deploy.sh                # Script de despliegue
```

## 🔧 Configuración

### Variables de Entorno Críticas

```bash
# Base de datos (REQUERIDO)
POSTGRES_PASSWORD=tu_password_seguro
DATABASE_URL=postgresql://user:pass@postgres:5432/db

# Aplicación
APP_PORT=8080
LOG_LEVEL=info

# Cache
CACHE_ENABLED=true
CACHE_SIZE_MB=50
```

### Perfiles de Servicios

El docker-compose incluye perfiles opcionales:

```bash
# Solo servicios básicos (app + postgres)
docker-compose up -d

# Con Redis para cache adicional
docker-compose --profile with-redis up -d

# Con Nginx como reverse proxy
docker-compose --profile with-nginx up -d

# Todos los servicios
docker-compose --profile with-redis --profile with-nginx up -d
```

## 🛠 Comandos de Gestión

### Scripts Disponibles

```bash
# Construcción de imágenes
./scripts/docker-build.sh                    # Imagen scratch
./scripts/docker-build.sh -t distroless      # Imagen distroless
./scripts/docker-build.sh -v 2.0.0 -l       # Versión específica + latest

# Despliegue y gestión
./scripts/deploy.sh build                    # Construir imágenes
./scripts/deploy.sh test                     # Ejecutar tests
./scripts/deploy.sh deploy                   # Despliegue completo
./scripts/deploy.sh status                   # Estado de servicios
./scripts/deploy.sh logs app                 # Ver logs
./scripts/deploy.sh restart app              # Reiniciar servicio
./scripts/deploy.sh backup                   # Backup base de datos
./scripts/deploy.sh health                   # Check de salud
```

### Comandos Docker Compose Directos

```bash
# Iniciar servicios
docker-compose -f docker-compose.production.yml up -d

# Ver logs en tiempo real
docker-compose -f docker-compose.production.yml logs -f app

# Escalar aplicación
docker-compose -f docker-compose.production.yml up -d --scale app=3

# Parar servicios
docker-compose -f docker-compose.production.yml down
```

## 🔍 Monitoreo y Logs

### Health Checks

```bash
# Check manual de salud
docker exec realty-core-app /inmobiliaria -health-check

# Estado de contenedores
docker-compose ps

# Uso de recursos
docker stats
```

### Logs

```bash
# Logs de aplicación
docker-compose logs app

# Logs de base de datos
docker-compose logs postgres

# Logs de Nginx (si está habilitado)
docker-compose logs nginx

# Logs con timestamps
docker-compose logs -t app
```

## 💾 Backup y Restore

### Backup Automático

```bash
# Backup manual
./scripts/deploy.sh backup

# El backup se guarda en: backups/realty_core_YYYYMMDD_HHMMSS.sql.gz
```

### Restauración

```bash
# Restaurar desde backup
./scripts/deploy.sh restore backups/realty_core_20240711_120000.sql.gz
```

### Backup Automatizado con Cron

```bash
# Agregar a crontab para backup diario a las 2 AM
0 2 * * * /path/to/realty-core/scripts/deploy.sh backup >> /var/log/realty-backup.log 2>&1
```

## 🔒 Seguridad

### Configuraciones de Seguridad

1. **Contenedores no-root**: Usa usuario `nonroot` en distroless
2. **Filesystem read-only**: Contenedores con filesystem de solo lectura
3. **No new privileges**: Previene escalación de privilegios
4. **Resource limits**: Límites de CPU y memoria configurados
5. **Health checks**: Monitoreo automático de salud

### SSL/TLS (Nginx)

```bash
# Configurar certificados SSL
mkdir -p nginx/ssl
# Copiar cert.pem y private.key a nginx/ssl/

# Habilitar HTTPS en nginx.conf (descomentar sección HTTPS)
# Reiniciar Nginx
docker-compose restart nginx
```

### Secrets Management

Para producción, usar Docker secrets:

```bash
# Crear secrets
echo "mi_password_seguro" | docker secret create postgres_password -
echo "redis_password" | docker secret create redis_password -

# Usar en docker-compose con external: true
```

## ⚡ Optimización de Performance

### Configuraciones Recomendadas

```bash
# Variables de entorno para optimización
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
CACHE_SIZE_MB=100

# Límites de recursos (docker-compose.yml)
deploy:
  resources:
    limits:
      memory: 512M
      cpus: '1.0'
```

### Escalado Horizontal

```bash
# Escalar aplicación a 3 instancias
docker-compose up -d --scale app=3

# Nginx balanceará automáticamente la carga
```

## 🚨 Troubleshooting

### Problemas Comunes

1. **Error de conexión a base de datos**
   ```bash
   # Verificar que PostgreSQL esté ejecutándose
   docker-compose ps postgres
   
   # Ver logs de PostgreSQL
   docker-compose logs postgres
   ```

2. **Aplicación no responde**
   ```bash
   # Check de salud
   ./scripts/deploy.sh health
   
   # Ver logs de aplicación
   docker-compose logs app
   ```

3. **Memoria insuficiente**
   ```bash
   # Ver uso de recursos
   docker stats
   
   # Ajustar límites en docker-compose.yml
   ```

### Logs de Debug

```bash
# Habilitar logs debug
echo "LOG_LEVEL=debug" >> .env.production
docker-compose restart app
```

## 📊 Métricas y Observabilidad

### Endpoints de Métricas

```bash
# Health check
curl http://localhost:8080/api/health

# Estadísticas de cache
curl http://localhost:8080/api/images/cache/stats

# Estadísticas generales
curl http://localhost:8080/api/properties/statistics
```

### Integración con Monitoring

Para integrar con Prometheus/Grafana:

```yaml
# Agregar al docker-compose.yml
services:
  prometheus:
    image: prom/prometheus
    # ... configuración
  
  grafana:
    image: grafana/grafana
    # ... configuración
```

## 🔄 Actualizaciones

### Rolling Updates

```bash
# Actualización sin downtime
./scripts/deploy.sh update

# O manualmente:
# 1. Backup
./scripts/deploy.sh backup

# 2. Build nueva imagen
./scripts/docker-build.sh

# 3. Deploy
docker-compose up -d --no-deps app
```

### Rollback

```bash
# Rollback a imagen anterior
docker tag realty-core:1.8.0 realty-core:latest
docker-compose up -d --no-deps app
```

## 📝 Mantenimiento

### Limpieza Regular

```bash
# Limpiar imágenes no utilizadas
docker image prune -f

# Limpiar contenedores parados
docker container prune -f

# Limpiar volúmenes no utilizados
docker volume prune -f

# Limpieza completa (CUIDADO en producción)
docker system prune -a -f
```

### Rotación de Logs

```bash
# Configurar logrotate para logs de Docker
sudo nano /etc/logrotate.d/docker

# Contenido:
/var/lib/docker/containers/*/*-json.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    sharedscripts
}
```

## 🆘 Soporte

Para problemas específicos:

1. Revisar logs: `./scripts/deploy.sh logs app`
2. Verificar salud: `./scripts/deploy.sh health`
3. Revisar configuración: verificar `.env.production`
4. Consultar documentación: `./scripts/deploy.sh --help`

---

**⚠️ Importante**: Siempre hacer backup antes de actualizaciones importantes y probar en ambiente de staging primero.