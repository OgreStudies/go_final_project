package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ogrestudies/go_final_project/config"
)

// Проверяет полученный от клиента Token
func verifyToken(token string) bool {
	//Парсим токен
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.TODOPassword()), nil
	})
	//При ошибке парсинга вернуть false
	if err != nil {
		return false
	}

	if !jwtToken.Valid {
		return false
	}

	//Получаем payload
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}
	//Котрольная сумма пароля
	checksumRaw := sha256.Sum256([]byte(config.TODOPassword()))
	checksum := hex.EncodeToString(checksumRaw[:])
	//Контрольная сумма из Payload
	checksumRcvRaw, ok := claims["checksum"]
	if !ok {
		return false
	}
	checksumRcv, ok := checksumRcvRaw.(string)
	if !ok {
		return false
	}
	//Если контрольные суммы не равны - ошибка
	if checksum != checksumRcv {
		return false
	}

	return true
}

// Проверка утентификации для запросов к API
func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := config.TODOPassword()
		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			// валидация и проверка JWT-токена
			valid := verifyToken(jwt)

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

// Проверка аутентификации для запросов к файловому серверу
// Если токен неверный - переадресация на страницу с вводом пароля
/*
func authHFS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := config.TODOPassword()

		if len(pass) > 0 {
			var jwt string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwt = cookie.Value
			}

			// валидация и проверки JWT-токена
			valid := verifyToken(jwt)
			//Если валидация не прошла - редирект на страницу с логином
			if !valid {
				r.URL.Path = "/login.html"
			}
		}
		next(w, r)
	})
}
*/
