package commands

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type Command interface {
	Validate(map[string]string) (Command, error)
	Execute(st storage.Storage) error
	Description() string
	GetName() string
}
