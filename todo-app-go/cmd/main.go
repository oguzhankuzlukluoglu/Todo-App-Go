package main

import (
    "log"
    "net/http"
    "github.com/go-chi/chi/v5"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    r := chi.NewRouter()

    // Sağlık kontrolü endpoint'i
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Uygulama çalışıyor!"))
    })

    log.Println("Server 8080 portunda çalışıyor...")
    http.ListenAndServe(":8080", r)
}
