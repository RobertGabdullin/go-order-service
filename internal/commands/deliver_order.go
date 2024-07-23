package commands

import (
	"errors"
	"strconv"
	"strings"
)

type DeliverOrder struct {
	Ords []int
}

func NewDeliverOrder(ords []int) DeliverOrder {
	return DeliverOrder{ords}
}

func (DeliverOrder) GetName() string {
	return "deliverOrd"
}

func (DeliverOrder) Description() string {
	return `Выдать заказ клиенту. На вход принимается список ID заказов (orders). 
	     Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты.
	     Все ID заказов должны принадлежать только одному клиенту.
	     Использование: deliverOrd -orders=[1,2,34]`
}

func convertToInt(in string) ([]int, error) {
	result := make([]int, 0)
	if len(in) < 2 || in[0] != '[' || in[len(in)-1] != ']' {
		return nil, errors.New("invalid number of flags")
	}

	parts := strings.Split(in[1:len(in)-1], ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		number, err := strconv.Atoi(part)
		if err != nil {
			return nil, errors.New("invalid flag value")
		}

		result = append(result, number)
	}

	return result, nil

}

func DeliverOrderAssignArgs(m map[string]string) (DeliverOrder, error) {
	if len(m) != 1 {
		return DeliverOrder{}, errors.New("invalid number of flags")
	}

	var ords []int
	var err error

	if elem, ok := m["orders"]; ok {
		ords, err = convertToInt(elem)
		if err != nil {
			return DeliverOrder{}, err
		}
	} else {
		return DeliverOrder{}, errors.New("missing orders flag")
	}

	return NewDeliverOrder(ords), nil
}
