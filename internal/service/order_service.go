package service

import (
	"errors"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/cache"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/metrics"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type OrderService struct {
	orderStorage   storage.TransactionalOrderStorage
	wrapperStorage storage.WrapperStorage
	cache          cache.Cache
	metrics        metrics.Metrics
}

func NewPostgresService(storage storage.TransactionalOrderStorage, wrapperStorage storage.WrapperStorage, cache cache.Cache, metrics metrics.Metrics) *OrderService {
	return &OrderService{
		orderStorage:   storage,
		wrapperStorage: wrapperStorage,
		cache:          cache,
		metrics:        metrics,
	}
}

func (s *OrderService) AddOrder(ord models.Order) error {
	return s.orderStorage.AddOrder(ord)
}

func (s *OrderService) ChangeStatus(id int, status, hash string) error {

	order, err := s.orderStorage.GetOrderById(id)
	if err != nil {
		return err
	}

	order.Status = status
	if status == "delivered" {
		order.DeliveredAt = time.Now()
	} else if status == "returned" {
		order.ReturnedAt = time.Now()
	} else {
		return errors.New("unknown status")
	}

	tx, err := s.orderStorage.BeginTransaction()
	if err != nil {
		return err
	}

	err = s.orderStorage.UpdateOrder(order)
	if err != nil {
		return err
	}

	err = s.orderStorage.UpdateHash(id, hash)
	if err != nil {
		s.orderStorage.RollbackTransaction(tx)
		return err
	}

	if status == "delivered" {
		s.metrics.IncIssuedOrders()
	}

	s.cache.InvalidateOrder(id)

	return s.orderStorage.CommitTransaction(tx)
}

func (s *OrderService) FindOrders(ids []int) ([]models.Order, error) {
	var orders []models.Order
	for _, id := range ids {
		order, err := s.getOrderFromCacheOrStorage(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *OrderService) ListOrders(recipient int) ([]models.Order, error) {
	orders, err := s.orderStorage.GetOrdersByRecipient(recipient)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		s.cache.SetOrder(order.Id, order)
	}
	return orders, nil
}

func (s *OrderService) GetReturns(offset, limit int) ([]models.Order, error) {
	orders, err := s.orderStorage.GetPaginatedOrdersByStatus("returned", offset, limit)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		s.cache.SetOrder(order.Id, order)
	}
	return orders, nil
}

func (s *OrderService) DeleteOrder(id int) error {
	err := s.orderStorage.DeleteOrder(id)
	if err != nil {
		return err
	}

	s.cache.InvalidateOrder(id)
	return nil
}

func (s *OrderService) GetWrapper(givenType string) (models.Wrapper, error) {
	wrappers, err := s.wrapperStorage.GetWrapperByType(givenType)
	if err != nil {
		return models.Wrapper{}, err
	}
	if len(wrappers) == 0 {
		return models.Wrapper{}, errors.New("wrapper not found")
	}
	return wrappers[0], nil
}

func (s *OrderService) getOrderFromCacheOrStorage(id int) (models.Order, error) {
	if order, found := s.cache.GetOrder(id); found {
		return order, nil
	}

	order, err := s.orderStorage.GetOrderById(id)
	if err != nil {
		return models.Order{}, err
	}

	s.cache.SetOrder(id, order)
	return order, nil
}
