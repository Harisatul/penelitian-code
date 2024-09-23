package order

import (
	"context"
	"github.com/riverqueue/river"
	"log/slog"
	"monolith/feature/shared"
	"monolith/pkg"
)

func (w *CancellationWorker) Work(ctx context.Context, job *river.Job[CancellationArgs]) error {
	var (
		lvState1       = shared.LogEventStateUpdateDB
		lfState1Status = "state_1_update_order_status"

		lf = []slog.Attr{
			pkg.LogEventName("CancelOrder"),
			pkg.LogEventPayload(job.Args),
		}
	)

	/*------------------------------------
	| Step 1 : Update order
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	err := updateOrderStatusToCancel(ctx, job.Args.ID)
	if err != nil {
		if err == errOrderNotFound || err == errTicketNotFound {
			return nil
		}

		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
		return err
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState1Status))

	pkg.LogInfoWithContext(ctx, "success cancel order", lf)
	return err
}
