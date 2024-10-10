FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o todo-app ./cmd/main.go

FROM alpine:3.18

WORKDIR /root/

COPY --from=builder /app/todo-app .

EXPOSE 8080

CMD ["./todo-app"]
