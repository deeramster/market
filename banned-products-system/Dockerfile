FROM golang:1.24.1-bookworm as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -v -o /app/banned-products-system ./cmd/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/banned-products-system /app/

COPY --from=builder /app/storage /app/storage
ENV BANNED_PRODUCTS_FILE=/app/storage/data/banned_products.json
RUN chmod +x /app/banned-products-system
ENTRYPOINT ["/app/banned-products-system"]
CMD ["list"]
