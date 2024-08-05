package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// Если err != nil - возвращает true, записывает в res ответ с соответствующим statusCode и сообщением об ошибке errString
// если errString - пустая строка, то сообщение об ошибке == err.Error()
// Если err == nil - возвращает false, ответ не формируется
func errResponceIfError(err error, res http.ResponseWriter, statusCode int, errString string) bool {

	if err != nil {
		res.WriteHeader(statusCode)
		if len(errString) != 0 {
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, errString)))
		} else {
			_, err = res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		}
		if err != nil {
			log.Output(1, err.Error())
		}
		return true

	}
	return false
}
