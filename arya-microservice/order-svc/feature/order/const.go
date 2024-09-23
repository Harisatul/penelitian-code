package order

import "time"

const (
	orderCancellationDuration      = 7 * time.Second
	CreateOrderTopic               = "order.create_order"
	CancelOrderTopic               = "order.cancel_order"
	CompleteOrderTopic             = "order.complete_order"
	NotifyCancelOrderTopic         = "scheduler.notify_cancel_order"
	NotifyCancelOrderConsumerGroup = "order.notify_cancel_order"
)
