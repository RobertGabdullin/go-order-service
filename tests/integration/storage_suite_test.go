package integration

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func (suite *StorageTestSuite) SetupTest() {
	suite.DB.Exec("TRUNCATE TABLE orders, wrappers RESTART IDENTITY;")
}

func (suite *StorageTestSuite) TestAddOrder() {

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

	err := suite.OrderStorage.AddOrder(order)
	assert.NoError(suite.T(), err)

	var retrievedOrder models.Order
	err = suite.DB.QueryRow("SELECT id, recipient, status, time_limit, delivered_at, returned_at, hash, weight, base_cost, wrapper FROM orders WHERE id = $1", order.Id).
		Scan(&retrievedOrder.Id, &retrievedOrder.Recipient, &retrievedOrder.Status, &retrievedOrder.Limit, &retrievedOrder.DeliveredAt, &retrievedOrder.ReturnedAt, &retrievedOrder.Hash, &retrievedOrder.Weight, &retrievedOrder.BasePrice, &retrievedOrder.Wrapper)
	assert.NoError(suite.T(), err)

	models.Normalize(&retrievedOrder, &order)

	assert.Equal(suite.T(), order, retrievedOrder)
}

func (suite *StorageTestSuite) TestUpdateOrder() {

	order := models.Order{
		Id:          1,
		Recipient:   1,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err := suite.OrderStorage.AddOrder(order)
	assert.NoError(suite.T(), err)

	order.Status = "updated"
	order.Limit = time.Now().Add(48 * time.Hour)

	err = suite.OrderStorage.UpdateOrder(order)
	assert.NoError(suite.T(), err)

	var retrievedOrder models.Order
	err = suite.DB.QueryRow("SELECT id, recipient, status, time_limit, delivered_at, returned_at, hash, weight, base_cost, wrapper FROM orders WHERE id = $1", order.Id).
		Scan(&retrievedOrder.Id, &retrievedOrder.Recipient, &retrievedOrder.Status, &retrievedOrder.Limit, &retrievedOrder.DeliveredAt, &retrievedOrder.ReturnedAt, &retrievedOrder.Hash, &retrievedOrder.Weight, &retrievedOrder.BasePrice, &retrievedOrder.Wrapper)
	assert.NoError(suite.T(), err)

	models.Normalize(&order, &retrievedOrder)

	assert.Equal(suite.T(), order, retrievedOrder)
}

func (suite *StorageTestSuite) TestDeleteOrder() {
	order := models.Order{
		Id:          1,
		Recipient:   1,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "pack",
	}

	err := suite.OrderStorage.AddOrder(order)
	assert.NoError(suite.T(), err)

	err = suite.OrderStorage.DeleteOrder(order.Id)
	assert.NoError(suite.T(), err)

	var count int
	err = suite.DB.QueryRow("SELECT COUNT(*) FROM orders WHERE id = $1", order.Id).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *StorageTestSuite) TestGetWrapperByType() {
	wrapper := models.NewWrapper(1, "pack", sql.NullInt64{Int64: 10, Valid: true}, 5)

	_, err := suite.DB.Exec("INSERT INTO wrappers (id, type, max_weight, markup) VALUES ($1, $2, $3, $4)", wrapper.Id, wrapper.Type, wrapper.MaxWeight.Int64, wrapper.Markup)
	assert.NoError(suite.T(), err)

	wrappers, err := suite.WrapperStorage.GetWrapperByType("pack")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(wrappers))

	assert.Equal(suite.T(), wrapper.Id, wrappers[0].Id)
	assert.Equal(suite.T(), wrapper.Type, wrappers[0].Type)
	assert.Equal(suite.T(), wrapper.MaxWeight, wrappers[0].MaxWeight)
	assert.Equal(suite.T(), wrapper.Markup, wrappers[0].Markup)
}

func (suite *StorageTestSuite) TestGetOrderById() {
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

	err := suite.OrderStorage.AddOrder(order)
	assert.NoError(suite.T(), err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	assert.NoError(suite.T(), err)

	models.Normalize(&order, &storedOrder)
	assert.Equal(suite.T(), order, storedOrder)
}

func (suite *StorageTestSuite) TestGetOrdersByRecipient() {
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

	err := suite.OrderStorage.AddOrder(order1)
	assert.NoError(suite.T(), err)
	err = suite.OrderStorage.AddOrder(order2)
	assert.NoError(suite.T(), err)

	orders, err := suite.OrderStorage.GetOrdersByRecipient(123)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(suite.T(), orders, order1)
	assert.Contains(suite.T(), orders, order2)
}

func (suite *StorageTestSuite) TestGetPaginatedOrdersByStatus() {
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

	err := suite.OrderStorage.AddOrder(order1)
	assert.NoError(suite.T(), err)
	err = suite.OrderStorage.AddOrder(order2)
	assert.NoError(suite.T(), err)

	orders, err := suite.OrderStorage.GetPaginatedOrdersByStatus("alive", 0, 2)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(suite.T(), orders, order1)
	assert.Contains(suite.T(), orders, order2)
}

func (suite *StorageTestSuite) TestUpdateHash() {
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

	err := suite.OrderStorage.AddOrder(order)
	assert.NoError(suite.T(), err)

	newHash := "newhash456"
	err = suite.OrderStorage.UpdateHash(order.Id, newHash)
	assert.NoError(suite.T(), err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), newHash, storedOrder.Hash)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
