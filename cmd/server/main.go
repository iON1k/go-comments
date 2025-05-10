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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

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
	http.ListenAndServe(":80", api.Router())
}
