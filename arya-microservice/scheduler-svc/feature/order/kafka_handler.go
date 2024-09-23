package order

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/riverqueue/river"
	"github.com/spf13/cast"
	"log/slog"
	"scheduler-svc/feature/shared"
	"scheduler-svc/pkg"
	"time"
)

type CreateOrderHandler struct {
}

func (*CreateOrderHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateInsertJobQueue
		lfState2Status = "state_2_scheduler_order_cancellation_status"

		lvState3       = shared.LogEventStateSetCache
		lfState3Status = "state_3_insert_order_status"

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
	| Step 2 : Schedule order cancellation
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	jobRes, err := queue.Insert(ctx, CancellationArgs{ID: payload.ID}, &river.InsertOpts{
		ScheduledAt: time.Now().Add(orderCancellationDuration),
	})
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Insert order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	err = cache.Set(ctx, fmt.Sprintf(orderCacheKey, payload.ID), jobRes.Job.ID, redis.KeepTTL).Err()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	pkg.LogInfoWithContext(ctx, "success create order", lf)
}

type CompleteOrderHandler struct {
}

func (*CompleteOrderHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) {
	var (
		lvState1       = shared.LogEventStateDecodeRequest
		lfState1Status = "state_1_decode_message_status"

		lvState2       = shared.LogEventStateFetchCache
		lfState2Status = "state_2_get_order_status"

		lvState3       = shared.LogEventStateCancelJobQueue
		lfState3Status = "state_3_cancel_job_status"

		lf = []slog.Attr{
			pkg.LogEventName("CompleteOrder"),
		}
	)

	/*------------------------------------
	| Step 1 : Decode request
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	var payload completeOrderMessage
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
	| Step 2 : Get order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState2))

	jobIdStr, err := cache.Get(ctx, fmt.Sprintf(orderCacheKey, payload.ID)).Result()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	jobId, err := cast.ToInt64E(jobIdStr)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState2Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState2Status))

	/*------------------------------------
	| Step 3 : Cancel job
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState3))

	_, err = queue.JobCancel(ctx, jobId)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState3Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState3Status))

	pkg.LogInfoWithContext(ctx, "success complete order", lf)
}
