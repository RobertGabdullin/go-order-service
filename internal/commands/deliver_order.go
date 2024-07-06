package commands

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/pkg/hash"
)

type deliverOrder struct {
	service service.StorageService
	ords    []int
}

func NewDeliverOrd(service service.StorageService) deliverOrder {
	return deliverOrder{service: service}
}

func SetDeliverOrd(service service.StorageService, ords []int) deliverOrder {
	return deliverOrder{service, ords}
}

func (deliverOrder) GetName() string {
	return "deliverOrd"
}

func (cur deliverOrder) Execute(mu *sync.Mutex) error {

	hash := hash.GenerateHash()

	mu.Lock()
	defer mu.Unlock()

	ords, err := cur.service.FindOrders(cur.ords)
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
		if elem.Expire.Before(time.Now()) {
			return errors.New("some orders is out of storage limit date")
		}
		tempErr := cur.service.ChangeStatus(elem.Id, "delivered", hash)
		if tempErr != nil {
			return tempErr
		}
	}

	return nil
}

func (deliverOrder) Description() string {
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

func (cmd deliverOrder) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 1 {
		return nil, errors.New("invalid number of flags")
	}

	var ords []int
	var err error

	if elem, ok := m["orders"]; ok {
		ords, err = convertToInt(elem)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("missing orders flag")
	}

	return SetDeliverOrd(cmd.service, ords), nil
}
