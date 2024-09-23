package ticket

import (
	"category-svc/feature/shared"
	"category-svc/pkg"
	"encoding/json"
	"log/slog"
	"net/http"
)

func HttpRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tickets", listUnreservedTickets)
}

func listUnreservedTickets(w http.ResponseWriter, r *http.Request) {
	var (
		lvState1       = shared.LogEventStateFetchCache
		lfState1Status = "state_1_fetch_cache_status"

		ctx = r.Context()

		lf = []slog.Attr{
			pkg.LogEventName("ListUnreservedTickets"),
		}
	)

	/*------------------------------------
	| Step 1 : Fetch cache
	* ----------------------------------*/
	lf = append(lf, pkg.LogEventState(lvState1))

	dataCache, err := cache.Get(ctx, listUnreservedTicketsCacheKey).Result()
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		shared.WriteInternalServerErrorResponse(w)
		pkg.LogErrorWithContext(ctx, err, lf)
		return
	}

	var ticketCache []listUnreservedTicketResponse
	err = json.Unmarshal([]byte(dataCache), &ticketCache)
	if err != nil {
		lf = append(lf, pkg.LogStatusFailed(lfState1Status))
		shared.WriteInternalServerErrorResponse(w)
		return
	}

	lf = append(lf, pkg.LogStatusSuccess(lfState1Status))

	response := []listUnreservedTicketResponse{
		{ID: 1, Name: "VIP", Price: 3_800_000},
		{ID: 2, Name: "PLATINUM", Price: 3_400_000},
		{ID: 3, Name: "CAT 1", Price: 2_900_000},
		{ID: 4, Name: "CAT 2", Price: 2_600_000},
		{ID: 5, Name: "CAT 3", Price: 2_100_000},
	}

	for _, ticket := range ticketCache {
		response[ticket.ID-1].Total = ticket.Total
	}

	shared.WriteSuccessResponse(w, http.StatusOK, response)
}
