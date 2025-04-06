FROM golang:1.20-alpine as builder

WORKDIR /app
COPY ./main.go .

RUN apk add --no-cache git \
    && go mod init example.com/app \
    && go mod tidy \
    && go build -o server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
ENV APP_PORT=8080
ENTRYPOINT ["./server"]
