package commands

import (
	"errors"
	"strconv"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
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

func (cur acceptReturn) GetName() string {
	return "acceptReturn"
}

func (cur acceptReturn) Execute(s storage.Storage) error {
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
			} else if elem.DeliviredAt.AddDate(0, 0, 2).Before(time.Now()) {
				return errors.New("the order can only be returned within two days after issue")
			} else {
				return s.ChangeStatus(cur.order, "returned")
			}

		}
	}
	return errors.New("order with such ids does not exist")
}

func (acceptReturn) Validate(m map[string]string) (Command, error) {
	if len(m) != 2 {
		return NewAcceptReturn(), errors.New("invalid number of flags")
	}
	var user, order int
	var err error
	ok := true
	for key, elem := range m {
		if key == "user" {
			user, err = strconv.Atoi(elem)
			if err != nil {
				ok = false
			}
		} else if key == "ord" {
			order, err = strconv.Atoi(elem)
			if err != nil {
				ok = false
			}
		} else {
			ok = false
		}
	}
	if ok {
		return SetAcceptReturn(user, order), nil
	}
	return NewAcceptReturn(), errors.New("invalid flag value")
}

func (acceptReturn) Description() string {
	return `Принять возврат от клиента. 
	     На вход принимается ID пользователя (user) и ID заказа (ord). 
	     Заказ может быть возвращен в течение двух дней с момента выдачи.
	     Использование: acceptReturn --user=1 --ord 1`
}
