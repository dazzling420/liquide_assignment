package redis

import (
	"liquide_assignment/internal/config"

	"github.com/redis/rueidis"
)

type RedisRepository struct {
	db     rueidis.Client
	config *config.Config
}

func InitRedisRepo(redisClient rueidis.Client, config *config.Config) RedisRepository {
	return RedisRepository{
		db:     redisClient,
		config: config,
	}
}
