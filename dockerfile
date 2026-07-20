# Step 1: Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Step 2: Run stage
FROM alpine:latest

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy built binary and assets
COPY --from=builder /app/server /app/server
COPY --from=builder /app/db/init_schema.sql /app/db/init_schema.sql

EXPOSE 7070

CMD ["/app/server"]
