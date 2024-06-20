package commands

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type getReturns struct {
	service service.StorageService
	offset  int
	limit   int
}

func NewGetReturns(service service.StorageService) getReturns {
	return getReturns{service: service}
}

func SetGetReturns(service service.StorageService, offset, limit int) getReturns {
	return getReturns{service, offset, limit}
}

func (getReturns) GetName() string {
	return "getReturns"
}

func (cur getReturns) Execute(mu *sync.Mutex) error {

	mu.Lock()
	ords, err := cur.service.GetReturns(cur.offset, cur.limit)
	mu.Unlock()

	if err != nil {
		return err
	}

	for i := range ords {
		fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s acceptedAt = %s\n", i+1, ords[i].Id, ords[i].Recipient, ords[i].Limit, ords[i].ReturnedAt)
	}

	return nil
}

func (getReturns) Description() string {
	return `Получить список возвратов. Можно указать отступ (offset) и максимальное количество строк вывода (limit). Строки отсортированы по дате возврата.
	     Использование: getReturns -offset=15 -limit=30`
}

func (cmd getReturns) AssignArgs(m map[string]string) (Command, error) {
	if len(m) > 2 {
		return nil, errors.New("invalid number of flags")
	}

	offset, limit := 0, -1
	var err error
	ok := true

	for key, elem := range m {
		if key == "offset" {
			offset, err = strconv.Atoi(elem)
			if err != nil || offset < 0 {
				ok = false
			}
		} else if key == "limit" {
			limit, err = strconv.Atoi(elem)
			if err != nil || limit < 0 {
				ok = false
			}
		} else {
			ok = false
		}
	}

	if ok {
		return SetGetReturns(cmd.service, offset, limit), nil
	}
	return nil, errors.New("invalid flag value")
}
