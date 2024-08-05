package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ogrestudies/go_final_project/internal/tasks"
)

// Обработчик запросов на манипуляцию с отдельной задачей
func TaskHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": //Получить задачу
		res.Header().Set("Content-Type", "application/json")
		taskIdReq := req.FormValue("id")
		taskId, err := strconv.ParseInt(taskIdReq, 10, 64)
		if errResponceIfError(err, res, http.StatusBadRequest, "ошибка идентификатора задачи") {
			return
		}

		storedTask, err := todoStorage.GetTaskById(taskId)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		//Преобразование task в json
		resp, err := json.Marshal(&storedTask)
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
		return
	case "PUT": //Обновить задачу
		res.Header().Set("Content-Type", "application/json")
		var taskData tasks.Task
		var buf bytes.Buffer
		//Чтение тела
		_, err := buf.ReadFrom(req.Body)
		if errResponceIfError(err, res, http.StatusInternalServerError, "") {
			return
		}

		//Преобразование тела в Task
		err = json.Unmarshal(buf.Bytes(), &taskData)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		//Добавление задачи
		taskId, err := strconv.ParseInt(taskData.ID, 10, 64)
		if errResponceIfError(err, res, http.StatusBadRequest, "ошибка идентификатора задачи") {
			return
		}
		//Проверка и коррекция данных для обновления
		checkedTask, err := taskData.TaskFieldCheckAndCorrect()
		if err != nil {
			errResponceIfError(fmt.Errorf("запись с id: %v ошибка формата данных для обновления: %w", checkedTask.ID, err), res, http.StatusInternalServerError, "")
			return
		}

		err = todoStorage.UpdateTask(taskId, &taskData)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, taskId)))
		if err != nil {
			log.Output(1, err.Error())
		}
		return

	case "POST": //Добавить задачу
		res.Header().Set("Content-Type", "application/json")
		var newTask tasks.Task
		var buf bytes.Buffer
		//Чтение тела
		_, err := buf.ReadFrom(req.Body)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		//Преобразование тела в Task
		err = json.Unmarshal(buf.Bytes(), &newTask)
		if errResponceIfError(err, res, http.StatusBadRequest, "") {
			return
		}

		//Проверка и коррекция данных для добаления
		checkedTask, err := newTask.TaskFieldCheckAndCorrect()
		if err != nil {
			errResponceIfError(fmt.Errorf("запись с id: %v ошибка формата данных для добавления: %w", checkedTask.ID, err), res, http.StatusInternalServerError, "")
			return
		}

		//Добавление задачи
		taskId, err := todoStorage.AddTask(&checkedTask)
		if errResponceIfError(err, res, http.StatusInternalServerError, "") {
			return
		}

		//Успешный ответ
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, taskId)))
		if err != nil {
			log.Output(1, err.Error())
		}

	case "DELETE": //Удаление задачи
		res.Header().Set("Content-Type", "application/json")

		taskIdReq := req.FormValue("id")

		taskId, err := strconv.ParseInt(taskIdReq, 10, 64)
		if errResponceIfError(err, res, http.StatusBadRequest, "ошибка идентификатора задачи") {
			return
		}

		err = todoStorage.DeleteTask(taskId)
		if errResponceIfError(err, res, http.StatusInternalServerError, "") {
			return
		}

		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte("{}"))
		if err != nil {
			log.Output(1, err.Error())
		}
		return
	}

}
