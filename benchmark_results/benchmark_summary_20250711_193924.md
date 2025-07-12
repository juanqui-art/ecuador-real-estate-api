# Reporte de Benchmark - Realty Core API

**Fecha:** 2025-07-11 19:39:24  
**Sistema:** Darwin 24.5.0  
**Go Version:** go version go1.24.4 darwin/arm64  

## Resumen Ejecutivo

Este reporte contiene los resultados de benchmarking para los endpoints críticos del sistema de gestión inmobiliaria.

### Endpoints Testeados

1. **List Properties** - `GET /api/properties`
2. **Search Properties Ranked** - `GET /api/properties/search/ranked`
3. **Filter Properties by Province** - `GET /api/properties/filter`
4. **Property Statistics** - `GET /api/properties/statistics`
5. **Paginated Properties** - `GET /api/pagination/properties`
6. **Image Statistics** - `GET /api/images/stats`
7. **User Statistics** - `GET /api/users/statistics`
8. **Agency Statistics** - `GET /api/agencies/statistics`
9. **Health Check** - `GET /api/health`

## Configuración de Pruebas

- **Requests por endpoint:** 100
- **Concurrencia:** 10
- **Timeout:** 30 segundos
- **URL Base:** http://localhost:8080

## Archivos de Resultados

- [benchmark_cpu_20250711_193123.txt](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/benchmark_cpu_20250711_193123.txt)
- [benchmark_cpu_20250711_193913.txt](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/benchmark_cpu_20250711_193913.txt)
- [benchmark_mem_20250711_193130.txt](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/benchmark_mem_20250711_193130.txt)
- [benchmark_mem_20250711_193919.txt](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/benchmark_mem_20250711_193919.txt)

## Archivos de Profiling

- [cpu_20250711_193123_cpu.pprof](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/cpu_20250711_193123_cpu.pprof)
- [cpu_20250711_193913_cpu.pprof](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/cpu_20250711_193913_cpu.pprof)
- [mem_20250711_193130_mem.pprof](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/mem_20250711_193130_mem.pprof)
- [mem_20250711_193919_mem.pprof](/Users/juanquizhpi/GolandProjects/realty-core/benchmark_results/mem_20250711_193919_mem.pprof)

## Análisis de Profiling

Para analizar los archivos de profiling:

```bash
# Para CPU profiling
go tool pprof cpu_timestamp_cpu.pprof

# Para memory profiling  
go tool pprof mem_timestamp_mem.pprof

# Para visualización web
go tool pprof -http=:8081 cpu_timestamp_cpu.pprof
```

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

