package config

import (
	"os"
	"strconv"
)

const defaultPort = "7540"
const defaultDb = "scheduler.db"
const defaultMaxVal = 50

// Возвращает сконфигурированный порт сервера. Если нет конфигурации порта, то вернёт порт по умолчанию
func TODOConfigPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		return defaultPort
	}

	return port
}

// Возвращает сконфигурированное имя базы данных. Если нет конфигурации имени, то вернёт имя по умолчанию
func TODODb() string {
	db := os.Getenv("TODO_DBFILE")
	if db == "" {
		return defaultDb
	}

	return db
}

// Возвращает сконфигурированное максимальное кол-во возвращаемых записей. Если нет конфигурации, то вернёт кол-во возвращаемых записей по умолчанию
func TODOTaskListMAX() int {
	maxVal, err := strconv.ParseInt(os.Getenv("TODO_TASK_LIST_MAX"), 10, 64)
	if err != nil {
		return defaultMaxVal
	}
	return int(maxVal)
}

// Возвращает пароль. Если пароль не задан вернёт пустую строку
func TODOPassword() string {
	return os.Getenv("TODO_PASSWORD")
}
