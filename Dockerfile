# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
  CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o koksmat-emit main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set environment variables
ENV PORT=8080

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/koksmat-emit .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./koksmat-emit","serve"]
