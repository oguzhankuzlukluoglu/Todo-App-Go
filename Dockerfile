# Go base image
FROM golang:1.22

# Çalışma dizinini ayarla
WORKDIR /app

# Go mod ve Go sum dosyalarını kopyala
COPY go.mod go.sum ./

# Modülleri indir
RUN go mod download

# Kaynak kodunu kopyala
COPY . .

# Uygulamayı derle
RUN go build -o todo-app ./cmd/main.go

# Uygulamayı başlat
CMD ["./todo-app"]
