package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ogrestudies/go_final_project/task"
)

// Обработчик запросов на манипуляцию с отдельной задачей
func taskHandle(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET": //Получить задачу
		res.Header().Set("Content-Type", "application/json")
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

		//Преобразование task в json
		resp, err := json.Marshal(&storedTask)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		//Формирование заголовка
		res.WriteHeader(http.StatusOK)
		//Запись тела
		res.Write(resp)
	case "PUT": //Обновить задачу
		res.Header().Set("Content-Type", "application/json")
		var taskData task.Task
		var buf bytes.Buffer
		//Чтение тела
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		//Преобразование тела в Task
		if err = json.Unmarshal(buf.Bytes(), &taskData); err != nil {

			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		//Добавление задачи
		taskId, err := strconv.ParseInt(taskData.ID, 10, 64)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			return
		}

		_, err = todoStorage.UpdateTask(taskId, &taskData)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, taskId)))

	case "POST": //Добавить задачу
		res.Header().Set("Content-Type", "application/json")
		var newTask task.Task
		var buf bytes.Buffer
		//Чтение тела
		_, err := buf.ReadFrom(req.Body)
		if err != nil {

			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		//Преобразование тела в Task
		if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {

			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		//Добавление задачи
		taskId, err := todoStorage.AddTask(newTask)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}

		//Успешный ответ
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(fmt.Sprintf(`{"id":"%d"}`, taskId)))

	case "DELETE": //Удаление задачи
		res.Header().Set("Content-Type", "application/json")

		taskIdReq := req.FormValue("id")

		taskId, err := strconv.ParseInt(taskIdReq, 10, 64)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			return
		}

		err = todoStorage.DeleteTask(taskId)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("{}"))
		return
	}

}
