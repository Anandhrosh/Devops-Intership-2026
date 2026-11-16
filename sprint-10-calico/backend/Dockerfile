
# Stage 1: Build the Go application
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY application/go.mod application/go.sum ./
RUN go mod download

# Copy application source code
COPY application/ .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /capstone-backend .

# Stage 2: Minimal non-root runtime
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /

COPY --from=builder /capstone-backend /capstone-backend

EXPOSE 3000

ENTRYPOINT ["/capstone-backend"]
