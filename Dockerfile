# Use the Go base image
FROM golang:1.23.0-alpine

# Set the working directory inside the container
WORKDIR /app

# Install bash in the Alpine image
RUN apk add --no-cache bash

# Copy all files to the working directory
COPY . .

# Ensure the wait-for-it.sh script is executable
RUN chmod +x wait-for-it.sh

# Download Go dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Set the entrypoint to use bash
ENTRYPOINT ["./wait-for-it.sh", "db:5432", "--", "/app/main"]
