# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Install ca-certificates for HTTPS requests during build
RUN apk add --no-cache ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Copy and setup corporate certificate if it exists (before go mod download)
COPY . /tmp/certcheck/
RUN if [ -f /tmp/certcheck/corp-ca.crt ]; then \
        cp /tmp/certcheck/corp-ca.crt /usr/local/share/ca-certificates/corp-ca.crt && \
        update-ca-certificates; \
    fi && rm -rf /tmp/certcheck

# Download dependencies (with certificate if needed)
RUN go mod download

# Copy only necessary source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY swagger/ ./swagger/
COPY migrations/ ./migrations/

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o main cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests at runtime
RUN apk --no-cache add ca-certificates

# Copy corporate certificate if it exists (optional)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the binary from builder
COPY --from=builder /build/main .

# Copy migration files from builder
COPY --from=builder /build/migrations ./migrations

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8000

CMD ["./main"]