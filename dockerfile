FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/api

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/server .
COPY web/ ./web
EXPOSE 8080
CMD ["./server"]
