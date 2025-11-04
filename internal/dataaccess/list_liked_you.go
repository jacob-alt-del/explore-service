package dataaccess

import "context"

func (r *Repository) ListLikedYou(
	ctx context.Context,
	recipientID string,
	paginationUnix int64,
	pageSize int,
) ([]Decision, error) {
	query, args := buildListLikedYouQuery(recipientID, paginationUnix, pageSize)

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

func buildListLikedYouQuery(recipientID string, paginationUnix int64, pageSize int) (string, []interface{}) {
	query := `
        SELECT actor_id, UNIX_TIMESTAMP(updated_at)
        FROM decisions
        WHERE recipient_id = ? AND liked = TRUE
    `
	args := []interface{}{recipientID}

	if paginationUnix > 0 {
		query += " AND updated_at < FROM_UNIXTIME(?)"
		args = append(args, paginationUnix)
	}

	query += " ORDER BY updated_at DESC LIMIT ?"
	args = append(args, pageSize+1) // +1 to check for next page

	return query, args
}
