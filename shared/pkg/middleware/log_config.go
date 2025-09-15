package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// LogConfig holds configuration for logging
type LogConfig struct {
	EnableColors    bool
	EnableTimestamp bool
	EnableRequestID bool
	LogLevel        LogLevel
	OutputFormat    string // "text" or "json"
}

// DefaultLogConfig returns default logging configuration
func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		EnableColors:    true,
		EnableTimestamp: true,
		EnableRequestID: true,
		LogLevel:        INFO,
		OutputFormat:    "text",
	}
}

// LoadLogConfigFromEnv loads logging configuration from environment variables
func LoadLogConfigFromEnv() *LogConfig {
	config := DefaultLogConfig()
	
	if os.Getenv("LOG_NO_COLOR") == "true" {
		config.EnableColors = false
	}
	
	if os.Getenv("LOG_NO_TIMESTAMP") == "true" {
		config.EnableTimestamp = false
	}
	
	if os.Getenv("LOG_NO_REQUEST_ID") == "true" {
		config.EnableRequestID = false
	}
	
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.LogLevel = LogLevel(strings.ToUpper(level))
	}
	
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		config.OutputFormat = format
	}
	
	return config
}

// GlobalLogConfig holds the global logging configuration
var GlobalLogConfig = LoadLogConfigFromEnv()

// LogWithConfig logs with specific configuration
func LogWithConfig(entry LogEntry, config *LogConfig) {
	if config == nil {
		config = GlobalLogConfig
	}
	
	// Check if we should log this level
	if !shouldLogLevel(entry.Level, config.LogLevel) {
		return
	}
	
	if config.OutputFormat == "json" {
		logJSON(entry)
	} else {
		logText(entry, config)
	}
}

// shouldLogLevel checks if the entry level should be logged based on config
func shouldLogLevel(entryLevel, configLevel LogLevel) bool {
	levels := map[LogLevel]int{
		DEBUG:   0,
		INFO:    1,
		WARN:    2,
		ERROR:   3,
		SUCCESS: 1,
	}
	
	entryPriority := levels[entryLevel]
	configPriority := levels[configLevel]
	
	return entryPriority >= configPriority
}

// logJSON logs in JSON format
func logJSON(entry LogEntry) {
	// Remove timestamp if disabled
	if !GlobalLogConfig.EnableTimestamp {
		entry.Timestamp = ""
	}
	
	// Remove request ID if disabled
	if !GlobalLogConfig.EnableRequestID {
		entry.RequestID = ""
	}
	
	jsonData, _ := json.Marshal(entry)
	fmt.Println(string(jsonData))
}

// logText logs in formatted text with colors
func logText(entry LogEntry, config *LogConfig) {
	if !config.EnableColors || os.Getenv("TERM") == "" || os.Getenv("NO_COLOR") != "" {
		logPlainText(entry, config)
		return
	}
	
	logColoredText(entry, config)
}

// logPlainText logs without colors
func logPlainText(entry LogEntry, config *LogConfig) {
	var parts []string
	
	if config.EnableTimestamp && entry.Timestamp != "" {
		parts = append(parts, fmt.Sprintf("[%s]", entry.Timestamp))
	}
	
	parts = append(parts, fmt.Sprintf("[%s]", entry.Level))
	
	if config.EnableRequestID && entry.RequestID != "" {
		parts = append(parts, fmt.Sprintf("[%s]", entry.RequestID))
	}
	
	if entry.Method != "" && entry.Path != "" {
		parts = append(parts, fmt.Sprintf("%s %s", entry.Method, entry.Path))
	}
	
	if entry.StatusCode > 0 {
		parts = append(parts, fmt.Sprintf("[%d]", entry.StatusCode))
	}
	
	if entry.Duration != "" {
		parts = append(parts, fmt.Sprintf("(%s)", entry.Duration))
	}
	
	if entry.UserID != nil {
		parts = append(parts, fmt.Sprintf("User:%v", entry.UserID))
	}
	
	parts = append(parts, entry.Message)
	
	fmt.Println(strings.Join(parts, " "))
}

// logColoredText logs with colors
func logColoredText(entry LogEntry, config *LogConfig) {
	color := getColorByLevel(entry.Level)
	levelStr := fmt.Sprintf("%s%s%s", color, entry.Level, colorReset)
	
	var message strings.Builder
	
	// Add timestamp
	if config.EnableTimestamp && entry.Timestamp != "" {
		message.WriteString(fmt.Sprintf("%s[%s]%s ", colorBold, entry.Timestamp, colorReset))
	}
	
	// Add level
	message.WriteString(levelStr)
	message.WriteString(" ")
	
	// Add request ID
	if config.EnableRequestID && entry.RequestID != "" {
		message.WriteString(fmt.Sprintf("%s[%s]%s ", colorPurple, entry.RequestID, colorReset))
	}
	
	// Add method and path
	if entry.Method != "" && entry.Path != "" {
		message.WriteString(fmt.Sprintf("%s%s %s%s ", colorCyan, entry.Method, entry.Path, colorReset))
	}
	
	// Add status code
	if entry.StatusCode > 0 {
		statusColor := colorGreen
		if entry.StatusCode >= 400 {
			statusColor = colorRed
		} else if entry.StatusCode >= 300 {
			statusColor = colorYellow
		}
		message.WriteString(fmt.Sprintf("%s[%d]%s ", statusColor, entry.StatusCode, colorReset))
	}
	
	// Add duration
	if entry.Duration != "" {
		message.WriteString(fmt.Sprintf("%s(%s)%s ", colorYellow, entry.Duration, colorReset))
	}
	
	// Add user info
	if entry.UserID != nil {
		message.WriteString(fmt.Sprintf("%sUser:%v%s ", colorGreen, entry.UserID, colorReset))
	}
	
	// Add message
	message.WriteString(entry.Message)
	
	fmt.Println(message.String())
}

// LogSuccess logs a success message
func LogSuccess(message string, extra ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05.000"),
		Level:     SUCCESS,
		Message:   message,
	}
	
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	
	LogWithConfig(entry, GlobalLogConfig)
}

// LogInfo logs an info message
func LogInfo(message string, extra ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05.000"),
		Level:     INFO,
		Message:   message,
	}
	
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	
	LogWithConfig(entry, GlobalLogConfig)
}

// LogWarn logs a warning message
func LogWarn(message string, extra ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05.000"),
		Level:     WARN,
		Message:   message,
	}
	
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	
	LogWithConfig(entry, GlobalLogConfig)
}

// LogError logs an error message
func LogError(message string, err error, extra ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05.000"),
		Level:     ERROR,
		Message:   message,
	}
	
	if err != nil {
		if entry.Extra == nil {
			entry.Extra = make(map[string]interface{})
		}
		entry.Extra["error"] = err.Error()
	}
	
	if len(extra) > 0 {
		if entry.Extra == nil {
			entry.Extra = make(map[string]interface{})
		}
		for k, v := range extra[0] {
			entry.Extra[k] = v
		}
	}
	
	LogWithConfig(entry, GlobalLogConfig)
}

// LogDebug logs a debug message
func LogDebug(message string, extra ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05.000"),
		Level:     DEBUG,
		Message:   message,
	}
	
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	
	LogWithConfig(entry, GlobalLogConfig)
}
