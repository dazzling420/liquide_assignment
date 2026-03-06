package report

type OrderBookResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
