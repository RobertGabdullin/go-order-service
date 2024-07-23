package commands

import (
	"errors"
	"strconv"
)

type GetReturns struct {
	Offset int
	Limit  int
}

func NewGetReturns(offset, limit int) GetReturns {
	return GetReturns{offset, limit}
}

func (GetReturns) GetName() string {
	return "getReturns"
}

func (GetReturns) Description() string {
	return `Получить список возвратов. Можно указать отступ (offset) и максимальное количество строк вывода (limit). Строки отсортированы по дате возврата.
	     Использование: getReturns -offset=15 -limit=30`
}

func GetReturnsAssignArgs(m map[string]string) (GetReturns, error) {
	if len(m) > 2 {
		return GetReturns{}, errors.New("invalid number of flags")
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
		return NewGetReturns(offset, limit), nil
	}
	return GetReturns{}, errors.New("invalid flag value")
}
