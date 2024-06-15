package cli

import (
	"errors"
	"fmt"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/parser"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
)

type CLI struct {
	storage storage.Storage
	parser  parser.Parser
}

func NewCLI(storage storage.Storage, parser parser.Parser) CLI {
	return CLI{storage, parser}
}

func (c CLI) Help() {
	fmt.Println("Утилита для управления ПВЗ. Для аргументов команд можно использовать следующие форматы: -word=x --word=x -word x --word x. В примерах будет использован только формат -word=x. Список команд:")
	commands := parser.GetCommands()
	for i, elem := range commands {
		fmt.Printf("%d) Команда: %s\n   Описание: %s\n", i+1, elem.GetName(), elem.Description())
	}
}

func validate(cmd string, args map[string]string) (commands.Command, error) {
	listCmd := parser.GetCommands()
	for _, elem := range listCmd {
		if elem.GetName() == cmd {
			return elem.Validate(args)
		}
	}
	return nil, errors.New("unknown command")
}

func (c CLI) Run(input string) error {

	cmdName, mapArgs, err := c.parser.Parse(input)
	if err != nil {
		return err
	}

	cmd, errValidate := validate(cmdName, mapArgs)
	if errValidate != nil {
		return errValidate
	}

	err = cmd.Execute(c.storage)
	if err != nil {
		return err
	}

	return c.storage.UpdateHash()

}
