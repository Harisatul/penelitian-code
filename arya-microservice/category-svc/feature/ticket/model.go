package ticket

type listUnreservedTicketResponse struct {
	ID    uint8  `json:"id"`
	Name  string `json:"name,omitempty"`
	Price uint32 `json:"price,omitempty"`
	Total uint16 `json:"total"`
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

type cancelOrderMessage struct {
	ID     uint64                   `json:"id"`
	Ticket cancelOrderTicketPayload `json:"ticket"`
}

type cancelOrderTicketPayload struct {
	Version uint32 `json:"version"`
}
