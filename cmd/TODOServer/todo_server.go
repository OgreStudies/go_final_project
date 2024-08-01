package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ogrestudies/go_final_project/internal/config"
	"github.com/ogrestudies/go_final_project/internal/handlers"
	"github.com/ogrestudies/go_final_project/internal/taskstorage"
)

func main() {
	var err error

	//Настройка логгера
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	//Подключение к хранилищу задач
	db, err := taskstorage.OpenStorage(config.TODODb())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Установка ссылки на хранилище для 'handlers'
	todoStorage := taskstorage.NewTaskstorage(db)
	handlers.SetStorage(&todoStorage)

	//Запуск сервера
	webDir := "./web"
	addrString := ":" + config.TODOConfigPort()
	fmt.Println(addrString)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/nextdate", handlers.NextDateHandle)
	mux.HandleFunc("/api/task", handlers.Auth(handlers.TaskHandle))
	mux.HandleFunc("/api/tasks", handlers.Auth(handlers.TasksHandle))
	mux.HandleFunc("/api/task/done", handlers.Auth(handlers.TasksDoneHandle))
	mux.HandleFunc("/api/signin", handlers.SignHandle)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(addrString, mux)
	if err != nil {
		panic(err)
	}
}
