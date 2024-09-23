package ticket

import (
	"category-svc/feature/shared"
	"category-svc/pkg"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type CreateOrderHandler struct {
}

func (*CreateOrderHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_order_status"

		lf = []slog.Attr{
			pkg.LogEventName("CreateOrder"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload createOrderMessage
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogWarnWithContext(ctx, "failed parse payload", err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payload),
	)

	/*------------------------------------
	| Step 2 : Update order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	ticket := ticketEntity{
		ID:      payload.Ticket.ID,
		OrderID: payload.ID,
		Version: payload.Ticket.Version - 1,
	}

	err = updateTicket(ctx, ticket)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))

		if err == errTicketNotFound {
			pkg.LogWarnWithContext(ctx, "ticket not found", err, lf)
			return
		}

		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	pkg.LogInfoWithContext(ctx, "success create order", lf)
}

type UpdateListUnreservedTicketHandler struct {
}

func (*UpdateListUnreservedTicketHandler) Handle(ctx context.Context, _ *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateFetchDB
		lfState1Status = "state_1_fetch_ticket_status"

		lvState2       = shared.LogEventStateSetCache
		lfState2Status = "state_2_set_cache_status"

		lf = []slog.Attr{
			pkg.LogEventName("UpdateListUnreservedTickets"),
		}
	)

	/*------------------------------------
	| Step 1 : Fetch ticket
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	categoryIds, total, err := countAvailableTicketByCategoryGroup(ctx)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	response := make([]listUnreservedTicketResponse, len(categoryIds))
	for i, categoryId := range categoryIds {
		response[i] = listUnreservedTicketResponse{
			ID:    categoryId,
			Total: total[i],
		}
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState1Status))

	/*------------------------------------
	| Step 2 : Set cache
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	bytes, err := json.Marshal(response)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	err = cache.Set(ctx, listUnreservedTicketsCacheKey, string(bytes), redis.KeepTTL).Err()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}
}

type CancelOrderHandler struct {
}

func (*CancelOrderHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_order_status"

		lf = []slog.Attr{
			pkg.LogEventName("CancelOrder"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload cancelOrderMessage
	err := json.Unmarshal(msg.Value, &payload)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf,
		pkg.LogStatusSuccess(lfState1Status),
		pkg.LogEventPayload(payload),
	)

	/*------------------------------------
	| Step 2 : Update order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	ticket := ticketEntity{
		OrderID: payload.ID,
		Version: payload.Ticket.Version - 1,
	}

	err = updateTicketOrderToNull(ctx, ticket)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))

		if err == errTicketNotFound {
			pkg.LogWarnWithContext(ctx, "ticket not found", err, lf)
			return
		}

		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	pkg.LogInfoWithContext(ctx, "success cancel order", lf)
}
