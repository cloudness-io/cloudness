# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=off go build -a -installsuffix cgo -ldflags="-w -s" -o placeholder .

# Final stage
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/placeholder /bin/sh

# Run the application
ENTRYPOINT ["/bin/sh"]
