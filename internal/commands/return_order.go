package commands

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type returnOrders struct {
	service service.StorageService
	order   int
}

func NewReturnOrd(service service.StorageService) returnOrders {
	return returnOrders{service: service}
}

func SetReturnOrd(service service.StorageService, order int) returnOrders {
	return returnOrders{service, order}
}

func (returnOrders) GetName() string {
	return "returnOrd"
}

func (returnOrders) Description() string {
	return `Вернуть заказ курьеру. На вход принимается ID заказа (order). Метод должен удалять заказ из вашего файла.
	     Можно вернуть только те заказы, у которых вышел срок хранения и если заказы не были выданы клиенту.
	     Использование: returnOrd -order=1`
}

func (cur returnOrders) Execute(mu *sync.Mutex) ([]models.Order, error) {

	mu.Lock()
	defer mu.Unlock()

	temp := []int{cur.order}
	ords, err := cur.service.FindOrders(temp)

	if err != nil {
		return nil, err
	}

	if len(ords) == 0 {
		return nil, errors.New("such order does not exist")
	}

	if ords[0].Status != "alive" && ords[0].Status != "returned" {
		return nil, errors.New("order is not at storage")
	}

	if ords[0].Expire.After(time.Now()) {
		return nil, errors.New("order should be out of storage limit")
	}

	err = cur.service.DeleteOrder(cur.order)
	if err != nil {
		return nil, err
	}

	return []models.Order{}, nil
}

func (cmd returnOrders) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 1 {
		return nil, errors.New("invalid number of flags")
	}

	var order int
	var err error

	if elem, ok := m["order"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil || order < 0 {
			return nil, errors.New("invalid value for order")
		}
	} else {
		return nil, errors.New("invalid flag name")
	}

	return SetReturnOrd(cmd.service, order), nil
}
