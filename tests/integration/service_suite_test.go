package integration

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type ServiceTestSuite struct {
	suite.Suite
	DB             *sql.DB
	OrderService   *service.OrderService
	OrderStorage   storage.OrderStorage
	WrapperStorage storage.WrapperStorage
}

func (suite *ServiceTestSuite) SetupSuite() {
	db, err := setupDB()
	if err != nil {
		suite.T().Fatalf("Failed to setup database: %v", err)
	}
	suite.DB = db
	orderStorage, err := storage.NewOrderStorage(dbConnStr)
	if err != nil {
		suite.T().Fatalf("Failed to create OrderStorage: %v", err)
	}
	wrapperStorage, err := storage.NewWrapperStorage(dbConnStr)
	if err != nil {
		suite.T().Fatalf("Failed to create WrapperStorage: %v", err)
	}
	suite.OrderService = service.NewPostgresService(orderStorage, wrapperStorage)
	suite.OrderStorage = orderStorage
	suite.WrapperStorage = wrapperStorage
}

func (suite *ServiceTestSuite) TearDownSuite() {
	teardownDB(suite.DB)
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.DB.Exec("TRUNCATE TABLE orders, wrappers RESTART IDENTITY;")
}

func (suite *ServiceTestSuite) TestAddOrder() {
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

	err := suite.OrderService.AddOrder(order)
	assert.NoError(suite.T(), err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	assert.NoError(suite.T(), err)

	models.Normalize(&storedOrder, &order)
	assert.Equal(suite.T(), order, storedOrder)
}

func (suite *ServiceTestSuite) TestChangeStatus() {
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

	err := suite.OrderService.AddOrder(order)
	assert.NoError(suite.T(), err)

	err = suite.OrderService.ChangeStatus(order.Id, "delivered", "newhash123")
	assert.NoError(suite.T(), err)

	updatedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "delivered", updatedOrder.Status)
	assert.NotZero(suite.T(), updatedOrder.DeliveredAt)
	assert.Equal(suite.T(), "newhash123", updatedOrder.Hash)
}

func (suite *ServiceTestSuite) TestFindOrders() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "alive",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Time{},
		ReturnedAt:  time.Time{},
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
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
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	assert.NoError(suite.T(), err)
	err = suite.OrderService.AddOrder(order2)
	assert.NoError(suite.T(), err)

	orders, err := suite.OrderService.FindOrders([]int{1, 2})
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(suite.T(), orders, order1)
	assert.Contains(suite.T(), orders, order2)
}

func (suite *ServiceTestSuite) TestListOrders() {
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
		Wrapper:     "none",
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
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	assert.NoError(suite.T(), err)
	err = suite.OrderService.AddOrder(order2)
	assert.NoError(suite.T(), err)

	orders, err := suite.OrderService.ListOrders(123)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), orders, 2)

	models.Normalize(&order1, &order2)
	for i := range orders {
		models.Normalize(&orders[i])
	}

	assert.Contains(suite.T(), orders, order1)
	assert.Contains(suite.T(), orders, order2)
}

func (suite *ServiceTestSuite) TestGetReturns() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "returned",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
	}

	order2 := models.Order{
		Id:          2,
		Recipient:   456,
		Status:      "returned",
		Limit:       time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	assert.NoError(suite.T(), err)
	err = suite.OrderService.AddOrder(order2)
	assert.NoError(suite.T(), err)

	returnedOrders, err := suite.OrderService.GetReturns(0, 2)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), returnedOrders, 2)

	models.Normalize(&order1, &order2)
	for i := range returnedOrders {
		models.Normalize(&returnedOrders[i])
	}

	assert.Contains(suite.T(), returnedOrders, order1)
	assert.Contains(suite.T(), returnedOrders, order2)
}

func (suite *ServiceTestSuite) TestDeleteOrder() {
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
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order)
	assert.NoError(suite.T(), err)

	err = suite.OrderService.DeleteOrder(order.Id)
	assert.NoError(suite.T(), err)

	_, err = suite.OrderStorage.GetOrderById(order.Id)
	assert.Error(suite.T(), err)
}

func (suite *ServiceTestSuite) TestGetWrapper() {
	wrapper := models.Wrapper{
		Id:        2,
		Type:      "pack",
		MaxWeight: sql.NullInt64{Int64: 20, Valid: true},
		Markup:    10,
	}

	_, err := suite.DB.Exec("INSERT INTO wrappers (id, type, max_weight, markup) VALUES ($1, $2, $3, $4)", wrapper.Id, wrapper.Type, wrapper.MaxWeight, wrapper.Markup)
	assert.NoError(suite.T(), err)

	storedWrapper, err := suite.OrderService.GetWrapper("pack")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), wrapper, storedWrapper)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
