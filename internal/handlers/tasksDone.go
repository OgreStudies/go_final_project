package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ogrestudies/go_final_project/internal/tasks"
)

// Обработчик запросов на выполнение задачи
func TasksDoneHandle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "POST":

		taskIdReq := req.FormValue("id")
		taskId, err := strconv.ParseInt(taskIdReq, 10, 64)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		storedTask, err := todoStorage.GetTaskById(taskId)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		//Если правило повторения пустое - удалить задачу
		if storedTask.Repeat == "" {
			err = todoStorage.DeleteTask(taskId)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
				if err != nil {
					log.Output(1, err.Error())
				}
				return
			}

		} else { //Если правило повторения не пустое - перенести задачу на следующую дату выполнения
			storedTask.Date, err = tasks.NextDate(time.Now(), storedTask.Date, storedTask.Repeat)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				_, err = res.Write([]byte(fmt.Sprintf(`{"error":"ошибка вычисления новой даты события с id: %v"}`, taskId)))
				if err != nil {
					log.Output(1, err.Error())
				}
				return
			}

			_, err = todoStorage.UpdateTask(taskId, &storedTask)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
				if err != nil {
					log.Output(1, err.Error())
				}
				return
			}
		}

		//Успешное завершение
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("{}"))
		if err != nil {
			log.Output(1, err.Error())
		}
		return

	}
}
