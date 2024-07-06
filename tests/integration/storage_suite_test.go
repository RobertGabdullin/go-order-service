//go:build integration

package integration

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type StorageTestSuite struct {
	suite.Suite
	DB             *sql.DB
	OrderStorage   *storage.PostgresOrderStorage
	WrapperStorage *storage.PostgresWrapperStorage
}

func (suite *StorageTestSuite) SetupSuite() {
	db, err := setupDB()
	if err != nil {
		suite.T().Fatalf("Failed to setup database: %v", err)
	}
	suite.DB = db
	suite.OrderStorage, err = storage.NewOrderStorage(dbConnStr)
	suite.WrapperStorage, err = storage.NewWrapperStorage(dbConnStr)
}

func (suite *StorageTestSuite) TearDownSuite() {
	teardownDB(suite.DB)
}

func (suite *StorageTestSuite) TearDownTest() {
	suite.DB.Exec("TRUNCATE TABLE orders RESTART IDENTITY;")
}

func (suite *StorageTestSuite) TestAddOrder() {
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
		Wrapper:     "pack",
	}

	err := suite.OrderStorage.AddOrder(order)
	suite.NoError(err)

	retrievedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)

	Normalize(&retrievedOrder, &order)

	suite.Equal(order, retrievedOrder)
}

func (suite *StorageTestSuite) TestUpdateOrder() {
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

	err := suite.OrderStorage.AddOrder(order)
	suite.NoError(err)

	order.Status = "updated"
	order.Expire = time.Now().Add(48 * time.Hour)

	err = suite.OrderStorage.UpdateOrder(order)
	suite.NoError(err)

	retrievedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)

	Normalize(&order, &retrievedOrder)

	suite.Equal(order, retrievedOrder)
}

func (suite *StorageTestSuite) TestDeleteOrder() {
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

	err := suite.OrderStorage.AddOrder(order)
	suite.NoError(err)

	err = suite.OrderStorage.DeleteOrder(order.Id)
	suite.NoError(err)

	var count int
	err = suite.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = $1", order.Id).Scan(&count)
	suite.NoError(err)
	suite.Equal(0, count)
}

func (suite *StorageTestSuite) TestGetWrapperByType() {
	wrapper := models.NewWrapper(1, "pack", sql.NullInt64{Int64: 10, Valid: true}, 5)

	wrappers, err := suite.WrapperStorage.GetWrapperByType("pack")
	suite.NoError(err)
	suite.Require().Len(wrappers, 1)

	suite.Equal(wrapper.Id, wrappers[0].Id)
	suite.Equal(wrapper.Type, wrappers[0].Type)
	suite.Equal(wrapper.MaxWeight, wrappers[0].MaxWeight)
	suite.Equal(wrapper.Markup, wrappers[0].Markup)
}

func (suite *StorageTestSuite) TestGetOrderById() {
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
		Wrapper:     "pack",
	}

	err := suite.OrderStorage.AddOrder(order)
	suite.NoError(err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)

	Normalize(&order, &storedOrder)
	suite.Equal(order, storedOrder)
}

func (suite *StorageTestSuite) TestGetOrdersByRecipient() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
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
		Id:          2,
		Recipient:   123,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "box",
	}

	err := suite.OrderStorage.AddOrder(order1)
	suite.NoError(err)
	err = suite.OrderStorage.AddOrder(order2)
	suite.NoError(err)

	orders, err := suite.OrderStorage.GetOrdersByRecipient(123)
	suite.NoError(err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}

	suite.ElementsMatch(orders, []models.Order{order1, order2})
}

func (suite *StorageTestSuite) TestGetPaginatedOrdersByStatus() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
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
		Id:          2,
		Recipient:   456,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "box",
	}

	err := suite.OrderStorage.AddOrder(order1)
	suite.NoError(err)
	err = suite.OrderStorage.AddOrder(order2)
	suite.NoError(err)

	orders, err := suite.OrderStorage.GetPaginatedOrdersByStatus("alive", 0, 2)
	suite.NoError(err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}

	suite.ElementsMatch(orders, []models.Order{order1, order2})
}

func (suite *StorageTestSuite) TestUpdateHash() {
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
		Wrapper:     "pack",
	}

	err := suite.OrderStorage.AddOrder(order)
	suite.NoError(err)

	newHash := "newhash456"
	err = suite.OrderStorage.UpdateHash(order.Id, newHash)
	suite.NoError(err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)
	suite.Equal(newHash, storedOrder.Hash)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
