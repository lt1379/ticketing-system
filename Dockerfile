# Stage 1: Build the Go application
FROM golang:1.23.3 AS builder
# Install necessary packages including tzdata
RUN apt-get update && apt-get install -y tzdata

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Stage 2: Create a minimal image with the compiled binary
FROM debian:bullseye-slim

# Install necessary packages including tzdata
RUN apt-get update && apt-get install -y tzdata

# Set the timezone environment variable
ENV TZ=Asia/Jakarta

# Install necessary packages (if any)
#RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main /main
COPY --from=builder /app/config.json /config.json

RUN mkdir "logs"

# Command to run the executable
CMD ["/main"]