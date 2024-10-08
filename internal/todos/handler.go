package todos

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// @Summary Create a new Todo
// @Description Bu endpoint, yeni bir todo eklemek için kullanılır.
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body Todo true "Todo nesnesi"
// @Success 201 {object} string "Todo başarıyla eklendi"
// @Failure 400 {string} string "Geçersiz veri"
// @Failure 500 {string} string "Todo eklenemedi"
// @Router /todos [post]
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

// @Summary Get all Todos
// @Description Bu endpoint tüm todo'ları listeler.
// @Tags todos
// @Produce  json
// @Success 200 {array} Todo "Todo'lar başarıyla listelendi"
// @Failure 500 {string} string "Todo'lar listelenemedi"
// @Router /todos [get]
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

// @Summary Get a Todo by ID
// @Description Bu endpoint, belirli bir todo'yu getirir.
// @Tags todos
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} Todo "Belirli todo başarıyla getirildi"
// @Failure 404 {string} string "Todo bulunamadı"
// @Failure 500 {string} string "Veritabanı hatası"
// @Router /todos/{id} [get]
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

// @Summary Update a Todo
// @Description Bu endpoint, belirli bir todo'yu günceller.
// @Tags todos
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param todo body Todo true "Todo nesnesi"
// @Success 200 {string} string "Todo başarıyla güncellendi"
// @Failure 400 {string} string "Geçersiz veri"
// @Failure 500 {string} string "Todo güncellenemedi"
// @Router /todos/{id} [put]
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

// @Summary Delete a Todo
// @Description Bu endpoint, belirli bir todo'yu siler.
// @Tags todos
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {string} string "Todo başarıyla silindi"
// @Failure 500 {string} string "Todo silinemedi"
// @Router /todos/{id} [delete]
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
