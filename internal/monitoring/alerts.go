package monitoring

import (
	"fmt"
	"sync"
	"time"
)

// AlertLevel represents the severity level of an alert
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert represents a monitoring alert
type Alert struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Level       AlertLevel             `json:"level"`
	Message     string                 `json:"message"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Resolved    bool                   `json:"resolved"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
	Tags        map[string]string      `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Level       AlertLevel                `json:"level"`
	Condition   func(*MetricsSnapshot) bool `json:"-"`
	Cooldown    time.Duration             `json:"cooldown"`
	Enabled     bool                      `json:"enabled"`
	Tags        map[string]string         `json:"tags"`
	lastTriggered time.Time
}

// AlertManager manages monitoring alerts and rules
type AlertManager struct {
	mutex       sync.RWMutex
	rules       map[string]*AlertRule
	activeAlerts map[string]*Alert
	alertHistory []*Alert
	maxHistory  int
	
	// Notification channels
	notifyChannels []AlertNotifier
}

// AlertNotifier interface for alert notification channels
type AlertNotifier interface {
	Notify(alert *Alert) error
	GetType() string
}

// NewAlertManager creates a new alert manager
func NewAlertManager() *AlertManager {
	am := &AlertManager{
		rules:        make(map[string]*AlertRule),
		activeAlerts: make(map[string]*Alert),
		alertHistory: make([]*Alert, 0),
		maxHistory:   1000,
		notifyChannels: make([]AlertNotifier, 0),
	}
	
	// Register default alert rules
	am.registerDefaultRules()
	
	return am
}

// registerDefaultRules registers default monitoring alert rules
func (am *AlertManager) registerDefaultRules() {
	// High error rate alert
	am.AddRule(&AlertRule{
		Name:        "high_error_rate",
		Description: "HTTP error rate is above 10%",
		Level:       AlertLevelWarning,
		Cooldown:    5 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "http"},
		Condition: func(metrics *MetricsSnapshot) bool {
			totalRequests := int64(0)
			errorRequests := int64(0)
			
			for _, httpMetric := range metrics.HTTP {
				totalRequests += httpMetric.Requests
			}
			
			// Calculate error rate (would need error metrics)
			// This is a simplified version
			if totalRequests > 100 {
				errorRate := float64(errorRequests) / float64(totalRequests) * 100
				return errorRate > 10
			}
			return false
		},
	})
	
	// High response time alert
	am.AddRule(&AlertRule{
		Name:        "high_response_time",
		Description: "Average response time is above 1000ms",
		Level:       AlertLevelWarning,
		Cooldown:    3 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "performance"},
		Condition: func(metrics *MetricsSnapshot) bool {
			for _, httpMetric := range metrics.HTTP {
				if httpMetric.AvgDuration > 1000 {
					return true
				}
			}
			return false
		},
	})
	
	// High memory usage alert
	am.AddRule(&AlertRule{
		Name:        "high_memory_usage",
		Description: "Memory usage is above 500MB",
		Level:       AlertLevelCritical,
		Cooldown:    2 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "system"},
		Condition: func(metrics *MetricsSnapshot) bool {
			return metrics.System.Memory > 500*1024*1024 // 500MB
		},
	})
	
	// High goroutine count alert
	am.AddRule(&AlertRule{
		Name:        "high_goroutine_count",
		Description: "Goroutine count is above 1000",
		Level:       AlertLevelWarning,
		Cooldown:    5 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "system"},
		Condition: func(metrics *MetricsSnapshot) bool {
			return metrics.System.Goroutines > 1000
		},
	})
	
	// Low cache hit rate alert
	am.AddRule(&AlertRule{
		Name:        "low_cache_hit_rate",
		Description: "Cache hit rate is below 80%",
		Level:       AlertLevelWarning,
		Cooldown:    10 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "cache"},
		Condition: func(metrics *MetricsSnapshot) bool {
			return metrics.Cache.HitRate < 80 && (metrics.Cache.Hits+metrics.Cache.Misses) > 100
		},
	})
	
	// Database connection alert
	am.AddRule(&AlertRule{
		Name:        "high_db_connections",
		Description: "Database connections are above 20",
		Level:       AlertLevelCritical,
		Cooldown:    1 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "database"},
		Condition: func(metrics *MetricsSnapshot) bool {
			return metrics.Database.Connections > 20
		},
	})
	
	// Database query duration alert
	am.AddRule(&AlertRule{
		Name:        "slow_db_queries",
		Description: "Average database query duration is above 100ms",
		Level:       AlertLevelWarning,
		Cooldown:    5 * time.Minute,
		Enabled:     true,
		Tags:        map[string]string{"category": "database"},
		Condition: func(metrics *MetricsSnapshot) bool {
			return metrics.Database.QueryDuration > 100
		},
	})
}

// AddRule adds a new alert rule
func (am *AlertManager) AddRule(rule *AlertRule) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.rules[rule.Name] = rule
}

// RemoveRule removes an alert rule
func (am *AlertManager) RemoveRule(name string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	delete(am.rules, name)
}

// EnableRule enables an alert rule
func (am *AlertManager) EnableRule(name string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if rule, exists := am.rules[name]; exists {
		rule.Enabled = true
	}
}

// DisableRule disables an alert rule
func (am *AlertManager) DisableRule(name string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if rule, exists := am.rules[name]; exists {
		rule.Enabled = false
	}
}

// AddNotifier adds an alert notifier
func (am *AlertManager) AddNotifier(notifier AlertNotifier) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	am.notifyChannels = append(am.notifyChannels, notifier)
}

// EvaluateRules evaluates all alert rules against current metrics
func (am *AlertManager) EvaluateRules(metrics *MetricsSnapshot) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	now := time.Now()
	
	for _, rule := range am.rules {
		if !rule.Enabled {
			continue
		}
		
		// Check cooldown
		if now.Sub(rule.lastTriggered) < rule.Cooldown {
			continue
		}
		
		// Evaluate condition
		if rule.Condition(metrics) {
			// Create alert
			alert := &Alert{
				ID:          fmt.Sprintf("%s_%d", rule.Name, now.Unix()),
				Name:        rule.Name,
				Level:       rule.Level,
				Message:     rule.Description,
				Description: am.generateAlertDescription(rule, metrics),
				Timestamp:   now,
				Resolved:    false,
				Tags:        rule.Tags,
				Metadata:    am.generateAlertMetadata(rule, metrics),
			}
			
			am.triggerAlert(alert, rule)
		}
	}
}

// triggerAlert triggers a new alert
func (am *AlertManager) triggerAlert(alert *Alert, rule *AlertRule) {
	// Update rule trigger time
	rule.lastTriggered = alert.Timestamp
	
	// Add to active alerts
	am.activeAlerts[alert.ID] = alert
	
	// Add to history
	am.addToHistory(alert)
	
	// Send notifications
	go am.notifyAlert(alert)
}

// generateAlertDescription generates a detailed description for an alert
func (am *AlertManager) generateAlertDescription(rule *AlertRule, metrics *MetricsSnapshot) string {
	switch rule.Name {
	case "high_response_time":
		maxDuration := float64(0)
		for _, httpMetric := range metrics.HTTP {
			if httpMetric.AvgDuration > maxDuration {
				maxDuration = httpMetric.AvgDuration
			}
		}
		return fmt.Sprintf("Maximum average response time: %.2fms", maxDuration)
		
	case "high_memory_usage":
		return fmt.Sprintf("Current memory usage: %d bytes (%.2f MB)", 
			metrics.System.Memory, float64(metrics.System.Memory)/1024/1024)
			
	case "high_goroutine_count":
		return fmt.Sprintf("Current goroutine count: %d", metrics.System.Goroutines)
		
	case "low_cache_hit_rate":
		return fmt.Sprintf("Current cache hit rate: %.2f%% (hits: %d, misses: %d)", 
			metrics.Cache.HitRate, metrics.Cache.Hits, metrics.Cache.Misses)
			
	case "high_db_connections":
		return fmt.Sprintf("Current database connections: %.0f", metrics.Database.Connections)
		
	case "slow_db_queries":
		return fmt.Sprintf("Average query duration: %.2fms", metrics.Database.QueryDuration)
		
	default:
		return rule.Description
	}
}

// generateAlertMetadata generates metadata for an alert
func (am *AlertManager) generateAlertMetadata(rule *AlertRule, metrics *MetricsSnapshot) map[string]interface{} {
	metadata := map[string]interface{}{
		"uptime":     metrics.Uptime.String(),
		"timestamp":  metrics.Timestamp,
	}
	
	// Add rule-specific metadata
	switch rule.Name {
	case "high_response_time":
		metadata["http_metrics"] = metrics.HTTP
	case "high_memory_usage", "high_goroutine_count":
		metadata["system_metrics"] = metrics.System
	case "low_cache_hit_rate":
		metadata["cache_metrics"] = metrics.Cache
	case "high_db_connections", "slow_db_queries":
		metadata["database_metrics"] = metrics.Database
	}
	
	return metadata
}

// notifyAlert sends notifications for an alert
func (am *AlertManager) notifyAlert(alert *Alert) {
	for _, notifier := range am.notifyChannels {
		if err := notifier.Notify(alert); err != nil {
			// Log notification error (would use structured logger in real implementation)
			fmt.Printf("Failed to send alert notification via %s: %v\n", notifier.GetType(), err)
		}
	}
}

// ResolveAlert resolves an active alert
func (am *AlertManager) ResolveAlert(alertID string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if alert, exists := am.activeAlerts[alertID]; exists {
		now := time.Now()
		alert.Resolved = true
		alert.ResolvedAt = &now
		
		// Remove from active alerts
		delete(am.activeAlerts, alertID)
		
		// Update in history
		am.updateInHistory(alert)
	}
}

// GetActiveAlerts returns all active alerts
func (am *AlertManager) GetActiveAlerts() []*Alert {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	alerts := make([]*Alert, 0, len(am.activeAlerts))
	for _, alert := range am.activeAlerts {
		alerts = append(alerts, alert)
	}
	
	return alerts
}

// GetAlertHistory returns alert history
func (am *AlertManager) GetAlertHistory(limit int) []*Alert {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	if limit <= 0 || limit > len(am.alertHistory) {
		limit = len(am.alertHistory)
	}
	
	// Return most recent alerts
	start := len(am.alertHistory) - limit
	return am.alertHistory[start:]
}

// GetRules returns all alert rules
func (am *AlertManager) GetRules() map[string]*AlertRule {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	rules := make(map[string]*AlertRule)
	for name, rule := range am.rules {
		rules[name] = rule
	}
	
	return rules
}

// addToHistory adds an alert to the history
func (am *AlertManager) addToHistory(alert *Alert) {
	am.alertHistory = append(am.alertHistory, alert)
	
	// Trim history if it exceeds max size
	if len(am.alertHistory) > am.maxHistory {
		am.alertHistory = am.alertHistory[len(am.alertHistory)-am.maxHistory:]
	}
}

// updateInHistory updates an alert in the history
func (am *AlertManager) updateInHistory(updatedAlert *Alert) {
	for i, alert := range am.alertHistory {
		if alert.ID == updatedAlert.ID {
			am.alertHistory[i] = updatedAlert
			break
		}
	}
}

// GetAlertSummary returns a summary of current alert status
func (am *AlertManager) GetAlertSummary() AlertSummary {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	summary := AlertSummary{
		Total:  len(am.activeAlerts),
		ByLevel: make(map[AlertLevel]int),
	}
	
	for _, alert := range am.activeAlerts {
		summary.ByLevel[alert.Level]++
	}
	
	return summary
}

// AlertSummary contains a summary of alert status
type AlertSummary struct {
	Total   int                  `json:"total"`
	ByLevel map[AlertLevel]int   `json:"by_level"`
}

// LogNotifier is a simple notifier that logs alerts
type LogNotifier struct{}

// NewLogNotifier creates a new log notifier
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// Notify logs the alert
func (ln *LogNotifier) Notify(alert *Alert) error {
	fmt.Printf("[ALERT] %s - %s: %s\n", alert.Level, alert.Name, alert.Message)
	return nil
}

// GetType returns the notifier type
func (ln *LogNotifier) GetType() string {
	return "log"
}

// Global alert manager instance
var globalAlertManager *AlertManager

// InitializeAlerts initializes the global alert manager
func InitializeAlerts() {
	globalAlertManager = NewAlertManager()
	
	// Add default log notifier
	globalAlertManager.AddNotifier(NewLogNotifier())
}

// GetGlobalAlertManager returns the global alert manager
func GetGlobalAlertManager() *AlertManager {
	return globalAlertManager
}