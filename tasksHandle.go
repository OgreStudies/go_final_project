package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogrestudies/go_final_project/config"
	"github.com/ogrestudies/go_final_project/task"
)

type Tasks struct {
	Tasks []task.Task `json:"tasks"`
}

// Обработчик запросов к списку задач
func tasksHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		res.Header().Set("Content-Type", "application/json")
		//Список ближайших задач
		var tasks Tasks
		var err error
		tasks.Tasks, err = todoStorage.GetLastTasks(config.TODOTaskListMAX(), req.FormValue("search"))

		if err != nil {

			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		//Преобразование tasks в json
		resp, err := json.Marshal(&tasks)
		if err != nil {

			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		//Формирование заголовка

		res.WriteHeader(http.StatusOK)
		//Запись тела
		res.Write(resp)

	}
}
