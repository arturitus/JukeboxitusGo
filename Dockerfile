# Start with the official Golang image from the Docker Hub
# FROM golang:1.23 AS builder
FROM golang:latest AS builder
# Continue with your Dockerfile instructions

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

ENV GOARCH=arm64
ENV GOOS=linux
# Build the Go app
RUN go build -o bin/main ./src

# Start a new stage from scratch
# FROM arm64v8/debian:bookworm-slim
FROM alpine:3.19

# Install required dependencies (if any)
# RUN apt-get update && apt-get install -y \
#     ca-certificates \
#     && rm -rf /var/lib/apt/lists/*

RUN apk add --no-cache ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bin/main /root/bin/main

# Ensure the binary is executable
RUN chmod +x /root/bin/main

# Command to run the executable
CMD ["/root/bin/main"]
