package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/joho/godotenv"
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
    // Load environment variables
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Connect to Entgo client
    client, err := ent.Open("mysql", os.Getenv("DB_USER")+":"+
        os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+
        os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME")+"?parseTime=True")

    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer client.Close()

    // Verify connection to the database
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("Failed creating schema resources: %v", err)
    }

    // Set up the router
    r := chi.NewRouter()

    // Swagger handler
    r.Get("/swagger/*", httpSwagger.WrapHandler)

    // Todo routes using Entgo client
    r.Post("/todos", todos.CreateTodo(client))
    r.Get("/todos", todos.GetTodos(client))
    r.Get("/todos/{id}", todos.GetTodoByID(client))
    r.Put("/todos/{id}", todos.UpdateTodo(client))
    r.Delete("/todos/{id}", todos.DeleteTodo(client))

    // Start server
    log.Println("Server " + os.Getenv("APP_PORT") + " portunda çalışıyor...")
    http.ListenAndServe(":"+os.Getenv("APP_PORT"), r)
}
