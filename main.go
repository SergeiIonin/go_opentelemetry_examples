package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cactus/go-statsd-client/v5/statsd"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Global StatSD client
var statsdClient statsd.Statter

// Global metric counter
var helloCounter metric.Int64Counter

// Initialize OpenTelemetry
func initOpenTelemetry() (*sdkmetric.MeterProvider, error) {
	ctx := context.Background()

	// Get OTLP endpoint from environment or use default
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	log.Printf("Configuring OpenTelemetry with endpoint: %s", endpoint)

	// Configure the OTLP exporter with retry and timeout options
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     10 * time.Second,
			MaxElapsedTime:  30 * time.Second,
		}),
		otlpmetricgrpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, err
	}

	// Create meter provider with a longer interval to reduce connection attempts
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))),
	)

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create a meter
	meter := meterProvider.Meter("gin-server")

	// Create instruments
	var err2 error
	helloCounter, err2 = meter.Int64Counter(
		"hello_requests_total",
		metric.WithDescription("Total number of hello requests"),
	)
	if err2 != nil {
		return nil, err2
	}

	return meterProvider, nil
}

// Initialize StatsD client
func initStatsd() (statsd.Statter, error) {
	// StatsD configuration from environment or use default
	statsdAddr := os.Getenv("STATSD_ADDR")
	if statsdAddr == "" {
		statsdAddr = "localhost:8125"
	}

	// Create a new StatsD client
	config := &statsd.ClientConfig{
		Address:       statsdAddr,
		Prefix:        "gin-app",
		UseBuffered:   true,
		FlushInterval: time.Second,
	}

	return statsd.NewClientWithConfig(config)
}

func main() {
	// Initialize StatsD
	var err error
	statsdClient, err = initStatsd()
	if err != nil {
		log.Fatalf("Failed to create StatsD client: %v", err)
	}
	defer statsdClient.Close()

	// Initialize OpenTelemetry
	mp, err := initOpenTelemetry()
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}
	defer func() {
		if err := mp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down meter provider: %v", err)
		}
	}()

	// Set up Gin in release mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Use the OpenTelemetry middleware
	r.Use(otelgin.Middleware("gin-server"))

	// Recovery middleware
	r.Use(gin.Recovery())

	// Define routes
	r.GET("/hello", func(c *gin.Context) {
		// Record metrics with OpenTelemetry
		helloCounter.Add(context.Background(), 1)

		// Record metrics with StatsD
		statsdClient.Inc("endpoint.hello.requests", 1, 1.0)

		// Return response
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, OpenTelemetry!",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
