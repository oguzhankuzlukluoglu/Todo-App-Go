package todos

import (
    "database/sql"
    "encoding/json"
    "net/http"
	"github.com/go-chi/chi/v5"
)


func CreateTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var todo Todo
        if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
            http.Error(w, "Geçersiz veri", http.StatusBadRequest)
            return
        }

        query := `INSERT INTO todos (title, description) VALUES (?, ?)`
        _, err := db.Exec(query, todo.Title, todo.Description)
        if err != nil {
            http.Error(w, "Todo eklenemedi", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("Todo başarıyla eklendi"))
    }
}
func GetTodos(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, title, description FROM todos")
        if err != nil {
            http.Error(w, "Todo'lar listelenemedi", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var todos []Todo
        for rows.Next() {
            var todo Todo
            if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
                http.Error(w, "Veriler okunamadı", http.StatusInternalServerError)
                return
            }
            todos = append(todos, todo)
        }

        if err := json.NewEncoder(w).Encode(todos); err != nil {
            http.Error(w, "Veriler encode edilemedi", http.StatusInternalServerError)
        }
    }
}
func GetTodoByID(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")

        var todo Todo
        query := "SELECT id, title, description FROM todos WHERE id = ?"
        err := db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Description)
        if err == sql.ErrNoRows {
            http.Error(w, "Todo bulunamadı", http.StatusNotFound)
            return
        } else if err != nil {
            http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
            return
        }

        if err := json.NewEncoder(w).Encode(todo); err != nil {
            http.Error(w, "Veri encode edilemedi", http.StatusInternalServerError)
        }
    }
}
func UpdateTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")

        var todo Todo
        if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
            http.Error(w, "Geçersiz veri", http.StatusBadRequest)
            return
        }

        query := "UPDATE todos SET title = ?, description = ? WHERE id = ?"
        _, err := db.Exec(query, todo.Title, todo.Description, id)
        if err != nil {
            http.Error(w, "Todo güncellenemedi", http.StatusInternalServerError)
            return
        }

        w.Write([]byte("Todo başarıyla güncellendi"))
    }
}
func DeleteTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")

        query := "DELETE FROM todos WHERE id = ?"
        _, err := db.Exec(query, id)
        if err != nil {
            http.Error(w, "Todo silinemedi", http.StatusInternalServerError)
            return
        }

        w.Write([]byte("Todo başarıyla silindi"))
    }
}

