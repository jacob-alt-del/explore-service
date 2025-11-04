package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_UpsertDecision_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// Expect the upsert query to be executed successfully
	mock.ExpectExec("INSERT INTO decisions").
		WithArgs("actor1", "recipient1", true).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = repo.UpsertDecision(ctx, "actor1", "recipient1", true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func Test_UpsertDecision_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("INSERT INTO decisions").
		WithArgs("actor1", "recipient1", false).
		WillReturnError(errors.New("insert failed"))

	ctx := context.Background()
	err = repo.UpsertDecision(ctx, "actor1", "recipient1", false)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func Test_CheckMutualLike_True(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// Simulate the SELECT returning a row (recipient liked actor back)
	rows := sqlmock.NewRows([]string{"liked"}).AddRow(1)

	mock.ExpectQuery("SELECT liked FROM decisions").
		WithArgs("recipient1", "actor1").
		WillReturnRows(rows)

	ctx := context.Background()
	mutual, err := repo.CheckMutualLike(ctx, "actor1", "recipient1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !mutual {
		t.Errorf("expected mutual like to be true, got false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func Test_CheckMutualLike_False_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// Simulate no results (recipient hasn’t liked actor)
	mock.ExpectQuery("SELECT liked FROM decisions").
		WithArgs("recipient1", "actor1").
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	mutual, err := repo.CheckMutualLike(ctx, "actor1", "recipient1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if mutual {
		t.Errorf("expected mutual like to be false, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func Test_CheckMutualLike_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// Simulate a real DB error (e.g., connection issue)
	mock.ExpectQuery("SELECT liked FROM decisions").
		WithArgs("recipient1", "actor1").
		WillReturnError(errors.New("query failed"))

	ctx := context.Background()
	mutual, err := repo.CheckMutualLike(ctx, "actor1", "recipient1")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if mutual {
		t.Errorf("expected mutual like to be false, got true")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
