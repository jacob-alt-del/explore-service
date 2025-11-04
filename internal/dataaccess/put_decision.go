package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *Repository) UpsertDecision(ctx context.Context, actorID, recipientID string, liked bool) error {
	const query = `
		INSERT INTO decisions (actor_id, recipient_id, liked)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			liked = VALUES(liked),
			updated_at = CURRENT_TIMESTAMP;
	`

	if _, err := r.db.ExecContext(ctx, query, actorID, recipientID, liked); err != nil {
		return fmt.Errorf("error upserting decision: %w", err)
	}
	return nil
}

func (r *Repository) CheckMutualLike(ctx context.Context, actorID, recipientID string) (bool, error) {
	const query = `
		SELECT liked FROM decisions
		WHERE actor_id = ? AND recipient_id = ? AND liked = TRUE;
	`

	err := r.db.QueryRowContext(ctx, query, recipientID, actorID).Scan(new(int))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // recipient hasn’t liked actor back :(
		}
		return false, fmt.Errorf("error checking mutual like: %w", err)
	}
	return true, nil
}
