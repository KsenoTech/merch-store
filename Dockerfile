# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Create a minimal container
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy wait-for-db script
COPY wait-for-db.sh .

# Install netcat for waiting for DB
RUN apk add --no-cache bash nc

# Expose port 8080
EXPOSE 8080

# Wait for DB and run the application
CMD ["./wait-for-db.sh", "./main"]