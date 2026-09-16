FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

FROM alpine:3.19

RUN apk add --no-cache ca-certificates bash

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations

COPY <<EOF /app/start.sh
#!/bin/bash
echo "Running migrations..."
./migrate
echo "Starting server..."
./server
EOF

RUN chmod +x /app/start.sh

EXPOSE 8080

CMD ["/app/start.sh"]
