package authentication

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/rueidis"
)

type RedisRepository interface {
	GetKey(key string) (string, error)
	HGetAll(key string) (map[string]string, error)
}

type MongoRepository interface {
}

func (s *service) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		tokenParsed, err := s.VerifyJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		platform := tokenParsed.Claims.(jwt.MapClaims)["platform"].(string)
		userId := tokenParsed.Claims.(jwt.MapClaims)["user_id"].(string)
		sessionId := tokenParsed.Claims.(jwt.MapClaims)["session_id"].(string)
		expireAt := tokenParsed.Claims.(jwt.MapClaims)["exp_at"].(time.Time)

		key := "user:session:" + userId + ":" + platform

		keyValueMap, err := s.redisRepo.HGetAll(key)
		if err != nil {
			if err == rueidis.Nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if _, ok := keyValueMap[sessionId]; !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if keyValueMap[sessionId] != token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if expireAt.Before(time.Now()) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		httpCTX := context.WithValue(r.Context(), "jSessionId", token)
		next.ServeHTTP(w, r.WithContext(httpCTX))
	})
}

func (s *service) VerifyJWT(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return s.config.SessionConfig.Secret, nil
	})
}
