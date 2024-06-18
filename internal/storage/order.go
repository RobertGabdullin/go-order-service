package storage

import "time"

type Order struct {
	Id          int
	Recipient   int
	Limit       time.Time
	DeliviredAt time.Time
	ReturnedAt  time.Time
	Status      string
}

func NewOrder(id, recipient int, limit time.Time, status string) Order {
	return Order{
		Id:          id,
		Recipient:   recipient,
		Limit:       limit,
		DeliviredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Status:      status,
	}
}
