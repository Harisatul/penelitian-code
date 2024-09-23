package order

type createOrderRequest struct {
	Email      string `json:"email"`
	CategoryID uint8  `json:"category_id"`
}

type createOrderResponse struct {
	ID     uint64 `json:"id"`
	Total  uint32 `json:"total"`
	VaCode string `json:"va_code"`
}

type paymentNotificationRequest struct {
	StatusCode        string `json:"status_code"`
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
}

type createOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket createOrderTicketPayload `json:"ticket"`
}

type createOrderTicketPayload struct {
	ID         uint32 `json:"id"`
	CategoryID uint8  `json:"category_id"`
	Version    uint32 `json:"version"`
}

type completeOrderMessage struct {
	ID uint64 `json:"id"`
}

type notifyCancelOrderMessage struct {
	ID uint64 `json:"id"`
}

type cancelOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket cancelOrderTicketPayload `json:"ticket"`
}

type cancelOrderTicketPayload struct {
	Version uint32 `json:"version"`
}
