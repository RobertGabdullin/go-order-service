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
	limit     time.Time
}

func NewAcceptOrd() AcceptOrder {
	return AcceptOrder{}
}

func SetAcceptOrd(rec, ord int, st time.Time) AcceptOrder {
	return AcceptOrder{rec, ord, st}
}

func (AcceptOrder) GetName() string {
	return "acceptOrd"
}

func (cur AcceptOrder) Execute(s storage.Storage) error {

	if cur.limit.Before(time.Now()) {
		return errors.New("storage time is out")
	}

	return s.AddOrder(storage.NewOrder(cur.order, cur.recipient, cur.limit, "alive"))

}

func (AcceptOrder) Description() string {
	return `Принять заказ от курьера. На вход принимается ID заказа (ord), ID получателя (user) и срок хранения (lim). 
	     Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение выдаст ошибку.
	     Использование: acceptOrd -user=1 -ord=1 -st=2024-06-05T10`
}

func (cmd AcceptOrder) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 3 {
		return nil, errors.New("invalid number of flags")
	}

	var user, order int
	var storage time.Time
	var err error

	if elem, ok := m["user"]; ok {
		user, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for user")
		}
	} else {
		return nil, errors.New("missing user flag")
	}

	if elem, ok := m["ord"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for ord")
		}
	} else {
		return nil, errors.New("missing ord flag")
	}

	if elem, ok := m["lim"]; ok {
		storage, err = time.Parse("2006-01-02T15", elem)
		if err != nil {
			return nil, errors.New("invalid value for lim")
		}
	} else {
		return nil, errors.New("missing st flag")
	}

	return SetAcceptOrd(user, order, storage), nil
}
