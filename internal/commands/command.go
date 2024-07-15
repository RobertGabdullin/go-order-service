package commands

import (
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type Command interface {
	AssignArgs(map[string]string) (Command, error)
	Execute(*sync.Mutex) ([]models.Order, error)
	Description() string
	GetName() string
}
