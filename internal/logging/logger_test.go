package logging

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	config := Config{
		Level:       InfoLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)

	assert.NotNil(t, logger)
	assert.Equal(t, InfoLevel, logger.level)
	assert.Equal(t, "test-service", logger.serviceName)
	assert.Equal(t, "1.0.0", logger.version)
}

func TestLogLevel_String(t *testing.T) {
	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{FatalLevel, "FATAL"},
		{LogLevel(99), "UNKNOWN"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, tc.level.String())
	}
}

func TestParseLogLevel(t *testing.T) {
	testCases := []struct {
		input    string
		expected LogLevel
	}{
		{"DEBUG", DebugLevel},
		{"debug", DebugLevel},
		{"INFO", InfoLevel},
		{"info", InfoLevel},
		{"WARN", WarnLevel},
		{"warn", WarnLevel},
		{"WARNING", WarnLevel},
		{"ERROR", ErrorLevel},
		{"error", ErrorLevel},
		{"FATAL", FatalLevel},
		{"fatal", FatalLevel},
		{"unknown", InfoLevel},
		{"", InfoLevel},
	}

	for _, tc := range testCases {
		result := ParseLogLevel(tc.input)
		assert.Equal(t, tc.expected, result, "Failed for input: %s", tc.input)
	}
}

func TestLogger_WithField(t *testing.T) {
	config := Config{
		Level:       InfoLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	newLogger := logger.WithField("key1", "value1")

	// Original logger should not have the field
	assert.Empty(t, logger.fields)

	// New logger should have the field
	assert.Equal(t, "value1", newLogger.fields["key1"])
}

func TestLogger_WithFields(t *testing.T) {
	config := Config{
		Level:       InfoLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}
	newLogger := logger.WithFields(fields)

	// Original logger should not have the fields
	assert.Empty(t, logger.fields)

	// New logger should have all the fields
	assert.Equal(t, "value1", newLogger.fields["key1"])
	assert.Equal(t, 123, newLogger.fields["key2"])
	assert.Equal(t, true, newLogger.fields["key3"])
}

func TestLogger_LogLevels(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       DebugLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	// Test Debug level
	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "debug message")
	assert.Contains(t, buf.String(), "DEBUG")
	buf.Reset()

	// Test Info level
	logger.Info("info message")
	assert.Contains(t, buf.String(), "info message")
	assert.Contains(t, buf.String(), "INFO")
	buf.Reset()

	// Test Warn level
	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "warn message")
	assert.Contains(t, buf.String(), "WARN")
	buf.Reset()

	// Test Error level with error
	logger.Error("error message", assert.AnError)
	assert.Contains(t, buf.String(), "error message")
	assert.Contains(t, buf.String(), "ERROR")
	assert.Contains(t, buf.String(), assert.AnError.Error())
	buf.Reset()
}

func TestLogger_LogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       WarnLevel, // Only WARN and above should be logged
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	// Debug and Info should not be logged
	logger.Debug("debug message")
	logger.Info("info message")
	assert.Empty(t, buf.String())

	// Warn should be logged
	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "warn message")
	buf.Reset()

	// Error should be logged
	logger.Error("error message", nil)
	assert.Contains(t, buf.String(), "error message")
}

func TestLogger_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       InfoLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	fields := map[string]interface{}{
		"user_id": "123",
		"action":  "test",
	}

	logger.Info("test message", fields)

	// Parse the JSON output
	var logEntry LogEntry
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify fields
	assert.Equal(t, "INFO", logEntry.Level)
	assert.Equal(t, "test message", logEntry.Message)
	assert.Equal(t, "test-service", logEntry.Service)
	assert.Equal(t, "1.0.0", logEntry.Version)
	assert.Equal(t, "123", logEntry.Fields["user_id"])
	assert.Equal(t, "test", logEntry.Fields["action"])
	assert.NotEmpty(t, logEntry.Timestamp)
	assert.NotEmpty(t, logEntry.Caller)
}

func TestLogger_HTTPRequest(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       InfoLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	duration := 150 * time.Millisecond
	logger.HTTPRequest("GET", "/api/test", 200, duration, "test-agent", "127.0.0.1")

	// Parse the JSON output
	var logEntry LogEntry
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify HTTP-specific fields
	assert.Equal(t, "GET", logEntry.Method)
	assert.Equal(t, "/api/test", logEntry.URL)
	assert.Equal(t, 200, logEntry.StatusCode)
	assert.Equal(t, int64(150), logEntry.Duration)
	assert.Equal(t, "test-agent", logEntry.UserAgent)
	assert.Equal(t, "127.0.0.1", logEntry.RemoteAddr)
}

func TestLogger_DatabaseQuery(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       DebugLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	duration := 50 * time.Millisecond
	logger.DatabaseQuery("SELECT * FROM properties WHERE id = $1", duration, 1)

	// Parse the JSON output
	var logEntry LogEntry
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify database-specific fields
	assert.Equal(t, "Database Query", logEntry.Message)
	assert.Contains(t, logEntry.Fields["query"], "SELECT * FROM properties")
	assert.Equal(t, float64(50), logEntry.Fields["duration_ms"])
	assert.Equal(t, float64(1), logEntry.Fields["rows_affected"])
}

func TestLogger_CacheOperation(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       DebugLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	duration := 5 * time.Millisecond
	logger.CacheOperation("GET", "property:123", true, duration)

	// Parse the JSON output
	var logEntry LogEntry
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify cache-specific fields
	assert.Equal(t, "Cache Operation", logEntry.Message)
	assert.Equal(t, "GET", logEntry.Fields["operation"])
	assert.Equal(t, "property:123", logEntry.Fields["key"])
	assert.Equal(t, true, logEntry.Fields["hit"])
	assert.Equal(t, float64(5), logEntry.Fields["duration_ms"])
}

func TestLogger_SecurityEvent(t *testing.T) {
	var buf bytes.Buffer
	
	config := Config{
		Level:       WarnLevel,
		ServiceName: "test-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&buf, "", 0)

	logger.SecurityEvent("Failed Login", "user123", "Invalid password")

	// Parse the JSON output
	var logEntry LogEntry
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	// Verify security-specific fields
	assert.Equal(t, "Security Event", logEntry.Message)
	assert.Equal(t, "user123", logEntry.UserID)
	assert.Equal(t, "Failed Login", logEntry.Fields["event"])
	assert.Equal(t, "Invalid password", logEntry.Fields["details"])
}

func TestTruncateQuery(t *testing.T) {
	// Short query should not be truncated
	shortQuery := "SELECT * FROM properties"
	result := truncateQuery(shortQuery)
	assert.Equal(t, shortQuery, result)

	// Long query should be truncated
	longQuery := strings.Repeat("SELECT * FROM properties WHERE condition = 'value' AND ", 10)
	result = truncateQuery(longQuery)
	assert.True(t, len(result) <= 203) // 200 + "..."
	assert.True(t, strings.HasSuffix(result, "..."))
}

func TestGlobalLogger(t *testing.T) {
	// Test setting and getting global logger
	config := Config{
		Level:       InfoLevel,
		ServiceName: "global-test",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	SetGlobalLogger(logger)

	retrieved := GetGlobalLogger()
	assert.Equal(t, logger, retrieved)

	// Test global convenience functions
	var buf bytes.Buffer
	logger.output = log.New(&buf, "", 0)

	Info("global info message")
	assert.Contains(t, buf.String(), "global info message")
}

// Benchmark tests
func BenchmarkLogger_Info(b *testing.B) {
	config := Config{
		Level:       InfoLevel,
		ServiceName: "bench-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&bytes.Buffer{}, "", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", map[string]interface{}{
			"iteration": i,
			"benchmark": true,
		})
	}
}

func BenchmarkLogger_WithFields(b *testing.B) {
	config := Config{
		Level:       InfoLevel,
		ServiceName: "bench-service",
		Version:     "1.0.0",
		Format:      "json",
	}

	logger := NewLogger(config)
	logger.output = log.New(&bytes.Buffer{}, "", 0)

	fields := map[string]interface{}{
		"user_id": "123",
		"action":  "test",
		"value":   42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		contextLogger := logger.WithFields(fields)
		contextLogger.Info("benchmark message")
	}
}