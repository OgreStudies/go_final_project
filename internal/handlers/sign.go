package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ogrestudies/go_final_project/internal/config"
)

// Обработчик запросов на аутентификацию
func SignHandle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	type authData struct {
		Password string `json:"password"`
	}

	var newAuthData authData
	var buf bytes.Buffer
	//Чтение тела
	_, err := buf.ReadFrom(req.Body)

	if errResponceIfError(err, res, http.StatusBadRequest, "") {
		return
	}

	//Преобразование тела в данные ауткнтификации
	err = json.Unmarshal(buf.Bytes(), &newAuthData)
	if errResponceIfError(err, res, http.StatusBadRequest, "") {
		return
	}

	//Проверка наличия пароля
	password := config.TODOPassword()
	if len(password) == 0 { //пароль отсутствует - возвращаем пустой токен
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(`{"token":""}`))
		if err != nil {
			log.Output(1, err.Error())
		}
		return
	} else { //пароль задан
		if newAuthData.Password != password { //Пароли не совпадают
			errResponceIfError(fmt.Errorf("неверный пароль"), res, http.StatusUnauthorized, "")
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
	if errResponceIfError(err, res, http.StatusUnauthorized, "ошибка создания токена") {
		return
	}

	//Отправка токена клиенту
	res.WriteHeader(http.StatusAccepted)
	_, err = res.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, signedToken)))
	if err != nil {
		log.Output(1, err.Error())
	}

}
