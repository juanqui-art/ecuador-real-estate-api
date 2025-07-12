package handlers

import (
	"encoding/json"
	"net/http"

	"realty-core/internal/logging"
	"realty-core/internal/middleware"
	"realty-core/internal/security"
)

// SecurityHandler handles security-related endpoints
type SecurityHandler struct {
	securityMiddleware *middleware.SecurityMiddleware
	passwordValidator  *security.PasswordValidator
	logger             *logging.Logger
}

// NewSecurityHandler creates a new security handler
func NewSecurityHandler(securityMiddleware *middleware.SecurityMiddleware) *SecurityHandler {
	return &SecurityHandler{
		securityMiddleware: securityMiddleware,
		passwordValidator:  security.NewPasswordValidator(),
		logger:             logging.GetGlobalLogger(),
	}
}

// SecurityMetrics returns current security metrics
func (sh *SecurityHandler) SecurityMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	metrics := sh.securityMiddleware.GetSecurityMetrics()
	loadStats := sh.securityMiddleware.GetLoadStats()
	
	response := SecurityMetricsResponse{
		SecurityMetrics: metrics,
		LoadStats:       loadStats,
		Timestamp:       metrics.LastReset,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if sh.logger != nil {
		sh.logger.Info("Security metrics requested", map[string]interface{}{
			"blocked_requests": metrics.BlockedRequests,
			"current_load":     loadStats.CurrentLoad,
			"request_count":    loadStats.RequestCount,
		})
	}
}

// SecurityMetricsResponse contains security metrics and load statistics
type SecurityMetricsResponse struct {
	SecurityMetrics security.SecurityMetricsSnapshot `json:"security_metrics"`
	LoadStats       security.LoadStats               `json:"load_stats"`
	Timestamp       interface{}                      `json:"timestamp"`
}

// ValidatePassword validates password strength
func (sh *SecurityHandler) ValidatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request PasswordValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	isValid, errors := sh.passwordValidator.ValidatePassword(request.Password)
	
	response := PasswordValidationResponse{
		IsValid: isValid,
		Errors:  errors,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if isValid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(response)
	
	if sh.logger != nil {
		sh.logger.Info("Password validation requested", map[string]interface{}{
			"is_valid":     isValid,
			"error_count":  len(errors),
		})
	}
}

// PasswordValidationRequest contains password validation request data
type PasswordValidationRequest struct {
	Password string `json:"password"`
}

// PasswordValidationResponse contains password validation results
type PasswordValidationResponse struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

// ValidateInput validates input for security threats
func (sh *SecurityHandler) ValidateInput(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request InputValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	validator := security.NewInputValidator()
	result := validator.ValidateInput(request.Input, request.Field)
	
	response := InputValidationResponse{
		IsValid:        result.IsValid,
		Field:          result.Field,
		Threats:        result.Threats,
		SanitizedInput: validator.SanitizeInput(request.Input),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if result.IsValid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(response)
	
	if sh.logger != nil {
		sh.logger.Info("Input validation requested", map[string]interface{}{
			"field":        request.Field,
			"is_valid":     result.IsValid,
			"threat_count": len(result.Threats),
		})
		
		if !result.IsValid {
			sh.logger.SecurityEvent(
				"Input Validation Failed",
				"",
				"Input validation detected threats",
				map[string]interface{}{
					"field":   request.Field,
					"threats": result.Threats,
				},
			)
		}
	}
}

// InputValidationRequest contains input validation request data
type InputValidationRequest struct {
	Input string `json:"input"`
	Field string `json:"field"`
}

// InputValidationResponse contains input validation results
type InputValidationResponse struct {
	IsValid        bool                     `json:"is_valid"`
	Field          string                   `json:"field"`
	Threats        []security.ThreatInfo    `json:"threats,omitempty"`
	SanitizedInput string                   `json:"sanitized_input"`
}

// SecurityStatus returns overall security status
func (sh *SecurityHandler) SecurityStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	metrics := sh.securityMiddleware.GetSecurityMetrics()
	loadStats := sh.securityMiddleware.GetLoadStats()
	
	// Determine security status based on metrics
	status := "healthy"
	alerts := make([]SecurityAlert, 0)
	
	// Check for high number of blocked requests
	if metrics.BlockedRequests > 100 {
		status = "warning"
		alerts = append(alerts, SecurityAlert{
			Level:   "warning",
			Message: "High number of blocked requests detected",
			Count:   metrics.BlockedRequests,
		})
	}
	
	// Check for high server load
	if loadStats.CurrentLoad > 80 {
		status = "critical"
		alerts = append(alerts, SecurityAlert{
			Level:   "critical",
			Message: "High server load detected",
			Value:   loadStats.CurrentLoad,
		})
	}
	
	// Check for high error rate
	if loadStats.ErrorRate > 20 {
		if status == "healthy" {
			status = "warning"
		}
		alerts = append(alerts, SecurityAlert{
			Level:   "warning",
			Message: "High error rate detected",
			Value:   loadStats.ErrorRate,
		})
	}
	
	response := SecurityStatusResponse{
		Status:    status,
		Alerts:    alerts,
		Metrics:   metrics,
		LoadStats: loadStats,
	}
	
	statusCode := http.StatusOK
	if status == "critical" {
		statusCode = http.StatusServiceUnavailable
	} else if status == "warning" {
		statusCode = http.StatusPartialContent
	}
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	
	if sh.logger != nil {
		sh.logger.Info("Security status requested", map[string]interface{}{
			"status":      status,
			"alert_count": len(alerts),
		})
	}
}

// SecurityAlert represents a security alert
type SecurityAlert struct {
	Level   string  `json:"level"`
	Message string  `json:"message"`
	Count   int64   `json:"count,omitempty"`
	Value   float64 `json:"value,omitempty"`
}

// SecurityStatusResponse contains security status information
type SecurityStatusResponse struct {
	Status    string                           `json:"status"`
	Alerts    []SecurityAlert                  `json:"alerts"`
	Metrics   security.SecurityMetricsSnapshot `json:"metrics"`
	LoadStats security.LoadStats               `json:"load_stats"`
}

// ThreatIntelligence returns threat intelligence information
func (sh *SecurityHandler) ThreatIntelligence(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	metrics := sh.securityMiddleware.GetSecurityMetrics()
	
	// Analyze threat patterns
	threatIntel := analyzeThreatPatterns(metrics)
	
	response := ThreatIntelligenceResponse{
		ThreatLevel:     threatIntel.Level,
		TopThreats:      threatIntel.TopThreats,
		SuspiciousIPs:   threatIntel.SuspiciousIPs,
		Recommendations: threatIntel.Recommendations,
		UpdatedAt:       metrics.LastReset,
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	
	if sh.logger != nil {
		sh.logger.Info("Threat intelligence requested", map[string]interface{}{
			"threat_level":    threatIntel.Level,
			"top_threat_count": len(threatIntel.TopThreats),
			"suspicious_ip_count": len(threatIntel.SuspiciousIPs),
		})
	}
}

// ThreatIntelligenceResponse contains threat intelligence information
type ThreatIntelligenceResponse struct {
	ThreatLevel     string              `json:"threat_level"`
	TopThreats      []ThreatPattern     `json:"top_threats"`
	SuspiciousIPs   []SuspiciousIP      `json:"suspicious_ips"`
	Recommendations []string            `json:"recommendations"`
	UpdatedAt       interface{}         `json:"updated_at"`
}

// ThreatPattern represents a threat pattern
type ThreatPattern struct {
	Type        string `json:"type"`
	Count       int64  `json:"count"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

// SuspiciousIP represents a suspicious IP address
type SuspiciousIP struct {
	IP          string `json:"ip"`
	Count       int64  `json:"count"`
	RiskLevel   string `json:"risk_level"`
	LastSeen    string `json:"last_seen"`
}

// ThreatIntelligence contains analyzed threat intelligence
type ThreatIntelligence struct {
	Level           string
	TopThreats      []ThreatPattern
	SuspiciousIPs   []SuspiciousIP
	Recommendations []string
}

// analyzeThreatPatterns analyzes security metrics to generate threat intelligence
func analyzeThreatPatterns(metrics security.SecurityMetricsSnapshot) ThreatIntelligence {
	intel := ThreatIntelligence{
		Level:           "low",
		TopThreats:      make([]ThreatPattern, 0),
		SuspiciousIPs:   make([]SuspiciousIP, 0),
		Recommendations: make([]string, 0),
	}
	
	// Analyze threat attempts
	for threatType, count := range metrics.ThreatAttempts {
		severity := "low"
		if count > 50 {
			severity = "high"
			intel.Level = "high"
		} else if count > 10 {
			severity = "medium"
			if intel.Level == "low" {
				intel.Level = "medium"
			}
		}
		
		intel.TopThreats = append(intel.TopThreats, ThreatPattern{
			Type:        threatType,
			Count:       count,
			Severity:    severity,
			Description: getThreatDescription(threatType),
		})
	}
	
	// Analyze suspicious IPs
	for ip, count := range metrics.SuspiciousIPs {
		riskLevel := "low"
		if count > 20 {
			riskLevel = "high"
		} else if count > 5 {
			riskLevel = "medium"
		}
		
		intel.SuspiciousIPs = append(intel.SuspiciousIPs, SuspiciousIP{
			IP:        ip,
			Count:     count,
			RiskLevel: riskLevel,
			LastSeen:  metrics.LastReset.Format("2006-01-02 15:04:05"),
		})
	}
	
	// Generate recommendations
	intel.Recommendations = generateRecommendations(metrics)
	
	return intel
}

// getThreatDescription returns a description for a threat type
func getThreatDescription(threatType string) string {
	descriptions := map[string]string{
		"sql_injection":         "SQL injection attempts detected",
		"xss":                   "Cross-site scripting attempts detected",
		"path_traversal":        "Path traversal attempts detected",
		"rate_limit_exceeded":   "Rate limit violations",
		"blocked_ip":            "Requests from blocked IP ranges",
		"suspicious_headers":    "Suspicious HTTP headers detected",
		"input_validation":      "Input validation failures",
	}
	
	if desc, exists := descriptions[threatType]; exists {
		return desc
	}
	
	return "Unknown threat pattern"
}

// generateRecommendations generates security recommendations based on metrics
func generateRecommendations(metrics security.SecurityMetricsSnapshot) []string {
	recommendations := make([]string, 0)
	
	if metrics.BlockedRequests > 50 {
		recommendations = append(recommendations, "Consider implementing additional IP blocking rules")
	}
	
	if len(metrics.SuspiciousIPs) > 10 {
		recommendations = append(recommendations, "Review and potentially block suspicious IP addresses")
	}
	
	if len(metrics.ThreatAttempts) > 5 {
		recommendations = append(recommendations, "Increase monitoring for detected threat patterns")
	}
	
	for threatType, count := range metrics.ThreatAttempts {
		if count > 20 {
			switch threatType {
			case "sql_injection":
				recommendations = append(recommendations, "Review database query parameterization")
			case "xss":
				recommendations = append(recommendations, "Implement additional output encoding")
			case "rate_limit_exceeded":
				recommendations = append(recommendations, "Consider adjusting rate limiting rules")
			}
		}
	}
	
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Security posture appears healthy")
	}
	
	return recommendations
}