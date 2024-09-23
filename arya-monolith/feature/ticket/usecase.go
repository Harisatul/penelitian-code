package ticket

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"monolith/feature/shared"
	"monolith/pkg"
)

func UpdateListUnreservedTickets(ctx context.Context) {
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
