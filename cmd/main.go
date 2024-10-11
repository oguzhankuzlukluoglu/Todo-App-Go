package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/swaggo/http-swagger"   // Swagger UI
    _ "todo-app-go/docs"               // Swagger docs
    "todo-app-go/ent"                  // Entgo Client
    "todo-app-go/internal/todos"
    _ "github.com/go-sql-driver/mysql" // MySQL Driver
)

// @title Todo App API
// @version 1.0
// @description This is a simple Todo App API.
// @host localhost:8080
// @BasePath /
func main() {
    // Entgo client bağlantısı
    client, err := ent.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+
        "@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=True")

    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer client.Close()

    // Veritabanı bağlantısını doğrula
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("Failed creating schema resources: %v", err)
    }

    // Router'ı ayarla
    r := chi.NewRouter()

    // Swagger endpoint'i
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    // Todo CRUD endpointleri
    r.Post("/todos", todos.CreateTodo(client))
    r.Get("/todos", todos.GetTodos(client))
    r.Get("/todos/{id}", todos.GetTodoByID(client))
    r.Put("/todos/{id}", todos.UpdateTodo(client))
    r.Delete("/todos/{id}", todos.DeleteTodo(client))

    // Sunucuyu başlat
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"  // Varsayılan olarak 8080 portu kullanılır
    }
    log.Println("Server " + port + " portunda çalışıyor...")
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
