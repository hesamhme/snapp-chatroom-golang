FROM golang:1.23 as builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o server ./cmd/server/main.go
RUN go build -o client ./cmd/client/main.go


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/server .
COPY --from=builder /app/client .



EXPOSE 8080


CMD ["./server"]
