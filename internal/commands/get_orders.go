package commands

import (
	"errors"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type getOrders struct {
	user  int
	count int
}

func NewGetOrds() getOrders {
	return getOrders{}
}

func (getOrders) Description() string {
	return `Получить список заказов. На вход принимается ID пользователя (user) как обязательный параметр.
	     Также можно указать опциональный параметр (count), который позволяет получить только последние N заказов.
	     Использование: getOrds -user=1 -count=5`
}

func SetGetOrds(user, count int) getOrders {
	return getOrders{user, count}
}

func (getOrders) GetName() string {
	return "getOrds"
}

func (cur getOrders) Execute(st service.StorageService) error {
	ords, err := st.ListOrders(cur.user)

	if err != nil {
		return err
	}

	cnt := 1

	for i := len(ords) - 1; i >= 0 && (cur.count == -1 || cur.count >= cnt); i-- {
		if ords[i].Status == "alive" {
			fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s\n", cnt, ords[i].Id, ords[i].Recipient, ords[i].Limit)
			cnt++
		}
	}

	return nil
}

func (cmd getOrders) AssignArgs(m map[string]string) (Command, error) {
	if len(m) < 1 || len(m) > 2 {
		return nil, errors.New("invalid number of flags")
	}

	user, count := 0, -1
	var err error

	if userStr, ok := m["user"]; ok {
		user, err = strconv.Atoi(userStr)
		if err != nil {
			return nil, errors.New("invalid value for user")
		}
	} else {
		return nil, errors.New("missing user flag")
	}

	if countStr, ok := m["count"]; ok {
		count, err = strconv.Atoi(countStr)
		if err != nil || count < 0 {
			return nil, errors.New("invalid value for count")
		}
	}

	return SetGetOrds(user, count), nil
}
