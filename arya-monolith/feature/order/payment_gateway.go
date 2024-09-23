package order

import (
	"context"
	"time"
)

func createVirtualAccountPayment(ctx context.Context, orderId uint64, amount uint32) (string, error) {
	// TODO: use payment gateway to create virtual account payment, ex: Midtrans
	time.Sleep(50 * time.Millisecond)
	return "RANDOM-CODE", nil
}
