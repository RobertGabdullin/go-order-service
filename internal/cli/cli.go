package cli

import (
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

func (c CLI) Run(input string) error {

	cmd, err := c.parser.Parse(input)

	if err != nil {
		return err
	}

	errExecute := cmd.Execute(c.storage)

	if errExecute != nil {
		return errExecute
	}

	return c.storage.UpdateHash()

}
