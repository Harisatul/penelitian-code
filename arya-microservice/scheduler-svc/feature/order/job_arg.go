package order

import (
	"github.com/riverqueue/river"
	"time"
)

type CancellationArgs struct {
	ID uint64 `json:"id"`
}

func (CancellationArgs) Kind() string                 { return "order_cancellation" }
func (CancellationArgs) InsertOpts() river.InsertOpts { return river.InsertOpts{MaxAttempts: 2} }

type CancellationWorker struct {
	river.WorkerDefaults[CancellationArgs]
}

func (w *CancellationWorker) Timeout(*river.Job[CancellationArgs]) time.Duration {
	return 3 * time.Second
}
