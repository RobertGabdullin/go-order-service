package commands

import (
	"errors"
	"fmt"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type help struct {
	commands []Command
}

func NewHelp() help {
	return help{make([]Command, 0)}
}

func (cur help) GetName() string {
	return "help"
}

func (help) Description() string {
	return `Вывести описание команд`
}

func (cur help) Execute(s storage.Storage) error {
	for i, elem := range cur.commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
	return nil
}

func (help) Validate(m map[string]string) (Command, error) {
	if len(m) != 0 {
		return NewHelp(), errors.New("invalid number of flags")
	}

	ans := NewHelp()
	ans.commands = append(ans.commands, NewHelp())
	ans.commands = append(ans.commands, NewAcceptOrd())
	ans.commands = append(ans.commands, NewAcceptReturn())
	ans.commands = append(ans.commands, NewDeliverOrd())
	ans.commands = append(ans.commands, NewGetOrds())
	ans.commands = append(ans.commands, NewGetReturns())
	ans.commands = append(ans.commands, NewReturnOrd())

	return ans, nil

}
