# Stage 1: Build the Go application
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install goose and sqlc
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./
# COPY *.go ./

# Generate SQL code with sqlc (assumes sqlc.yaml is properly configured)
RUN sqlc generate -f ./db/sqlc.yaml

# Run goose migrations (assuming you have a db/migrations directory)
RUN goose -dir ./db/migrations sqlite3 ./db/minesweeper.db up

# Build the Go app
# RUN CGO_ENABLED=0 GOOS=linux go build -o /server-build
RUN CGO_ENABLED=0 GOOS=linux go build -o main /app

# Stage 2: A minimal image to run the Go app
FROM alpine:latest

# Install SQLite and necessary ca-certificates (if your Go app makes HTTPS requests)
RUN apk --no-cache add sqlite-libs ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy SQLite database file with migrations applied
COPY --from=builder /app/db/minesweeper.db ./db/minesweeper.db

# Copy any necessary static files or templates
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/dist ./dist

# Copy migrations directory if you want to run migrations at startup
COPY --from=builder /app/db/migrations ./db/migrations

# Expose port (adjust according to your Go app configuration)
EXPOSE 8080

# Command to run the executable
CMD ["/server-build"]
