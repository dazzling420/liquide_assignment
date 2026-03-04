package jwttoken

import (
	"fmt"
	"liquide_assignment/internal/config"
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
	ss, err := token.SignedString(s.config.SessionConfig.Secret)
	if err != nil {
		return "", err
	}
	return ss, err
}

func (s *service) GetExpiry(tokenString string) (time.Time, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.config.SessionConfig.Secret, nil
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
