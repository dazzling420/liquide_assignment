package main

import (
	"fmt"
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/http/rest"
	jwttoken "liquide_assignment/internal/jwt_token"
	"liquide_assignment/internal/logger"
	"liquide_assignment/internal/service/authentication"
	"liquide_assignment/internal/service/login"
	"liquide_assignment/internal/service/order"
	"liquide_assignment/internal/service/report"
	"liquide_assignment/internal/storage/database"
	mongodb "liquide_assignment/internal/storage/database/mongoDb"
	"liquide_assignment/internal/storage/database/redis"
	"net/http"
	"os"
	"strconv"
)

func main() {
	config := config.GetConfig()
	loggerService := logger.Init(config.LoggerConfig)

	loggerService.Info("=========================== LOGS START HERE ===========================")
	loggerService.Info("Logger initialized")

	// Mongo Connections
	mongoClient, ctx, err := database.MongoConnect(config.MongoDbConfig)
	if err != nil {
		loggerService.Error("Failed to connect to MongoDB!! ", err)
		panic("Failed to connect to MongoDB!!")
	}
	loggerService.Info("Mongo client initialized")

	// Redis Connection
	redisClient, err := database.RedisConnect(config.SMRedisConfig)
	if err != nil {
		loggerService.Error("Failed to connect to Redis!! ", err)
		panic("Failed to connect to Redis!!")
	}
	loggerService.Info("Redis client initialized")

	// Creating Repos
	mongoRepo := mongodb.InitMongoRepo(mongoClient, ctx, config)
	redisRepo := redis.InitRedisRepo(redisClient, config)
	jwtRepo := jwttoken.InitJWTTokenService(config)

	// Creating Services
	loginService := login.InitLoginService(loggerService, config, redisRepo, mongoRepo, jwtRepo)
	orderService := order.InitOrderService(loggerService, config, redisRepo, mongoRepo, jwtRepo)
	reportService := report.InitReportService(loggerService, config, redisRepo, mongoRepo, jwtRepo)
	authService := authentication.InitAuthService(loggerService, config, redisRepo, mongoRepo, jwtRepo)

	router := rest.InitHandlerNew(authService, loginService, orderService, reportService)

	loggerService.Info("Initializing done!!!!")

	loggerService.Info("Application listening on " + strconv.Itoa(config.AppConfig.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.AppConfig.Port), router); err != nil {
		loggerService.Error("Unable to start router", err.Error())
		os.Exit(0)
	}
}
