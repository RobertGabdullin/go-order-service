package commands

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type Command interface {
	AssignArgs(map[string]string) (Command, error)
	Execute(st service.StorageService) error
	Description() string
	GetName() string
}
