FROM golang:1.22.0 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

CMD ["./main"]