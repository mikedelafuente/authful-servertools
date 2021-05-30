package config

import (
	"log"
	"os"
	"strings"
	"sync"
)

var configOnce sync.Once
var configInstance *AuthfulConfig

func GetConfig() *AuthfulConfig {
	configOnce.Do(func() {
		var err error

		configInstance, err = getConfigInstanceFromEnvironment()

		if err != nil {
			log.Println(err)
			panic(err)
		}
	})

	return configInstance
}

func (s AuthfulConfig) GetLogLevel() string {
	return s.logLevel
}

// If no value is set in the environmental variable "AUTHFUL_LOG_LEVEL" then "ERROR" is returned
func getLogLevel() string {
	logLevel := os.Getenv("AUTHFUL_LOG_LEVEL")
	if len(logLevel) == 0 {
		logLevel = "ERROR"
	}

	return strings.TrimSpace(strings.ToUpper(logLevel))
}

func getConfigInstanceFromEnvironment() (*AuthfulConfig, error) {
	log.Println("Loading config from environment")

	var myConfig *AuthfulConfig = &AuthfulConfig{
		logLevel: getLogLevel(),
	}

	switch logLevel := myConfig.logLevel; {
	case logLevel == "VERBOSE" || logLevel == "ALL":
		myConfig.LogVerbose = true
		fallthrough
	case logLevel == "DEBUG":
		myConfig.LogDebug = true
		fallthrough
	case logLevel == "INFO":
		myConfig.LogInfo = true
		fallthrough
	case logLevel == "WARN":
		myConfig.LogWarn = true
		fallthrough
	case logLevel == "ERROR":
		myConfig.LogError = true
	case logLevel == "FATAL":
		myConfig.LogFatal = true
	default: // SAME AS "OFF"
		// Do nothing
	}

	return myConfig, nil
}
