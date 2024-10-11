# Build aşaması
FROM golang:1.23 AS builder

WORKDIR /app

# Gerekli modülleri yükle
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodunu kopyala ve binary'yi build et
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./cmd/main.go

# Release aşaması
FROM alpine:3.18

WORKDIR /root/

# Build edilen binary'yi kopyala
COPY --from=builder /app/todo-app .

# Binary'nin çalıştırılabilir olup olmadığını kontrol et
RUN chmod +x /root/todo-app

# Portu expose et
EXPOSE 8080

# Uygulamayı başlat
CMD ["./todo-app"]
