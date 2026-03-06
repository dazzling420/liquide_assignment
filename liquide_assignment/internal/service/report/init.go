package report

import (
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
)

type ReportService interface {
	GetLogger() logger.Service
	GetConfig() *config.Config
	GetJWTRepoInstance() jwttoken.Service

	GetOrderBook(claims map[string]interface{}) (*OrderBookResponse, error)
}

type service struct {
	logger    logger.Service
	config    *config.Config
	redisRepo RedisRepository
	mongoRepo MongoRepository
	jwtRepo   jwttoken.Service
}

func InitReportService(logger logger.Service, config *config.Config, redisRepo RedisRepository, mongoRepo MongoRepository, jwtRepo jwttoken.Service) *service {
	return &service{
		logger,
		config,
		redisRepo,
		mongoRepo,
		jwtRepo,
	}
}
