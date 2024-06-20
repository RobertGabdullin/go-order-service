package cli

import (
	"errors"
	"fmt"
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
)

type CLI struct {
	parser   parser.Parser
	mu       *sync.Mutex
	commands []commands.Command
}

func NewCLI(storage service.StorageService, parser parser.Parser) CLI {
	return CLI{
		parser: parser,
		mu:     new(sync.Mutex),
		commands: []commands.Command{
			commands.NewAcceptOrd(storage),
			commands.NewAcceptReturn(storage),
			commands.NewDeliverOrd(storage),
			commands.NewGetOrds(storage),
			commands.NewGetReturns(storage),
			commands.NewReturnOrd(storage),
		},
	}
}

func (c CLI) Help() {
	fmt.Println("Утилита для управления ПВЗ. Для аргументов команд можно использовать следующие форматы: -word=x --word=x -word x --word x. В примерах будет использован только формат -word=x. Список команд:")
	commands := c.commands
	for i, elem := range commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
}

func (c CLI) Find(cmd string) (commands.Command, error) {
	listCmd := c.commands
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

	cmd, errFind := c.Find(cmdName)
	if errFind != nil {
		return errFind
	}

	cmd, errAssign := cmd.AssignArgs(mapArgs)
	if errAssign != nil {
		return errAssign
	}

	return cmd.Execute(c.mu)
}
