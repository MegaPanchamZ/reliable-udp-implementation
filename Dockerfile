# === Build Stage ===
# Use the official Go image as the builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the sender executable
RUN go build -o /bin/sender -ldflags="-s -w" ./cmd/sender

# Build the receiver executable
RUN go build -o /bin/receiver -ldflags="-s -w" ./cmd/receiver

# === Final Stage ===
# Use a minimal alpine image for the final container
FROM alpine:latest

# Install tools for testing (e.g., diff)
RUN apk add --no-cache coreutils

# Copy only the compiled binaries from the builder stage
COPY --from=builder /bin/sender /bin/sender
COPY --from=builder /bin/receiver /bin/receiver

WORKDIR /data
