FROM golang:1.23.5 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat-server ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/chat-server .

EXPOSE 8080
CMD ["./chat-server"]
