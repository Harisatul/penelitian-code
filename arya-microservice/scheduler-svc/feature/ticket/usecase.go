package ticket

import (
	"context"
	"log/slog"
	"scheduler-svc/feature/shared"
	"scheduler-svc/pkg"
)

func UpdateListUnreservedTickets(ctx context.Context) {
	var (
		lvState1       = shared.LogEventStatePublishMessage
		lfState1Status = "state_1_publish_message_status"

		lf = []slog.Attr{
			pkg.LogEventName("UpdateListUnreservedTickets"),
		}
	)

	/*------------------------------------
	| Step 1 : Publish message
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	err := pkg.PublishMessage(kp, notifyUpdateListUnreservedTicketsTopic, "")
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		pkg.LogErrorWithContext(ctx, err, lf)
	}
}
