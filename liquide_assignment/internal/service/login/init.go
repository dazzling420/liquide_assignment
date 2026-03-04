package login

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"net/http"
)

type LoginService interface {
	GetLogger() logger.Service
	GetConfig() *config.Config

	Login(request *json.Decoder) (*Response, error)
	Signup(request *json.Decoder) (*Response, error)
	CheckRateLimited(request *http.Request) (int64, error)
	AddRateLimiting(request *http.Request)
	GetJWTClaimsMap(userId, platform, deviceId string) map[string]interface{}
	GetActiveSessions(userId, platform string) (map[string]string, error)
	DeleteExpiredAndOldSessions(activeSessionsMap map[string]string, userId, platform string) ([]tokenWithExp, error)
}

type service struct {
	logger    logger.Service
	config    *config.Config
	redisRepo RedisRepository
	mongoRepo MongoRepository
	jwtRepo   jwttoken.Service
}

func InitLoginService(logger logger.Service, config *config.Config, redisRepo RedisRepository, mongoRepo MongoRepository, jwtRepo jwttoken.Service) *service {
	return &service{
		logger,
		config,
		redisRepo,
		mongoRepo,
		jwtRepo,
	}
}
