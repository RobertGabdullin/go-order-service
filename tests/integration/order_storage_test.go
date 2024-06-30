package integration

import (
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

func TestOrderStorage_AddOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err = orderStorage.AddOrder(order)
	assert.NoError(t, err)

	storedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)

	models.Normalize(&order, &storedOrder)

	assert.Equal(t, order, storedOrder)
}

func TestOrderStorage_UpdateOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err = orderStorage.AddOrder(order)
	require.NoError(t, err)

	order.Status = "delivered"
	order.Hash = "newhash123"

	err = orderStorage.UpdateOrder(order)
	assert.NoError(t, err)

	updatedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)

	models.Normalize(&order, &updatedOrder)

	assert.Equal(t, order, updatedOrder)
}

func TestOrderStorage_DeleteOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err = orderStorage.AddOrder(order)
	require.NoError(t, err)

	err = orderStorage.DeleteOrder(order.Id)
	assert.NoError(t, err)

	_, err = orderStorage.GetOrderById(order.Id)
	assert.Error(t, err)
}

func TestOrderStorage_GetOrderById(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err = orderStorage.AddOrder(order)
	require.NoError(t, err)

	storedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)

	models.Normalize(&order, &storedOrder)

	assert.Equal(t, order, storedOrder)
}

func TestOrderStorage_GetOrdersByRecipient(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order1 := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	order2 := models.Order{
		Id:          2,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "box",
	}

	err = orderStorage.AddOrder(order1)
	require.NoError(t, err)
	err = orderStorage.AddOrder(order2)
	require.NoError(t, err)

	orders, err := orderStorage.GetOrdersByRecipient(123)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(t, orders, order1)
	assert.Contains(t, orders, order2)
}

func TestOrderStorage_GetPaginatedOrdersByStatus(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order1 := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	order2 := models.Order{
		Id:          2,
		Recipient:   456,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "box",
	}

	err = orderStorage.AddOrder(order1)
	require.NoError(t, err)
	err = orderStorage.AddOrder(order2)
	require.NoError(t, err)

	orders, err := orderStorage.GetPaginatedOrdersByStatus("alive", 0, 2)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(t, orders, order1)
	assert.Contains(t, orders, order2)
}

func TestOrderStorage_UpdateHash(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err = orderStorage.AddOrder(order)
	require.NoError(t, err)

	newHash := "newhash123"
	err = orderStorage.UpdateHash(order.Id, newHash)
	assert.NoError(t, err)

	updatedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)
	assert.Equal(t, newHash, updatedOrder.Hash)
}
