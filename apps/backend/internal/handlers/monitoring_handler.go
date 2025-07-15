package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"realty-core/internal/logging"
	"realty-core/internal/monitoring"
)

// MonitoringHandler handles monitoring and metrics endpoints
type MonitoringHandler struct {
	metricsCollector *monitoring.MetricsCollector
	alertManager     *monitoring.AlertManager
	logger           *logging.Logger
}

// NewMonitoringHandler creates a new monitoring handler
func NewMonitoringHandler(metricsCollector *monitoring.MetricsCollector, alertManager *monitoring.AlertManager) *MonitoringHandler {
	return &MonitoringHandler{
		metricsCollector: metricsCollector,
		alertManager:     alertManager,
		logger:           logging.GetGlobalLogger(),
	}
}

// GetMetrics returns current application metrics
func (mh *MonitoringHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Update system metrics before getting snapshot
	mh.metricsCollector.UpdateSystemMetrics()
	
	snapshot := mh.metricsCollector.GetMetricsSnapshot()
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(snapshot)
	
	if mh.logger != nil {
		mh.logger.Info("Metrics snapshot requested", map[string]interface{}{
			"uptime_seconds":    snapshot.Uptime.Seconds(),
			"memory_mb":         float64(snapshot.System.Memory) / 1024 / 1024,
			"goroutines":        snapshot.System.Goroutines,
			"cache_hit_rate":    snapshot.Cache.HitRate,
		})
	}
}

// GetPrometheusMetrics returns metrics in Prometheus format
func (mh *MonitoringHandler) GetPrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	// Update system metrics
	mh.metricsCollector.UpdateSystemMetrics()
	snapshot := mh.metricsCollector.GetMetricsSnapshot()
	
	// Generate Prometheus format output
	output := mh.generatePrometheusOutput(snapshot)
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

// generatePrometheusOutput generates Prometheus-formatted metrics
func (mh *MonitoringHandler) generatePrometheusOutput(snapshot monitoring.MetricsSnapshot) string {
	var output string
	
	// System metrics
	output += "# HELP realty_core_uptime_seconds Total uptime in seconds\n"
	output += "# TYPE realty_core_uptime_seconds counter\n"
	output += "realty_core_uptime_seconds " + strconv.FormatFloat(snapshot.Uptime.Seconds(), 'f', 2, 64) + "\n\n"
	
	output += "# HELP realty_core_memory_bytes Current memory usage in bytes\n"
	output += "# TYPE realty_core_memory_bytes gauge\n"
	output += "realty_core_memory_bytes " + strconv.FormatInt(snapshot.System.Memory, 10) + "\n\n"
	
	output += "# HELP realty_core_goroutines Current number of goroutines\n"
	output += "# TYPE realty_core_goroutines gauge\n"
	output += "realty_core_goroutines " + strconv.Itoa(snapshot.System.Goroutines) + "\n\n"
	
	// Database metrics
	output += "# HELP realty_core_db_connections Current database connections\n"
	output += "# TYPE realty_core_db_connections gauge\n"
	output += "realty_core_db_connections " + strconv.FormatFloat(snapshot.Database.Connections, 'f', 0, 64) + "\n\n"
	
	output += "# HELP realty_core_db_queries_total Total database queries\n"
	output += "# TYPE realty_core_db_queries_total counter\n"
	output += "realty_core_db_queries_total " + strconv.FormatInt(snapshot.Database.Queries, 10) + "\n\n"
	
	output += "# HELP realty_core_db_query_duration_ms Average database query duration\n"
	output += "# TYPE realty_core_db_query_duration_ms gauge\n"
	output += "realty_core_db_query_duration_ms " + strconv.FormatFloat(snapshot.Database.QueryDuration, 'f', 2, 64) + "\n\n"
	
	// Cache metrics
	output += "# HELP realty_core_cache_hits_total Total cache hits\n"
	output += "# TYPE realty_core_cache_hits_total counter\n"
	output += "realty_core_cache_hits_total " + strconv.FormatInt(snapshot.Cache.Hits, 10) + "\n\n"
	
	output += "# HELP realty_core_cache_misses_total Total cache misses\n"
	output += "# TYPE realty_core_cache_misses_total counter\n"
	output += "realty_core_cache_misses_total " + strconv.FormatInt(snapshot.Cache.Misses, 10) + "\n\n"
	
	output += "# HELP realty_core_cache_hit_rate Cache hit rate percentage\n"
	output += "# TYPE realty_core_cache_hit_rate gauge\n"
	output += "realty_core_cache_hit_rate " + strconv.FormatFloat(snapshot.Cache.HitRate, 'f', 2, 64) + "\n\n"
	
	// Business metrics
	output += "# HELP realty_core_properties_total Total number of properties\n"
	output += "# TYPE realty_core_properties_total gauge\n"
	output += "realty_core_properties_total " + strconv.FormatInt(snapshot.Business.Properties, 10) + "\n\n"
	
	output += "# HELP realty_core_images_total Total number of images\n"
	output += "# TYPE realty_core_images_total gauge\n"
	output += "realty_core_images_total " + strconv.FormatInt(snapshot.Business.Images, 10) + "\n\n"
	
	output += "# HELP realty_core_users_total Total number of users\n"
	output += "# TYPE realty_core_users_total gauge\n"
	output += "realty_core_users_total " + strconv.FormatInt(snapshot.Business.Users, 10) + "\n\n"
	
	output += "# HELP realty_core_agencies_total Total number of agencies\n"
	output += "# TYPE realty_core_agencies_total gauge\n"
	output += "realty_core_agencies_total " + strconv.FormatInt(snapshot.Business.Agencies, 10) + "\n\n"
	
	// HTTP metrics
	for endpoint, metrics := range snapshot.HTTP {
		sanitized := sanitizeMetricName(endpoint)
		
		output += "# HELP realty_core_http_requests_total Total HTTP requests for " + endpoint + "\n"
		output += "# TYPE realty_core_http_requests_total counter\n"
		output += "realty_core_http_requests_total{endpoint=\"" + endpoint + "\"} " + strconv.FormatInt(metrics.Requests, 10) + "\n\n"
		
		output += "# HELP realty_core_http_duration_ms HTTP request duration for " + endpoint + "\n"
		output += "# TYPE realty_core_http_duration_ms gauge\n"
		output += "realty_core_http_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"avg\"} " + strconv.FormatFloat(metrics.AvgDuration, 'f', 2, 64) + "\n"
		output += "realty_core_http_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"p95\"} " + strconv.FormatFloat(metrics.P95Duration, 'f', 2, 64) + "\n"
		output += "realty_core_http_duration_ms{endpoint=\"" + endpoint + "\",quantile=\"p99\"} " + strconv.FormatFloat(metrics.P99Duration, 'f', 2, 64) + "\n\n"
		
		_ = sanitized // Use variable to avoid unused warning
	}
	
	return output
}

// GetAlerts returns current active alerts
func (mh *MonitoringHandler) GetAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	alerts := mh.alertManager.GetActiveAlerts()
	summary := mh.alertManager.GetAlertSummary()
	
	response := AlertsResponse{
		Summary: summary,
		Alerts:  alerts,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if mh.logger != nil {
		mh.logger.Info("Active alerts requested", map[string]interface{}{
			"total_alerts":     summary.Total,
			"critical_alerts": summary.ByLevel[monitoring.AlertLevelCritical],
			"warning_alerts":  summary.ByLevel[monitoring.AlertLevelWarning],
		})
	}
}

// AlertsResponse contains alerts information
type AlertsResponse struct {
	Summary monitoring.AlertSummary `json:"summary"`
	Alerts  []*monitoring.Alert     `json:"alerts"`
}

// GetAlertHistory returns alert history
func (mh *MonitoringHandler) GetAlertHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Parse limit parameter
	limit := 50 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	history := mh.alertManager.GetAlertHistory(limit)
	
	response := AlertHistoryResponse{
		Limit:   limit,
		Count:   len(history),
		History: history,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if mh.logger != nil {
		mh.logger.Info("Alert history requested", map[string]interface{}{
			"limit":         limit,
			"returned_count": len(history),
		})
	}
}

// AlertHistoryResponse contains alert history information
type AlertHistoryResponse struct {
	Limit   int                 `json:"limit"`
	Count   int                 `json:"count"`
	History []*monitoring.Alert `json:"history"`
}

// GetAlertRules returns configured alert rules
func (mh *MonitoringHandler) GetAlertRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	rules := mh.alertManager.GetRules()
	
	// Convert to response format (excluding the function field)
	rulesResponse := make(map[string]AlertRuleResponse)
	for name, rule := range rules {
		rulesResponse[name] = AlertRuleResponse{
			Name:        rule.Name,
			Description: rule.Description,
			Level:       rule.Level,
			Cooldown:    rule.Cooldown,
			Enabled:     rule.Enabled,
			Tags:        rule.Tags,
		}
	}
	
	response := AlertRulesResponse{
		Count: len(rulesResponse),
		Rules: rulesResponse,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if mh.logger != nil {
		mh.logger.Info("Alert rules requested", map[string]interface{}{
			"rule_count": len(rulesResponse),
		})
	}
}

// AlertRuleResponse contains alert rule information for API response
type AlertRuleResponse struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Level       monitoring.AlertLevel `json:"level"`
	Cooldown    time.Duration         `json:"cooldown"`
	Enabled     bool                  `json:"enabled"`
	Tags        map[string]string     `json:"tags"`
}

// AlertRulesResponse contains alert rules information
type AlertRulesResponse struct {
	Count int                            `json:"count"`
	Rules map[string]AlertRuleResponse  `json:"rules"`
}

// UpdateAlertRule updates an alert rule configuration
func (mh *MonitoringHandler) UpdateAlertRule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Extract rule name from URL path
	ruleName := r.URL.Path[len("/api/monitoring/rules/"):]
	if ruleName == "" {
		http.Error(w, "Rule name required", http.StatusBadRequest)
		return
	}
	
	var request AlertRuleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	// Update rule
	if request.Enabled != nil {
		if *request.Enabled {
			mh.alertManager.EnableRule(ruleName)
		} else {
			mh.alertManager.DisableRule(ruleName)
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "updated",
		"rule":    ruleName,
	})
	
	if mh.logger != nil {
		mh.logger.Info("Alert rule updated", map[string]interface{}{
			"rule_name": ruleName,
			"enabled":   request.Enabled,
		})
	}
}

// AlertRuleUpdateRequest contains alert rule update data
type AlertRuleUpdateRequest struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// GetDashboard returns monitoring dashboard data
func (mh *MonitoringHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Update system metrics
	mh.metricsCollector.UpdateSystemMetrics()
	
	// Get current metrics snapshot
	snapshot := mh.metricsCollector.GetMetricsSnapshot()
	
	// Get alert summary
	alertSummary := mh.alertManager.GetAlertSummary()
	
	// Calculate additional dashboard metrics
	totalRequests := int64(0)
	avgResponseTime := float64(0)
	for _, httpMetric := range snapshot.HTTP {
		totalRequests += httpMetric.Requests
		if httpMetric.AvgDuration > avgResponseTime {
			avgResponseTime = httpMetric.AvgDuration
		}
	}
	
	dashboard := DashboardResponse{
		Timestamp: time.Now(),
		Uptime:    snapshot.Uptime,
		System: SystemDashboard{
			Memory:     snapshot.System.Memory,
			MemoryMB:   float64(snapshot.System.Memory) / 1024 / 1024,
			CPU:        snapshot.System.CPU,
			Goroutines: snapshot.System.Goroutines,
		},
		Performance: PerformanceDashboard{
			TotalRequests:   totalRequests,
			AvgResponseTime: avgResponseTime,
			CacheHitRate:    snapshot.Cache.HitRate,
			DBConnections:   snapshot.Database.Connections,
			DBQueryTime:     snapshot.Database.QueryDuration,
		},
		Business: BusinessDashboard{
			Properties: snapshot.Business.Properties,
			Images:     snapshot.Business.Images,
			Users:      snapshot.Business.Users,
			Agencies:   snapshot.Business.Agencies,
		},
		Alerts: AlertsDashboard{
			Total:    alertSummary.Total,
			Critical: alertSummary.ByLevel[monitoring.AlertLevelCritical],
			Warning:  alertSummary.ByLevel[monitoring.AlertLevelWarning],
			Info:     alertSummary.ByLevel[monitoring.AlertLevelInfo],
		},
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboard)
	
	if mh.logger != nil {
		mh.logger.Info("Monitoring dashboard requested", map[string]interface{}{
			"uptime_hours":     snapshot.Uptime.Hours(),
			"memory_mb":        dashboard.System.MemoryMB,
			"total_requests":   totalRequests,
			"active_alerts":    alertSummary.Total,
		})
	}
}

// DashboardResponse contains monitoring dashboard data
type DashboardResponse struct {
	Timestamp   time.Time             `json:"timestamp"`
	Uptime      time.Duration         `json:"uptime"`
	System      SystemDashboard       `json:"system"`
	Performance PerformanceDashboard  `json:"performance"`
	Business    BusinessDashboard     `json:"business"`
	Alerts      AlertsDashboard       `json:"alerts"`
}

// SystemDashboard contains system metrics for dashboard
type SystemDashboard struct {
	Memory     int64   `json:"memory_bytes"`
	MemoryMB   float64 `json:"memory_mb"`
	CPU        float64 `json:"cpu_percent"`
	Goroutines int     `json:"goroutines"`
}

// PerformanceDashboard contains performance metrics for dashboard
type PerformanceDashboard struct {
	TotalRequests   int64   `json:"total_requests"`
	AvgResponseTime float64 `json:"avg_response_time_ms"`
	CacheHitRate    float64 `json:"cache_hit_rate"`
	DBConnections   float64 `json:"db_connections"`
	DBQueryTime     float64 `json:"db_query_time_ms"`
}

// BusinessDashboard contains business metrics for dashboard
type BusinessDashboard struct {
	Properties int64 `json:"properties"`
	Images     int64 `json:"images"`
	Users      int64 `json:"users"`
	Agencies   int64 `json:"agencies"`
}

// AlertsDashboard contains alert metrics for dashboard
type AlertsDashboard struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	Warning  int `json:"warning"`
	Info     int `json:"info"`
}

// sanitizeMetricName sanitizes metric names for Prometheus format
func sanitizeMetricName(name string) string {
	// Simple sanitization - replace invalid characters
	sanitized := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			sanitized += string(r)
		} else {
			sanitized += "_"
		}
	}
	return sanitized
}