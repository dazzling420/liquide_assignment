package report

import (
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"liquide_assignment/internal/service/order"
)

type RedisRepository interface {
	CheckRateLimited(key string) (int64, error)
	AddRateLimiting(key string, limit int)
}

type MongoRepository interface {
	GetOrderBook(userId string) ([]order.OrderRequest, error)
}

func (s *service) GetLogger() logger.Service {
	return s.logger
}

func (s *service) GetConfig() *config.Config {
	return s.config
}

func (s *service) GetJWTRepoInstance() jwttoken.Service {
	return s.jwtRepo
}

func (s *service) GetOrderBook(claims map[string]interface{}) (*OrderBookResponse, error) {
	userId := claims["user_id"].(string)
	resp, err := s.mongoRepo.GetOrderBook(userId)
	if err != nil {
		s.logger.Error("Failed to fetch order book", err.Error())
		return nil, config.Wrap(err, config.ErrInternalServerErrorMongo)
	}
	return &OrderBookResponse{Message: "Order book fetched successfully", Data: resp}, nil
}
