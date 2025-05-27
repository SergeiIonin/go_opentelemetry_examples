# Go Application with Gin and OpenTelemetry

This is a simple Go web application built with the Gin framework that demonstrates OpenTelemetry integration with StatsD metrics.

## Project Structure

- `main.go` - Main application code with Gin HTTP server and OpenTelemetry instrumentation
- `Dockerfile` - Docker image definition for the Go application
- `docker-compose.yml` - Docker Compose configuration to run the app with an OpenTelemetry collector
- `otel-collector-config.yaml` - Configuration for the OpenTelemetry collector

## Features

- REST API with Gin framework
- `/hello` endpoint on port 8888
- OpenTelemetry instrumentation for metrics and traces
- StatsD metrics collection
- Docker and Docker Compose setup

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

## Environment Variables

The application can be configured using the following environment variables:

- `PORT` - The port the application listens on (default: 8888)
- `OTEL_EXPORTER_OTLP_ENDPOINT` - The endpoint of the OpenTelemetry collector (default: otel-collector:4317)
- `STATSD_ADDR` - The address of the StatsD server (default: otel-collector:8125)
