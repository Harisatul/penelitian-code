package order

import "time"

type orderEntity struct {
	ID         uint64
	CategoryID uint8
	Email      string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
