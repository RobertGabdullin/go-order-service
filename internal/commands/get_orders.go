package commands

import (
	"errors"
	"strconv"
)

type GetOrders struct {
	User  int
	Count int
}

func (GetOrders) GetName() string {
	return "getOrds"
}

func (GetOrders) Description() string {
	return `Получить список заказов. На вход принимается ID пользователя (user) как обязательный параметр.
	     Также можно указать опциональный параметр (count), который позволяет получить только последние N заказов.
	     Использование: getOrds -user=1 -count=5`
}

func NewGetOrders(user, count int) GetOrders {
	return GetOrders{user, count}
}

func GetOrdersAssignArgs(m map[string]string) (GetOrders, error) {
	if len(m) < 1 || len(m) > 2 {
		return GetOrders{}, errors.New("invalid number of flags")
	}

	user, count := 0, -1
	var err error

	if userStr, ok := m["user"]; ok {
		user, err = strconv.Atoi(userStr)
		if err != nil {
			return GetOrders{}, errors.New("invalid value for user")
		}
	} else {
		return GetOrders{}, errors.New("missing user flag")
	}

	if countStr, ok := m["count"]; ok {
		count, err = strconv.Atoi(countStr)
		if err != nil || count < 0 {
			return GetOrders{}, errors.New("invalid value for count")
		}
	}

	return NewGetOrders(user, count), nil
}
