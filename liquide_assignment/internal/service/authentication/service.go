package authentication

import (
	"context"
	"errors"
	"fmt"
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
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
			err := config.Wrap(errors.New("Authorization header required"), config.ErrInvalidToken)
			err.Message = "Authorization header required"
			response.HandleError(w, err)
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(authHeader, prefix) {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		tokenParsed, err := s.VerifyJWT(token)
		if err != nil {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		platform := tokenParsed.Claims.(jwt.MapClaims)["platform"].(string)
		userId := tokenParsed.Claims.(jwt.MapClaims)["user_id"].(string)
		sessionId := tokenParsed.Claims.(jwt.MapClaims)["session_id"].(string)
		expireAt := tokenParsed.Claims.(jwt.MapClaims)["exp_at"].(float64)

		key := "user:session:" + userId + ":" + platform

		keyValueMap, err := s.redisRepo.HGetAll(key)
		if err != nil {
			if err == rueidis.Nil {
				response.HandleError(w, config.ErrInvalidToken)
				return
			} else {
				response.HandleError(w, config.ErrInternalServerErrorRedis)
				return
			}
		}

		if _, ok := keyValueMap[sessionId]; !ok {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		if keyValueMap[sessionId] != token {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		if time.Unix(int64(expireAt), 0).Before(time.Now()) {
			response.HandleError(w, config.ErrInvalidToken)
			return
		}

		httpCTX := context.WithValue(r.Context(), "session_id", token)
		next.ServeHTTP(w, r.WithContext(httpCTX))
	})
}

func (s *service) VerifyJWT(tokenString string) (*jwt.Token, error) {

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(s.config.SessionConfig.Secret), nil
	})
}
