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
	weight    int
	basePrice int
	wrapper   string
}

func NewAcceptOrd(service service.StorageService) AcceptOrder {
	return AcceptOrder{service: service, wrapper: "none"}
}

func SetAcceptOrd(service service.StorageService, rec, ord, weight, basePrice int, expire time.Time, wrapper string) AcceptOrder {
	return AcceptOrder{service, rec, ord, expire, weight, basePrice, wrapper}
}

func (AcceptOrder) GetName() string {
	return "acceptOrd"
}

func (cur AcceptOrder) Execute(mu *sync.Mutex) error {
	if cur.expire.Before(time.Now()) {
		return errors.New("storage time is out")
	}

	wrapper, err := cur.service.GetWrapper(cur.wrapper)
	if err != nil {
		return err
	}

	if wrapper.MaxWeight.Valid && cur.weight > int(wrapper.MaxWeight.Int64) {
		return errors.New("order weight exceeds the maximum limit for the chosen wrapper")
	}

	hash := hash.GenerateHash()

	mu.Lock()
	defer mu.Unlock()

	return cur.service.AddOrder(models.NewOrder(cur.order, cur.recipient, cur.expire, "alive", hash, cur.weight, cur.basePrice, wrapper.Type))
}

func (AcceptOrder) Description() string {
	return `Принять заказ от курьера. На вход принимается ID заказа (order), ID получателя (user), вес товара (weight), базовая цена (basePrice), срок хранения (expire), и опционально обертка (wrapper). 
		Заказ нельзя принять дважды. Если срок хранения в прошлом, приложение выдаст ошибку.
		Использование: acceptOrd -user=1 -order=1 -weight=5 -basePrice=100 -expire=2024-06-05T10 -wrapper=pack`
}

func (cmd AcceptOrder) AssignArgs(m map[string]string) (Command, error) {
	if len(m) < 5 || len(m) > 6 {
		return nil, errors.New("invalid number of flags")
	}

	var user, order, weight, basePrice int
	var expire time.Time
	var wrapper string
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

	if elem, ok := m["weight"]; ok {
		weight, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for weight")
		}
	} else {
		return nil, errors.New("missing weight flag")
	}

	if elem, ok := m["basePrice"]; ok {
		basePrice, err = strconv.Atoi(elem)
		if err != nil {
			return nil, errors.New("invalid value for basePrice")
		}
	} else {
		return nil, errors.New("missing basePrice flag")
	}

	if elem, ok := m["expire"]; ok {
		expire, err = time.Parse("2006-01-02T15", elem)
		if err != nil {
			return nil, errors.New("invalid value for expire")
		}
	} else {
		return nil, errors.New("missing expire flag")
	}

	if elem, ok := m["wrapper"]; ok {
		wrapper = elem
	} else {
		wrapper = "none"
	}

	return SetAcceptOrd(cmd.service, user, order, weight, basePrice, expire, wrapper), nil
}
