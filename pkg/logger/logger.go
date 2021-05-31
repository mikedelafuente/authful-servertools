package logger

import (
	"context"
	"log"

	"github.com/mikedelafuente/authful-servertools/pkg/config"
)

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

// What a Terrible Failure: Report a condition that should never happen. The error will
// always be logged at level Fatal/Failure with the call stack. Depending on system
// configuration, a report may be sent to the SDK developer and/or the process may be
// terminated immediately with an error dialog.
// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "FATAL" to see these logs.
func Fatal(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogFatal {
		log.Printf("[AUTHFUL FATAL] %v \n", v...)
	}
}

func Println(v ...interface{}) {
	log.Println(v...)
}

// Use this when you suspect something shady is going on. You may not be completely in full
// on error mode, but maybe you recovered from some unexpected behavior. Basically, use
// this to log stuff you didn't expect to happen but isn't necessarily an error. Kind of
// like a "hey, this happened, and it's weird, we should look into it."
// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "WARN" to see these logs.
func Warn(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogWarn {
		log.Printf("[AUTHFUL WARN] %v \n", v...)
	}

}

// This is for when bad stuff happens. Use this tag in places like inside a catch
// statement. You know that an error has occurred and therefore you're logging an error.
// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "ERROR" to see these logs.
func Error(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogError {
		log.Printf("[AUTHFUL ERROR] %v \n", v...)
	}
}

// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "DEBUG" to see these logs.
func Debug(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogDebug {
		log.Printf("[AUTHFUL DEBUG] %v \n", v...)
	}
}

// Use this to post useful information to the log. For example: that you have successfully connected to a server.
// Basically use it to report successes.
// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "INFO" to see these logs.
func Info(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogInfo {
		log.Printf("[AUTHFUL INFO] %v \n", v...)
	}
}

// Use this when you want to go absolutely nuts with your logging. If for some reason
// you've decided to log every little thing in a particular part of your app, use the
// Verbose tag.
// Set environmental variable 'AUTHFUL_LOG_LEVEL' to at least "VERBOSE" or "ALL" to see these logs.
func Verbose(ctx context.Context, v ...interface{}) {
	if config.GetConfig().LogVerbose {
		log.Printf("[AUTHFUL VERBOSE] %v \n", v...)
	}
}

// If no value is set in the environmental variable "AUTHFUL_LOG_LEVEL" then "ERROR" is returned. Returns the same output as config.GetConfig().GetLogLevel()
func GetLogLevel() string {
	return config.GetConfig().GetLogLevel()
}
