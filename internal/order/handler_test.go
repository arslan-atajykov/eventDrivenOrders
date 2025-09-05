// package order

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	_ "github.com/lib/pq"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateOrderHandler(t *testing.T) {
// 	db := setupTestDB(t)
// 	repo := NewRepository(db)
// 	producer := &MockProducer{} // fake Kafka for testing
// 	handler := NewHandler(repo, producer)

// 	// build request
// 	body := []byte(`{"customer":"Bob"}`)
// 	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
// 	w := httptest.NewRecorder()

// 	handler.CreateOrder(w, req)

// 	resp := w.Result()
// 	defer resp.Body.Close()

// 	require.Equal(t, http.StatusCreated, resp.StatusCode)

// 	var o Order
// 	err := json.NewDecoder(resp.Body).Decode(&o)
// 	require.NoError(t, err)
// 	require.Equal(t, "Bob", o.Customer)
// 	require.Equal(t, "new", o.Status)
// 	require.NotZero(t, o.ID)
// }

// // MockProducer fakes Kafka producer
// type MockProducer struct{}

// func (m *MockProducer) PublishOrder(ctx context.Context, o *Order) error {
// 	// do nothing
// 	return nil
// }
