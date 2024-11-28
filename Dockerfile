# Use a Golang base image as the base image for building the application
FROM golang:alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy the Go module files to the container
COPY go.mod go.sum ./

# Install Go dependencies
RUN go mod download

# Copy the entire project code to the container (including the src directory)
COPY . .

# Build the Go application
RUN go build -o app main.go

# Use a lightweight base image for the final runtime image
FROM alpine:latest

# Install necessary dependencies for HTTP requests (e.g., certificates)
RUN apk add --no-cache ca-certificates

# Set the working directory in the container for runtime
WORKDIR /app

# Copy the compiled Go binary from the builder stage to the final image
COPY --from=builder /app/app /app/

# Copy the src directory (which contains static files like index.html)
COPY --from=builder /app/src /app/src

# Expose the port your application will listen on
EXPOSE 3000

# Run the Go application
CMD ["/app/app"]
