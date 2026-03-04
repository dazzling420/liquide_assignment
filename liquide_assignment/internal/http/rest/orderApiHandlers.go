package rest

import (
	"liquide_assignment/internal/config"
	"liquide_assignment/internal/response"
	"liquide_assignment/internal/service/order"
	"net/http"
)

func CreateOrder(os order.OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inputRequest := getRequestDetails(r)
		resp, err := os.CreateOrder(inputRequest)
		if err != nil {
			response.HandleError(w, err.(config.Errors))
			return
		}

		apiResponse := response.NewAPIResponse(resp.Message, config.Success, http.StatusOK)
		response.SendResponse(w, apiResponse)
	}
}
