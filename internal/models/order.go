package models

import "time"

type Order struct {
	Id          int
	Recipient   int
	Limit       time.Time
	DeliveredAt time.Time
	ReturnedAt  time.Time
	Status      string
}

func NewOrder(id, recipient int, limit time.Time, status string) Order {
	return Order{
		Id:          id,
		Recipient:   recipient,
		Limit:       limit,
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Status:      status,
	}
}
