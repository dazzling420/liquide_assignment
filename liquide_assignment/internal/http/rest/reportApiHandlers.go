package rest

import (
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
	"liquide_assignment/internal/service/report"
	"net/http"
)

func GetOrderBook(rs report.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := rs.GetJWTRepoInstance().GetInstance().GetParsedJWT(w, r)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}
		orderBook, err := rs.GetOrderBook(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		apiResponse := response.NewAPIResponse(orderBook.Message, config.Success, http.StatusOK)
		apiResponse.Data = orderBook.Data
		response.SendResponse(w, apiResponse)
	}
}
