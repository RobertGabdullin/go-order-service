package commands

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/pkg/hash"
)

type AcceptOrder struct {
	service   service.StorageService
	recipient int
	order     int
	expire    time.Time
}

func NewAcceptOrd(service service.StorageService) AcceptOrder {
	return AcceptOrder{service: service}
}

func SetAcceptOrd(service service.StorageService, rec, ord int, st time.Time) AcceptOrder {
	return AcceptOrder{service, rec, ord, st}
}

func (AcceptOrder) GetName() string {
	return "acceptOrd"
}

func (cur AcceptOrder) Execute(mu *sync.Mutex) error {

	if cur.expire.Before(time.Now()) {
		return errors.New("storage time is out")
	}

	hash := hash.GenerateHash()

	mu.Lock()
	defer mu.Unlock()

	return cur.service.AddOrder(models.NewOrder(cur.order, cur.recipient, cur.expire, "alive", hash))

}

func (AcceptOrder) Description() string {
	return `Принять заказ от курьера. На вход принимается ID заказа (order), ID получателя (user) и срок хранения (expire). 
	     Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение выдаст ошибку.
	     Использование: acceptOrd -user=1 -order=1 -expire=2024-06-05T10`
}

func (cmd AcceptOrder) AssignArgs(m map[string]string) (Command, error) {
	if len(m) != 3 {
		return nil, errors.New("invalid number of flags")
	}

	var user, order int
	var storage time.Time
	var err error

	if elem, ok := m["user"]; ok {
		user, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for user")
		}
	} else {
		return nil, errors.New("missing user flag")
	}

	if elem, ok := m["order"]; ok {
		order, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for order")
		}
	} else {
		return nil, errors.New("missing order flag")
	}

	if elem, ok := m["expire"]; ok {
		storage, err = time.Parse("2006-01-02T15", elem)
		if err != nil {
			return nil, errors.New("invalid value for expire")
		}
	} else {
		return nil, errors.New("missing expire flag")
	}

	return SetAcceptOrd(cmd.service, user, order, storage), nil
}
