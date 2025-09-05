package order

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgress://user:pass@localhost:5432/ordersdb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestCreateOrder(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	o := &Order{
		Customer: "TestCustomer",
		Status:   "new",
	}

	err := repo.CreateOrder(context.Background(), o)
	require.NoError(t, err)
	require.NotZero(t, o.ID)
	require.WithinDuration(t, time.Now().UTC(), o.CreatedAt, 2*time.Second)
}
