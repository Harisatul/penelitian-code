package order

import (
	"context"
	"encoding/json"
	"github.com/riverqueue/river"
	"log/slog"
	"scheduler-svc/feature/shared"
	"scheduler-svc/pkg"
)

func (w *CancellationWorker) Work(ctx context.Context, job *river.Job[CancellationArgs]) error {
	var (
		lvState1       = shared.LogEventStatePublishMessage
		lfState1Status = "state_1_publish_message_status"

		lf = []slog.Attr{
			pkg.LogEventName("CancelOrder"),
			pkg.LogEventPayload(job.Args),
		}
	)

	/*------------------------------------
	| Step 1 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	message := notifyCancelOrderMessage{
		ID: job.Args.ID,
	}

	messageByte, err := json.Marshal(message)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return err
	}

	err = pkg.PublishMessage(kp, NotifyCancelOrderTopic, string(messageByte))
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return err
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState1Status))

	pkg.LogInfoWithContext(ctx, "success cancel order", lf)
	return nil
}
