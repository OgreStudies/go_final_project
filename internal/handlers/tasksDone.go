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
		if errResponceIfError(err, res, http.StatusBadRequest, "ошибка идентификатора задачи") {
			return
		}

		storedTask, err := todoStorage.GetTaskById(taskId)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		//Если правило повторения пустое - удалить задачу
		if storedTask.Repeat == "" {
			err = todoStorage.DeleteTask(taskId)
			if errResponceIfError(err, res, http.StatusInternalServerError, "") {
				return
			}

		} else { //Если правило повторения не пустое - перенести задачу на следующую дату выполнения
			storedTask.Date, err = tasks.NextDate(time.Now(), storedTask.Date, storedTask.Repeat)
			if errResponceIfError(err, res, http.StatusInternalServerError, fmt.Sprintf("ошибка вычисления новой даты события с id: %v", taskId)) {
				return
			}

			//Проверка и коррекция данных для обновления
			checkedTask, err := storedTask.TaskFieldCheckAndCorrect()
			if err != nil {
				errResponceIfError(fmt.Errorf("запись с id: %v ошибка формата данных для обновления: %w", checkedTask.ID, err), res, http.StatusInternalServerError, "")
				return
			}

			err = todoStorage.UpdateTask(taskId, &checkedTask)
			if errResponceIfError(err, res, http.StatusInternalServerError, "") {
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
