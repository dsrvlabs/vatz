# Start from the official Golang image based on Debian
FROM golang:1.20-buster as builder

# Set the Current Working Directory inside the container
WORKDIR /app/vatz

# Copy the entire project into the container
COPY . .

# Run the Makefile to build the application
RUN make build

# Run vatz init to initialize the application (this should generate the binary)
RUN ./vatz init

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/vatz/vatz ./vatz

# Expose port 8080 (or whatever port your application uses)
EXPOSE 8080

# Command to start the application, expecting a config.yaml file to be provided at runtime
CMD ["./vatz", "start", "--config", "/config/config.yaml"]