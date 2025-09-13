FROM golang:1.21-alpine AS builder

# Install dependencies for building
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tailscale-exporter .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates curl

# Install Tailscale
RUN curl -fsSL https://tailscale.com/install.sh | sh

# Create non-root user
RUN addgroup -g 1001 -S tailscale && \
    adduser -u 1001 -S tailscale -G tailscale

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tailscale-exporter .

# Change ownership
RUN chown tailscale:tailscale tailscale-exporter

# Switch to non-root user
USER tailscale

# Expose port
EXPOSE 9090

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:9090/metrics || exit 1

# Run the exporter
CMD ["./tailscale-exporter"]

