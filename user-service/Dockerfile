FROM golang:1.23-alpine

# Install Git
RUN apk add --no-cache git

WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o main .

EXPOSE 8082

CMD ["./main"]