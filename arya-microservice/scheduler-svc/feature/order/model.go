package order

type createOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket createOrderTicketPayload `json:"ticket"`
}

type createOrderTicketPayload struct {
	ID         uint32 `json:"id"`
	CategoryID uint8  `json:"category_id"`
	Version    uint32 `json:"version"`
}

type notifyCancelOrderMessage struct {
	ID uint64 `json:"id"`
}

type completeOrderMessage struct {
	ID uint64 `json:"id"`
}
