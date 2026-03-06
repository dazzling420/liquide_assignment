package rest

import (
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
	"liquide_assignment/internal/service/login"
	"net/http"
	"strconv"
)

func Signup(ls login.LoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RateLimitRemainingSec, err := ls.CheckRateLimited(r)
		if err != nil {
			ls.GetLogger().Error("Failed to check rate limit", err)
			err := config.Wrap(err, config.ErrInternalServerErrorRedis)
			response.HandleError(w, err)
			return
		}
		if RateLimitRemainingSec > 0 {
			ls.GetLogger().Error("Too many requests")
			err := config.ErrTooManyRequests
			err.Message = "Too many requests, Please retry after " + strconv.FormatInt(RateLimitRemainingSec, 10) + " seconds"
			response.HandleError(w, err)
			return
		}
		ls.AddRateLimiting(r)
		inputRequest := getRequestDetails(r)
		resp, err := ls.Signup(inputRequest)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}
		apiResponse := response.NewAPIResponse(resp.Message, config.Success, http.StatusOK)
		response.SendResponse(w, apiResponse)
	}
}

func Login(ls login.LoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RateLimitRemainingSec, err := ls.CheckRateLimited(r)
		if err != nil {
			ls.GetLogger().Error("Failed to check rate limit", err)
			err := config.Wrap(err, config.ErrInternalServerErrorRedis)
			response.HandleError(w, err)
			return
		}
		if RateLimitRemainingSec > 0 {
			ls.GetLogger().Error("Too many requests")
			err := config.ErrTooManyRequests
			err.Message = "Too many requests, Please retry after " + strconv.FormatInt(RateLimitRemainingSec, 10) + " seconds"
			response.HandleError(w, err)
			return
		}
		ls.AddRateLimiting(r)
		inputRequest := getRequestDetails(r)
		resp, err := ls.Login(inputRequest)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}
		apiResponse := response.NewAPIResponse(resp.Message, config.Success, http.StatusOK)
		apiResponse.Data = resp.Data
		response.SendResponse(w, apiResponse)
	}
}
