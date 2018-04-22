package config

// Config encapsulates the configuration for the process.
type Config struct {
	Verbose bool
}

var instance *Config

// Get gets the configuration instance.
func Get() *Config {
	if instance == nil {
		instance = &Config{}
	}
	return instance
}
