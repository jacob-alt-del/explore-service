package dataaccess

import "context"

func (r *Repository) ListNewLikedYou(
	ctx context.Context,
	recipientID string,
	paginationUnix int64,
	pageSize int,
) ([]Decision, error) {
	query, args := buildListNewLikedYouQuery(recipientID, paginationUnix, pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Decision
	for rows.Next() {
		var d Decision
		if err := rows.Scan(&d.ActorID, &d.UpdatedAtUnix); err != nil {
			return nil, err
		}
		results = append(results, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func buildListNewLikedYouQuery(recipientID string, paginationUnix int64, pageSize int) (string, []interface{}) {
	query := `
		SELECT d1.actor_id, UNIX_TIMESTAMP(d1.updated_at) AS unix_timestamp
		FROM decisions AS d1
		LEFT JOIN decisions AS d2
			ON d1.actor_id = d2.recipient_id
			AND d1.recipient_id = d2.actor_id
		WHERE d1.recipient_id = ?
		  AND d1.liked = TRUE
		  AND (d2.liked IS NULL OR d2.liked = FALSE)
    `
	args := []interface{}{recipientID}

	if paginationUnix > 0 {
		query += " AND d1.updated_at < FROM_UNIXTIME(?)"
		args = append(args, paginationUnix)
	}

	query += " ORDER BY d1.updated_at DESC LIMIT ?;"
	args = append(args, pageSize+1) // +1 to check for next page

	return query, args
}
