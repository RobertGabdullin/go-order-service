package commands

import (
	"errors"
	"strconv"
)

type AcceptReturn struct {
	User  int
	Order int
}

func NewAcceptReturn(user, order int) AcceptReturn {
	return AcceptReturn{user, order}
}

func (AcceptReturn) GetName() string {
	return "acceptReturn"
}

func (AcceptReturn) Description() string {
	return `Принять возврат от клиента. 
		 На вход принимается ID пользователя (user) и ID заказа (order). 
		 Заказ может быть возвращен в течение двух дней с момента выдачи.
		 Использование: acceptReturn -user=1 -order=1`
}

func AcceptReturnAssignArgs(m map[string]string) (AcceptReturn, error) {
	if len(m) != 2 {
		return AcceptReturn{}, errors.New("invalid number of flags")
	}

	var user, order int
	var err error

	if userStr, ok := m["user"]; ok {
		user, err = strconv.Atoi(userStr)
		if err != nil {
			return AcceptReturn{}, errors.New("invalid value for user")
		}
	} else {
		return AcceptReturn{}, errors.New("missing user flag")
	}

	if orderStr, ok := m["order"]; ok {
		order, err = strconv.Atoi(orderStr)
		if err != nil {
			return AcceptReturn{}, errors.New("invalid value for order")
		}
	} else {
		return AcceptReturn{}, errors.New("missing order flag")
	}

	return NewAcceptReturn(user, order), nil
}
