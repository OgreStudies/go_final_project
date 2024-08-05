package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ogrestudies/go_final_project/internal/config"
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
func Auth(next http.HandlerFunc) http.HandlerFunc {
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
