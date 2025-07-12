package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	
	"github.com/pkg/profile"
)

// BenchmarkConfig configura los parámetros del benchmark
type BenchmarkConfig struct {
	BaseURL       string
	Requests      int
	Concurrency   int
	Timeout       time.Duration
	ProfileMode   string // "cpu", "mem", "block", "mutex"
}

// EndpointTest define un test específico para un endpoint
type EndpointTest struct {
	Name     string
	Method   string
	Path     string
	Body     string
	Headers  map[string]string
	Expected int // Status code esperado
}

// BenchmarkResult almacena los resultados del benchmark
type BenchmarkResult struct {
	EndpointName    string
	TotalRequests   int
	SuccessRequests int
	FailedRequests  int
	TotalTime       time.Duration
	AvgTime         time.Duration
	MinTime         time.Duration
	MaxTime         time.Duration
	RequestsPerSec  float64
	MemUsageMB      float64
}

// CriticalEndpoints define los endpoints más importantes para benchmarking
var CriticalEndpoints = []EndpointTest{
	{
		Name:     "List Properties",
		Method:   "GET",
		Path:     "/api/properties",
		Expected: 200,
	},
	{
		Name:     "Search Properties Ranked",
		Method:   "GET", 
		Path:     "/api/properties/search/ranked?q=casa&limit=10",
		Expected: 200,
	},
	{
		Name:     "Filter Properties by Province",
		Method:   "GET",
		Path:     "/api/properties/filter?province=Guayas",
		Expected: 200,
	},
	{
		Name:     "Property Statistics",
		Method:   "GET",
		Path:     "/api/properties/statistics",
		Expected: 200,
	},
	{
		Name:     "Paginated Properties",
		Method:   "GET",
		Path:     "/api/pagination/properties?page=1&page_size=20",
		Expected: 200,
	},
	{
		Name:     "Image Statistics",
		Method:   "GET",
		Path:     "/api/images/stats",
		Expected: 200,
	},
	{
		Name:     "User Statistics",
		Method:   "GET",
		Path:     "/api/users/statistics",
		Expected: 200,
	},
	{
		Name:     "Agency Statistics",
		Method:   "GET",
		Path:     "/api/agencies/statistics",
		Expected: 200,
	},
	{
		Name:     "Health Check",
		Method:   "GET",
		Path:     "/api/health",
		Expected: 200,
	},
}

func main() {
	config := BenchmarkConfig{
		BaseURL:     "http://localhost:8080",
		Requests:    100,
		Concurrency: 10,
		Timeout:     30 * time.Second,
		ProfileMode: "cpu", // Por defecto CPU profiling
	}

	// Verificar argumentos de línea de comandos
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cpu":
			config.ProfileMode = "cpu"
		case "mem":
			config.ProfileMode = "mem"
		case "block":
			config.ProfileMode = "block"
		case "mutex":
			config.ProfileMode = "mutex"
		default:
			fmt.Printf("Modo de profiling no válido: %s. Usando CPU por defecto.\n", os.Args[1])
		}
	}

	fmt.Printf("🚀 Iniciando Benchmark de Performance - Modo: %s\n", config.ProfileMode)
	fmt.Printf("⚙️  Configuración: %d requests, %d concurrency, timeout %v\n", 
		config.Requests, config.Concurrency, config.Timeout)
	fmt.Println("=" + fmt.Sprintf("%80s", "="))

	// Configurar profiling según el modo seleccionado
	var prof interface{ Stop() }
	switch config.ProfileMode {
	case "cpu":
		prof = profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	case "mem":
		prof = profile.Start(profile.MemProfile, profile.ProfilePath("."))
	case "block":
		prof = profile.Start(profile.BlockProfile, profile.ProfilePath("."))
	case "mutex":
		prof = profile.Start(profile.MutexProfile, profile.ProfilePath("."))
	}
	defer prof.Stop()

	// Verificar conectividad del servidor
	if !checkServerHealth(config.BaseURL) {
		log.Fatal("❌ Servidor no disponible en " + config.BaseURL)
	}

	fmt.Println("✅ Servidor conectado correctamente")
	fmt.Println()

	var allResults []BenchmarkResult
	
	// Ejecutar benchmark para cada endpoint crítico
	for _, endpoint := range CriticalEndpoints {
		fmt.Printf("🔄 Testeando: %s\n", endpoint.Name)
		
		result := benchmarkEndpoint(config, endpoint)
		allResults = append(allResults, result)
		
		printEndpointResult(result)
		fmt.Println()
	}

	// Reporte final
	printFinalReport(allResults)
	
	// Información del sistema
	printSystemInfo()
}

func checkServerHealth(baseURL string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baseURL + "/api/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func benchmarkEndpoint(config BenchmarkConfig, endpoint EndpointTest) BenchmarkResult {
	client := &http.Client{Timeout: config.Timeout}
	url := config.BaseURL + endpoint.Path
	
	var results []time.Duration
	successCount := 0
	failedCount := 0
	
	// Obtener memoria inicial
	var memStart runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)
	
	startTime := time.Now()
	
	// Ejecutar requests
	for i := 0; i < config.Requests; i++ {
		requestStart := time.Now()
		
		req, err := http.NewRequest(endpoint.Method, url, nil)
		if err != nil {
			failedCount++
			continue
		}
		
		// Añadir headers si existen
		for key, value := range endpoint.Headers {
			req.Header.Set(key, value)
		}
		
		resp, err := client.Do(req)
		if err != nil {
			failedCount++
			continue
		}
		
		// Leer y descartar el body para simular uso real
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		
		requestDuration := time.Since(requestStart)
		results = append(results, requestDuration)
		
		if resp.StatusCode == endpoint.Expected {
			successCount++
		} else {
			failedCount++
		}
	}
	
	totalTime := time.Since(startTime)
	
	// Obtener memoria final
	var memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memEnd)
	
	// Calcular estadísticas
	avgTime := calculateAverage(results)
	minTime := calculateMin(results)
	maxTime := calculateMax(results)
	requestsPerSec := float64(config.Requests) / totalTime.Seconds()
	memUsage := float64(memEnd.Alloc-memStart.Alloc) / 1024 / 1024 // MB
	
	return BenchmarkResult{
		EndpointName:    endpoint.Name,
		TotalRequests:   config.Requests,
		SuccessRequests: successCount,
		FailedRequests:  failedCount,
		TotalTime:       totalTime,
		AvgTime:         avgTime,
		MinTime:         minTime,
		MaxTime:         maxTime,
		RequestsPerSec:  requestsPerSec,
		MemUsageMB:      memUsage,
	}
}

func calculateAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func calculateMin(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	min := durations[0]
	for _, d := range durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}

func calculateMax(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	max := durations[0]
	for _, d := range durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

func printEndpointResult(result BenchmarkResult) {
	fmt.Printf("  📊 Resultados para %s:\n", result.EndpointName)
	fmt.Printf("     ✅ Success: %d/%d (%.1f%%)\n", 
		result.SuccessRequests, result.TotalRequests, 
		float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("     ⏱️  Tiempo promedio: %v\n", result.AvgTime)
	fmt.Printf("     ⚡ Requests/seg: %.2f\n", result.RequestsPerSec)
	fmt.Printf("     🧠 Memoria usada: %.2f MB\n", result.MemUsageMB)
	fmt.Printf("     📈 Min/Max: %v / %v\n", result.MinTime, result.MaxTime)
}

func printFinalReport(results []BenchmarkResult) {
	fmt.Println("📈 REPORTE FINAL DE PERFORMANCE")
	fmt.Println("=" + fmt.Sprintf("%80s", "="))
	
	var totalRequests, totalSuccess int
	var totalTime time.Duration
	var avgRequestsPerSec float64
	
	fmt.Printf("%-25s | %8s | %8s | %10s | %8s\n", 
		"Endpoint", "Success%", "Avg Time", "Req/Sec", "Mem MB")
	fmt.Println(fmt.Sprintf("%80s", "-"))
	
	for _, result := range results {
		successRate := float64(result.SuccessRequests) / float64(result.TotalRequests) * 100
		fmt.Printf("%-25s | %7.1f%% | %8v | %8.2f | %6.2f\n",
			truncateString(result.EndpointName, 25),
			successRate,
			result.AvgTime,
			result.RequestsPerSec,
			result.MemUsageMB)
		
		totalRequests += result.TotalRequests
		totalSuccess += result.SuccessRequests
		totalTime += result.TotalTime
		avgRequestsPerSec += result.RequestsPerSec
	}
	
	fmt.Println(fmt.Sprintf("%80s", "-"))
	fmt.Printf("TOTALES: %d requests, %.1f%% success, %.2f avg req/sec\n",
		totalRequests,
		float64(totalSuccess)/float64(totalRequests)*100,
		avgRequestsPerSec/float64(len(results)))
}

func printSystemInfo() {
	fmt.Println("\n💻 INFORMACIÓN DEL SISTEMA")
	fmt.Println("=" + fmt.Sprintf("%80s", "="))
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("Arch: %s\n", runtime.GOARCH)
	fmt.Printf("CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Go Version: %s\n", runtime.Version())
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memoria asignada: %.2f MB\n", float64(m.Alloc)/1024/1024)
	fmt.Printf("Total asignaciones: %.2f MB\n", float64(m.TotalAlloc)/1024/1024)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}