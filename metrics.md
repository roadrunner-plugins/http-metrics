# Grafana Queries - Complete Guide with Configuration

This guide provides ready-to-use PromQL queries for all Enhanced HTTP Metrics, including panel configuration details (
legend, min step, units).

---

## 1. Traffic Overview Metrics

### 1.1 Total Requests Per Second (RPS)

**Query:**

```promql
sum(rate(rr_http_requests_by_endpoint_total[5m]))
```

**Configuration:**

- **Legend:** `Total RPS`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Graph
- **Description:** Overall request rate across all endpoints

---

### 1.2 RPS by Endpoint

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Graph or Bar gauge
- **Description:** Top 10 endpoints by request rate

---

### 1.3 RPS by HTTP Method

**Query:**

```promql
sum by (method) (rate(rr_http_requests_by_endpoint_total[5m]))
```

**Configuration:**

- **Legend:** `{{method}}`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Graph or Pie chart
- **Description:** Request distribution by HTTP method (GET, POST, etc.)

---

### 1.4 RPS by Status Code

**Query:**

```promql
sum by (status) (rate(rr_http_requests_by_endpoint_total[5m]))
```

**Configuration:**

- **Legend:** `Status {{status}}`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Graph (stacked area)
- **Description:** Request rate grouped by status code

---

### 1.5 Success Rate (2xx responses)

**Query:**

```promql
sum(rate(rr_http_requests_by_endpoint_total{status=~"2.."}[5m])) / sum(rate(rr_http_requests_by_endpoint_total[5m])) * 100
```

**Configuration:**

- **Legend:** `Success Rate`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Graph or Gauge
- **Thresholds:** Red < 95%, Yellow 95-99%, Green > 99%
- **Description:** Percentage of successful (2xx) requests

---

### 1.6 Current Queue Size

**Query:**

```promql
rr_http_requests_queue
```

**Configuration:**

- **Legend:** `Queue Size`
- **Min Step:** `5s`
- **Unit:** `short`
- **Panel Type:** Graph
- **Description:** Number of requests currently waiting in queue

---

## 2. Performance Metrics

### 2.1 Average Request Duration

**Query:**

```promql
rate(rr_http_request_duration_seconds_sum[5m]) / rate(rr_http_request_duration_seconds_count[5m])
```

**Configuration:**

- **Legend:** `Avg Duration`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** Mean request duration across all requests

---

### 2.2 Request Duration by Endpoint

**Query:**

```promql
sum by (endpoint) (rate(rr_http_request_duration_seconds_sum[5m])) / sum by (endpoint) (rate(rr_http_request_duration_seconds_count[5m]))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** Average duration per endpoint

---

### 2.3 P50 Latency (Median)

**Query:**

```promql
histogram_quantile(0.50, sum by (le) (rate(rr_http_request_duration_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `P50 (Median)`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** 50th percentile latency (median response time)

---

### 2.4 P95 Latency

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_request_duration_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `P95`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Thresholds:** Green < 1s, Yellow 1-5s, Red > 5s
- **Description:** 95th percentile latency (95% of requests complete faster)

---

### 2.5 P99 Latency

**Query:**

```promql
histogram_quantile(0.99, sum by (le) (rate(rr_http_request_duration_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `P99`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** 99th percentile latency (99% of requests complete faster)

---

### 2.6 P95 Latency by Endpoint

**Query:**

```promql
histogram_quantile(0.95, sum by (endpoint, le) (rate(rr_http_request_duration_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph or Table
- **Description:** 95th percentile per endpoint

---

### 2.7 Queue Time (Average)

**Query:**

```promql
rate(rr_http_queue_time_seconds_sum[5m]) / rate(rr_http_queue_time_seconds_count[5m])
```

**Configuration:**

- **Legend:** `Queue Time`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** Average time requests spend waiting in queue

---

### 2.8 Queue Time P95

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_queue_time_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `Queue Time P95`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Thresholds:** Green < 0.1s, Yellow 0.1-0.5s, Red > 0.5s
- **Description:** 95th percentile queue wait time

---

### 2.9 Processing Time (Average)

**Query:**

```promql
rate(rr_http_processing_time_seconds_sum[5m]) / rate(rr_http_processing_time_seconds_count[5m])
```

**Configuration:**

- **Legend:** `Processing Time`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** Average PHP worker processing time

---

### 2.10 Processing Time P95

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_processing_time_seconds_bucket[5m])))
```

**Configuration:**

- **Legend:** `Processing Time P95`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph
- **Description:** 95th percentile processing time

---

### 2.11 Performance Breakdown (Stacked)

**Query 1 (Queue Time):**

```promql
sum(rate(rr_http_queue_time_seconds_sum[5m])) / sum(rate(rr_http_queue_time_seconds_count[5m]))
```

**Query 2 (Processing Time):**

```promql
sum(rate(rr_http_processing_time_seconds_sum[5m])) / sum(rate(rr_http_processing_time_seconds_count[5m]))
```

**Configuration:**

- **Legend:**
    - Query 1: `Queue Time`
    - Query 2: `Processing Time`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph (stacked area)
- **Description:** Visual breakdown of where time is spent

---

### 2.12 Queue vs Processing Time Ratio

**Query:**

```promql
(rate(rr_http_queue_time_seconds_sum[5m]) / rate(rr_http_queue_time_seconds_count[5m])) / (rate(rr_http_request_duration_seconds_sum[5m]) / rate(rr_http_request_duration_seconds_count[5m])) * 100
```

**Configuration:**

- **Legend:** `Queue Time %`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Graph
- **Description:** Percentage of total time spent waiting in queue

---

## 3. Endpoint Analysis

### 3.1 Top 10 Slowest Endpoints (by P95)

**Query:**

```promql
topk(10, histogram_quantile(0.95, sum by (endpoint, le) (rate(rr_http_duration_by_endpoint_seconds_bucket[15m]))))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `30s`
- **Unit:** `seconds (s)`
- **Panel Type:** Bar gauge (horizontal) or Table
- **Description:** Slowest endpoints ranked by 95th percentile

---

### 3.2 Top 10 Slowest Endpoints (by Average)

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_sum[5m])) / sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_count[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `seconds (s)`
- **Panel Type:** Bar gauge or Table
- **Description:** Slowest endpoints by mean duration

---

### 3.3 Top 10 Most Trafficked Endpoints

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Bar gauge or Table
- **Description:** Endpoints with highest request rate

---

### 3.4 Endpoint Performance Table

**Query 1 (RPS):**

```promql
sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m]))
```

**Query 2 (Avg Duration):**

```promql
sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_sum[5m])) / sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_count[5m]))
```

**Query 3 (P95 Duration):**

```promql
histogram_quantile(0.95, sum by (endpoint, le) (rate(rr_http_duration_by_endpoint_seconds_bucket[5m])))
```

**Query 4 (Error %):**

```promql
sum by (endpoint) (rate(rr_http_errors_total[5m])) / sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m])) * 100
```

**Configuration:**

- **Legend:** N/A (Table columns)
- **Min Step:** `15s`
- **Unit:**
    - Query 1: `requests/sec (reqps)`
    - Query 2: `seconds (s)`
    - Query 3: `seconds (s)`
    - Query 4: `percent (0-100)`
- **Panel Type:** Table
- **Column Names:** `Endpoint`, `RPS`, `Avg Duration`, `P95 Duration`, `Error Rate %`
- **Description:** Comprehensive endpoint performance overview

---

### 3.5 Endpoint Duration Heatmap

**Query:**

```promql
sum by (le, endpoint) (rate(rr_http_duration_by_endpoint_seconds_bucket[5m]))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `30s`
- **Unit:** N/A (heatmap)
- **Panel Type:** Heatmap
- **Description:** Visual distribution of response times across endpoints

---

## 4. Error Tracking

### 4.1 Total Error Rate

**Query:**

```promql
sum(rate(rr_http_errors_total[5m]))
```

**Configuration:**

- **Legend:** `Errors/sec`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Description:** Total errors per second (all types)

---

### 4.2 Error Rate Percentage

**Query:**

```promql
sum(rate(rr_http_errors_total[5m])) / sum(rate(rr_http_requests_by_endpoint_total[5m])) * 100
```

**Configuration:**

- **Legend:** `Error Rate`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Graph or Gauge
- **Thresholds:** Green < 1%, Yellow 1-5%, Red > 5%
- **Description:** Percentage of requests that result in errors

---

### 4.3 Error Rate by Type

**Query:**

```promql
sum by (type) (rate(rr_http_errors_total[5m]))
```

**Configuration:**

- **Legend:** `{{type}}`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph (stacked) or Pie chart
- **Description:** Errors grouped by classification (client_error, server_error, timeout, no_workers)

---

### 4.4 Error Rate by Status Code

**Query:**

```promql
sum by (status) (rate(rr_http_errors_total[5m]))
```

**Configuration:**

- **Legend:** `HTTP {{status}}`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph or Table
- **Description:** Errors grouped by specific HTTP status code

---

### 4.5 Error Rate by Endpoint

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_errors_total[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Bar gauge or Table
- **Description:** Endpoints with highest error rate

---

### 4.6 Most Error-Prone Endpoints (Percentage)

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_errors_total[5m])) / sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m])) * 100)
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Bar gauge or Table
- **Description:** Endpoints with highest error percentage

---

### 4.7 4xx vs 5xx Errors

**Query 1 (Client Errors):**

```promql
sum(rate(rr_http_errors_total{status=~"4.."}[5m]))
```

**Query 2 (Server Errors):**

```promql
sum(rate(rr_http_errors_total{status=~"5.."}[5m]))
```

**Configuration:**

- **Legend:**
    - Query 1: `4xx (Client Errors)`
    - Query 2: `5xx (Server Errors)`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph (stacked area)
- **Description:** Comparison of client vs server errors

---

### 4.8 Error Rate by Endpoint and Type (Heatmap)

**Query:**

```promql
sum by (endpoint, type) (rate(rr_http_errors_total[5m]))
```

**Configuration:**

- **Legend:** N/A (heatmap)
- **Min Step:** `30s`
- **Unit:** `errors/sec`
- **Panel Type:** Heatmap
- **Description:** Visual correlation between endpoints and error types

---

### 4.9 No Free Workers Errors

**Query:**

```promql
sum(rate(rr_http_no_free_workers_total[5m]))
```

**Configuration:**

- **Legend:** `No Workers Available`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Thresholds:** Any value > 0 is critical
- **Description:** Rate of requests rejected due to worker pool exhaustion

---

### 4.10 Timeout Errors

**Query:**

```promql
sum(rate(rr_http_errors_total{type="timeout"}[5m]))
```

**Configuration:**

- **Legend:** `Timeouts`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Description:** Rate of timeout errors (408, 504)

---

### 4.11 Client Errors (4xx)

**Query:**

```promql
sum(rate(rr_http_errors_total{type="client_error"}[5m]))
```

**Configuration:**

- **Legend:** `Client Errors (4xx)`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Description:** Rate of client-side errors

---

### 4.12 Server Errors (5xx)

**Query:**

```promql
sum(rate(rr_http_errors_total{type="server_error"}[5m]))
```

**Configuration:**

- **Legend:** `Server Errors (5xx)`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Description:** Rate of server-side errors

---

### 4.13 Specific Status Code Rates

**Query (404 Not Found):**

```promql
sum(rate(rr_http_requests_by_endpoint_total{status="404"}[5m]))
```

**Query (500 Internal Server Error):**

```promql
sum(rate(rr_http_requests_by_endpoint_total{status="500"}[5m]))
```

**Query (503 Service Unavailable):**

```promql
sum(rate(rr_http_requests_by_endpoint_total{status="503"}[5m]))
```

**Configuration:**

- **Legend:** `HTTP {{status}}`
- **Min Step:** `15s`
- **Unit:** `errors/sec`
- **Panel Type:** Graph
- **Description:** Track specific problematic status codes

---

## 5. Worker Pool Health

### 5.1 Active Workers (Current)

**Query:**

```promql
rr_http_active_workers
```

**Configuration:**

- **Legend:** `Active Workers`
- **Min Step:** `5s`
- **Unit:** `short`
- **Panel Type:** Graph or Stat
- **Description:** Number of workers currently processing requests

---

### 5.2 Idle Workers (Current)

**Query:**

```promql
rr_http_idle_workers
```

**Configuration:**

- **Legend:** `Idle Workers`
- **Min Step:** `5s`
- **Unit:** `short`
- **Panel Type:** Graph or Stat
- **Description:** Number of workers available and waiting

---

### 5.3 Worker Utilization Percentage

**Query:**

```promql
rr_http_worker_utilization_percent
```

**Configuration:**

- **Legend:** `Utilization`
- **Min Step:** `5s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Gauge or Graph
- **Thresholds:** Green < 70%, Yellow 70-90%, Red > 90%
- **Description:** Worker pool utilization percentage

---

### 5.4 Total Workers

**Query:**

```promql
rr_http_active_workers + rr_http_idle_workers
```

**Configuration:**

- **Legend:** `Total Workers`
- **Min Step:** `5s`
- **Unit:** `short`
- **Panel Type:** Stat
- **Description:** Total number of workers in pool

---

### 5.5 Active vs Idle Workers (Stacked)

**Query 1:**

```promql
rr_http_active_workers
```

**Query 2:**

```promql
rr_http_idle_workers
```

**Configuration:**

- **Legend:**
    - Query 1: `Active`
    - Query 2: `Idle`
- **Min Step:** `5s`
- **Unit:** `short`
- **Panel Type:** Graph (stacked area)
- **Description:** Visual representation of worker pool state

---

### 5.6 Worker Utilization Over Time

**Query:**

```promql
avg_over_time(rr_http_worker_utilization_percent[5m])
```

**Configuration:**

- **Legend:** `Avg Utilization (5m)`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Graph
- **Description:** Smoothed worker utilization trend

---

### 5.7 Peak Worker Utilization

**Query:**

```promql
max_over_time(rr_http_worker_utilization_percent[1h])
```

**Configuration:**

- **Legend:** `Peak Utilization (1h)`
- **Min Step:** `1m`
- **Unit:** `percent (0-100)`
- **Panel Type:** Graph or Stat
- **Description:** Maximum utilization observed in last hour

---

### 5.8 Average Queue Length Over Time

**Query:**

```promql
avg_over_time(rr_http_requests_queue[5m])
```

**Configuration:**

- **Legend:** `Avg Queue Size`
- **Min Step:** `15s`
- **Unit:** `short`
- **Panel Type:** Graph
- **Description:** Average number of requests in queue

---

### 5.9 Maximum Queue Length

**Query:**

```promql
max_over_time(rr_http_requests_queue[1h])
```

**Configuration:**

- **Legend:** `Max Queue Size (1h)`
- **Min Step:** `1m`
- **Unit:** `short`
- **Panel Type:** Graph or Stat
- **Description:** Peak queue size in last hour

---

## 6. Request/Response Sizes

### 6.1 Average Request Size

**Query:**

```promql
rate(rr_http_request_size_bytes_sum[5m]) / rate(rr_http_request_size_bytes_count[5m])
```

**Configuration:**

- **Legend:** `Avg Request Size`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph
- **Description:** Mean request body size

---

### 6.2 Average Response Size

**Query:**

```promql
rate(rr_http_response_size_bytes_sum[5m]) / rate(rr_http_response_size_bytes_count[5m])
```

**Configuration:**

- **Legend:** `Avg Response Size`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph
- **Description:** Mean response body size

---

### 6.3 Request Size P95

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_request_size_bytes_bucket[5m])))
```

**Configuration:**

- **Legend:** `Request Size P95`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph
- **Description:** 95th percentile request size

---

### 6.4 Response Size P95

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_response_size_bytes_bucket[5m])))
```

**Configuration:**

- **Legend:** `Response Size P95`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph
- **Description:** 95th percentile response size

---

### 6.5 Request Size by Endpoint

**Query:**

```promql
sum by (endpoint) (rate(rr_http_request_size_bytes_sum[5m])) / sum by (endpoint) (rate(rr_http_request_size_bytes_count[5m]))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph or Table
- **Description:** Average request size per endpoint

---

### 6.6 Response Size by Endpoint

**Query:**

```promql
sum by (endpoint) (rate(rr_http_response_size_bytes_sum[5m])) / sum by (endpoint) (rate(rr_http_response_size_bytes_count[5m]))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Graph or Table
- **Description:** Average response size per endpoint

---

### 6.7 Largest Request Endpoints

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_request_size_bytes_sum[5m])) / sum by (endpoint) (rate(rr_http_request_size_bytes_count[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Bar gauge or Table
- **Description:** Endpoints with largest average request size

---

### 6.8 Largest Response Endpoints

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_response_size_bytes_sum[5m])) / sum by (endpoint) (rate(rr_http_response_size_bytes_count[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes (IEC)`
- **Panel Type:** Bar gauge or Table
- **Description:** Endpoints with largest average response size

---

### 6.9 Total Bandwidth In (Request Rate)

**Query:**

```promql
sum(rate(rr_http_request_size_bytes_sum[5m]))
```

**Configuration:**

- **Legend:** `Inbound Bandwidth`
- **Min Step:** `15s`
- **Unit:** `bytes/sec (Bps)`
- **Panel Type:** Graph
- **Description:** Total bytes per second received

---

### 6.10 Total Bandwidth Out (Response Rate)

**Query:**

```promql
sum(rate(rr_http_response_size_bytes_sum[5m]))
```

**Configuration:**

- **Legend:** `Outbound Bandwidth`
- **Min Step:** `15s`
- **Unit:** `bytes/sec (Bps)`
- **Panel Type:** Graph
- **Description:** Total bytes per second sent

---

### 6.11 Bandwidth by Endpoint (Inbound)

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_request_size_bytes_sum[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes/sec (Bps)`
- **Panel Type:** Bar gauge or Table
- **Description:** Request bandwidth per endpoint

---

### 6.12 Bandwidth by Endpoint (Outbound)

**Query:**

```promql
topk(10, sum by (endpoint) (rate(rr_http_response_size_bytes_sum[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `bytes/sec (Bps)`
- **Panel Type:** Bar gauge or Table
- **Description:** Response bandwidth per endpoint

---

### 6.13 Request/Response Size Ratio

**Query:**

```promql
(rate(rr_http_response_size_bytes_sum[5m]) / rate(rr_http_response_size_bytes_count[5m])) / (rate(rr_http_request_size_bytes_sum[5m]) / rate(rr_http_request_size_bytes_count[5m]))
```

**Configuration:**

- **Legend:** `Response/Request Ratio`
- **Min Step:** `15s`
- **Unit:** `short` (ratio)
- **Panel Type:** Graph
- **Description:** How many times larger responses are compared to requests

---

## 7. Advanced Analytics

### 7.1 Request Rate Trend (Hour over Hour)

**Query:**

```promql
sum(rate(rr_http_requests_by_endpoint_total[1h])) / sum(rate(rr_http_requests_by_endpoint_total[1h] offset 24h))
```

**Configuration:**

- **Legend:** `HoH Change`
- **Min Step:** `5m`
- **Unit:** `short` (ratio)
- **Panel Type:** Graph or Stat
- **Description:** Current hour traffic vs same hour yesterday (1.0 = same, 2.0 = double)

---

### 7.2 Performance Degradation Detection

**Query:**

```promql
histogram_quantile(0.95, sum by (le) (rate(rr_http_request_duration_seconds_bucket[5m]))) / histogram_quantile(0.95, sum by (le) (rate(rr_http_request_duration_seconds_bucket[5m] offset 1h)))
```

**Configuration:**

- **Legend:** `P95 Degradation`
- **Min Step:** `15s`
- **Unit:** `short` (ratio)
- **Panel Type:** Graph
- **Thresholds:** Green < 1.2x, Yellow 1.2-1.5x, Red > 1.5x
- **Description:** Current P95 vs 1 hour ago (1.0 = no change, 2.0 = twice as slow)

---

### 7.3 Throughput per Worker

**Query:**

```promql
sum(rate(rr_http_requests_by_endpoint_total[5m])) / (rr_http_active_workers + rr_http_idle_workers)
```

**Configuration:**

- **Legend:** `RPS per Worker`
- **Min Step:** `15s`
- **Unit:** `requests/sec (reqps)`
- **Panel Type:** Graph
- **Description:** Average requests handled per worker

---

### 7.4 Error Burst Detection

**Query:**

```promql
sum(rate(rr_http_errors_total[1m])) > 2 * avg_over_time(sum(rate(rr_http_errors_total[1m]))[10m:1m])
```

**Configuration:**

- **Legend:** `Error Burst`
- **Min Step:** `15s`
- **Unit:** `bool` (0 or 1)
- **Panel Type:** Graph (binary)
- **Description:** Detects sudden spikes in errors (>2x baseline)

---

### 7.5 Slowest Hour of Day

**Query:**

```promql
avg_over_time((rate(rr_http_request_duration_seconds_sum[1h]) / rate(rr_http_request_duration_seconds_count[1h]))[24h:1h])
```

**Configuration:**

- **Legend:** `Hourly Avg Latency`
- **Min Step:** `1h`
- **Unit:** `seconds (s)`
- **Panel Type:** Graph (24 hour view)
- **Description:** Average latency by hour over 24 hours

---

### 7.6 Endpoint Performance Correlation

**Query:**

```promql
(sum by (endpoint) (rate(rr_http_errors_total[5m])) / sum by (endpoint) (rate(rr_http_requests_by_endpoint_total[5m]))) * (sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_sum[5m])) / sum by (endpoint) (rate(rr_http_duration_by_endpoint_seconds_count[5m])))
```

**Configuration:**

- **Legend:** `{{endpoint}}`
- **Min Step:** `15s`
- **Unit:** `short`
- **Panel Type:** Graph or Table
- **Description:** Correlation score: slow + error-prone endpoints have higher values

---

### 7.7 Server Uptime

**Query:**

```promql
rr_http_uptime_seconds
```

**Configuration:**

- **Legend:** `Uptime`
- **Min Step:** `1m`
- **Unit:** `seconds (s)` or `duration (s)`
- **Panel Type:** Stat
- **Description:** How long server has been running

---

### 7.8 Request Distribution by HTTP Method (Pie)

**Query:**

```promql
sum by (method) (rate(rr_http_requests_by_endpoint_total[5m])) / sum(rate(rr_http_requests_by_endpoint_total[5m])) * 100
```

**Configuration:**

- **Legend:** `{{method}}`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Pie chart
- **Description:** Percentage breakdown of requests by HTTP method

---

### 7.9 Status Code Distribution (Pie)

**Query:**

```promql
sum by (status) (rate(rr_http_requests_by_endpoint_total[5m])) / sum(rate(rr_http_requests_by_endpoint_total[5m])) * 100
```

**Configuration:**

- **Legend:** `{{status}}`
- **Min Step:** `15s`
- **Unit:** `percent (0-100)`
- **Panel Type:** Pie chart
- **Description:** Percentage breakdown of responses by status code

---

## 8. Dashboard Layout Recommendations

### Row 1: Key Metrics Overview (4 panels)

1. **Total RPS** - Stat panel
2. **P95 Latency** - Stat panel with threshold colors
3. **Error Rate %** - Gauge with thresholds
4. **Worker Utilization** - Gauge with thresholds

### Row 2: Traffic & Performance (2 panels)

1. **RPS by Endpoint** - Graph (time series)
2. **Latency Percentiles** - Graph (P50, P95, P99)

### Row 3: Performance Breakdown (2 panels)

1. **Queue vs Processing Time** - Stacked area graph
2. **Top Slowest Endpoints** - Bar gauge

### Row 4: Error Analysis (2 panels)

1. **Error Rate by Type** - Stacked area graph
2. **Most Error-Prone Endpoints** - Table

### Row 5: Worker Pool Health (2 panels)

1. **Active vs Idle Workers** - Stacked area graph
2. **Queue Size** - Graph with threshold line

### Row 6: Bandwidth Analysis (2 panels)

1. **Request/Response Sizes** - Graph (dual Y-axis)
2. **Top Bandwidth Consumers** - Table

### Row 7: Detailed Endpoint Table (1 panel)

1. **Endpoint Performance Table** - Table with multiple queries

---

## 9. Unit Reference Guide

### Standard Grafana Units

**Time Units:**

- `seconds (s)` - for durations
- `milliseconds (ms)` - for sub-second timings
- `duration (s)` - auto-formats (1m 30s, 2h 15m, etc.)

**Rate Units:**

- `requests/sec (reqps)` - for request rates
- `errors/sec` - for error rates
- `bytes/sec (Bps)` - for bandwidth

**Size Units:**

- `bytes (IEC)` - auto-formats (KB, MB, GB) using 1024 base
- `bytes (SI)` - auto-formats using 1000 base

**Percentage:**

- `percent (0-100)` - displays as 95%
- `percentunit (0.0-1.0)` - displays 0.95 as 95%

**Count:**

- `short` - auto-formats large numbers (1K, 1M)
- `none` - raw number

**Boolean:**

- `bool` - 0 or 1
- `bool_yes_no` - displays as Yes/No
- `bool_on_off` - displays as On/Off

---

## 10. Common Threshold Configurations

### Latency Thresholds

```
Green: < 1s
Yellow: 1-5s
Red: > 5s
```

### Error Rate Thresholds

```
Green: < 1%
Yellow: 1-5%
Red: > 5%
```

### Worker Utilization Thresholds

```
Green: < 70%
Yellow: 70-90%
Red: > 90%
```

### Queue Time Thresholds

```
Green: < 100ms
Yellow: 100-500ms
Red: > 500ms
```

### Success Rate Thresholds

```
Red: < 95%
Yellow: 95-99%
Green: > 99%
```