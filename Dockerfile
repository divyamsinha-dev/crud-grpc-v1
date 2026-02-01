
FROM golang:alpine AS builder


RUN apk update && apk add --no-cache git

WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .


RUN go build -o main ./server/main.go


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/main .





CMD ["./main"]

