FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN CGO_ENABLED=0 go build -o redis-server ./cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/redis-server /redis-server
EXPOSE 6379
CMD ["/redis-server"]