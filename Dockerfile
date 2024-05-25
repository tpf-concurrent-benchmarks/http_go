# Stage 1: Build the Go application
FROM golang:1.22 AS builder
LABEL stage=builder

WORKDIR /app

COPY ./src/ .
RUN go mod download
RUN go get http_go

RUN CGO_ENABLED=0 GOOS=linux go build -o main.exe

# Stage 2: Create a minimal image with the application
FROM alpine

# Install necessary packages
RUN apk add --no-cache bash

# Create a non-root user and group
RUN addgroup -S nonrootgroup && adduser -S nonrootuser -G nonrootgroup

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app /app

# Ensure the non-root user has ownership of the app directory
RUN chown -R nonrootuser:nonrootgroup /app

# Switch to the non-root user
USER nonrootuser

# Run the application
CMD ["./main.exe"]
