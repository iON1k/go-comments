package main

import (
	"comments/pkg/api"
	"comments/pkg/storage/postgres"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем файл окружения
	godotenv.Load()

	db_conn := os.Getenv("DB")
	if db_conn == "" {
		log.Fatal("No environment for DB")
	}

	// Создаем подключение к БД
	store, err := postgres.New(db_conn)
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	// Запускаем API
	api := api.New(store)
	log.Print("Server is starting...")
	http.ListenAndServe(":8080", api.Router())
	log.Print("Server has been stopped.")
}
