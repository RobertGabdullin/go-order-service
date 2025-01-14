package models

import "time"

type Order struct {
	Id          int
	Recipient   int
	Expire      time.Time
	DeliveredAt time.Time
	ReturnedAt  time.Time
	Status      string
	Hash        string
	BasePrice   int
	Weight      int
	Wrapper     string
}

func NewOrder(id, recipient int, limit time.Time, status, hash string, price, weight int, wrapper string) Order {
	return Order{
		Id:          id,
		Recipient:   recipient,
		Expire:      limit,
		DeliveredAt: time.Now(),
		ReturnedAt:  time.Now(),
		Status:      status,
		Hash:        hash,
		Weight:      weight,
		BasePrice:   price,
		Wrapper:     wrapper,
	}
}
