package order

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func isUserHasOrder(ctx context.Context, email string) (isExist bool, err error) {
	q := `SELECT EXISTS(SELECT 1 FROM orders WHERE email = $1 AND status IN ('created', 'completed'));`

	err = db.QueryRow(ctx, q, email).Scan(&isExist)
	if err != nil {
		return
	}

	return
}

func insertOrder(ctx context.Context, order orderEntity) (ticketId, version uint32, err error) {
	q := `
		WITH selected_ticket AS (
			SELECT id, version
			FROM tickets
			WHERE category_id = $1 AND order_id IS NULL
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		),

		inserted_order AS (
			INSERT INTO orders (id, category_id, email)
			VALUES ($2, $1, $3)
		)

		UPDATE tickets
		SET order_id = $2,
			version = version + 1
		WHERE tickets.id = (SELECT id FROM selected_ticket)
		RETURNING tickets.id, tickets.version;
	`

	err = db.QueryRow(ctx, q, order.CategoryID, order.ID, order.Email).Scan(&ticketId, &version)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errTicketNotFound
		} else if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			err = errTicketNotFound
		}
	}
	return
}

func updateOrderStatusToSuccess(ctx context.Context, id uint64) error {
	q := `UPDATE orders SET status = 'completed', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'created';`

	res, err := db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errOrderNotFound
	}

	return nil
}

func updateOrderStatusToCancel(ctx context.Context, id uint64) (version uint32, err error) {
	q1 := `UPDATE orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'created';`
	q2 := `UPDATE tickets SET order_id = NULL, version = version + 1 WHERE order_id = $1 RETURNING version;`

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

	err = br.QueryRow().Scan(&version)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = errTicketNotFound
		}
	}

	return
}
