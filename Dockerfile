# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o dhcp-api-server ./main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/dhcp-api-server ./
COPY dhcpd.conf.template ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./dhcp-api-server"]