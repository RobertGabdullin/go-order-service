package cli

import (
	"errors"
	"fmt"
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
	"gitlab.ozon.dev/r_gabdullin/homework-1/pkg/hash"
)

type CLI struct {
	storage storage.Storage
	parser  parser.Parser
	mu      *sync.Mutex
}

func NewCLI(storage storage.Storage, parser parser.Parser) CLI {
	return CLI{storage, parser, new(sync.Mutex)}
}

func (c CLI) Help() {
	fmt.Println("Утилита для управления ПВЗ. Для аргументов команд можно использовать следующие форматы: -word=x --word=x -word x --word x. В примерах будет использован только формат -word=x. Список команд:")
	commands := parser.GetCommands()
	for i, elem := range commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
}

func Find(cmd string) (commands.Command, error) {
	listCmd := parser.GetCommands()
	for _, elem := range listCmd {
		if elem.GetName() == cmd {
			return elem, nil
		}
	}
	return nil, errors.New("unknown command")
}

func (c CLI) Run(input string) error {

	cmdName, mapArgs, err := c.parser.Parse(input)
	if err != nil {
		return err
	}

	cmd, errFind := Find(cmdName)
	if errFind != nil {
		return errFind
	}

	cmd, errAssign := cmd.AssignArgs(mapArgs)
	if errAssign != nil {
		return errAssign
	}

	hash := hash.GenerateHash()

	c.mu.Lock()
	defer c.mu.Unlock()

	err = cmd.Execute(c.storage)
	if err != nil {
		return err
	}

	return c.storage.UpdateHash(hash)

}
