package commands

import (
	"errors"
	"strconv"
	"time"
)

type AcceptOrder struct {
	Recipient int
	Order     int
	Expire    time.Time
	Weight    int
	BasePrice int
	Wrapper   string
}

func NewAcceptOrder(rec, ord, weight, basePrice int, expire time.Time, wrapper string) AcceptOrder {
	return AcceptOrder{rec, ord, expire, weight, basePrice, wrapper}
}

func (AcceptOrder) GetName() string {
	return "acceptOrd"
}

func (AcceptOrder) Description() string {
	return `Принять заказ от курьера. На вход принимается ID заказа (order), ID получателя (user), вес товара (weight), базовая цена (basePrice), срок хранения (expire), и опционально обертка (wrapper). 
		Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение выдаст ошибку.
		Использование: acceptOrd -user=1 -order=1 -weight=5 -basePrice=100 -expire=2024-06-05T10 -wrapper=pack`
}

func AcceptOrderAssignArgs(m map[string]string) (AcceptOrder, error) {
	if len(m) < 5 || len(m) > 6 {
		return AcceptOrder{}, errors.New("invalid number of flags")
	}

	var user, order, weight, basePrice int
	var expire time.Time
	var wrapper string
	var err error

	if elem, ok := m["user"]; ok {
		user, err = strconv.Atoi(elem)
		if err != nil {
			return AcceptOrder{}, errors.New("invalid value for user")
		}
	} else {
		return AcceptOrder{}, errors.New("missing user flag")
	}

	if elem, ok := m["order"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil {
			return AcceptOrder{}, errors.New("invalid value for order")
		}
	} else {
		return AcceptOrder{}, errors.New("missing order flag")
	}

	if elem, ok := m["weight"]; ok {
		weight, err = strconv.Atoi(elem)
		if err != nil {
			return AcceptOrder{}, errors.New("invalid value for weight")
		}
	} else {
		return AcceptOrder{}, errors.New("missing weight flag")
	}

	if elem, ok := m["basePrice"]; ok {
		basePrice, err = strconv.Atoi(elem)
		if err != nil {
			return AcceptOrder{}, errors.New("invalid value for basePrice")
		}
	} else {
		return AcceptOrder{}, errors.New("missing basePrice flag")
	}

	if elem, ok := m["expire"]; ok {
		expire, err = time.Parse("2006-01-02T15", elem)
		if err != nil {
			return AcceptOrder{}, errors.New("invalid value for expire")
		}
	} else {
		return AcceptOrder{}, errors.New("missing expire flag")
	}

	if elem, ok := m["wrapper"]; ok {
		wrapper = elem
	} else {
		wrapper = "none"
	}

	return NewAcceptOrder(user, order, weight, basePrice, expire, wrapper), nil
}
