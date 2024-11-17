# Use the official Golang image as the base image
FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for dependency resolution
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o bank-ifsc main.go

# Expose the port your application will run on
EXPOSE 3000

# Command to run the application
CMD ["./bank-ifsc"]