FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

#RUN go mod tidy
# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a small image for the runtime
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for potential HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/app .

# Expose the application port
EXPOSE 8888

# Set environment variables with defaults
ENV PORT=8888
ENV OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
ENV STATSD_ADDR=otel-collector:8125

# Run the binary
CMD ["./app"]
