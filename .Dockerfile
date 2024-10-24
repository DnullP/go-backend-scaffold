FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o comment_service ./apps/comment/main.go

RUN go build -o user_service ./apps/user/main.go

# comment
FROM alpine:latest AS comment_service
WORKDIR /app
COPY --from=builder /app/comment_service .
CMD ["./comment_service"]

# user
FROM alpine:latest AS user_service
WORKDIR /app
COPY --from=builder /app/user_service .
CMD ["./user_service"]
