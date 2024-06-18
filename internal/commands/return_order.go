package commands

import (
	"errors"
	"strconv"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type returnOrders struct {
	order int
}

func NewReturnOrd() returnOrders {
	return returnOrders{}
}

func SetReturnOrd(order int) returnOrders {
	return returnOrders{order}
}

func (returnOrders) GetName() string {
	return "returnOrd"
}

func (returnOrders) Description() string {
	return `Вернуть заказ курьеру. На вход принимается ID заказа. Метод должен удалять заказ из вашего файла.
	     Можно вернуть только те заказы, у которых вышел срок хранения и если заказы не были выданы клиенту.
	     Использование: returnOrd -ord=1`
}

func (cur returnOrders) Execute(st storage.Storage) error {

	temp := make([]int, 0)
	temp = append(temp, cur.order)

	ords, err := st.FindOrders(temp)

	if err != nil {
		return err
	}

	if len(ords) == 0 {
		return errors.New("such order does not exist")
	}

	if ords[0].Status != "alive" && ords[0].Status != "returned" {
		return errors.New("order is not at storage")
	}

	if ords[0].Limit.After(time.Now()) {
		return errors.New("order should be out of storage limit")
	}

	return st.ChangeStatus(ords[0].Id, "deleted")
}

func (cmd returnOrders) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 1 {
		return nil, errors.New("invalid number of flags")
	}

	var order int
	var err error

	if elem, ok := m["ord"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil || order < 0 {
			return nil, errors.New("invalid value for ord")
		}
	} else {
		return nil, errors.New("invalid flag name")
	}

	return SetReturnOrd(order), nil
}
