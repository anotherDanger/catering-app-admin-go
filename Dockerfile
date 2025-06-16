FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

RUN apk --no-cache add -ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

ENTRYPOINT [ "./app" ]