package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func connectToDB() {
	var err error
	for {
		dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
		db, err = sql.Open("mysql", dsn)
		if err == nil && db.Ping() == nil {
			log.Println("Veritabanına başarıyla bağlanıldı!")
			break
		}
		log.Println("Veritabanına bağlanılamadı, tekrar deneniyor...")
		time.Sleep(5 * time.Second) 
	}
}

func migrate() {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE
		);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}
	log.Println("Todos tablosu başarıyla oluşturuldu!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectToDB() 
	defer db.Close()

	migrate()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Uygulama çalışıyor ve MySQL'e bağlı!"))
	})

	log.Println("Server " + os.Getenv("APP_PORT") + " portunda çalışıyor...")
	http.ListenAndServe(":"+os.Getenv("APP_PORT"), r)
}
