package order

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"net/http"
)

type OrderService interface {
	GetLogger() logger.Service
	GetConfig() *config.Config
	GetJWTRepoInstance() jwttoken.Service

	CreateOrder(request *json.Decoder, claims map[string]interface{}) (*OrderResponse, error)
	CheckRateLimited(request *http.Request) (int64, error)
	AddRateLimiting(request *http.Request)
}

type service struct {
	logger    logger.Service
	config    *config.Config
	redisRepo RedisRepository
	mongoRepo MongoRepository
	jwtRepo   jwttoken.Service
}

func InitOrderService(logger logger.Service, config *config.Config, redisRepo RedisRepository, mongoRepo MongoRepository, jwtRepo jwttoken.Service) *service {
	return &service{
		logger,
		config,
		redisRepo,
		mongoRepo,
		jwtRepo,
	}
}
