package models

import "time"

type Order struct {
	Id          int
	Recipient   int
	Limit       time.Time
	DeliveredAt time.Time
	ReturnedAt  time.Time
	Status      string
	Hash        string
	TotalPrice  int
	Weight      int
	Wrapper     string
}

func NewOrder(id, recipient int, limit time.Time, status, hash string, price, weight int, wrapper string) Order {
	return Order{
		Id:          id,
		Recipient:   recipient,
		Limit:       limit,
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Status:      status,
		Hash:        hash,
		Weight:      weight,
		TotalPrice:  price,
		Wrapper:     wrapper,
	}
}
