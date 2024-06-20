package commands

import (
	"sync"
)

type Command interface {
	AssignArgs(map[string]string) (Command, error)
	Execute(*sync.Mutex) error
	Description() string
	GetName() string
}
