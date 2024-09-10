# Stage 1: Build the Go application
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install goose and sqlc
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install Node.js and npm for Tailwind CSS
RUN apt-get update && apt-get install -y nodejs npm

# Copy go.mod, go.sum, and package.json files
COPY go.mod go.sum package.json package-lock.json ./

# Download all dependencies for Go and Node.js. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
RUN npm install

# Copy the source from the current directory to the Working Directory inside the container
COPY . ./

# Generate SQL code with sqlc (assumes sqlc.yaml is properly configured)
RUN sqlc generate -f ./db/sqlc.yaml

# Run goose migrations (assuming you have a db/migrations directory)
RUN goose -dir ./db/migrations sqlite3 ./db/minesweeper.db up

# Build the Tailwind CSS file
RUN npx tailwindcss -i ./main.css -o ./dist/tailwind.css

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

# Copy any necessary static files, templates, and CSS
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/dist ./dist

# Copy migrations directory if you want to run migrations at startup
COPY --from=builder /app/db/migrations ./db/migrations

# Expose port (adjust according to your Go app configuration)
EXPOSE 8080

# Command to run the executable
# CMD ["/server-build"]
CMD ["./main"]
