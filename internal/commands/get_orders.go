package commands

import (
	"errors"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type getOrders struct {
	user  int
	count int
}

func NewGetOrds() getOrders {
	return getOrders{}
}

func (getOrders) Description() string {
	return `Получить список заказов. На вход принимается ID пользователя как обязательный параметр и опциональные параметры.
	     Параметры позволяют получать только последние N заказов или заказы клиента, находящиеся в нашем ПВЗ
	     Использование: getOrds -user=1 -count=5`
}

func SetGetOrds(user, count int) getOrders {
	return getOrders{user, count}
}

func (cur getOrders) GetName() string {
	return "getOrds"
}

func (cur getOrders) Execute(st storage.Storage) error {
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

func (getOrders) Validate(m map[string]string) (Command, error) {
	if len(m) < 1 || len(m) > 2 {
		return NewGetOrds(), errors.New("invalid number of flags")
	}

	user, count := 0, -1
	ok := true
	var err error

	for key, elem := range m {
		if key == "user" {
			user, err = strconv.Atoi(elem)
			if err != nil {
				ok = false
			}
		} else if key == "count" {
			count, err = strconv.Atoi(elem)
			if err != nil || count < 0 {
				ok = false
			}
		} else {
			ok = false
		}
	}

	if ok {
		return SetGetOrds(user, count), nil
	}

	return NewGetOrds(), errors.New("invalid flag")
}
