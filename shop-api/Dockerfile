FROM golang:1.24.1-bookworm as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=1 GOOS=linux go build -o /shop-api ./cmd/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /shop-api /app/
COPY config.json /app/config/
COPY products.json /app/data/
RUN mkdir -p /app/certs
RUN chmod +x /app/shop-api
ENTRYPOINT ["/app/shop-api"]
