# Start from the official lightweight Go image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the Go application into an executable named 'graphene-api'
RUN go build -o graphene-api .

# Expose the port your Gin router is running on
EXPOSE 8080

# Command to run the executable when the container starts
CMD ["./graphene-api"]