FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy source code
COPY server/*.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gosync-server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/gosync-server .

# Create data directory
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

CMD ["./gosync-server"]
