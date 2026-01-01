package main

import (
	"fmt"
	"os"
)

// Logger provides structured logging with colors
type Logger struct {
	verbose bool
}

// ANSI color codes
const (
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[0;34m"
	colorReset  = "\033[0m"
)

// NewLogger creates a new logger
func NewLogger(verbose bool) *Logger {
	return &Logger{verbose: verbose}
}

// Error prints an error message
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorRed+"❌ "+format+colorReset+"\n", args...)
}

// Warn prints a warning message
func (l *Logger) Warn(format string, args ...interface{}) {
	fmt.Printf(colorYellow+"⚠️  "+format+colorReset+"\n", args...)
}

// Info prints an info message
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Success prints a success message
func (l *Logger) Success(format string, args ...interface{}) {
	fmt.Printf(colorGreen+"✔ "+format+colorReset+"\n", args...)
}

// Step prints a step completion message
func (l *Logger) Step(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s  "+colorGreen+"✔"+colorReset+"\n", msg)
}

// Debug prints a debug message (only if verbose)
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.verbose {
		fmt.Printf(colorBlue+"[DEBUG] "+format+colorReset+"\n", args...)
	}
}

// Section prints a section header
func (l *Logger) Section(title string) {
	fmt.Println()
	fmt.Println(colorBlue + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + colorReset)
	fmt.Printf(colorBlue+"  %s"+colorReset+"\n", title)
	fmt.Println(colorBlue + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" + colorReset)
}
