package dataaccess

import "context"

func (r *Repository) CountLikedYou(ctx context.Context, recipientID string) (uint64, error) {
	const query = `
        SELECT COUNT(*) s
        FROM decisions 
        WHERE recipient_id = ? AND liked = TRUE;
    `

	var count uint64
	err := r.db.QueryRowContext(ctx, query, recipientID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
