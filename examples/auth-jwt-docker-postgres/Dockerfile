# Use the official Go image as a base
FROM golang:1.21-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
# CGO_ENABLED=0 is important for static binaries, useful in Alpine
# -o specifies the output file name
RUN CGO_ENABLED=0 go build -o /auth-app ./main.go

# --- Start a new stage for a smaller final image ---
FROM alpine:latest

WORKDIR /root/

# Copy the compiled application from the builder stage
COPY --from=builder /auth-app .

# Copy the go.mod and go.sum from the builder stage to allow replace directive to work
# This is a workaround for the replace directive in examples. In a real app,
# you'd typically vendor dependencies or build from the root module.
COPY --from=builder /app/go.mod /app/go.sum /app/

# Expose the port the app runs on
EXPOSE 8082

# Command to run the executable
CMD ["./auth-app"]
