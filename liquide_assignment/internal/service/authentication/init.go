package authentication

import (
	"liquide_assignment/internal/config"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	ValidateSession(next http.Handler) http.Handler
	VerifyJWT(tokenString string) (*jwt.Token, error)
}

type service struct {
	logger    logger.Service
	config    *config.Config
	redisRepo RedisRepository
	mongoRepo MongoRepository
	jwtRepo   jwttoken.Service
}

func InitAuthService(logger logger.Service, config *config.Config, redisRepo RedisRepository, mongoRepo MongoRepository, jwtRepo jwttoken.Service) *service {
	return &service{
		logger,
		config,
		redisRepo,
		mongoRepo,
		jwtRepo,
	}
}
