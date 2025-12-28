FROM golang:1.25.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o account-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/account-service .

EXPOSE 8000
CMD ["./account-service"]