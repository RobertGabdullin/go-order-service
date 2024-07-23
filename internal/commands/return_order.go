package commands

import (
	"errors"
	"strconv"
)

type ReturnOrder struct {
	Order int
}

func NewReturnOrder(order int) ReturnOrder {
	return ReturnOrder{order}
}

func (ReturnOrder) GetName() string {
	return "returnOrd"
}

func (ReturnOrder) Description() string {
	return `Вернуть заказ курьеру. На вход принимается ID заказа (order). Метод должен удалять заказ из вашего файла.
	     Можно вернуть только те заказы, у которых вышел срок хранения и если заказы не были выданы клиенту.
	     Использование: returnOrd -order=1`
}

func ReturnOrderAssignArgs(m map[string]string) (ReturnOrder, error) {
	if len(m) != 1 {
		return ReturnOrder{}, errors.New("invalid number of flags")
	}

	var order int
	var err error

	if elem, ok := m["order"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil || order < 0 {
			return ReturnOrder{}, errors.New("invalid value for order")
		}
	} else {
		return ReturnOrder{}, errors.New("invalid flag name")
	}

	return NewReturnOrder(order), nil
}
