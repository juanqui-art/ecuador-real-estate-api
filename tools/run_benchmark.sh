#!/bin/bash

# run_benchmark.sh - Script para ejecutar benchmarks de performance del sistema

set -e

# Configuración
BENCHMARK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$BENCHMARK_DIR")"
RESULTS_DIR="$PROJECT_ROOT/benchmark_results"
SERVER_URL="http://localhost:8080"
SERVER_PID=""

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para logging
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Función para cleanup
cleanup() {
    if [ ! -z "$SERVER_PID" ]; then
        log "Deteniendo servidor (PID: $SERVER_PID)..."
        kill $SERVER_PID 2>/dev/null || true
        wait $SERVER_PID 2>/dev/null || true
    fi
}

# Registrar cleanup para exit
trap cleanup EXIT

# Función para verificar si el servidor está corriendo
check_server() {
    local max_attempts=30
    local attempt=1
    
    log "Verificando conectividad del servidor..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s -f "$SERVER_URL/api/health" > /dev/null 2>&1; then
            success "Servidor disponible en $SERVER_URL"
            return 0
        fi
        
        log "Intento $attempt/$max_attempts - Esperando servidor..."
        sleep 2
        ((attempt++))
    done
    
    error "Servidor no disponible después de $max_attempts intentos"
    return 1
}

# Función para iniciar el servidor si no está corriendo
start_server() {
    if curl -s -f "$SERVER_URL/api/health" > /dev/null 2>&1; then
        success "Servidor ya está corriendo"
        return 0
    fi
    
    log "Iniciando servidor de desarrollo..."
    
    # Cambiar al directorio del proyecto
    cd "$PROJECT_ROOT"
    
    # Verificar que el binario existe
    if [ ! -f "./server" ]; then
        log "Compilando servidor..."
        if ! go build -o server ./cmd/server; then
            error "Error compilando servidor"
            return 1
        fi
    fi
    
    # Iniciar servidor en background
    log "Ejecutando servidor en background..."
    ./server > server.log 2>&1 &
    SERVER_PID=$!
    
    log "Servidor iniciado con PID: $SERVER_PID"
    
    # Verificar que el servidor está respondiendo
    if check_server; then
        return 0
    else
        error "Error iniciando servidor"
        return 1
    fi
}

# Función para ejecutar benchmark
run_benchmark() {
    local profile_mode=$1
    local timestamp=$(date +'%Y%m%d_%H%M%S')
    local result_file="$RESULTS_DIR/benchmark_${profile_mode}_${timestamp}.txt"
    
    log "Ejecutando benchmark con perfil: $profile_mode"
    
    # Crear directorio de resultados si no existe
    mkdir -p "$RESULTS_DIR"
    
    # Cambiar al directorio tools para ejecutar
    cd "$BENCHMARK_DIR"
    
    # Ejecutar benchmark y guardar resultado
    if go run benchmark.go "$profile_mode" | tee "$result_file"; then
        success "Benchmark $profile_mode completado. Resultados en: $result_file"
        
        # Mover archivos de profiling si existen
        for prof_file in cpu.pprof mem.pprof block.pprof mutex.pprof; do
            if [ -f "$prof_file" ]; then
                mv "$prof_file" "$RESULTS_DIR/${profile_mode}_${timestamp}_${prof_file}"
                log "Archivo de profiling movido: $RESULTS_DIR/${profile_mode}_${timestamp}_${prof_file}"
            fi
        done
        
        return 0
    else
        error "Error ejecutando benchmark $profile_mode"
        return 1
    fi
}

# Función para generar reporte comparativo
generate_report() {
    local report_file="$RESULTS_DIR/benchmark_summary_$(date +'%Y%m%d_%H%M%S').md"
    
    log "Generando reporte comparativo..."
    
    cat > "$report_file" << EOF
# Reporte de Benchmark - Realty Core API

**Fecha:** $(date +'%Y-%m-%d %H:%M:%S')  
**Sistema:** $(uname -s) $(uname -r)  
**Go Version:** $(go version)  

## Resumen Ejecutivo

Este reporte contiene los resultados de benchmarking para los endpoints críticos del sistema de gestión inmobiliaria.

### Endpoints Testeados

1. **List Properties** - \`GET /api/properties\`
2. **Search Properties Ranked** - \`GET /api/properties/search/ranked\`
3. **Filter Properties by Province** - \`GET /api/properties/filter\`
4. **Property Statistics** - \`GET /api/properties/statistics\`
5. **Paginated Properties** - \`GET /api/pagination/properties\`
6. **Image Statistics** - \`GET /api/images/stats\`
7. **User Statistics** - \`GET /api/users/statistics\`
8. **Agency Statistics** - \`GET /api/agencies/statistics\`
9. **Health Check** - \`GET /api/health\`

## Configuración de Pruebas

- **Requests por endpoint:** 100
- **Concurrencia:** 10
- **Timeout:** 30 segundos
- **URL Base:** $SERVER_URL

## Archivos de Resultados

EOF

    # Listar archivos de resultados recientes
    find "$RESULTS_DIR" -name "benchmark_*.txt" -mtime -1 | sort | while read file; do
        echo "- [$(basename "$file")]($file)" >> "$report_file"
    done
    
    cat >> "$report_file" << EOF

## Archivos de Profiling

EOF

    # Listar archivos de profiling recientes
    find "$RESULTS_DIR" -name "*.pprof" -mtime -1 | sort | while read file; do
        echo "- [$(basename "$file")]($file)" >> "$report_file"
    done
    
    cat >> "$report_file" << EOF

## Análisis de Profiling

Para analizar los archivos de profiling:

\`\`\`bash
# Para CPU profiling
go tool pprof cpu_timestamp_cpu.pprof

# Para memory profiling  
go tool pprof mem_timestamp_mem.pprof

# Para visualización web
go tool pprof -http=:8081 cpu_timestamp_cpu.pprof
\`\`\`

## Métricas Clave

- **Success Rate:** Porcentaje de requests exitosos
- **Avg Time:** Tiempo promedio de respuesta
- **Req/Sec:** Requests por segundo
- **Memory Usage:** Uso de memoria durante las pruebas

## Recomendaciones

1. **Endpoints con > 100ms tiempo promedio:** Requieren optimización
2. **Success rate < 95%:** Investigar errores
3. **Memory usage > 50MB:** Revisar memory leaks
4. **Req/sec < 50:** Optimizar rendimiento

EOF

    success "Reporte generado: $report_file"
}

# Función principal
main() {
    log "🚀 Iniciando suite de benchmarks para Realty Core"
    log "=============================================="
    
    # Verificar dependencias
    if ! command -v go &> /dev/null; then
        error "Go no está instalado"
        exit 1
    fi
    
    if ! command -v curl &> /dev/null; then
        error "curl no está instalado"
        exit 1
    fi
    
    # Crear directorio de resultados
    mkdir -p "$RESULTS_DIR"
    
    # Iniciar servidor si no está corriendo
    if ! start_server; then
        error "No se pudo iniciar el servidor"
        exit 1
    fi
    
    # Ejecutar benchmarks con diferentes tipos de profiling
    local profiles=("cpu" "mem")
    local failed=0
    
    for profile in "${profiles[@]}"; do
        if ! run_benchmark "$profile"; then
            ((failed++))
            warning "Benchmark $profile falló"
        fi
        
        # Pausa entre benchmarks para estabilización
        log "Pausa de 5 segundos entre benchmarks..."
        sleep 5
    done
    
    # Generar reporte final
    generate_report
    
    if [ $failed -eq 0 ]; then
        success "🎉 Todos los benchmarks completados exitosamente"
        log "📊 Resultados disponibles en: $RESULTS_DIR"
    else
        warning "⚠️  $failed benchmark(s) fallaron. Revisar logs."
    fi
    
    # Mostrar archivos generados
    log "\n📁 Archivos generados:"
    ls -la "$RESULTS_DIR" | grep "$(date +'%Y%m%d')" || true
}

# Verificar argumentos
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  --help, -h     Mostrar esta ayuda"
    echo "  --no-server    No iniciar servidor automáticamente"
    echo ""
    echo "Este script ejecuta benchmarks de performance para el API de Realty Core."
    echo "Genera perfiles de CPU y memoria, y crea reportes detallados."
    exit 0
fi

if [ "$1" = "--no-server" ]; then
    log "Modo --no-server: asumiendo que el servidor ya está corriendo"
    if ! check_server; then
        error "Servidor no disponible. Inicie el servidor manualmente o ejecute sin --no-server"
        exit 1
    fi
else
    # Ejecutar función principal
    main
fi