package service

import (
	"errors"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type OrderService struct {
	storage storage.TransactionalStorage
}

func NewPostgresService(storage storage.TransactionalStorage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) AddOrder(ord models.Order) error {
	return s.storage.AddOrder(ord)
}

func (s *OrderService) ChangeStatus(id int, status, hash string) error {

	order, err := s.storage.GetOrderById(id)
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

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return err
	}

	err = s.storage.UpdateOrder(order)
	if err != nil {
		return err
	}

	err = s.storage.UpdateHash(id, hash)
	if err != nil {
		s.storage.RollbackTransaction(tx)
		return err
	}

	return s.storage.CommitTransaction(tx)
}

func (s *OrderService) FindOrders(ids []int) ([]models.Order, error) {

	var orders []models.Order
	for _, id := range ids {
		order, err := s.storage.GetOrderById(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) ListOrders(recipient int) ([]models.Order, error) {
	return s.storage.GetOrdersByRecipient(recipient)
}

func (s *OrderService) GetReturns(offset, limit int) ([]models.Order, error) {
	return s.storage.GetPaginatedOrdersByStatus("returned", offset, limit)
}

func (s *OrderService) DeleteOrder(id int) error {
	return s.storage.DeleteOrder(id)
}
