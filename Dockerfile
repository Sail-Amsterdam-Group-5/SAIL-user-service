FROM golang:1.22.0 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]