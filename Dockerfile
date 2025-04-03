FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
ENV APP_PORT=8080
CMD ["./server"]
