# ================================
# Realty Core API - Makefile
# Sistema de Gestión Inmobiliaria
# ================================

# Variables
BINARY_NAME=realty-api
BUILD_DIR=bin
MAIN_PATH=./cmd/server
COVERAGE_FILE=coverage.out

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run test clean deps lint format check

# Comando por defecto
.DEFAULT_GOAL := help

## ================================
## 🏗️  BUILD COMMANDS
## ================================

## build: Construir el binario de la aplicación
build:
	@echo "$(GREEN)🔨 Construyendo aplicación...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd apps/backend && go build -o ../../$(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Binario creado: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## build-dev: Build con flags de desarrollo
build-dev:
	@echo "$(BLUE)🔨 Build modo desarrollo...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd apps/backend && go build -race -o ../../$(BUILD_DIR)/$(BINARY_NAME)-dev $(MAIN_PATH)
	@echo "$(GREEN)✅ Build dev completado$(NC)"

## build-prod: Build optimizado para producción
build-prod:
	@echo "$(BLUE)🔨 Build modo producción...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@cd apps/backend && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ../../$(BUILD_DIR)/$(BINARY_NAME)-prod $(MAIN_PATH)
	@echo "$(GREEN)✅ Build producción completado$(NC)"

## ================================
## 🏃 RUN COMMANDS
## ================================

## run: Ejecutar la aplicación en modo desarrollo
run:
	@echo "$(GREEN)🚀 Ejecutando servidor...$(NC)"
	@cd apps/backend && go run $(MAIN_PATH)

## run-prod: Ejecutar binario de producción
run-prod: build-prod
	@echo "$(GREEN)🚀 Ejecutando servidor (producción)...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME)-prod

## ================================
## 🧪 TEST COMMANDS
## ================================

## test: Ejecutar todos los tests
test:
	@echo "$(YELLOW)🧪 Ejecutando todos los tests...$(NC)"
	@cd apps/backend && go test ./... -v

## test-short: Ejecutar tests rápidos (sin integración)
test-short:
	@echo "$(YELLOW)⚡ Ejecutando tests rápidos...$(NC)"
	@cd apps/backend && go test ./... -short -v

## test-cache: Tests específicos del sistema de cache
test-cache:
	@echo "$(YELLOW)🗃️  Ejecutando tests de cache...$(NC)"
	@cd apps/backend && go test ./internal/cache/... -v

## test-images: Tests del sistema de imágenes
test-images:
	@echo "$(YELLOW)🖼️  Ejecutando tests de imágenes...$(NC)"
	@cd apps/backend && go test ./internal/storage/... ./internal/processors/... -v

## test-properties: Tests del CRUD de propiedades
test-properties:
	@echo "$(YELLOW)🏠 Ejecutando tests de propiedades...$(NC)"
	@cd apps/backend && go test ./internal/domain/... ./internal/service/... ./internal/repository/... -v

## test-handlers: Tests de handlers HTTP
test-handlers:
	@echo "$(YELLOW)🌐 Ejecutando tests de handlers...$(NC)"
	@cd apps/backend && go test ./internal/handlers/... -v

## test-coverage: Ejecutar tests con reporte de cobertura
test-coverage:
	@echo "$(YELLOW)📊 Generando reporte de cobertura...$(NC)"
	@cd apps/backend && go test ./... -coverprofile=../../$(COVERAGE_FILE)
	@cd apps/backend && go tool cover -html=../../$(COVERAGE_FILE) -o ../../coverage.html
	@echo "$(GREEN)✅ Reporte generado: coverage.html$(NC)"

## test-bench: Ejecutar benchmarks del cache
test-bench:
	@echo "$(YELLOW)⚡ Ejecutando benchmarks...$(NC)"
	@cd apps/backend && go test ./internal/cache/... -bench=. -benchmem

## ================================
## 🔍 QUALITY COMMANDS  
## ================================

## lint: Ejecutar linter
lint:
	@echo "$(BLUE)🔍 Ejecutando linter...$(NC)"
	@cd apps/backend && go vet ./...
	@echo "$(GREEN)✅ Linting completado$(NC)"

## format: Formatear código
format:
	@echo "$(BLUE)📝 Formateando código...$(NC)"
	@cd apps/backend && go fmt ./...
	@echo "$(GREEN)✅ Formato aplicado$(NC)"

## check: Verificación completa (format + lint + test)
check: format lint test-short
	@echo "$(GREEN)✅ Verificación completa exitosa$(NC)"

## check-full: Verificación completa con todos los tests
check-full: format lint test
	@echo "$(GREEN)✅ Verificación completa con todos los tests$(NC)"

## ================================
## 📦 DEPENDENCY COMMANDS
## ================================

## deps: Descargar dependencias
deps:
	@echo "$(BLUE)📦 Descargando dependencias...$(NC)"
	@cd apps/backend && go mod download
	@pnpm install
	@echo "$(GREEN)✅ Dependencias descargadas$(NC)"

## deps-update: Actualizar dependencias
deps-update:
	@echo "$(BLUE)🔄 Actualizando dependencias...$(NC)"
	@cd apps/backend && go get -u ./... && go mod tidy
	@pnpm update
	@echo "$(GREEN)✅ Dependencias actualizadas$(NC)"

## deps-tidy: Limpiar dependencias no utilizadas
deps-tidy:
	@echo "$(BLUE)🧹 Limpiando dependencias...$(NC)"
	@cd apps/backend && go mod tidy
	@echo "$(GREEN)✅ Dependencias limpiadas$(NC)"

## ================================
## 🗃️  DATABASE COMMANDS
## ================================

## db-connect: Conectar a PostgreSQL local
db-connect:
	@echo "$(BLUE)🐘 Conectando a PostgreSQL local...$(NC)"
	@psql -h localhost -p 5433 -U juanquizhpi -d inmobiliaria_db

## db-status: Verificar estado de PostgreSQL local
db-status:
	@echo "$(BLUE)📋 Estado de PostgreSQL:$(NC)"
	@psql -h localhost -p 5433 -U juanquizhpi -d inmobiliaria_db -c "SELECT version();"

## ================================
## 🧹 CLEANUP COMMANDS
## ================================

## clean: Limpiar archivos generados
clean:
	@echo "$(YELLOW)🧹 Limpiando archivos generados...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE) coverage.html
	@echo "$(GREEN)✅ Limpieza completada$(NC)"

## clean-cache: Limpiar cache de Go
clean-cache:
	@echo "$(YELLOW)🗑️  Limpiando cache de Go...$(NC)"
	@cd apps/backend && go clean -cache
	@echo "$(GREEN)✅ Cache limpiado$(NC)"

## ================================
## 📈 DEVELOPMENT WORKFLOWS
## ================================

## dev: Workflow completo de desarrollo
dev: clean deps format lint test-short build-dev
	@echo "$(GREEN)🎉 Workflow de desarrollo completado$(NC)"

## ci: Workflow de integración continua
ci: deps check-full build
	@echo "$(GREEN)🎉 Workflow CI completado$(NC)"

## release: Workflow de release
release: clean deps check-full test-coverage build-prod
	@echo "$(GREEN)🎉 Release preparado$(NC)"

## ================================
## 🗄️ DATABASE MIGRATION COMMANDS
## ================================

## migrate-up: Aplicar todas las migraciones pendientes
migrate-up:
	@echo "$(BLUE)🆙 Aplicando migraciones...$(NC)"
	@tools/scripts/migrate.sh up

## migrate-down: Revertir última migración (o N migraciones)
migrate-down:
	@echo "$(YELLOW)⬇️  Revirtiendo migración...$(NC)"
	@tools/scripts/migrate.sh down $(N)

## migrate-version: Ver versión actual de migraciones
migrate-version:
	@echo "$(BLUE)📊 Versión actual de migraciones:$(NC)"
	@tools/scripts/migrate.sh version

## migrate-create: Crear nueva migración (usar NAME=nombre_migracion)
migrate-create:
	@echo "$(GREEN)📝 Creando nueva migración: $(NAME)$(NC)"
	@tools/scripts/migrate.sh create $(NAME)

## migrate-force: Forzar versión específica (usar VERSION=numero)
migrate-force:
	@echo "$(RED)⚠️  Forzando versión $(VERSION) (peligroso!)$(NC)"
	@tools/scripts/migrate.sh force $(VERSION)

## migrate-validate: Validar conexión a base de datos
migrate-validate:
	@echo "$(BLUE)🔍 Validando conexión a base de datos...$(NC)"
	@tools/scripts/migrate.sh validate

## db-setup: Setup completo de base de datos (migraciones + datos de prueba)
db-setup: migrate-up
	@echo "$(GREEN)🎉 Base de datos configurada exitosamente$(NC)"

## db-reset: Resetear base de datos (WARNING: destructivo)
db-reset:
	@echo "$(RED)⚠️  RESETEO DESTRUCTIVO - Presiona Ctrl+C para cancelar$(NC)"
	@sleep 3
	@tools/scripts/migrate.sh down 999
	@tools/scripts/migrate.sh up
	@echo "$(GREEN)🔄 Base de datos reseteada$(NC)"

## ================================
## 📚 DOCUMENTATION COMMANDS
## ================================

## sync-docs: Sincronizar toda la documentación desde PROGRESS.md
sync-docs:
	@echo "$(BLUE)📚 Sincronizando documentación...$(NC)"
	@echo "$(YELLOW)⚠️  Herramienta sync-docs no disponible en nueva estructura$(NC)"

## validate-docs: Validar consistencia de documentación
validate-docs:
	@echo "$(YELLOW)🔍 Validando consistencia de documentación...$(NC)"
	@echo "$(YELLOW)⚠️  Herramienta sync-docs no disponible en nueva estructura$(NC)"

## check-docs: Verificar estado actual de documentación
check-docs:
	@echo "$(BLUE)📋 Estado actual de documentación:$(NC)"
	@echo "$(GREEN)📁 Documentación organizada en /docs/$(NC)"
	@ls -la docs/

## fix-docs: Forzar sincronización y validación completa
fix-docs: sync-docs validate-docs
	@echo "$(GREEN)🎉 Documentación corregida y validada$(NC)"

## ================================
## 📊 PROJECT INFO
## ================================

## info: Mostrar información del proyecto
info:
	@echo "$(BLUE)📊 Información del Proyecto$(NC)"
	@echo "$(YELLOW)Nombre:$(NC) Realty Core API"
	@echo "$(YELLOW)Versión Go:$(NC) $$(go version)"
	@echo "$(YELLOW)Archivos Go:$(NC) $$(find . -name '*.go' | wc -l)"
	@echo "$(YELLOW)Tests:$(NC) $$(find . -name '*_test.go' | wc -l)"
	@echo "$(YELLOW)Líneas de código:$(NC) $$(find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -1)"

## status: Estado actual del desarrollo
status:
	@echo "$(BLUE)📈 Estado del Desarrollo$(NC)"
	@echo "$(YELLOW)Funcionalidades completadas:$(NC)"
	@echo "  ✅ CRUD Propiedades"
	@echo "  ✅ PostgreSQL FTS"
	@echo "  ✅ Sistema de Imágenes"
	@echo "  ✅ Cache LRU"
	@echo "  ✅ Sistema de Paginación"
	@echo "$(YELLOW)Próximas funcionalidades:$(NC)"
	@echo "  🔄 Sistema de Usuarios"
	@echo "  🔄 Dashboard Analytics"

## ================================
## 🆘 HELP
## ================================

## help: Mostrar esta ayuda
help:
	@echo "$(GREEN)================================"
	@echo "🏠 Realty Core API - Makefile"
	@echo "Sistema de Gestión Inmobiliaria"
	@echo "================================$(NC)"
	@echo ""
	@echo "$(YELLOW)📋 Comandos Disponibles:$(NC)"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'
	@echo ""
	@echo "$(BLUE)💡 Ejemplos de uso:$(NC)"
	@echo "  make dev          # Desarrollo completo"
	@echo "  make test-cache   # Solo tests de cache"
	@echo "  make ci          # Pipeline CI"
	@echo "  make build-prod  # Build producción"
	@echo ""
	@echo "$(GREEN)🚀 Para empezar: make dev$(NC)"