FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o comment_service ./apps/comment/main.go
RUN go build -o user_service ./apps/user/main.go
RUN go build -o gateway ./apps/gateway/main.go

# user
FROM alpine:latest AS server
WORKDIR /app
COPY --from=builder /app/user_service .
COPY --from=builder /app/comment_service .
COPY --from=builder /app/gateway .
COPY --from=builder /config/common-config.yaml ./config/