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

func (s *OrderService) ChangeStatus(id int, status string) error {

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return err
	}

	order, err := s.storage.GetOrderById(id)
	if err != nil {
		s.storage.RollbackTransaction(tx)
		return err
	}

	order.Status = status
	if status == "delivered" {
		order.DeliveredAt = time.Now()
	} else if status == "returned" {
		order.ReturnedAt = time.Now()
	} else {
		s.storage.RollbackTransaction(tx)
		return errors.New("unknown status")
	}

	err = s.storage.UpdateOrder(order)
	if err != nil {
		s.storage.RollbackTransaction(tx)
		return err
	}
	return s.storage.CommitTransaction(tx)
}

func (s *OrderService) FindOrders(ids []int) ([]models.Order, error) {

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	for _, id := range ids {
		order, err := s.storage.GetOrderById(id)
		if err != nil {
			s.storage.RollbackTransaction(tx)
			return nil, err
		}
		orders = append(orders, order)
	}
	s.storage.CommitTransaction(tx)
	return orders, nil
}

func (s *OrderService) ListOrders(recipient int) ([]models.Order, error) {
	return s.storage.GetOrdersByRecipient(recipient)
}

func (s *OrderService) GetReturns() ([]models.Order, error) {
	return s.storage.GetOrdersByStatus("returned")
}

func (s *OrderService) UpdateHash(hash string) error {
	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return err
	}

	exists, err := s.storage.HashExists()
	if err != nil {
		s.storage.RollbackTransaction(tx)
		return err
	}

	if exists {
		err := s.storage.UpdateHash(hash)
		if err != nil {
			s.storage.RollbackTransaction(tx)
			return err
		}
		s.storage.CommitTransaction(tx)
		return nil
	}

	err = s.storage.InsertHash(hash)
	if err != nil {
		s.storage.RollbackTransaction(tx)
		return err
	}

	s.storage.CommitTransaction(tx)
	return nil
}

func (s *OrderService) DeleteOrder(id int) error {
	return s.storage.DeleteOrder(id)
}
