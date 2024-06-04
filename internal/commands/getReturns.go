package commands

import (
	"errors"
	"fmt"
	"strconv"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type getReturns struct {
	last int
}

func NewGetReturns() getReturns {
	return getReturns{}
}

func SetGetReturns(last int) getReturns {
	return getReturns{last}
}

func (cur getReturns) GetName() string {
	return "getReturns"
}

func (cur getReturns) Execute(st storage.Storage) error {
	ords, err := st.GetReturns()
	if err != nil {
		return err
	}

	cnt := 1

	for i := len(ords) - 1; i >= 0 && len(ords)-i <= cur.last; i-- {
		fmt.Printf("%d) orderID = %d recipientID = %d storedUntil = %s\n", cnt, ords[i].Id, ords[i].Recipient, ords[i].Limit)
		cnt++
	}

	return nil
}

func (getReturns) Description() string {
	return `Получить список возвратов. По умолчанию возвращает последние 10. Может быть указан с помощью флага last.
	     Использование: getReturns --last=15`
}

func (getReturns) Validate(m map[string]string) (Command, error) {
	if len(m) > 1 {
		return NewGetReturns(), errors.New("invalid number of flags")
	}

	last := 10
	var err error
	ok := true

	for key, elem := range m {
		if key == "last" {
			last, err = strconv.Atoi(elem)
			if err != nil || last < 0 {
				ok = false
			}
		} else {
			ok = false
		}
	}

	if ok {
		return SetGetReturns(last), nil
	}
	return NewGetReturns(), errors.New("invalid flag value")
}
