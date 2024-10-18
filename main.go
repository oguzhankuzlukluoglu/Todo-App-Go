package main

import (
	"context"
	"log"
	"net/http"
	"todo-app-go/ent"
	"todo-app-go/ent/user"
	"todo-app-go/internal/todos"

	_ "todo-app-go/docs"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
)

func CasbinMiddleware(client *ent.Client, enforcer *casbin.Enforcer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, _, _ := r.BasicAuth()

			currentUser, err := client.User.
				Query().
				Where(user.UsernameEQ(username)).
				Only(context.Background())
			if err != nil {
				http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
				return
			}

			allowed, err := enforcer.Enforce(currentUser.Username, r.URL.Path, r.Method)
			if err != nil || !allowed {
				http.Error(w, "Yetki yok", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// @title Todo App API
// @version 1.0
// @description This is a simple Todo App API.
// @host localhost:8080
// @BasePath /
func main() {

	// Casbin enforcer'ı oluştur
	modelPath := "model.conf"
	policyPath := "policy.csv"

	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		log.Fatalf("Casbin enforcer yapılandırması başarısız: %v", err)
	}

	if err := e.LoadPolicy(); err != nil {
		log.Fatalf("Casbin politikaları yüklenemedi: %v", err)
	}

	log.Println("Casbin yapılandırması başarıyla yüklendi")
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

	r := chi.NewRouter()
	CasbinMiddleware(client, &casbin.Enforcer{})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/users/register", todos.RegisterUser(client))
	r.Post("/users/login", todos.LoginUser(client))
	r.Get("/users", todos.GetUsers(client))

	r.Post("/todos", todos.CreateTodo(client))
	r.Get("/todos", todos.GetTodos(client))
	r.Get("/todos/{id}", todos.GetTodoByID(client))
	r.Put("/todos/{id}", todos.UpdateTodo(client))
	r.Delete("/todos/{id}", todos.DeleteTodo(client))

	log.Println("Sunucu 8080 portunda çalışıyor...")
	http.ListenAndServe(":8080", r)
}
