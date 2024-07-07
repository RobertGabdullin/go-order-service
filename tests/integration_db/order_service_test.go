//go:build integration

package integration

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

func TestOrderService_AddOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

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

	err = orderService.AddOrder(order)
	assert.NoError(t, err)

	storedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)

	Normalize(&order, &storedOrder)

	assert.Equal(t, order, storedOrder)

	orderService.DeleteOrder(1)

}

func TestOrderService_ChangeStatus(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

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

	err = orderService.AddOrder(order)
	require.NoError(t, err)

	err = orderService.ChangeStatus(order.Id, "delivered", "newhash123")
	assert.NoError(t, err)

	updatedOrder, err := orderStorage.GetOrderById(order.Id)
	assert.NoError(t, err)
	assert.Equal(t, "delivered", updatedOrder.Status)
	assert.NotZero(t, updatedOrder.DeliveredAt)
	assert.Equal(t, "newhash123", updatedOrder.Hash)

	orderService.DeleteOrder(2)
}

func TestOrderService_FindOrders(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

	order1 := models.Order{
		Id:          3,
		Recipient:   3,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Time{},
		ReturnedAt:  time.Time{},
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
	}

	order2 := models.Order{
		Id:          4,
		Recipient:   4,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "none",
	}

	err = orderService.AddOrder(order1)
	require.NoError(t, err)

	err = orderService.AddOrder(order2)
	require.NoError(t, err)

	orders, err := orderService.FindOrders([]int{order1.Id, order2.Id})
	assert.NoError(t, err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}
	assert.ElementsMatch(t, orders, []models.Order{order1, order2})

	orderService.DeleteOrder(3)
	orderService.DeleteOrder(4)
}

func TestOrderService_ListOrders(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

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
		Wrapper:     "none",
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
		Wrapper:     "none",
	}

	err = orderService.AddOrder(order1)
	require.NoError(t, err)

	err = orderService.AddOrder(order2)
	require.NoError(t, err)

	orders, err := orderService.ListOrders(5)
	assert.NoError(t, err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}
	assert.ElementsMatch(t, orders, []models.Order{order1, order2})

	orderService.DeleteOrder(5)
	orderService.DeleteOrder(6)
}

func TestOrderService_GetReturns(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

	order1 := models.Order{
		Id:          7,
		Recipient:   7,
		Status:      "returned",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
	}

	order2 := models.Order{
		Id:          8,
		Recipient:   8,
		Status:      "returned",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "none",
	}

	err = orderService.AddOrder(order1)
	require.NoError(t, err)

	err = orderService.AddOrder(order2)
	require.NoError(t, err)

	orders, err := orderService.GetReturns(0, 2)
	assert.NoError(t, err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}
	assert.ElementsMatch(t, orders, []models.Order{order1, order2})

	orderService.DeleteOrder(7)
	orderService.DeleteOrder(8)
}

func TestOrderService_DeleteOrder(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

	order := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
	}

	err = orderService.AddOrder(order)
	require.NoError(t, err)

	err = orderService.DeleteOrder(order.Id)
	assert.NoError(t, err)

	_, err = orderStorage.GetOrderById(order.Id)
	assert.Error(t, err)
}

func TestOrderService_GetWrapper(t *testing.T) {
	db, err := setupDB()
	require.NoError(t, err)
	defer teardownDB(db)

	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	require.NoError(t, err)
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	require.NoError(t, err)
	orderService := service.NewPostgresService(orderStorage, wrapperStorage)

	wrapper := models.Wrapper{
		Id:        1,
		Type:      "pack",
		MaxWeight: sql.NullInt64{Int64: 10, Valid: true},
		Markup:    5,
	}

	storedWrapper, err := orderService.GetWrapper("pack")
	assert.NoError(t, err)
	assert.Equal(t, wrapper, storedWrapper)
}
