FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tailscale-exporter ./cmd/tailscale-exporter

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN addgroup -g 1001 -S tailscale-exporter && \
  adduser -u 1001 -S tailscale-exporter -G tailscale-exporter

WORKDIR /app

COPY --from=builder /app/tailscale-exporter .

RUN chown tailscale-exporter:tailscale-exporter tailscale-exporter

USER tailscale

EXPOSE 9090

# Run the exporter
CMD ["./tailscale-exporter"]

