package order

import "liquide_assignment/internal/config"

type OrderRequest struct {
	UserId      string             `json:"_"`
	Name        string             `json:"name"`
	ISIN        string             `json:"isin"`
	Quantity    int                `json:"quantity"`
	Price       float32            `json:"price"`
	OrderType   string             `json:"order_type"`
	OrderStatus config.OrderStatus `json:"-"`
	OrderId     string             `json:"_"`
}

type OrderResponse struct {
	Message string `json:"message"`
}

func (o *OrderRequest) Validate() bool {
	if o.Name == "" || o.ISIN == "" || o.Quantity == 0 || o.Price == 0 {
		return false
	}
	return true
}
