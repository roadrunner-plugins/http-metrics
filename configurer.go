package prometheus

// Configurer interface for reading plugin configuration
type Configurer interface {
	// UnmarshalKey reads configuration section into provided structure
	UnmarshalKey(name string, out interface{}) error
	// Has checks if configuration section exists
	Has(name string) bool
}

// configKey is the configuration section name for this plugin
const configKey = "http_metrics"
