package ticket

import "context"

func updateTicket(ctx context.Context, ticket ticketEntity) error {
	q := `UPDATE tickets SET order_id = $3, version = version + 1 WHERE id = $1 AND version = $2 AND order_id IS NULL`

	res, err := db.Exec(ctx, q, ticket.ID, ticket.Version, ticket.OrderID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errTicketNotFound
	}

	return nil
}

func updateTicketOrderToNull(ctx context.Context, ticket ticketEntity) error {
	q := `UPDATE tickets SET order_id = NULL, version = version + 1 WHERE order_id = $1 AND version = $2;`

	res, err := db.Exec(ctx, q, ticket.OrderID, ticket.Version)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errTicketNotFound
	}
	return nil
}

func countAvailableTicketByCategoryGroup(ctx context.Context) (categoryIds []uint8, total []uint16, err error) {
	q := `SELECT category_id, COUNT(*) FROM tickets WHERE order_id IS NULL GROUP BY category_id;`

	rows, err := db.Query(ctx, q)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var categoryId uint8
		var count uint16

		err = rows.Scan(&categoryId, &count)
		if err != nil {
			return
		}

		categoryIds = append(categoryIds, categoryId)
		total = append(total, count)
	}

	return
}
