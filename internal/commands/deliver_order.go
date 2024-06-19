package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type deliverOrder struct {
	ords []int
}

func NewDeliverOrd() deliverOrder {
	return deliverOrder{}
}

func SetDeliverOrd(ords []int) deliverOrder {
	return deliverOrder{ords}
}

func (deliverOrder) GetName() string {
	return "deliverOrd"
}

func (cur deliverOrder) Execute(st service.StorageService) error {
	ords, err := st.FindOrders(cur.ords)

	if err != nil {
		return err
	}

	temp := make(map[int]bool)
	for _, elem := range ords {
		temp[elem.Recipient] = true
	}
	if len(temp) > 1 {
		return errors.New("list of orders should belong only to one person")
	}

	for _, elem := range ords {
		if elem.Status != "alive" {
			return errors.New("some orders are not available")
		}
		if elem.Limit.Before(time.Now()) {
			return errors.New("some orders is out of storage limit date")
		}
		tempErr := st.ChangeStatus(elem.Id, "delivered")
		if tempErr != nil {
			return tempErr
		}
	}

	return nil
}

func (deliverOrder) Description() string {
	return `Выдать заказ клиенту. На вход принимается список ID заказов (ords). 
	     Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты.
	     Все ID заказов должны принадлежать только одному клиенту.
	     Использование: deliverOrd -ords=[1,2,34]`
}

func convertToInt(in string) ([]int, error) {
	result := make([]int, 0)
	if len(in) < 2 || in[0] != '[' || in[len(in)-1] != ']' {
		return nil, errors.New("invalid flag value")
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

func (cmd deliverOrder) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 1 {
		return nil, errors.New("invalid number of flags")
	}

	var ords []int
	var err error

	if elem, ok := m["ords"]; ok {
		ords, err = convertToInt(elem)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("missing ords flag")
	}

	return SetDeliverOrd(ords), nil
}
