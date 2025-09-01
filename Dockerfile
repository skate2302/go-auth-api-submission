# syntax=docker/dockerfile:1

# Stage 1: The builder stage
# We use a specific Go version to build our application.
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies.
# This is done in a separate step to leverage Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application.
# -o /app/main specifies the output file.
# CGO_ENABLED=0 is important for creating a static binary.
# -ldflags "-s -w" makes the binary smaller.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Stage 2: The final stage
# We use a minimal 'scratch' image which contains nothing but our application.
# This makes the final image very small and secure.
FROM scratch

# Set the working directory
WORKDIR /app

# Copy the built application binary from the 'builder' stage
COPY --from=builder /app/main .

# Expose the port that our application will run on
EXPOSE 8080

# The command to run when the container starts
ENTRYPOINT ["/app/main"]