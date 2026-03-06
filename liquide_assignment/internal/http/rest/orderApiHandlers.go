package rest

import (
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
	"liquide_assignment/internal/service/order"
	"net/http"
	"strconv"
)

func CreateOrder(os order.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RateLimitRemainingSec, err := os.CheckRateLimited(r)
		if err != nil {
			os.GetLogger().Error("Failed to check rate limit", err)
			err := config.Wrap(err, config.ErrInternalServerErrorRedis)
			response.HandleError(w, err)
			return
		}
		if RateLimitRemainingSec > 0 {
			os.GetLogger().Error("Too many requests")
			err := config.ErrTooManyRequests
			err.Message = "Too many requests, Please retry after " + strconv.FormatInt(RateLimitRemainingSec, 10) + " seconds"
			response.HandleError(w, err)
			return
		}
		os.AddRateLimiting(r)
		claims, err := os.GetJWTRepoInstance().GetInstance().GetParsedJWT(w, r)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}
		inputRequest := getRequestDetails(r)
		resp, err := os.CreateOrder(inputRequest, claims)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}

		apiResponse := response.NewAPIResponse(resp.Message, config.Success, http.StatusOK)
		response.SendResponse(w, apiResponse)
	}
}
