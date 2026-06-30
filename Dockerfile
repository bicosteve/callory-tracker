# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

# Copy source code and UI files
COPY . .

# Build the web application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o web ./cmd/web

# Production stage
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/web .
COPY --from=builder /app/ui ./ui


EXPOSE 4001

# Run the application
CMD ["./web"]
