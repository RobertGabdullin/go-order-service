package commands

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/pkg/hash"
)

type acceptReturn struct {
	service service.StorageService
	user    int
	order   int
}

func NewAcceptReturn(service service.StorageService) acceptReturn {
	return acceptReturn{service: service}
}

func SetAcceptReturn(service service.StorageService, user, order int) acceptReturn {
	return acceptReturn{service, user, order}
}

func (acceptReturn) GetName() string {
	return "acceptReturn"
}

func (cur acceptReturn) Execute(mu *sync.Mutex) error {

	hash := hash.GenerateHash()

	mu.Lock()
	defer mu.Unlock()

	temp := []int{cur.order}
	ords, err := cur.service.FindOrders(temp)

	if err != nil {
		return err
	}

	for _, elem := range ords {
		if elem.Id != cur.order {
			continue
		}
		if elem.Status != "delivered" {
			return errors.New("such an order has never been issued")
		}
		if elem.DeliveredAt.AddDate(0, 0, 2).Before(time.Now()) {
			return errors.New("the order can only be returned within two days after issue")
		}
		return cur.service.ChangeStatus(cur.order, "returned", hash)
	}

	return errors.New("order with such ids does not exist")
}

func (cmd acceptReturn) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 2 {
		return nil, errors.New("invalid number of flags")
	}

	var user, order int
	var err error

	if userStr, ok := m["user"]; ok {
		user, err = strconv.Atoi(userStr)
		if err != nil {
			return nil, errors.New("invalid value for user")
		}
	} else {
		return nil, errors.New("missing user flag")
	}

	if orderStr, ok := m["order"]; ok {
		order, err = strconv.Atoi(orderStr)
		if err != nil {
			return nil, errors.New("invalid value for order")
		}
	} else {
		return nil, errors.New("missing order flag")
	}

	return SetAcceptReturn(cmd.service, user, order), nil
}

func (acceptReturn) Description() string {
	return `Принять возврат от клиента. 
	     На вход принимается ID пользователя (user) и ID заказа (order). 
	     Заказ может быть возвращен в течение двух дней с момента выдачи.
	     Использование: acceptReturn -user=1 -order=1`
}
