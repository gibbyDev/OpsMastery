FROM golang:1.23-alpine AS builder

WORKDIR /app

# Initialize a new Go module inside the container
RUN go mod init opsmastery

# Copy the rest of the files
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o /app/main .

FROM alpine:latest

WORKDIR /root/

# Copy the built application and .env file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
