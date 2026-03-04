package order

import "encoding/json"

type RedisRepository interface {
}

type MongoRepository interface {
}

func (s *service) CreateOrder(request *json.Decoder) (OrderResponse, error) {
	return OrderResponse{}, nil
}
