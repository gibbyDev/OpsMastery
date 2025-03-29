FROM golang:1.22-alpine AS builder

WORKDIR /app

# Initialize a new Go module inside the container
RUN go mod init opsmastery

# Download dependencies (if any are added later)
RUN go mod tidy

# Copy the rest of the files
COPY . .

# Copy the .env file into the container
COPY .env /app/.env

# Build the application
RUN go build -o /app/main .

FROM alpine:latest

WORKDIR /root/

# Copy the built application and .env file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
