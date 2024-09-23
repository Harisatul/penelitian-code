package order

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
	"order-svc/feature/shared"
	"order-svc/pkg"
)

type CancelOrderHandler struct {
}

func (*CancelOrderHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateUpdateDB
		lfState2Status = "state_2_update_db_status"

		lvState3       = shared.LogEventStateKafkaPublish
		lfState3Status = "state_3_publish_message_status"

		lf = []slog.Attr{
			pkg.LogEventName("CancelOrder"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload notifyCancelOrderMessage
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
	| Step 2 : Update order to cancel
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	version, err := updateOrderStatusToCancel(ctx, payload.ID)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))

		if err == errOrderNotFound || err == errTicketNotFound {
			pkg.LogInfoWithContext(ctx, "order not found", lf)
			return
		}

		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	message, err := json.Marshal(cancelOrderMessage{
		ID:     payload.ID,
		Ticket: cancelOrderTicketPayload{Version: version},
	})
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	err = pkg.PublishMessage(kp, CancelOrderTopic, string(message))

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	pkg.LogInfoWithContext(ctx, "success cancel order", lf)
}
