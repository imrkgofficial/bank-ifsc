# Step 1: Use a Golang base image
FROM golang:1.23.3 AS builder

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy go.mod and go.sum files to install dependencies
COPY go.mod go.sum ./

# Step 4: Install the dependencies (go mod tidy)
RUN go mod tidy

# Step 5: Copy the rest of the application files into the container
COPY . .

# Step 6: Run the tests
RUN go test ./...

# Step 7: Build the Go application
RUN go build -o bank-ifsc main.go

# Step 8: Create a minimal runtime image (optional for production deployment)
FROM golang:1.23.3 AS runtime

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bank-ifsc .

# Step 9: Define the command to run the application
CMD ["./bank-ifsc"]

# Step 10: Optionally, you can add deployment steps
# This could involve deploying the binary with Docker, Portainer, or other methods
# Here, we use a simple echo command to simulate the deployment:
RUN echo "Deploying application..."
