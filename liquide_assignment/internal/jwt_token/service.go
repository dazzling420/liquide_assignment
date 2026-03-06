package jwttoken

import (
	"errors"
	"fmt"
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *service) GetInstance() *service {
	return s
}

func (s *service) IssueToken(reqClaims map[string]interface{}) (string, error) {

	// Create the Claims
	claims := &jwt.MapClaims{
		"user_id":    reqClaims["user_id"].(string),
		"platform":   reqClaims["platform"].(string),
		"device_id":  reqClaims["device_id"].(string),
		"session_id": reqClaims["session_id"].(string),
		"in_at":      reqClaims["in_at"].(int64),
		"exp_at":     reqClaims["exp_at"].(int64),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(s.config.SessionConfig.Secret)
	ss, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return ss, err
}

func (s *service) GetExpiry(tokenString string) (time.Time, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SessionConfig.Secret), nil
	})

	if err != nil {
		return time.Time{}, config.ErrSomethingWentWrong
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp_at"].(float64); ok {
			return time.Unix(int64(exp), 0), nil
		}
	}

	return time.Time{}, fmt.Errorf("no exp found")
}

func (s *service) GetParsedJWT(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := config.Wrap(errors.New("Authorization header required"), config.ErrInvalidToken)
		err.Message = "Authorization header required"
		response.HandleError(w, err)
		return nil, err
	}

	const prefix = "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		response.HandleError(w, config.ErrInvalidToken)
		return nil, config.ErrInvalidToken
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		response.HandleError(w, config.ErrInvalidToken)
		return nil, config.ErrInvalidToken
	}

	tokenParsed, err := s.VerifyJWT(token)
	if err != nil {
		response.HandleError(w, config.ErrInvalidToken)
		return nil, config.ErrInvalidToken
	}

	platform := tokenParsed.Claims.(jwt.MapClaims)["platform"].(string)
	userId := tokenParsed.Claims.(jwt.MapClaims)["user_id"].(string)
	sessionId := tokenParsed.Claims.(jwt.MapClaims)["session_id"].(string)
	expireAt := tokenParsed.Claims.(jwt.MapClaims)["exp_at"].(float64)

	var claims = map[string]interface{}{
		"platform":   platform,
		"user_id":    userId,
		"session_id": sessionId,
		"exp_at":     expireAt,
	}

	return claims, nil
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
