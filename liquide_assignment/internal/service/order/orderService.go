package order

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/logger"
)

type OrderService interface {
	CreateOrder(request *json.Decoder) (OrderResponse, error)
}

type service struct {
	logger    logger.Service
	config    *config.Config
	redisRepo RedisRepository
	mongoRepo MongoRepository
}

func InitOrderService(logger logger.Service, config *config.Config, redisRepo RedisRepository, mongoRepo MongoRepository) *service {
	return &service{
		logger,
		config,
		redisRepo,
		mongoRepo,
	}
}
