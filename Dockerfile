# Use the official Golang image
FROM golang:latest

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for caching deps)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the remaining source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Expose app port
EXPOSE 8080

# Start the app
CMD ["./main"]
