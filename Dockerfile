# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o task-manager-app .

# Stage 2: Create a lightweight image for execution
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/task-manager-app .

# Copy migrations folder
COPY migrations ./migrations

# Expose port and run the application
EXPOSE 8080
CMD ["./task-manager-app"]