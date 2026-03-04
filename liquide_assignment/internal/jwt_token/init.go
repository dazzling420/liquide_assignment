package jwttoken

import (
	"liquide_assignment/internal/config"
	"time"
)

type Service interface {
	GetInstance() *service
	IssueToken(claims map[string]interface{}) (string, error)
	GetExpiry(token string) (time.Time, error)
	//GetClaim(key string) (string, error)
}

type service struct {
	config *config.Config
}

func InitJWTTokenService(config *config.Config) *service {
	return &service{
		config: config,
	}
}
