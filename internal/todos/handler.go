
package todos
import (
    "context"
    "encoding/json"
    "net/http"
    "strconv" 
    "github.com/go-chi/chi/v5"
    "todo-app-go/ent"
    "todo-app-go/ent/user"  // User modeli ve sorguları için gerekli

)

// @Summary Create a new Todo
// @Description Bu endpoint, yeni bir todo eklemek için kullanılır.
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body ent.Todo true "Todo nesnesi"
// @Success 201 {string} string "Todo başarıyla eklendi"
// @Failure 400 {string} string "Geçersiz veri"
// @Failure 500 {string} string "Todo eklenemedi"
// @Router /todos [post]
func CreateTodo(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var todo ent.Todo
        if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
            http.Error(w, "Geçersiz veri", http.StatusBadRequest)
            return
        }

        _, err := client.Todo.
            Create().
            SetTitle(todo.Title).
            SetDescription(todo.Description).
            Save(context.Background())
        if err != nil {
            http.Error(w, "Todo eklenemedi", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("Todo başarıyla eklendi"))
    }
}

// @Summary Get all Todos
// @Description Bu endpoint, tüm todo'ları listeler.
// @Tags todos
// @Produce  json
// @Success 200 {array} ent.Todo "Todo'lar başarıyla listelendi"
// @Failure 500 {string} string "Todo'lar listelenemedi"
// @Router /todos [get]
func GetTodos(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        todos, err := client.Todo.Query().All(context.Background())
        if err != nil {
            http.Error(w, "Todo'lar listelenemedi", http.StatusInternalServerError)
            return
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
// @Success 200 {object} ent.Todo "Belirli todo başarıyla getirildi"
// @Failure 400 {string} string "Geçersiz ID"
// @Failure 404 {string} string "Todo bulunamadı"
// @Failure 500 {string} string "Veritabanı hatası"
// @Router /todos/{id} [get]
func GetTodoByID(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)  // string'i integer'a çeviriyoruz
        if err != nil {
            http.Error(w, "Geçersiz ID", http.StatusBadRequest)
            return
        }

        todo, err := client.Todo.Get(context.Background(), id)
        if err != nil {
            http.Error(w, "Todo bulunamadı", http.StatusNotFound)
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
// @Param todo body ent.Todo true "Todo nesnesi"
// @Success 200 {string} string "Todo başarıyla güncellendi"
// @Failure 400 {string} string "Geçersiz veri"
// @Failure 500 {string} string "Todo güncellenemedi"
// @Router /todos/{id} [put]
func UpdateTodo(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)  // string'i integer'a çeviriyoruz
        if err != nil {
            http.Error(w, "Geçersiz ID", http.StatusBadRequest)
            return
        }

        var todo ent.Todo
        if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
            http.Error(w, "Geçersiz veri", http.StatusBadRequest)
            return
        }

        _, err = client.Todo.
            UpdateOneID(id).
            SetTitle(todo.Title).
            SetDescription(todo.Description).
            Save(context.Background())
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
// @Failure 400 {string} string "Geçersiz ID"
// @Failure 500 {string} string "Todo silinemedi"
// @Router /todos/{id} [delete]
func DeleteTodo(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := chi.URLParam(r, "id")
        id, err := strconv.Atoi(idStr)  // string'i integer'a çeviriyoruz
        if err != nil {
            http.Error(w, "Geçersiz ID", http.StatusBadRequest)
            return
        }

        err = client.Todo.DeleteOneID(id).Exec(context.Background())
        if err != nil {
            http.Error(w, "Todo silinemedi", http.StatusInternalServerError)
            return
        }

        w.Write([]byte("Todo başarıyla silindi"))
    }
}
// RegisterUser godoc
// @Summary Register a new user
// @Description This endpoint registers a new user with a role.
// @Tags User
// @Accept  json
// @Produce  json
// @Param registerUser body object{username=string,password=string,role=string} true "User Registration Data"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /users/register [post]
func RegisterUser(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
            Role     string `json:"role"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        _, err := client.User.
            Create().
            SetUsername(req.Username).
            SetPassword(req.Password).
            SetRole(req.Role).
            Save(context.Background())

        if err != nil {
            http.Error(w, "User could not be registered", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        w.Write([]byte("User registered successfully"))
    }
}

// LoginUser godoc
// @Summary Log in a user
// @Description This endpoint logs in a user with their username and password.
// @Tags User
// @Accept  json
// @Produce  json
// @Param loginUser body object{username=string,password=string} true "User Login Data"
// @Success 200 {object} ent.User "Logged in successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /users/login [post]
func LoginUser(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }

        user, err := client.User.
            Query().
            Where(user.UsernameEQ(req.Username), user.PasswordEQ(req.Password)).
            Only(context.Background())

        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(user)
    }
}
// @Summary Get all User
// @Description Bu endpoint, tüm user'ları listeler.
// @Tags User
// @Produce  json
// @Success 200 {array} ent.Todo "user'lar başarıyla listelendi"
// @Failure 500 {string} string "user'lar listelenemedi"
// @Router /users [get]
func GetUsers(client *ent.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := client.User.Query().All(context.Background())
        if err != nil {
            http.Error(w, "user'lar listelenemedi", http.StatusInternalServerError)
            return
        }

        if err := json.NewEncoder(w).Encode(users); err != nil {
            http.Error(w, "Veriler encode edilemedi", http.StatusInternalServerError)
        }
    }
}