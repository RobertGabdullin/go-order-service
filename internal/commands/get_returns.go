package commands

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type getReturns struct {
	offset int
	limit  int
}

func NewGetReturns() getReturns {
	return getReturns{}
}

func SetGetReturns(offset, limit int) getReturns {
	return getReturns{offset, limit}
}

func (getReturns) GetName() string {
	return "getReturns"
}

func (cur getReturns) Execute(st service.StorageService) error {
	ords, err := st.GetReturns()
	if err != nil {
		return err
	}

	sort.Slice(ords, func(i, j int) bool {
		return ords[i].ReturnedAt.Before(ords[j].ReturnedAt)
	})

	cnt := 1
	if cur.limit == -1 {
		cur.limit = len(ords)
	}

	for i := cur.offset; i < len(ords) && cnt <= cur.limit; i++ {
		fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s acceptedAt = %s\n", cnt, ords[i].Id, ords[i].Recipient, ords[i].Limit, ords[i].ReturnedAt)
		cnt++
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
		return SetGetReturns(offset, limit), nil
	}
	return nil, errors.New("invalid flag value")
}
