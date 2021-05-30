package config

type AuthfulConfig struct {
	LogFatal   bool `json:"log_fatal"`
	LogError   bool `json:"log_error"`
	LogWarn    bool `json:"log_warn"`
	LogInfo    bool `json:"log_info"`
	LogDebug   bool `json:"log_debug"`
	LogVerbose bool `json:"log_verbose"`
	logLevel   string
}
