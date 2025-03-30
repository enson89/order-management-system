# Use multi-architecture build to support amd64 and arm64
FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

# Set the working directory
WORKDIR /app

# Install required dependencies for librdkafka and build tools
RUN apt-get update && apt-get install -y gcc librdkafka-dev pkg-config

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy application files
COPY . .

# Accept build-time arguments for target OS and architecture
ARG TARGETOS
ARG TARGETARCH

# Set LIBRDKAFKA_PATH so that the dynamic build finds the correct librdkafka installation
ENV LIBRDKAFKA_PATH=/usr

# Build the application for the target architecture using dynamic linking.
# The "-tags dynamic" flag forces confluentinc/confluent-kafka-go to link dynamically.
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=1 go build -tags dynamic -o main ./cmd/app/main.go

# Use a Debian Bookworm-based slim image for the final stage, which provides GLIBC >= 2.34
FROM --platform=$TARGETPLATFORM debian:bookworm-slim

# Set the working directory in the final image
WORKDIR /root/

# Update and install librdkafka runtime (Bookworm includes a newer glibc)
RUN apt-get update && apt-get install -y librdkafka1 && rm -rf /var/lib/apt/lists/*

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy the config directory from the builder stage into the final image
COPY --from=builder /app/config ./config

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]