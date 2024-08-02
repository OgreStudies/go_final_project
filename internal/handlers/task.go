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

		//Преобразование task в json
		resp, err := json.Marshal(&storedTask)
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
		return
	case "PUT": //Обновить задачу
		res.Header().Set("Content-Type", "application/json")
		var taskData tasks.Task
		var buf bytes.Buffer
		//Чтение тела
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}
		//Преобразование тела в Task
		if err = json.Unmarshal(buf.Bytes(), &taskData); err != nil {

			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		//Добавление задачи
		taskId, err := strconv.ParseInt(taskData.ID, 10, 64)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		_, err = todoStorage.UpdateTask(taskId, &taskData)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
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
		if err != nil {

			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}
		//Преобразование тела в Task
		if err = json.Unmarshal(buf.Bytes(), &newTask); err != nil {

			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}
		//Добавление задачи
		taskId, err := todoStorage.AddTask(newTask)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
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

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			_, err = res.Write([]byte(`{"error":"Ошибка идентификатора задачи"}`))
			if err != nil {
				log.Output(1, err.Error())
			}
			return
		}

		err = todoStorage.DeleteTask(taskId)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
			if err != nil {
				log.Output(1, err.Error())
			}
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
