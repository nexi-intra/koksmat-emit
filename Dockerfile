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

# Run Go tests with coverage
RUN go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out > coverage.txt

# Optionally, enforce a minimum coverage percentage (e.g., 80%)
# Uncomment the following lines to enforce coverage threshold
# RUN COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//') && \
#     if [ $(echo "$COVERAGE < 80" | bc) -eq 1 ]; then \
#         echo "Coverage ($COVERAGE%) is below the required threshold (80%)"; \
#         exit 1; \
#     fi

# Display the coverage summary in build logs
RUN cat coverage.txt

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
