package prometheus

import "time"

// Config represents the prometheus plugin configuration
type Config struct {
	// Enabled controls whether enhanced metrics are collected
	Enabled bool `mapstructure:"enabled"`

	// EndpointPatterns configures endpoint pattern matching and grouping
	EndpointPatterns EndpointPatternsConfig `mapstructure:"endpoint_patterns"`

	// CollectSizes enables request/response body size tracking
	CollectSizes bool `mapstructure:"collect_sizes"`

	// CollectQueueTime enables queue time vs processing time breakdown
	CollectQueueTime bool `mapstructure:"collect_queue_time"`

	// CollectWorkerInfo enables worker pool health metrics
	CollectWorkerInfo bool `mapstructure:"collect_worker_info"`

	// DurationBuckets defines histogram buckets for duration metrics (in seconds)
	DurationBuckets []float64 `mapstructure:"duration_buckets"`

	// SizeBuckets defines histogram buckets for size metrics (in bytes)
	SizeBuckets []float64 `mapstructure:"size_buckets"`
}

// EndpointPatternsConfig configures endpoint pattern matching
type EndpointPatternsConfig struct {
	// Enabled controls whether endpoint pattern matching is active
	Enabled bool `mapstructure:"enabled"`

	// MaxPatterns limits the number of unique endpoint patterns to prevent cardinality explosion
	MaxPatterns int `mapstructure:"max_patterns"`

	// Rules defines regex patterns for grouping similar endpoints
	Rules []PatternRule `mapstructure:"rules"`

	// CacheSize defines the LRU cache size for pattern matching results
	CacheSize int `mapstructure:"cache_size"`
}

// PatternRule defines a single pattern matching rule
type PatternRule struct {
	// Pattern is the regex pattern to match against request paths
	Pattern string `mapstructure:"pattern"`

	// Name is the label value to use when pattern matches
	Name string `mapstructure:"name"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled: true,
		EndpointPatterns: EndpointPatternsConfig{
			Enabled:     true,
			MaxPatterns: 100,
			Rules:       []PatternRule{},
			CacheSize:   10000,
		},
		CollectSizes:      true,
		CollectQueueTime:  true,
		CollectWorkerInfo: false, // Disabled by default as it requires HTTP plugin integration
		DurationBuckets: []float64{
			0.001, // 1ms
			0.005, // 5ms
			0.01,  // 10ms
			0.05,  // 50ms
			0.1,   // 100ms
			0.5,   // 500ms
			1.0,   // 1s
			5.0,   // 5s
			10.0,  // 10s
		},
		SizeBuckets: []float64{
			1024,     // 1KB
			10240,    // 10KB
			102400,   // 100KB
			1048576,  // 1MB
			10485760, // 10MB
		},
	}
}

// WorkerPoolStats represents worker pool statistics
// This interface allows decoupling from specific HTTP plugin implementation
type WorkerPoolStats struct {
	Active      int
	Idle        int
	Total       int
	UpdatedAt   time.Time
	Utilization float64 // Percentage 0-100
}
