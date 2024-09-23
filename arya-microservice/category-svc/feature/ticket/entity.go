package ticket

type ticketEntity struct {
	ID         uint32 `json:"id"`
	CategoryID uint8  `json:"category_id"`
	OrderID    uint64 `json:"order_id"`
	Version    uint32 `json:"version"`
}
