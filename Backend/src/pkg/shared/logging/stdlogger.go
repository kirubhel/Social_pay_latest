package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[36m"
	colorGray   = "\033[37m"
	colorBold   = "\033[1m"
)

// LogLevel represents logging levels
type LogLevel struct {
	Name  string
	Color string
}

var (
	debugLevel = LogLevel{Name: "DEBUG", Color: colorGray}
	infoLevel  = LogLevel{Name: "INFO", Color: colorBlue}
	warnLevel  = LogLevel{Name: "WARN", Color: colorYellow}
	errorLevel = LogLevel{Name: "ERROR", Color: colorRed}
	fatalLevel = LogLevel{Name: "FATAL", Color: colorRed}
)

// StdLogger wraps the standard logger to implement the Logger interface
type StdLogger struct {
	mu      sync.Mutex
	logger  *log.Logger
	prefix  string
	useJSON bool
}

// NewStdLogger creates a new StdLogger
func NewStdLogger(prefix string) Logger {
	// Remove all time-related flags, we'll handle timestamp formatting ourselves
	return &StdLogger{
		logger:  log.New(os.Stdout, "", 0), // Changed from log.Ldate | log.Ltime | log.Lmicroseconds
		prefix:  prefix,
		useJSON: false,
	}
}

// formatFields formats the fields map in a more readable way
func formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	// Get sorted keys for consistent output
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build the fields string with proper indentation
	var pairs []string
	maxKeyLength := 0

	// First pass: determine max key length for alignment
	for _, k := range keys {
		if len(k) > maxKeyLength {
			maxKeyLength = len(k)
		}
	}

	// Second pass: format values with proper alignment
	for _, k := range keys {
		v := fields[k]
		// Handle different types of values
		var valueStr string
		switch val := v.(type) {
		case error:
			valueStr = fmt.Sprintf("%q", val.Error())
		case string:
			valueStr = fmt.Sprintf("%q", val)
		case nil:
			valueStr = "null"
		case map[string]interface{}:
			// Pretty format nested maps
			b, _ := json.MarshalIndent(val, "    ", "    ")
			valueStr = string(b)
		case []interface{}:
			// Pretty format arrays
			b, _ := json.MarshalIndent(val, "    ", "    ")
			valueStr = string(b)
		default:
			// Use JSON marshaling for complex types
			b, err := json.MarshalIndent(val, "    ", "    ")
			if err != nil {
				valueStr = fmt.Sprintf("%v", val)
			} else {
				valueStr = string(b)
			}
		}

		// Format key-value pair with proper alignment
		formattedPair := fmt.Sprintf(
			"%-*s %s %s",
			maxKeyLength,
			colorBold+k+colorReset,
			colorGray+"="+colorReset,
			valueStr,
		)
		pairs = append(pairs, formattedPair)
	}

	// Join all pairs with newlines and proper indentation
	return "\n    " + strings.Join(pairs, "\n    ")
}

// formatLog formats the log message with timestamp, level, and fields
func (l *StdLogger) formatLog(level LogLevel, msg string, fields map[string]interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	if l.useJSON {
		logData := map[string]interface{}{
			"timestamp": timestamp,
			"level":     level.Name,
			"message":   msg,
			"prefix":    l.prefix,
		}
		if len(fields) > 0 {
			logData["fields"] = fields
		}
		jsonBytes, _ := json.MarshalIndent(logData, "", "    ")
		return string(jsonBytes)
	}

	// Format with colors and proper spacing
	var builder strings.Builder

	// Timestamp
	builder.WriteString(colorGray)
	builder.WriteString(timestamp)
	builder.WriteString(colorReset)
	builder.WriteString(" ")

	// Level
	builder.WriteString(level.Color)
	builder.WriteString(fmt.Sprintf("[%-5s]", level.Name))
	builder.WriteString(colorReset)
	builder.WriteString(" ")

	// Prefix
	if l.prefix != "" {
		builder.WriteString(colorBold)
		builder.WriteString(l.prefix)
		builder.WriteString(colorReset)
		builder.WriteString(" ")
	}

	// Message
	builder.WriteString(msg)

	// Fields
	if formattedFields := formatFields(fields); formattedFields != "" {
		builder.WriteString(formattedFields)
	}

	return builder.String()
}

// Debug logs a debug message
func (l *StdLogger) Debug(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, l.formatLog(debugLevel, msg, fields))
}

// Info logs an info message
func (l *StdLogger) Info(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, l.formatLog(infoLevel, msg, fields))
}

// Warn logs a warning message
func (l *StdLogger) Warn(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, l.formatLog(warnLevel, msg, fields))
}

// Error logs an error message
func (l *StdLogger) Error(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, l.formatLog(errorLevel, msg, fields))
}

// Fatal logs a fatal message and exits
func (l *StdLogger) Fatal(msg string, fields map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Output(2, l.formatLog(fatalLevel, msg, fields))
	os.Exit(1)
}

// SetJSON enables or disables JSON formatting
func (l *StdLogger) SetJSON(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.useJSON = enabled
}
