package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ogrestudies/go_final_project/config"
	"github.com/ogrestudies/go_final_project/storage"
	"github.com/ogrestudies/go_final_project/task"
)

// Указатель на хранилище данных
var todoStorage task.TaskStorage

func main() {
	var err error

	//Настройка логгера
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(file)

	//Подключение к хранилищу задач
	db, err := storage.OpenStorage(config.TODODb())

	if err != nil {
		panic(err)
	}
	defer db.Close()
	todoStorage = task.NewTaskstorage(db)

	//Запуск сервера
	webDir := "./web"
	addrString := ":" + config.TODOConfigPort()
	fmt.Println(addrString)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/nextdate", nextDateHandle)
	mux.HandleFunc("/api/task", auth(taskHandle))
	mux.HandleFunc("/api/tasks", auth(tasksHandle))
	mux.HandleFunc("/api/task/done", auth(tasksDoneHandle))
	mux.HandleFunc("/api/signin", signHandle)

	//Данная реализация файлового сервера предполагает, что клиент должен отсылать token всегда, включая запрос файйлов
	//Однако реализация тестов финального задания не предполагает такого поведения
	//hfs := http.FileServer(http.Dir(webDir))
	//mux.HandleFunc("/", authHFS(hfs.ServeHTTP))

	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	err = http.ListenAndServe(addrString, nil)
	if err != nil {
		panic(err)
	}
}
