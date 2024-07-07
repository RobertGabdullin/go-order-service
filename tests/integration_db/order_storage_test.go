//go:build integration

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
		Recipient:   1,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	Normalize(&order, &storedOrder)

	assert.Equal(t, order, storedOrder)

	orderStorage.DeleteOrder(1)
}

func TestOrderStorage_UpdateOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          2,
		Recipient:   2,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	Normalize(&order, &updatedOrder)

	assert.Equal(t, order, updatedOrder)

	orderStorage.DeleteOrder(2)
}

func TestOrderStorage_DeleteOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          3,
		Recipient:   3,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	orderStorage.DeleteOrder(3)
}

func TestOrderStorage_GetOrderById(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          4,
		Recipient:   4,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	Normalize(&order, &storedOrder)

	assert.Equal(t, order, storedOrder)

	orderStorage.DeleteOrder(4)
}

func TestOrderStorage_GetOrdersByRecipient(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order1 := models.Order{
		Id:          5,
		Recipient:   5,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	order2 := models.Order{
		Id:          6,
		Recipient:   5,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	orders, err := orderStorage.GetOrdersByRecipient(5)
	assert.NoError(t, err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}
	assert.ElementsMatch(t, orders, []models.Order{order1, order2})

	orderStorage.DeleteOrder(5)
	orderStorage.DeleteOrder(6)
}

func TestOrderStorage_GetPaginatedOrdersByStatus(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order1 := models.Order{
		Id:          7,
		Recipient:   7,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	order2 := models.Order{
		Id:          8,
		Recipient:   8,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}

	assert.ElementsMatch(t, orders, []models.Order{order1, order2})

	orderStorage.DeleteOrder(7)
	orderStorage.DeleteOrder(8)
}

func TestOrderStorage_UpdateHash(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)

	order := models.Order{
		Id:          9,
		Recipient:   9,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
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

	orderStorage.DeleteOrder(9)
}
