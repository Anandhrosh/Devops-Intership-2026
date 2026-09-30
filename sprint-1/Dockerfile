# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/server main.go

# Create a custom non-root user
RUN printf 'appuser:x:10001:10001:Application User:/nonexistent:/sbin/nologin\n' \
    > /out/passwd

# Stage 2: Minimal production image
FROM scratch

WORKDIR /app

# Copy the compiled application
COPY --from=builder /out/server /app/server

# Copy the custom user information
COPY --from=builder /out/passwd /etc/passwd

# Run as the custom non-root user
USER appuser

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s \
    --start-period=5s --retries=3 \
    CMD ["/app/server", "healthcheck"]

CMD ["/app/server"]