# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.27.2-bookworm AS builder

WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go untuk Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /roomly-api .


# ==========================================
# STAGE 2: RUNTIME
# ==========================================
FROM debian:bookworm-slim AS runtime

WORKDIR /app

# Install timezone data dan SSL certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

# Set timezone Jakarta
ENV TZ=Asia/Jakarta

# Copy executable dari stage builder
COPY --from=builder /roomly-api ./roomly-api

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./roomly-api"]