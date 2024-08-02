package handlers

import (
	"encoding/json"
	"fmt"
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

		if err != nil {

			res.WriteHeader(http.StatusInternalServerError)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		//Преобразование tasks в json
		resp, err := json.Marshal(&tasks)
		if err != nil {

			res.WriteHeader(http.StatusInternalServerError)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
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
