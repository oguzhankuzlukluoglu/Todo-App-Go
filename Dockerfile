# Build aşaması
FROM golang:1.23 AS builder

WORKDIR /app

# Gerekli modülleri yükle
COPY go.mod go.sum ./ 
RUN go mod download

# Kaynak kodunu ve gerekli dosyaları kopyala
COPY . ./

# Binary'yi build et
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./main.go

# Release aşaması
FROM alpine:3.18

WORKDIR /root/

# Build edilen binary'yi kopyala
COPY --from=builder /app/todo-app .

# Casbin için gerekli dosyaları kopyala
COPY --from=builder /app/model.conf /app/policy.csv ./

# Swagger için gerekli dokümanları kopyala
COPY --from=builder /app/docs ./docs  

# Binary'nin ve gerekli dosyaların izinlerini ayarla
RUN chmod +x /root/todo-app && \
    chmod 644 /root/model.conf /root/policy.csv

# Portu expose et
EXPOSE 8080

# Uygulamayı başlat
CMD ["./todo-app"]