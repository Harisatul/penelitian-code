package order

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func isUserHasOrder(ctx context.Context, email string) (isExist bool, err error) {
	q := `SELECT EXISTS(SELECT 1 FROM orders WHERE email = $1 AND status IN ('created', 'completed'));`

	err = db.QueryRow(ctx, q, email).Scan(&isExist)
	if err != nil {
		return
	}

	return
}

func insertOrder(ctx context.Context, tx pgx.Tx, order orderEntity) (err error) {
	q := `
		WITH selected_ticket AS (
			SELECT id
			FROM tickets
			WHERE category_id = $1 AND order_id IS NULL
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		),

		inserted_order AS (
			INSERT INTO orders (id, category_id, job_id, email)
			VALUES ($2, $1, $3, $4)
		)

		UPDATE tickets
		SET order_id = $2
		WHERE tickets.id = (SELECT id FROM selected_ticket);
	`

	res, err := tx.Exec(ctx, q, order.CategoryID, order.ID, order.JobID, order.Email)
	if err != nil {
		return
	}

	if res.RowsAffected() == 0 {
		err = errTicketNotFound
	}

	return
}

func updateOrderStatusToSuccess(ctx context.Context, id uint64) (jobId int64, err error) {
	q := `UPDATE orders SET status = 'completed', updated_at = CURRENT_TIMESTAMP 
              WHERE id = $1 AND status = 'created' RETURNING job_id;`

	err = db.QueryRow(ctx, q, id).Scan(&jobId)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errOrderNotFound
		}
	}

	return
}

func updateOrderStatusToCancel(ctx context.Context, id uint64) (err error) {
	q1 := `UPDATE orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'created';`
	q2 := `UPDATE tickets SET order_id = NULL WHERE order_id = $1;`

	batch := &pgx.Batch{}
	batch.Queue(q1, id)
	batch.Queue(q2, id)

	br := db.SendBatch(ctx, batch)
	defer br.Close()

	cmd, err := br.Exec()
	if err != nil {
		return
	}

	if cmd.RowsAffected() == 0 {
		err = errOrderNotFound
		return
	}

	cmd, err = br.Exec()
	if err != nil {
		return
	}

	if cmd.RowsAffected() == 0 {
		err = errTicketNotFound
	}

	return
}
