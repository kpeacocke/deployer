# Multi-stage build for gh-deployer
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 deployer

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o gh-deployer .

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/gh-deployer /usr/local/bin/gh-deployer

# Use non-root user
USER deployer

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/gh-deployer"]

# Default command
CMD ["--help"]

# Labels
LABEL org.opencontainers.image.title="gh-deployer" \
      org.opencontainers.image.description="GitHub release deployer with blue/green deployment" \
      org.opencontainers.image.url="https://github.com/kpeacocke/deployer" \
      org.opencontainers.image.source="https://github.com/kpeacocke/deployer" \
      org.opencontainers.image.licenses="MIT"
