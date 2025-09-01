FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy from root project (2 level up)
COPY ../../go.mod ../../go.sum ./
COPY ../../ .

# Build gateway
RUN CGO_ENABLED=0 go build -o gateway-app ./cmd/gateway

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/gateway-app /gateway-app
EXPOSE 8080
CMD ["/gateway-app"]