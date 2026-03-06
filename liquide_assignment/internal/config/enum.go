package config

type Status string
type OrderType string
type OrderStatus string

const (
	Success Status = "success"
	Failure Status = "failure"
)

const (
	Buy  OrderType = "BUY"
	Sell OrderType = "SELL"
)

const (
	New          OrderStatus = "NEW"
	Pending      OrderStatus = "PENDING"
	Modified     OrderStatus = "MODIFIED"
	Cancelled    OrderStatus = "CANCELLED"
	Completed    OrderStatus = "COMPLETED"
	Failed       OrderStatus = "FAILED"
	OrderSuccess OrderStatus = "SUCCESS"
)
