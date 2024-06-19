package commands

import (
	"errors"
	"strconv"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type acceptReturn struct {
	user  int
	order int
}

func NewAcceptReturn() acceptReturn {
	return acceptReturn{}
}

func SetAcceptReturn(user, order int) acceptReturn {
	return acceptReturn{user, order}
}

func (acceptReturn) GetName() string {
	return "acceptReturn"
}

func (cur acceptReturn) Execute(s service.StorageService) error {
	temp := make([]int, 0)
	temp = append(temp, cur.order)
	ords, err := s.FindOrders(temp)

	if err != nil {
		return err
	}

	for _, elem := range ords {
		if elem.Id == cur.order {
			if elem.Status != "delivered" {
				return errors.New("such an order has never been issued")
			} else if elem.DeliveredAt.AddDate(0, 0, 2).Before(time.Now()) {
				return errors.New("the order can only be returned within two days after issue")
			} else {
				return s.ChangeStatus(cur.order, "returned")
			}

		}
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

	if orderStr, ok := m["ord"]; ok {
		order, err = strconv.Atoi(orderStr)
		if err != nil {
			return nil, errors.New("invalid value for order")
		}
	} else {
		return nil, errors.New("missing ord flag")
	}

	return SetAcceptReturn(user, order), nil
}

func (acceptReturn) Description() string {
	return `Принять возврат от клиента. 
	     На вход принимается ID пользователя (user) и ID заказа (ord). 
	     Заказ может быть возвращен в течение двух дней с момента выдачи.
	     Использование: acceptReturn -user=1 -ord=1`
}
