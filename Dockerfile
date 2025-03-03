FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY .env .env

EXPOSE 14704

CMD ["./main"]
