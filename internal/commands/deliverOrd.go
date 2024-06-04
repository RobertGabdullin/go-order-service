package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type deliverOrd struct {
	ords []int
}

func NewDeliverOrd() deliverOrd {
	return deliverOrd{}
}

func SetDeliverOrd(ords []int) deliverOrd {
	return deliverOrd{ords}
}

func (cur deliverOrd) GetName() string {
	return "deliverOrd"
}

func (cur deliverOrd) Execute(st storage.Storage) error {
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

func (deliverOrd) Description() string {
	return `Выдать заказ клиенту. На вход принимается список ID заказов. 
	     Можно выдавать только те заказы, которые были приняты от курьера и чей срок хранения меньше текущей даты.
	     Все ID заказов должны принадлежать только одному клиенту.
	     Использование: deliverOrd --ords=[1,2,34]`
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

func (deliverOrd) Validate(m map[string]string) (Command, error) {
	if len(m) != 1 {
		return NewDeliverOrd(), errors.New("invalid number of arguments")
	}

	var ords []int
	var err error

	for key, elem := range m {
		if key == "ords" {
			ords, err = convertToInt(elem)
			if err != nil {
				return NewDeliverOrd(), err
			}
		} else {
			return NewDeliverOrd(), errors.New("invalid flag name")
		}
	}

	return SetDeliverOrd(ords), nil

}
