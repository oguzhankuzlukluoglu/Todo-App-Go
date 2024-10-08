package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"

    _ "todo-app-go/docs" // Swagger docs
    "todo-app-go/internal/todos"
    httpSwagger "github.com/swaggo/http-swagger" // Swagger UI
)

var db *sql.DB

func connectToDB() {
    var err error
    dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Veritabanına bağlanılamadı: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Veritabanına bağlanılamadı: %v", err)
    }
    log.Println("Veritabanına başarıyla bağlanıldı!")
}

// @title Todo App API
// @version 1.0
// @description This is a simple Todo App API.
// @host localhost:8080
// @BasePath /
func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    connectToDB()

    r := chi.NewRouter()

    r.Get("/swagger/*", httpSwagger.WrapHandler)

    r.Post("/todos", todos.CreateTodo(db))
    r.Get("/todos", todos.GetTodos(db))
    r.Get("/todos/{id}", todos.GetTodoByID(db))
    r.Put("/todos/{id}", todos.UpdateTodo(db))
    r.Delete("/todos/{id}", todos.DeleteTodo(db))

    log.Println("Server " + os.Getenv("APP_PORT") + " portunda çalışıyor...")
    http.ListenAndServe(":"+os.Getenv("APP_PORT"), r)
}
