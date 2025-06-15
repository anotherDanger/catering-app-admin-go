FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

RUN apk --no-cache add -ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]