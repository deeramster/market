FROM golang:1.24.1-bookworm as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o spark-app .


FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y bash libc6 krb5-user && rm -rf /var/lib/apt/lists/*
COPY --from=builder /app/spark-app /app/spark-app
COPY certs/ /app/certs/
RUN mkdir -p /data
WORKDIR /app

CMD ["/app/spark-app"]
