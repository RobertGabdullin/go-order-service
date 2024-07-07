//go:build integration

package integration

import (
	"database/sql"
	"testing"
	"time"

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

func (suite *ServiceTestSuite) TearDownTest() {
	suite.DB.Exec("TRUNCATE TABLE orders RESTART IDENTITY;")
}

func (suite *ServiceTestSuite) TestAddOrder() {
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

	err := suite.OrderService.AddOrder(order)
	suite.NoError(err)

	storedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)

	Normalize(&storedOrder, &order)
	suite.Equal(order, storedOrder)
}

func (suite *ServiceTestSuite) TestChangeStatus() {
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

	err := suite.OrderService.AddOrder(order)
	suite.NoError(err)

	err = suite.OrderService.ChangeStatus(order.Id, "delivered", "newhash123")
	suite.NoError(err)

	updatedOrder, err := suite.OrderStorage.GetOrderById(order.Id)
	suite.NoError(err)
	suite.Equal("delivered", updatedOrder.Status)
	suite.NotZero(updatedOrder.DeliveredAt)
	suite.Equal("newhash123", updatedOrder.Hash)
}

func (suite *ServiceTestSuite) TestFindOrders() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
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
		Id:          2,
		Recipient:   456,
		Status:      "alive",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	suite.NoError(err)
	err = suite.OrderService.AddOrder(order2)
	suite.NoError(err)

	orders, err := suite.OrderService.FindOrders([]int{1, 2})
	suite.NoError(err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}

	suite.ElementsMatch(orders, []models.Order{order1, order2})

	orders, err = suite.OrderService.FindOrders([]int{3, 4})
	suite.Nil(orders)
	suite.ErrorContains(err, "no rows in result set")

}

func (suite *ServiceTestSuite) TestListOrders() {
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
		Wrapper:     "none",
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
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	suite.NoError(err)
	err = suite.OrderService.AddOrder(order2)
	suite.NoError(err)

	orders, err := suite.OrderService.ListOrders(123)
	suite.NoError(err)

	Normalize(&order1, &order2)
	for i := range orders {
		Normalize(&orders[i])
	}

	suite.ElementsMatch(orders, []models.Order{order1, order2})

	orders, err = suite.OrderService.ListOrders(234)
	suite.NoError(err)
	suite.ElementsMatch(orders, []models.Order{})

}

func (suite *ServiceTestSuite) TestGetReturns() {
	order1 := models.Order{
		Id:          1,
		Recipient:   123,
		Status:      "returned",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now().Add(12 * time.Hour),
		Hash:        "hash123",
		Weight:      10,
		BasePrice:   1000,
		Wrapper:     "none",
	}

	order2 := models.Order{
		Id:          2,
		Recipient:   456,
		Status:      "returned",
		Expire:      time.Now().Add(24 * time.Hour),
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Hash:        "hash456",
		Weight:      15,
		BasePrice:   1500,
		Wrapper:     "none",
	}

	err := suite.OrderService.AddOrder(order1)
	suite.NoError(err)
	err = suite.OrderService.AddOrder(order2)
	suite.NoError(err)

	returnedOrders, err := suite.OrderService.GetReturns(0, 2)
	suite.NoError(err)

	Normalize(&order1, &order2)
	for i := range returnedOrders {
		Normalize(&returnedOrders[i])
	}
	suite.ElementsMatch(returnedOrders, []models.Order{order1, order2})

	returnedOrders, err = suite.OrderService.GetReturns(0, 1)
	suite.NoError(err)
	for i := range returnedOrders {
		Normalize(&returnedOrders[i])
	}
	suite.ElementsMatch(returnedOrders, []models.Order{order2})

}

func (suite *ServiceTestSuite) TestDeleteOrder() {
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

	err := suite.OrderService.AddOrder(order)
	suite.NoError(err)

	err = suite.OrderService.DeleteOrder(order.Id)
	suite.NoError(err)

	_, err = suite.OrderStorage.GetOrderById(order.Id)
	suite.ErrorContains(err, "no rows in result set")
}

func (suite *ServiceTestSuite) TestGetWrapper() {
	wrapper := models.Wrapper{
		Id:        1,
		Type:      "pack",
		MaxWeight: sql.NullInt64{Int64: 10, Valid: true},
		Markup:    5,
	}

	storedWrapper, err := suite.OrderService.GetWrapper("pack")
	suite.NoError(err)
	suite.Equal(wrapper, storedWrapper)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
