FROM golang:1.25-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o anxiety-relief-tgbot-go ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache tzdata ca-certificates
WORKDIR /root/
COPY --from=builder /app/anxiety-relief-tgbot-go /anxiety-relief-tgbot-go
ENTRYPOINT ["/anxiety-relief-tgbot-go"]
