# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /workspace

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the controller
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager controller.go

# Final stage
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy the binary from builder stage
COPY --from=builder /workspace/manager .

USER 65532:65532

ENTRYPOINT ["/manager"]