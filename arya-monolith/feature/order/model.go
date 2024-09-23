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
