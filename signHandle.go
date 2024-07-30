package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ogrestudies/go_final_project/config"
)

// Обработчик запросов на аутентификацию
func signHandle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	type authData struct {
		Password string `json:"password"`
	}

	var newAuthData authData
	var buf bytes.Buffer
	//Чтение тела
	_, err := buf.ReadFrom(req.Body)
	if err != nil {

		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		return
	}
	//Преобразование тела в данные ауткнтификации
	if err = json.Unmarshal(buf.Bytes(), &newAuthData); err != nil {

		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		return
	}

	//Проверка наличия пароля
	password := config.TODOPassword()
	if len(password) == 0 { //пароль отсутствует - возвращаем пустой токен
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"token":""}`))
		return
	} else { //пароль задан
		if newAuthData.Password != password { //Пароли не совпадают
			res.WriteHeader(http.StatusUnauthorized)
			res.Write([]byte(`{"error": "Неверный пароль"}`))
			return
		}
	}

	//Сгенерить токен
	// создаём payload - контрольная сумма пароля
	checksum := sha256.Sum256([]byte(password))
	claims := jwt.MapClaims{
		"checksum": hex.EncodeToString(checksum[:]),
	}

	// создаём jwt и указываем payload
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// получаем подписанный токен
	signedToken, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write([]byte(`{"error": "Ошибка создания токена"}`))
		return
	}

	//Отправка токена клиенту
	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, signedToken)))

}
