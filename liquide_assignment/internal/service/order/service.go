package order

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"net/http"
)

type RedisRepository interface {
	CheckRateLimited(key string) (int64, error)
	AddRateLimiting(key string, limit int)
}

type MongoRepository interface {
	AddOrderEntry(orderRequest OrderRequest) (bool, error)
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

func (s *service) CreateOrder(request *json.Decoder, claims map[string]interface{}) (*OrderResponse, error) {
	var orderRequest OrderRequest
	err := request.Decode(&orderRequest)
	if err != nil {
		s.logger.Error("Invalid request", err.Error())
		err = config.Wrap(err, config.ErrInvalidRequest)
		return nil, err
	}
	if !orderRequest.Validate() {
		s.logger.Error("Invalid order Request")
		return nil, config.ErrInvalidRequest
	}

	orderRequest.OrderId = generateOrderId()
	orderRequest.UserId = claims["user_id"].(string)
	orderRequest.OrderStatus = config.New
	_, err = s.mongoRepo.AddOrderEntry(orderRequest)
	if err != nil {
		s.logger.Error("Failed to add order entry", err.Error())
		return nil, config.Wrap(err, config.ErrDatabaseUpdateErrorMongo)
	}
	return &OrderResponse{Message: "Order created successfully"}, nil
}

func (s *service) CheckRateLimited(request *http.Request) (int64, error) {
	key := "ORDER_RATE_LIMIT:" + request.RemoteAddr
	return s.redisRepo.CheckRateLimited(key)
}

func (s *service) AddRateLimiting(request *http.Request) {
	key := "ORDER_RATE_LIMIT:" + request.RemoteAddr
	s.redisRepo.AddRateLimiting(key, s.config.RateLimitConfig.OrderService)
}
