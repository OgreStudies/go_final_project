package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ogrestudies/go_final_project/internal/tasks"
)

// Обработчик запросов к списку задач
func TasksHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		res.Header().Set("Content-Type", "application/json")
		//Список ближайших задач
		var tasks tasks.Tasks
		var err error
		tasks.Tasks, err = todoStorage.GetLastTasks(req.FormValue("search"))
		if errResponceIfError(err, res, http.StatusInternalServerError, "") {
			return
		}

		//Преобразование tasks в json
		resp, err := json.Marshal(&tasks)
		if errResponceIfError(err, res, http.StatusInternalServerError, "") {
			return
		}
		//Формирование заголовка

		res.WriteHeader(http.StatusOK)
		//Запись тела
		_, err = res.Write(resp)
		if err != nil {
			log.Output(1, err.Error())
		}

	}
}
