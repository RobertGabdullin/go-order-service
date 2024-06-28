package service

import (
	"errors"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type OrderService struct {
	orderStorage   storage.TransactionalOrderStorage
	wrapperStorage storage.WrapperStorage
}

func NewPostgresService(storage storage.TransactionalOrderStorage, wrapperStorage storage.WrapperStorage) *OrderService {
	return &OrderService{orderStorage: storage, wrapperStorage: wrapperStorage}
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

	return s.orderStorage.CommitTransaction(tx)
}

func (s *OrderService) FindOrders(ids []int) ([]models.Order, error) {

	var orders []models.Order
	for _, id := range ids {
		order, err := s.orderStorage.GetOrderById(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) ListOrders(recipient int) ([]models.Order, error) {
	return s.orderStorage.GetOrdersByRecipient(recipient)
}

func (s *OrderService) GetReturns(offset, limit int) ([]models.Order, error) {
	return s.orderStorage.GetPaginatedOrdersByStatus("returned", offset, limit)
}

func (s *OrderService) DeleteOrder(id int) error {
	return s.orderStorage.DeleteOrder(id)
}

func (s *OrderService) GetWrapper(givenType string) (models.Wrapper, error) {
	wrappers, err := s.wrapperStorage.GetWrapperByType(givenType)
	if err != nil {
		return models.Wrapper{}, err
	}
	if len(wrappers) == 0 {
		return models.Wrapper{}, errors.New("undefined type")
	}
	return wrappers[0], nil
}
