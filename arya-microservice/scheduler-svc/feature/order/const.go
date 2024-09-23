package order

import "time"

const (
	orderCancellationDuration  = 7 * time.Second
	orderCacheKey              = "scheduler.order_%d"
	NotifyCancelOrderTopic     = "scheduler.notify_cancel_order"
	CreateOrderTopic           = "order.create_order"
	CreateOrderConsumerGroup   = "scheduler.create_order"
	CompleteOrderTopic         = "order.complete_order"
	CompleteOrderConsumerGroup = "scheduler.complete_order"
)
