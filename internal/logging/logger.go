package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity level of log messages
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	Version     string                 `json:"version"`
	Method      string                 `json:"method,omitempty"`
	URL         string                 `json:"url,omitempty"`
	StatusCode  int                    `json:"status_code,omitempty"`
	Duration    int64                  `json:"duration_ms,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	RemoteAddr  string                 `json:"remote_addr,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StackTrace  string                 `json:"stack_trace,omitempty"`
	Caller      string                 `json:"caller,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// Logger represents the structured logger
type Logger struct {
	level       LogLevel
	serviceName string
	version     string
	output      *log.Logger
	fields      map[string]interface{}
}

// Config represents logger configuration
type Config struct {
	Level       LogLevel
	ServiceName string
	Version     string
	Format      string // "json" or "text"
}

// NewLogger creates a new structured logger
func NewLogger(config Config) *Logger {
	output := log.New(os.Stdout, "", 0)
	
	return &Logger{
		level:       config.Level,
		serviceName: config.ServiceName,
		version:     config.Version,
		output:      output,
		fields:      make(map[string]interface{}),
	}
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:       l.level,
		serviceName: l.serviceName,
		version:     l.version,
		output:      l.output,
		fields:      make(map[string]interface{}),
	}
	
	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	
	// Add new field
	newLogger.fields[key] = value
	
	return newLogger
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:       l.level,
		serviceName: l.serviceName,
		version:     l.version,
		output:      l.output,
		fields:      make(map[string]interface{}),
	}
	
	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	
	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	
	return newLogger
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, message, fields...)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, message, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, message, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, err error, fields ...map[string]interface{}) {
	if l.level <= ErrorLevel {
		entry := l.createLogEntry(ErrorLevel, message, fields...)
		if err != nil {
			entry.Error = err.Error()
			entry.StackTrace = getStackTrace(2) // Skip this function and log()
		}
		l.writeLog(entry)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, err error, fields ...map[string]interface{}) {
	entry := l.createLogEntry(FatalLevel, message, fields...)
	if err != nil {
		entry.Error = err.Error()
		entry.StackTrace = getStackTrace(2)
	}
	l.writeLog(entry)
	os.Exit(1)
}

// HTTPRequest logs HTTP request information
func (l *Logger) HTTPRequest(method, url string, statusCode int, duration time.Duration, userAgent, remoteAddr string, fields ...map[string]interface{}) {
	entry := l.createLogEntry(InfoLevel, "HTTP Request", fields...)
	entry.Method = method
	entry.URL = url
	entry.StatusCode = statusCode
	entry.Duration = duration.Milliseconds()
	entry.UserAgent = userAgent
	entry.RemoteAddr = remoteAddr
	
	l.writeLog(entry)
}

// DatabaseQuery logs database query information
func (l *Logger) DatabaseQuery(query string, duration time.Duration, rowsAffected int64, fields ...map[string]interface{}) {
	entry := l.createLogEntry(DebugLevel, "Database Query", fields...)
	entry.Fields["query"] = truncateQuery(query)
	entry.Fields["duration_ms"] = duration.Milliseconds()
	entry.Fields["rows_affected"] = rowsAffected
	
	l.writeLog(entry)
}

// CacheOperation logs cache operations
func (l *Logger) CacheOperation(operation, key string, hit bool, duration time.Duration, fields ...map[string]interface{}) {
	entry := l.createLogEntry(DebugLevel, "Cache Operation", fields...)
	entry.Fields["operation"] = operation
	entry.Fields["key"] = key
	entry.Fields["hit"] = hit
	entry.Fields["duration_ms"] = duration.Milliseconds()
	
	l.writeLog(entry)
}

// SecurityEvent logs security-related events
func (l *Logger) SecurityEvent(event, userID, details string, fields ...map[string]interface{}) {
	entry := l.createLogEntry(WarnLevel, "Security Event", fields...)
	entry.UserID = userID
	entry.Fields["event"] = event
	entry.Fields["details"] = details
	
	l.writeLog(entry)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, fields ...map[string]interface{}) {
	entry := l.createLogEntry(level, message, fields...)
	l.writeLog(entry)
}

// createLogEntry creates a new log entry with common fields
func (l *Logger) createLogEntry(level LogLevel, message string, fields ...map[string]interface{}) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level.String(),
		Message:   message,
		Service:   l.serviceName,
		Version:   l.version,
		Caller:    getCaller(3), // Skip createLogEntry, log, and the public method
		Fields:    make(map[string]interface{}),
	}
	
	// Add logger context fields
	for k, v := range l.fields {
		entry.Fields[k] = v
	}
	
	// Add additional fields from method call
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			entry.Fields[k] = v
		}
	}
	
	return entry
}

// writeLog outputs the log entry
func (l *Logger) writeLog(entry *LogEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to standard logging if JSON marshaling fails
		l.output.Printf("ERROR: Failed to marshal log entry: %v", err)
		return
	}
	
	l.output.Printf("%s", string(jsonData))
}

// getCaller returns the caller function name and line
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	
	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]
	
	// Get just the function name, not the full package path
	funcName := fn.Name()
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}
	
	return fmt.Sprintf("%s:%d:%s", filename, line, funcName)
}

// getStackTrace returns a formatted stack trace
func getStackTrace(skip int) string {
	var traces []string
	
	for i := skip; i < skip+10; i++ { // Limit to 10 frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		
		// Get just the filename
		parts := strings.Split(file, "/")
		filename := parts[len(parts)-1]
		
		// Get just the function name
		funcName := fn.Name()
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}
		
		traces = append(traces, fmt.Sprintf("%s:%d:%s", filename, line, funcName))
	}
	
	return strings.Join(traces, " -> ")
}

// truncateQuery truncates long SQL queries for logging
func truncateQuery(query string) string {
	const maxLength = 200
	
	// Remove extra whitespace
	cleaned := strings.Join(strings.Fields(query), " ")
	
	if len(cleaned) <= maxLength {
		return cleaned
	}
	
	return cleaned[:maxLength] + "..."
}

// ParseLogLevel parses a string log level
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN", "WARNING":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Global logger instance
var globalLogger *Logger

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	return globalLogger
}

// Global convenience functions
func Debug(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(message, fields...)
	}
}

func Info(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Info(message, fields...)
	}
}

func Warn(message string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(message, fields...)
	}
}

func Error(message string, err error, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Error(message, err, fields...)
	}
}

func Fatal(message string, err error, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(message, err, fields...)
	}
}

func HTTPRequest(method, url string, statusCode int, duration time.Duration, userAgent, remoteAddr string, fields ...map[string]interface{}) {
	if globalLogger != nil {
		globalLogger.HTTPRequest(method, url, statusCode, duration, userAgent, remoteAddr, fields...)
	}
}