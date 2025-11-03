package prometheus

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	rrcontext "github.com/roadrunner-server/context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	jprop "go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	pluginName string = "http_metrics"
	namespace  string = "rr_http"

	// should be in sync with the http/handler.go constants
	noWorkers string = "No-Workers"
	trueStr   string = "true"
)

type Plugin struct {
	// Core components
	writersPool sync.Pool
	prop        propagation.TextMapPropagator
	stopCh      chan struct{}

	// Configuration
	cfg             Configurer
	config          *Config
	endpointMatcher *EndpointMatcher

	// Existing metrics
	queueSize       prometheus.Gauge
	noFreeWorkers   *prometheus.CounterVec
	requestCounter  *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	uptime          *prometheus.CounterVec

	// NEW: Phase 1 metrics - Performance breakdown
	queueTime      *prometheus.HistogramVec
	processingTime *prometheus.HistogramVec

	// NEW: Phase 1 metrics - Request/Response sizes
	requestSize  *prometheus.HistogramVec
	responseSize *prometheus.HistogramVec

	// NEW: Phase 1 metrics - Endpoint-level tracking
	requestsByEndpoint *prometheus.CounterVec
	durationByEndpoint *prometheus.HistogramVec

	// NEW: Phase 1 metrics - Error classification
	errorsByType *prometheus.CounterVec

	// NEW: Phase 1 metrics - Worker pool health
	activeWorkers     prometheus.Gauge
	idleWorkers       prometheus.Gauge
	workerUtilization prometheus.Gauge
}

func (p *Plugin) Init(cfg Configurer) error {
	// Initialize default configuration
	p.config = DefaultConfig()

	// Store configurer for potential future use
	p.cfg = cfg

	// Try to load configuration from http_metrics section
	if cfg != nil && cfg.Has(configKey) {
		if err := cfg.UnmarshalKey(configKey, p.config); err != nil {
			return err
		}
	}

	// Initialize writers pool
	p.writersPool = sync.Pool{
		New: func() any {
			wr := new(writer)
			wr.code = -1
			return wr
		},
	}

	p.stopCh = make(chan struct{}, 1)

	// Initialize endpoint matcher
	var err error
	p.endpointMatcher, err = NewEndpointMatcher(p.config.EndpointPatterns)
	if err != nil {
		return err
	}

	// Initialize existing metrics
	p.queueSize = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "requests_queue",
		Help:      "Total number of queued requests.",
	})

	p.noFreeWorkers = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "no_free_workers_total",
		Help:      "Total number of NoFreeWorkers occurrences.",
	}, nil)

	p.requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "request_total",
		Help:      "Total number of handled http requests after server restart.",
	}, []string{"status"})

	p.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "request_duration_seconds",
			Help:      "HTTP request duration.",
			Buckets:   p.config.DurationBuckets,
		},
		[]string{"status"},
	)

	p.uptime = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "uptime_seconds",
		Help:      "Uptime in seconds",
	}, nil)

	// Initialize NEW metrics - Performance breakdown
	if p.config.CollectQueueTime {
		p.queueTime = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "queue_time_seconds",
				Help:      "Time request spent waiting in queue before being picked up by a worker.",
				Buckets:   p.config.DurationBuckets,
			},
			[]string{"method", "endpoint"},
		)

		p.processingTime = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "processing_time_seconds",
				Help:      "Time spent processing the request by PHP worker.",
				Buckets:   p.config.DurationBuckets,
			},
			[]string{"method", "endpoint"},
		)
	}

	// Initialize NEW metrics - Request/Response sizes
	if p.config.CollectSizes {
		p.requestSize = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "request_size_bytes",
				Help:      "HTTP request body size in bytes.",
				Buckets:   p.config.SizeBuckets,
			},
			[]string{"method", "endpoint"},
		)

		p.responseSize = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "response_size_bytes",
				Help:      "HTTP response body size in bytes.",
				Buckets:   p.config.SizeBuckets,
			},
			[]string{"method", "endpoint", "status"},
		)
	}

	// Initialize NEW metrics - Endpoint-level tracking
	p.requestsByEndpoint = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "requests_by_endpoint_total",
			Help:      "Total number of HTTP requests by endpoint pattern.",
		},
		[]string{"method", "endpoint", "status"},
	)

	p.durationByEndpoint = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "duration_by_endpoint_seconds",
			Help:      "HTTP request duration by endpoint pattern.",
			Buckets:   p.config.DurationBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Initialize NEW metrics - Error classification
	p.errorsByType = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "errors_total",
			Help:      "Total number of HTTP errors classified by type.",
		},
		[]string{"type", "endpoint", "status"},
	)

	// Initialize NEW metrics - Worker pool health
	if p.config.CollectWorkerInfo {
		p.activeWorkers = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "active_workers",
			Help:      "Number of workers currently processing requests.",
		})

		p.idleWorkers = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "idle_workers",
			Help:      "Number of idle workers available to process requests.",
		})

		p.workerUtilization = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "worker_utilization_percent",
			Help:      "Worker pool utilization percentage (0-100).",
		})
	}

	p.prop = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, jprop.Jaeger{})

	return nil
}

func (p *Plugin) Serve() chan error {
	errCh := make(chan error, 1)

	// Existing uptime ticker
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-p.stopCh:
				return
			case <-ticker.C:
				p.uptime.With(nil).Inc()
			}
		}
	}()

	return errCh
}

func (p *Plugin) Stop(context.Context) error {
	close(p.stopCh)
	return nil
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle OpenTelemetry tracing (existing logic)
		if val, ok := r.Context().Value(rrcontext.OtelTracerNameKey).(string); ok {
			tp := trace.SpanFromContext(r.Context()).TracerProvider()
			ctx, span := tp.Tracer(val, trace.WithSchemaURL(semconv.SchemaURL),
				trace.WithInstrumentationVersion(otelhttp.Version())).
				Start(r.Context(), pluginName, trace.WithSpanKind(trace.SpanKindServer))
			defer span.End()

			// inject
			p.prop.Inject(ctx, propagation.HeaderCarrier(r.Header))
			r = r.WithContext(ctx)
		}

		// Record arrival time
		arrivalTime := time.Now()

		// Get writer from pool and initialize timing
		rrWriter := p.getWriter(w)
		defer p.putWriter(rrWriter)

		rrWriter.arrivalTime = arrivalTime
		rrWriter.queueStart = time.Now()
		rrWriter.requestSize = r.ContentLength

		// Track queue size
		p.queueSize.Inc()

		// Processing starts when worker picks up request
		rrWriter.processStart = time.Now()

		// Execute request
		next.ServeHTTP(rrWriter, r)

		processEnd := time.Now()

		// Extract request metadata
		endpoint := "other"
		if p.config.EndpointPatterns.Enabled {
			endpoint = p.endpointMatcher.Match(r.URL.Path)
		}
		method := r.Method
		status := strconv.Itoa(rrWriter.code)

		// Calculate timings
		queueTime := rrWriter.processStart.Sub(rrWriter.queueStart)
		processingTime := processEnd.Sub(rrWriter.processStart)
		totalTime := processEnd.Sub(rrWriter.arrivalTime)

		// Create label sets for metrics
		endpointLabels := prometheus.Labels{
			"method":   method,
			"endpoint": endpoint,
		}

		fullLabels := prometheus.Labels{
			"method":   method,
			"endpoint": endpoint,
			"status":   status,
		}

		// Record existing metrics
		p.requestCounter.With(prometheus.Labels{"status": status}).Inc()
		p.requestDuration.With(prometheus.Labels{"status": status}).Observe(totalTime.Seconds())

		// Record NEW metrics - Performance breakdown
		if p.config.CollectQueueTime {
			p.queueTime.With(endpointLabels).Observe(queueTime.Seconds())
			p.processingTime.With(endpointLabels).Observe(processingTime.Seconds())
		}

		// Record NEW metrics - Request/Response sizes
		if p.config.CollectSizes {
			if rrWriter.requestSize > 0 {
				p.requestSize.With(endpointLabels).Observe(float64(rrWriter.requestSize))
			}
			if rrWriter.bytesWritten > 0 {
				p.responseSize.With(fullLabels).Observe(float64(rrWriter.bytesWritten))
			}
		}

		// Record NEW metrics - Endpoint-level tracking
		p.requestsByEndpoint.With(fullLabels).Inc()
		p.durationByEndpoint.With(endpointLabels).Observe(totalTime.Seconds())

		// Record NEW metrics - Error classification
		if isErrorStatus(rrWriter.code) {
			errorType := string(classifyError(rrWriter.code, w.Header()))
			p.errorsByType.With(prometheus.Labels{
				"type":     errorType,
				"endpoint": endpoint,
				"status":   status,
			}).Inc()
		}

		// Handle no workers case (existing logic)
		if w.Header().Get(noWorkers) == trueStr {
			p.noFreeWorkers.With(nil).Inc()
		}

		p.queueSize.Dec()
	})
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) MetricsCollector() []prometheus.Collector {
	collectors := []prometheus.Collector{
		// Existing metrics
		p.requestCounter,
		p.requestDuration,
		p.queueSize,
		p.noFreeWorkers,
		p.uptime,
		// NEW: Endpoint-level metrics (always enabled)
		p.requestsByEndpoint,
		p.durationByEndpoint,
		p.errorsByType,
	}

	// Add conditional metrics
	if p.config.CollectQueueTime {
		collectors = append(collectors, p.queueTime, p.processingTime)
	}

	if p.config.CollectSizes {
		collectors = append(collectors, p.requestSize, p.responseSize)
	}

	if p.config.CollectWorkerInfo {
		collectors = append(collectors, p.activeWorkers, p.idleWorkers, p.workerUtilization)
	}

	return collectors
}

func (p *Plugin) getWriter(w http.ResponseWriter) *writer {
	wr := p.writersPool.Get().(*writer)
	wr.w = w
	return wr
}

func (p *Plugin) putWriter(w *writer) {
	w.reset()
	p.writersPool.Put(w)
}

// updateWorkerMetrics queries worker pool state and updates metrics
// NOTE: This requires integration with HTTP plugin's worker pool
// Implementation depends on RoadRunner's internal API for accessing worker pool state
func (p *Plugin) updateWorkerMetrics() {
	// TODO: Implement worker pool state access
	// This requires extending the HTTP plugin to expose worker pool statistics
	// or using RoadRunner's internal interfaces to query worker state
	//
	// Example pseudo-code:
	// stats := p.getWorkerPoolStats()
	// if stats != nil {
	//     p.activeWorkers.Set(float64(stats.Active))
	//     p.idleWorkers.Set(float64(stats.Idle))
	//     p.workerUtilization.Set(stats.Utilization)
	// }
}
