package main

import (
    "context"
    "log"
    "net/http"


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
   

    // Entgo veritabanı bağlantısı
    client, err := ent.Open("mysql", "todo_user:user_password@tcp(mysql:3306)/todoapp?parseTime=True")
    if err != nil {
        log.Fatalf("Veritabanı bağlantısı başarısız: %v", err)
    }
    defer client.Close()

    // Veritabanı şemasını oluştur
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("Şema oluşturulamadı: %v", err)
    }

    // Router oluştur
    r := chi.NewRouter()

    // Casbin middleware'i ekle

    // Swagger handler'ı ekle
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    // Todo rotaları
    r.Post("/users/register", todos.RegisterUser(client))
    r.Post("/users/login", todos.LoginUser(client))
    r.Get("/users",todos.GetUsers(client))
    

    r.Post("/todos", todos.CreateTodo(client))
    r.Get("/todos", todos.GetTodos(client))
    r.Get("/todos/{id}", todos.GetTodoByID(client))
    r.Put("/todos/{id}", todos.UpdateTodo(client))
    r.Delete("/todos/{id}", todos.DeleteTodo(client))

    // Sunucuyu başlat
    log.Println("Sunucu 8080 portunda çalışıyor...")
    http.ListenAndServe(":8080", r)
}
