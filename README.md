# Go Application with Gin and OpenTelemetry

This is a simple Go web application built with the Gin framework that demonstrates OpenTelemetry integration with StatsD metrics, Prometheus for metrics collection, and Grafana for visualization.

## Project Structure

- `main.go` - Main application code with Gin HTTP server and OpenTelemetry instrumentation
- `Dockerfile` - Docker image definition for the Go application
- `docker-compose.yml` - Docker Compose configuration to run the app with an OpenTelemetry collector, Prometheus, and Grafana
- `otel-collector-config.yaml` - Configuration for the OpenTelemetry collector
- `prometheus.yml` - Configuration for Prometheus to scrape metrics from the OpenTelemetry collector
- `grafana-provisioning/` - Provisioning configuration for Grafana
  - `datasources/` - Datasource configurations
  - `dashboards/` - Dashboard configurations and JSON definitions

## Features

- REST API with Gin framework
- `/hello` endpoint on port 8888
- OpenTelemetry instrumentation for metrics and traces
- StatsD metrics collection
- Complete observability stack with:
  - OpenTelemetry Collector for collecting and processing telemetry data
  - Prometheus for metrics storage
  - Grafana for metrics visualization with pre-configured dashboards
- Docker and Docker Compose setup for easy deployment

## Running the Application

```bash
# Build and run with Docker Compose
docker-compose up --build
```

## Testing the Application

Once the application is running, you can test it with:

```bash
curl http://localhost:8888/hello
```

## Observing Metrics

The metrics are exported to the OpenTelemetry collector and can be viewed in various ways:

- StatsD metrics are received on port 8125
- OTLP metrics are available via gRPC on port 4317
- Prometheus metrics are exposed on port 8889
- Prometheus UI is available at http://localhost:9090
- Grafana dashboards are available at http://localhost:3000
  - Default login credentials: admin/admin
  - Pre-configured dashboard "Go Application Metrics" is automatically provisioned

### Available Metrics

The application collects and visualizes the following metrics:

- `ginapp_hello_requests_total` - Counter of requests to the /hello endpoint
- `ginapp_http_server_request_duration_seconds` - Histogram of request durations
- `ginapp_http_server_request_body_size_bytes` - Histogram of request body sizes
- `ginapp_http_server_response_body_size_bytes` - Histogram of response body sizes

## Environment Variables

The application can be configured using the following environment variables:

- `PORT` - The port the application listens on (default: 8888)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - The endpoint of the OpenTelemetry collector (default: otel-collector:4317)
- `STATSD_ADDR` - The address of the StatsD server (default: otel-collector:8125)

## Component Configuration

### OpenTelemetry Collector

The OpenTelemetry collector is configured in `otel-collector-config.yaml` with:

- OTLP receivers on ports 4317 (gRPC) and 4318 (HTTP)
- StatsD receiver on port 8125
- Memory limiter and batch processors
- Debug exporter for logging
- Prometheus exporter on port 8889
- Health check, pprof, and zPages extensions

### Prometheus

Prometheus is configured in `prometheus.yml` to scrape metrics from the OpenTelemetry collector every 5 seconds.

### Grafana

Grafana is set up with:

- Automatic provisioning of Prometheus data source
- Pre-configured dashboards for application metrics
- Persistence using Docker volumes

## Architecture

The system architecture consists of:

1. Go application with OpenTelemetry SDK
2. OpenTelemetry collector for collecting and processing telemetry data
3. Prometheus for metrics storage
4. Grafana for metrics visualization

The data flows as follows:
- Go app → OpenTelemetry collector (OTLP/StatsD) → Prometheus → Grafana
