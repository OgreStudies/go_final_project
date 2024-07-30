package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ogrestudies/go_final_project/task"
)

// Обработчик запросов на выполнение задачи
func tasksDoneHandle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	switch req.Method {
	case "POST":

		taskIdReq := req.FormValue("id")
		taskId, err := strconv.ParseInt(taskIdReq, 10, 64)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			return
		}

		storedTask, err := todoStorage.GetTaskById(taskId)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		//Если правило повторения пустое - удалить задачу
		if storedTask.Repeat == "" {
			err = todoStorage.DeleteTask(taskId)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
				return
			}

		} else { //Если правило повторения не пустое - перенести задачу на следующую дату выполнения
			storedTask.Date, err = task.NextDate(time.Now(), storedTask.Date, storedTask.Repeat)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(fmt.Sprintf(`{"error":"ошибка вычисления новой даты события с id: %v"}`, taskId)))
				return
			}

			_, err = todoStorage.UpdateTask(taskId, &storedTask)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
				return
			}
		}

		//Успешное завершение
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
		return

	}
}
