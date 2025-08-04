# ---- Stage 1: Build ----
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code AND the .env file
COPY . .
COPY .env .env

# Build the Go binary
RUN go build -o server .

# ---- Stage 2: Run (Minimal Image) ----
FROM alpine:latest

# Create a non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
# Copy .env file from builder
COPY --from=builder /app/.env .env

# Set executable permissions
RUN chmod +x ./server

# Set the user
USER appuser

# Expose the port
EXPOSE 8080

# Start the Go server
ENTRYPOINT ["./server"]