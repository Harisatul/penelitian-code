package ticket

const (
	listUnreservedTicketsCacheKey            = "category.list_unreserved_tickets"
	CreateOrderTopic                         = "order.create_order"
	CreateOrderConsumerGroup                 = "category.create_order"
	NotifyUpdateListUnreservedTicketsTopic   = "scheduler.notify_update_list_unreserved_tickets"
	UpdateListUnreservedTicketsConsumerGroup = "category.notify_update_list_unreserved_tickets"
	CancelOrderTopic                         = "order.cancel_order"
	CancelOrderConsumerGroup                 = "category.cancel_order"
)
