package commands

import (
	"errors"
	"strconv"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type AcceptOrder struct {
	recipient int
	order     int
	storage   time.Time
}

func NewAcceptOrd() AcceptOrder {
	return AcceptOrder{}
}

func SetAcceptOrd(rec, ord int, st time.Time) AcceptOrder {
	return AcceptOrder{rec, ord, st}
}

func (cur AcceptOrder) GetName() string {
	return "acceptOrd"
}

func (cur AcceptOrder) Execute(s storage.Storage) error {

	if cur.storage.Before(time.Now()) {
		return errors.New("storage time is out")
	}

	return s.AddOrder(storage.NewOrder(cur.order, cur.recipient, cur.storage, "alive"))

}

func (AcceptOrder) Description() string {
	return `Принять заказ от курьера. На вход принимается ID заказа, ID получателя и срок хранения. 
	     Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение выдаст ошибку.
	     Использование: acceptOrd --user=1 -ord 1 -st=2024-06-05T10`
}

func (AcceptOrder) Validate(m map[string]string) (Command, error) {
	if len(m) != 3 {
		return NewAcceptOrd(), errors.New("wrong number of arguments")
	}
	var order, rec int
	var storage time.Time
	var err error
	ok := true
	for key, elem := range m {
		if key == "user" {
			rec, err = strconv.Atoi(elem)
			if err != nil {
				ok = false
			}
		} else if key == "ord" {
			order, err = strconv.Atoi(elem)
			if err != nil {
				ok = false
			}
		} else if key == "st" {
			storage, err = time.Parse("2006-01-02T15", elem)
			if err != nil {
				ok = false
			}
		} else {
			return NewAcceptOrd(), errors.New("unknown flag")
		}
	}
	if ok {
		return SetAcceptOrd(rec, order, storage), nil
	}
	return NewAcceptOrd(), errors.New("invalid flag format")
}
