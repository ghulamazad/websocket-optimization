# Step 1: Build the Go application
FROM golang:1.23.1-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all Go dependencies
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o websocket-server ./cmd/main.go

# Step 2: Run the application
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /app/websocket-server .

# Expose port 8080 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./websocket-server"]
