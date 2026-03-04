package login

import (
	"encoding/json"
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/logger"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/redis/rueidis"
)

type RedisRepository interface {
	AddSession(key string, session, sessionToken string) error
	CheckRateLimited(key string) (int64, error)
	AddRateLimiting(key string, limit int)
	GetActiveSessions(key string) (map[string]string, error)
	DeleteSession(key, session string) error
}

type MongoRepository interface {
	AddUser(signupRequest SignupRequest) error
	GetUser(loginRequest LoginRequest) (SignupRequest, error)
	DoesUserExist(pan string) bool
}

func (s *service) GetLogger() logger.Service {
	return s.logger
}

func (s *service) GetConfig() *config.Config {
	return s.config
}

func (s *service) Signup(request *json.Decoder) (*Response, error) {
	var signupRequest SignupRequest
	err := request.Decode(&signupRequest)
	if err != nil {
		s.logger.Error("Invalid request", err.Error())
		err = config.Wrap(err, config.ErrInvalidRequest)
		return nil, err
	}
	isUserValid := checkUserValidity(&signupRequest)
	if !isUserValid {
		s.logger.Error("Invalid user")
		return nil, config.ErrInvalidUser
	}

	if s.mongoRepo.DoesUserExist(signupRequest.Pan) {
		s.logger.Error("User already exists")
		return nil, config.ErrUserAlreadyExists
	}

	signupRequest.UserId = uuid.New().String()

	signupRequest.Password, err = HashPassword(signupRequest.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", err)
		err = config.Wrap(err, config.ErrPasswordHashError)
		return nil, err
	}
	err = s.mongoRepo.AddUser(signupRequest)
	if err != nil {
		s.logger.Error("Failed to add user to database", err)
		err = config.Wrap(err, config.ErrDatabaseUpdateErrorMongo)
		return nil, err
	}
	return &Response{Message: "User added Successfully", Status: config.Success, ErrorCode: ""}, nil
}

func (s *service) Login(request *json.Decoder) (*Response, error) {
	var loginRequest LoginRequest
	err := request.Decode(&loginRequest)
	if err != nil {
		s.logger.Error("Invalid request", err.Error())
		err = config.Wrap(err, config.ErrInvalidRequest)
		return nil, err
	}

	userDetails, err := s.mongoRepo.GetUser(loginRequest)
	if err != nil {
		s.logger.Error("Failed to get user from database", err)
		return nil, config.ErrInvalidUser
	}

	if !CheckPassword(userDetails.Password, loginRequest.Password) {
		s.logger.Error("Invalid Password")
		return nil, config.ErrInvalidPassword
	}

	activeSessionsMap, err := s.GetActiveSessions(userDetails.UserId, loginRequest.Platform)
	if err != nil {
		s.logger.Error("Failed to get active sessions from database", err)
		err = config.Wrap(err, config.ErrDatabaseUpdateErrorMongo)
		return nil, err
	}

	if activeSessionsMap != nil {
		if len(activeSessionsMap) >= s.config.SessionConfig.SessionsPerPlatform {
			list, err := s.DeleteExpiredAndOldSessions(activeSessionsMap, userDetails.UserId, loginRequest.Platform)
			if err != nil {
				s.logger.Error("Failed to delete expired and old sessions", err)
				return nil, err
			}

			for len(list) >= s.config.SessionConfig.SessionsPerPlatform {
				key := "user:session:" + userDetails.UserId + ":" + loginRequest.Platform
				s.redisRepo.DeleteSession(key, list[0].Token)
				list = list[1:]
			}
		}
	}

	//create session
	claimsMap := s.GetJWTClaimsMap(userDetails.UserId, loginRequest.Platform, loginRequest.DeviceId)
	sessionToken, err := s.jwtRepo.IssueToken(claimsMap)
	if err != nil {
		s.logger.Error("Failed to issue token", err)
		err = config.Wrap(err, config.ErrSessionIssueError)
		e := err.(config.Errors)
		e.Message = "Failed to issue token, Please try again after some time"
		return nil, e
	}

	key := "user:session:" + userDetails.UserId + ":" + loginRequest.Platform
	err = s.redisRepo.AddSession(key, claimsMap["session_id"].(string), sessionToken)
	if err != nil {
		s.logger.Error("Failed to add session to database", err)
		return nil, config.ErrDatabaseUpdateErrorRedis
	}
	return &Response{Message: "User Logged In Successfully", Status: config.Success, ErrorCode: ""}, nil
}

func (s *service) CheckRateLimited(request *http.Request) (int64, error) {
	key := "LOGIN_RATE_LIMIT:" + request.RemoteAddr
	return s.redisRepo.CheckRateLimited(key)
}

func (s *service) AddRateLimiting(request *http.Request) {
	key := "LOGIN_RATE_LIMIT:" + request.RemoteAddr
	s.redisRepo.AddRateLimiting(key, s.config.RateLimitConfig.LoginService)
}

func (s *service) GetJWTClaimsMap(userId, platform, deviceId string) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    userId,
		"platform":   platform,
		"device_id":  deviceId,
		"session_id": generateReqToken(),
		"in_at":      time.Now().Unix(),
		"exp_at":     time.Now().Add(time.Duration(s.config.SessionConfig.ExpiryTime) * time.Minute).Unix(),
	}
}

func (s *service) GetActiveSessions(userId, platform string) (map[string]string, error) {
	sessionMap, err := s.redisRepo.GetActiveSessions("user:session:" + userId + ":" + platform)
	if err != nil {
		if err == rueidis.Nil {
			return nil, nil
		}
		s.logger.Error("Failed to get active sessions from database", err)
		err = config.Wrap(err, config.ErrInternalServerErrorRedis)
		return nil, err
	}

	return sessionMap, nil
}

func (s *service) DeleteExpiredAndOldSessions(activeSessionsMap map[string]string, userId, platform string) ([]tokenWithExp, error) {
	var list []tokenWithExp

	for session, jwtToken := range activeSessionsMap {
		key := "user:session:" + userId + ":" + platform
		expiryTime, err := s.jwtRepo.GetExpiry(jwtToken)
		if err != nil {
			s.logger.Error("Failed to get expiry time for token "+jwtToken, err)
			s.redisRepo.DeleteSession(key, session)
			continue
		}
		if expiryTime.Before(time.Now()) {
			s.redisRepo.DeleteSession(key, session)
		}
		list = append(list, tokenWithExp{jwtToken, expiryTime})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Exp.Before(list[j].Exp)
	})
	return list, nil
}
