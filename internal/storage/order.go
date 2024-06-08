package storage

import "time"

type Order struct {
	Id          int
	Recipient   int
	Limit       time.Time
	DeliviredAt time.Time
	AcceptedAt  time.Time
	Status      string
}

func NewOrder(id, recipient int, limit time.Time, status string) Order {
	return Order{id, recipient, limit, time.Now(), time.Now(), status}
}
